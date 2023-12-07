package ssh

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
	"github.com/ydb-platform/ydb-rolling-restart/pkg/rolling/service"
)

const (
	ServiceTypeStorage = "storage"
	ServiceTypeDynamic = "dynamic"
)

var (
	ServiceTypes = []string{ServiceTypeStorage, ServiceTypeDynamic}
)

func init() {
	service.OptionsMap[name] = &opts{}
}

type opts struct {
	ServiceName string
	ServiceType string
	SSH         sshOptions
}

type sshOptions struct {
	Command  string
	User     string
	Args     []string
	Logging  bool
	PoolSize int
}

func (o *opts) DefineFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.ServiceName, "service-name", "", "ydb-server-storage",
		"The service name on target host.")
	fs.StringVarP(&o.ServiceType, "service-type", "", ServiceTypeStorage,
		fmt.Sprintf("The service type on target host. Available choices: %s", util.JoinStrings(ServiceTypes, ", ")))
	fs.StringVarP(&o.SSH.Command, "ssh-command", "", "ssh", "SSH command")
	fs.StringVarP(&o.SSH.User, "ssh-user", "", "", "SSH username")
	fs.StringArrayVarP(&o.SSH.Args, "ssh-args", "", nil,
		"SSH command arguments, this overrides any other parameters")
	fs.BoolVarP(&o.SSH.Logging, "ssh-logging", "", false,
		"Enable ssh command output logging")
	fs.IntVarP(&o.SSH.PoolSize, "ssh-pool-size", "", 10, "SSH processes pool size")
}

func (o *opts) Validate() error {
	if len(o.ServiceName) == 0 {
		return fmt.Errorf("empty service name specified")
	}
	if !util.Contains(ServiceTypes, o.ServiceType) {
		return fmt.Errorf("invalid service type specified: %s", o.ServiceType)
	}
	if len(o.SSH.User) == 0 {
		return fmt.Errorf("empty ssh user specified")
	}
	if o.SSH.PoolSize < 1 {
		return fmt.Errorf("empty ssh pool size specified: %d, must be >= 1", o.SSH.PoolSize)
	}

	return nil
}
