package options

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
)

var (
	CMSAvailabilityModes = []string{"max, keep, force"}
)

const (
	CMSDefaultApiTimeoutSeconds = 60
	CMSDefaultRetryWaitTime     = 60
)

type CMS struct {
	AuthUser          string
	AuthPasswordFile  string
	AvailabilityMode  string
	ApiTimeoutSeconds int
	RetryWaitSeconds  int
}

func (cms *CMS) DefineFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&cms.AuthUser, "cms-auth-user", "", "rolling-restart",
		"Specify username which will be used by CMS")
	fs.StringVarP(&cms.AuthPasswordFile, "cms-auth-password-file", "", "",
		"Specify path to authentication file for cms user")
	fs.StringVarP(&cms.AvailabilityMode, "cms-availability-mode", "", "max",
		fmt.Sprintf("CMS Availability mode (%+v)", CMSAvailabilityModes))
	fs.IntVarP(&cms.ApiTimeoutSeconds, "cms-api-timeout-seconds", "", CMSDefaultApiTimeoutSeconds,
		"CMS API response timeout in seconds")
	fs.IntVarP(&cms.RetryWaitSeconds, "cms-wait-time-seconds", "", CMSDefaultRetryWaitTime,
		"CMS retry time in seconds")
}

func (cms *CMS) Validate() error {
	if len(cms.AuthUser) == 0 {
		return fmt.Errorf("empty auth user")
	}
	if len(cms.AuthPasswordFile) != 0 {
		if _, err := os.Stat(cms.AuthPasswordFile); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("auth password file not exists: %v", err)
		}
	}
	if !util.Contains(CMSAvailabilityModes, cms.AvailabilityMode) {
		return fmt.Errorf("invalid availability mode specified: %v, use one of: %+v", cms.AvailabilityMode, CMSAvailabilityModes)
	}
	if cms.ApiTimeoutSeconds < 0 {
		return fmt.Errorf("invalid value specified: %d", cms.ApiTimeoutSeconds)
	}
	if cms.RetryWaitSeconds < 0 {
		return fmt.Errorf("invalid value specified: %d", cms.RetryWaitSeconds)
	}

	return nil
}
