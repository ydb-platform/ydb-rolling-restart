package rolling

import (
	"fmt"
	"strings"
	"time"

	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/pkg/cms"
)

type Rolling struct {
	logger *zap.SugaredLogger
	cms    *cms.CMSClient
	opts   *Options
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

	nodesToRestart := service.Filter(
		FilterNodeParams{
			Service:         r.opts.Service,
			AllTenants:      tenants,
			AllNodes:        nodes,
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

	r.logger.Infof("Maintenance task id: %s", task.TaskUid)
	return r.loop(task)
}

func (r *Rolling) RestartPrevious() error {
	service, err := r.createService()
	if err != nil {
		return err
	}

	_ = service
	// todo:
	//  1. find previous restart
	//  2. run loop

	return nil
}

func (r *Rolling) loop(result *Ydb_Maintenance.MaintenanceTaskResult) error {
	const (
		defaultDelay = time.Second * 30
	)

	var (
		delay = defaultDelay
		err   error
	)

	for {
		if result != nil {
			r.logResult(result)
		}

		// action can be performed
		if result == nil || result.RetryAfter != nil {
			// calculate delay relative to current time
			delay = result.RetryAfter.AsTime().Sub(time.Now().UTC())
			if defaultDelay < delay {
				delay = defaultDelay
			}
		} else {
			// process action groups & use default delay
			ok := r.next(result)
			if !ok {
				r.logger.Infof("Processing completed")
				break
			}
			delay = defaultDelay
		}

		r.logger.Infof("Wait for delay: %s", delay)
		time.Sleep(delay)

		r.logger.Infof("Refresh maintenance task")
		result, err = r.cms.RefreshMaintenanceTask(result.TaskUid)
		if err != nil {
			r.logger.Warnf("Failed to refresh maintenance task: %+v", err)
		}
	}

	return nil
}

func (r *Rolling) next(result *Ydb_Maintenance.MaintenanceTaskResult) bool {
	// todo: get available actions from result & perform restart on each node

	return true
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

func (r *Rolling) logResult(result *Ydb_Maintenance.MaintenanceTaskResult) {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Uid: %s\n", result.TaskUid))

	if result.RetryAfter != nil {
		sb.WriteString(fmt.Sprintf("Retry after: %s\n", result.RetryAfter.AsTime().Format(time.DateTime)))
	}

	for _, gs := range result.ActionGroupStates {
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
