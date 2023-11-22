package cms

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ydb-platform/ydb-go-genproto/protos/Ydb_Operations"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/durationpb"

	_ "github.com/ydb-platform/ydb-go-genproto/draft/Ydb_Maintenance_V1"

	"github.com/ydb-platform/ydb-rolling-restart/pkg/options"
)

const (
	BufferSize = 32 << 20
)

type Factory struct {
	grpc options.GRPC
	auth options.CMSAuth
	user string
}

func NewConnectionFactory(
	grpc options.GRPC,
	auth options.CMSAuth,
	user string,
) *Factory {
	return &Factory{
		grpc: grpc,
		auth: auth,
		user: user,
	}
}

func (f Factory) Context() (context.Context, context.CancelFunc) {
	ctx, cf := context.WithTimeout(context.Background(), time.Second*time.Duration(f.grpc.TimeoutSeconds))

	t, err := f.auth.Token()
	if err != nil {
		zap.S().Warnf("Failed to load auth token: %v", err)
		return ctx, cf
	}

	return metadata.AppendToOutgoingContext(ctx,
		"x-ydb-auth-ticket", t.Secret,
		"authorization", t.Token()), cf
}

func (f Factory) OperationParams() *Ydb_Operations.OperationParams {
	return &Ydb_Operations.OperationParams{
		OperationMode:    Ydb_Operations.OperationParams_SYNC,
		OperationTimeout: durationpb.New(time.Duration(f.grpc.TimeoutSeconds) * time.Second),
		CancelAfter:      durationpb.New(time.Duration(f.grpc.TimeoutSeconds) * time.Second),
	}
}

func (f Factory) Connection() (*grpc.ClientConn, error) {
	cr, err := f.Credentials()
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %v", err)
	}

	return grpc.Dial(f.Endpoint(),
		grpc.WithTransportCredentials(cr),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(BufferSize),
			grpc.MaxCallRecvMsgSize(BufferSize)))
}

func (f Factory) Credentials() (credentials.TransportCredentials, error) {
	if !f.grpc.Secure {
		return insecure.NewCredentials(), nil
	}

	return credentials.NewClientTLSFromFile(f.grpc.RootCA, "")
}

func (f Factory) Endpoint() string {
	addrIndex := rand.Intn(
		len(f.grpc.Addr),
	)
	return fmt.Sprintf("%s:%d", f.grpc.Addr[addrIndex], f.grpc.Port)
}

func (f Factory) User() string {
	return f.user
}
