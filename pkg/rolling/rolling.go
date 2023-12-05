package rolling

import (
	"fmt"
	"strings"
	"time"

	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/cms"
)

type Rolling struct {
	logger *zap.SugaredLogger
	cms    *cms.CMSClient
	opts   *Options

	// internal state available during restart
	service Service
	nodes   []*Ydb_Maintenance.Node
	tenants []string
}

const (
	RestartTaskPrefix = "rolling_restart"
	RestartTaskUid    = RestartTaskPrefix + "_001"
)

func New(cms *cms.CMSClient, logger *zap.SugaredLogger, opts *Options) *Rolling {
	return &Rolling{
		cms:    cms,
		logger: logger,
		opts:   opts,
	}
}

func (r *Rolling) Restart() error {
	if err := r.prepareState(); err != nil {
		return err
	}

	nodesToRestart := r.service.Filter(
		FilterNodeParams{
			Service:         r.opts.Service,
			AllTenants:      r.tenants,
			AllNodes:        r.nodes,
			SelectedTenants: r.opts.Tenants,
			SelectedNodeIds: r.opts.Nodes,
		},
	)
	taskParams := cms.MaintenanceTaskParams{
		TaskUid:          RestartTaskUid,
		AvailAbilityMode: r.opts.AvailabilityMode(),
		Duration:         r.opts.RestartDuration(),
		Nodes:            nodesToRestart,
	}
	task, err := r.cms.CreateMaintenanceTask(taskParams)
	if err != nil {
		return fmt.Errorf("failed to create maintenance task: %+v", err)
	}

	r.logger.Infof("Maintenance task id: %s", task.GetTaskUid())
	return r.loop(task)
}

func (r *Rolling) RestartPrevious() error {
	if err := r.prepareState(); err != nil {
		return err
	}

	result, err := r.cms.GetMaintenanceTask(RestartTaskUid)
	if err != nil {
		return fmt.Errorf("failed to get maintenance task with id: %s, err: %+v", RestartTaskUid, err)
	}

	return r.loop(result)
}

func (r *Rolling) loop(task cms.MaintenanceTask) error {
	const (
		defaultDelay = time.Second * 30
	)

	var (
		err    error
		delay  time.Duration
		taskId = task.GetTaskUid()
	)

	r.logger.Infof("Maintenance task processing loop started")
	for {
		delay = defaultDelay

		if task != nil {
			r.logTask(task)

			if task.GetRetryAfter() != nil {
				retryTime := task.GetRetryAfter().AsTime()
				r.logger.Debugf("Task has retry after attribute: %s", retryTime.Format(time.DateTime))

				if retryDelay := retryTime.Sub(time.Now().UTC()); defaultDelay < retryDelay {
					delay = defaultDelay
				}
			}

			r.logger.Info("Processing task action group states")
			if completed := r.process(task); completed {
				break
			}
		}

		r.logger.Infof("Wait next %s delay", delay)
		time.Sleep(delay)

		r.logger.Infof("Refresh maintenance task with id: %s", taskId)
		task, err = r.cms.RefreshMaintenanceTask(taskId)
		if err != nil {
			r.logger.Warnf("Failed to refresh maintenance task: %+v", err)
		}
	}

	r.logger.Infof("Maintenance task processing loop completed")
	return nil
}

func (r *Rolling) process(task cms.MaintenanceTask) bool {
	performed := util.FilterBy(task.GetActionGroupStates(),
		func(gs *Ydb_Maintenance.ActionGroupStates) bool {
			return gs.ActionStates[0].Status == Ydb_Maintenance.ActionState_ACTION_STATUS_PERFORMED
		},
	)
	ids := make([]*Ydb_Maintenance.ActionUid, 0, len(performed))

	for _, gs := range performed {
		var (
			as     = gs.ActionStates[0]
			lock   = as.Action.GetLockAction()
			nodeId = lock.Scope.GetNodeId()
		)

		nodes := util.FilterBy(r.nodes,
			func(node *Ydb_Maintenance.Node) bool {
				return node.NodeId == nodeId
			},
		)

		// todo: drain node

		if err := r.service.RestartNode(nodes[0]); err != nil {
			// todo: failed to restart node ?
		}

		ids = append(ids, as.ActionUid)
	}

	result, err := r.cms.CompleteAction(ids)
	if err != nil {
		r.logger.Warnf("Failed to complete action: %+v", err)
		return false
	}
	r.logCompleteResult(result)

	// completed when all actions marked as completed
	return len(task.GetActionGroupStates()) == len(result.ActionStatuses)
}

func (r *Rolling) prepareState() error {
	service, err := r.createService()
	if err != nil {
		return err
	}

	tenants, err := r.cms.Tenants()
	if err != nil {
		return fmt.Errorf("failed to list avaialble tenants: %+v", err)
	}

	nodes, err := r.cms.Nodes()
	if err != nil {
		return fmt.Errorf("failed to list avaialble nodes: %+v", err)
	}

	r.service = service
	r.nodes = nodes
	r.tenants = tenants

	return nil
}

func (r *Rolling) createService() (Service, error) {
	factory := ServiceFactoryMap[r.opts.Service]
	service, err := factory(ServiceOptionsMap[r.opts.Service])
	if err != nil {
		return nil, err
	}

	if err := service.Prepare(); err != nil {
		return nil, err
	}

	return service, nil
}

func (r *Rolling) logTask(task cms.MaintenanceTask) {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Uid: %s\n", task.GetTaskUid()))

	if task.GetRetryAfter() != nil {
		sb.WriteString(fmt.Sprintf("Retry after: %s\n", task.GetRetryAfter().AsTime().Format(time.DateTime)))
	}

	for _, gs := range task.GetActionGroupStates() {
		as := gs.ActionStates[0]
		sb.WriteString(fmt.Sprintf("  Lock on node %d ", as.Action.GetLockAction().Scope.GetNodeId()))
		if as.Status == Ydb_Maintenance.ActionState_ACTION_STATUS_PERFORMED {
			sb.WriteString(fmt.Sprintf("PERFORMED, until: %s", as.Deadline.AsTime().Format(time.DateTime)))
		} else {
			sb.WriteString("PENDING")
		}
		sb.WriteString("\n")
	}
	r.logger.Debugf("Maintenance task result:\n%s", sb.String())
}

func (r *Rolling) logCompleteResult(result *Ydb_Maintenance.ManageActionResult) {
	if result == nil {
		return
	}

	sb := strings.Builder{}

	for _, status := range result.ActionStatuses {
		sb.WriteString(fmt.Sprintf("  Action: %s, status: %s", status.ActionUid, status.Status))
	}

	r.logger.Debugf("Manage action result:\n%s", sb.String())
}
