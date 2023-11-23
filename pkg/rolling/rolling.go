package rolling

import (
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-rolling-restart/pkg/cms"
)

type Rolling struct {
	logger *zap.SugaredLogger
	cms    *cms.CMSClient
}

func New(cms *cms.CMSClient, logger *zap.SugaredLogger) *Rolling {
	return &Rolling{
		cms:    cms,
		logger: logger,
	}
}

func (r *Rolling) Restart() error {

	// todo:
	//  1. filter out nodes by tenant/specified nodes
	//  2. invoke CreateMaintenanceTask
	//  3. run loop

	return nil
}

func (r *Rolling) RestartPrevious() error {

	// todo:
	//  1. find previous restart
	//  2. run loop

	return nil
}

func (r *Rolling) loop() {
	for r.next() {
		// todo: sleep if required
	}
}

func (r *Rolling) next() bool {
	return true
}
