package options

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/ydb-platform/ydb-rolling-restart/internal/util"
)

var (
	CMSAvailabilityModes = []string{"max", "keep", "force"}
	CMSAuths             = map[string]CMSAuth{
		"none": &CMSAuthNone{},
		"env":  &CMSAuthEnv{},
		"file": &CMSAuthFile{},
		"iam":  &CMSAuthIAM{},
	}
)

const (
	CMSDefaultRetryWaitTime    = 60
	CMSDefaultAvailAbilityMode = "max"
	CMSDefaultAuthType         = "none"
)

type CMS struct {
	Auth             CMSAuth
	AuthType         string
	User             string
	AvailabilityMode string
	RetryWaitSeconds int
}

func (cms *CMS) DefineFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&cms.AuthType, "cms-auth-type", "", CMSDefaultAuthType,
		fmt.Sprintf("CMS Authentication types: %+v", util.Keys(CMSAuths)))

	for _, auth := range CMSAuths {
		auth.DefineFlags(fs)
	}

	fs.StringVarP(&cms.User, "cms-user", "", "rolling-restart",
		"CMS user that will be used for restart")
	fs.StringVarP(&cms.AvailabilityMode, "cms-availability-mode", "", CMSDefaultAvailAbilityMode,
		fmt.Sprintf("CMS Availability mode (%+v)", CMSAvailabilityModes))
	fs.IntVarP(&cms.RetryWaitSeconds, "cms-wait-time-seconds", "", CMSDefaultRetryWaitTime,
		"CMS retry time in seconds")
}

func (cms *CMS) Validate() error {
	if !util.Contains(util.Keys(CMSAuths), cms.AuthType) {
		return fmt.Errorf("invalid auth type specified: %s, use one of: %+v", cms.AuthType, util.Keys(CMSAuths))
	}
	if len(cms.User) == 0 {
		return fmt.Errorf("empty auth user")
	}

	cms.Auth = CMSAuths[cms.AuthType]
	if err := cms.Auth.Validate(); err != nil {
		return err
	}

	if !util.Contains(CMSAvailabilityModes, cms.AvailabilityMode) {
		return fmt.Errorf("invalid availability mode specified: %v, use one of: %+v", cms.AvailabilityMode, CMSAvailabilityModes)
	}
	if cms.RetryWaitSeconds < 0 {
		return fmt.Errorf("invalid value specified: %d", cms.RetryWaitSeconds)
	}

	return nil
}
