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
	service, err := r.service()
	if err != nil {
		return err
	}

	if err := service.Prepare(); err != nil {
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

	// todo:
	//  1. find previous restart
	//  2. run loop

	return nil
}

func (r *Rolling) loop(result *Ydb_Maintenance.MaintenanceTaskResult) error {
	for r.next(result) {
		// todo: sleep if required
	}

	return nil
}

func (r *Rolling) next(result *Ydb_Maintenance.MaintenanceTaskResult) bool {
	const (
		delay = time.Second * 30
	)

	r.logResult(result)
	r.logger.Infof("Waiting locks for %s:", delay)
	time.Sleep(delay)
	return true
}

func (r *Rolling) service() (Service, error) {
	factory := ServiceFactoryMap[r.opts.Service]
	service, err := factory(ServiceOptionsMap[r.opts.Service])
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (r *Rolling) logResult(result *Ydb_Maintenance.MaintenanceTaskResult) {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Uid: %s\n", result.TaskUid))

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
