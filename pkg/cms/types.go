package cms

import (
	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"
	"google.golang.org/protobuf/types/known/durationpb"
)

type MaintenanceTaskParams struct {
	TaskUid          string
	AvailAbilityMode Ydb_Maintenance.AvailabilityMode
	Duration         *durationpb.Duration // todo: duration (?)
	Nodes            []*Ydb_Maintenance.Node
}
