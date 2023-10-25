package options

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

const (
	DefaultCMSAuthEnvVar      = "YDB_TOKEN"
	DefaultCMSAuthIAMEndpoint = "iam.api.cloud.yandex.net"
)

type (
	CMSAuth interface {
		Options
		Token() (AuthToken, error)
	}
	AuthToken struct {
		Type   string
		Secret string
	}
	CMSAuthNone struct{}
	CMSAuthEnv  struct {
		Name string

		t AuthToken
	}
	CMSAuthFile struct {
		Filename string

		t AuthToken
	}
	CMSAuthIAM struct {
		KeyFilename string
		Endpoint    string
	}
)

func (t AuthToken) Token() string {
	if t.Type == "" {
		return t.Secret
	}
	return fmt.Sprintf("%s %s", t.Type, t.Secret)
}

func (an CMSAuthNone) DefineFlags(_ *pflag.FlagSet) {}
func (an CMSAuthNone) Validate() error              { return nil }
func (an CMSAuthNone) Token() (AuthToken, error)    { return AuthToken{}, nil }

func (ae CMSAuthEnv) DefineFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&ae.Name, "cms-auth-env-name", "", DefaultCMSAuthEnvVar,
		"CMS Authentication environment variable name (type: env)")
}
func (ae CMSAuthEnv) Validate() error {
	if len(ae.Name) == 0 {
		return fmt.Errorf("auth env variable name empty")
	}
	return nil
}
func (ae CMSAuthEnv) Token() (AuthToken, error) {
	if ae.t.Secret != "" {
		return ae.t, nil
	}

	zap.S().Debugf("Read auth token from %s variable", ae.Name)
	val, ok := os.LookupEnv(ae.Name)
	if !ok {
		return AuthToken{}, fmt.Errorf("failed to lookup variable: %s", ae.Name)
	}
	return AuthToken{
		Type:   "OAuth",
		Secret: val,
	}, nil
}

func (af CMSAuthFile) DefineFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&af.Filename, "cms-auth-file-token", "", "",
		"CMS Authentication file token name (type: file)")
}
func (af CMSAuthFile) Validate() error {
	if len(af.Filename) != 0 {
		if _, err := os.Stat(af.Filename); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("auth password file not exists: %v", err)
		}
	}
	return nil
}
func (af CMSAuthFile) Token() (AuthToken, error) {
	if af.t.Secret != "" {
		return af.t, nil
	}

	zap.S().Debugf("Read auth token from %s file", af.Filename)
	b, err := os.ReadFile(af.Filename)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed to read token file: %v", err)
	}
	return AuthToken{
		Type:   "OAuth",
		Secret: string(b),
	}, nil
}

func (at CMSAuthIAM) DefineFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&at.KeyFilename, "cms-auth-iam-key-file", "", "",
		"CMS Authentication iam key file path (type: iam)")
	fs.StringVarP(&at.Endpoint, "cms-auth-iam-endpoint", "", DefaultCMSAuthIAMEndpoint,
		"CMS Authentication iam endpoint (type: iam)")
}
func (at CMSAuthIAM) Validate() error {
	if len(at.KeyFilename) != 0 {
		if _, err := os.Stat(at.KeyFilename); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("auth iam key file %s not exists: %v", err)
		}
	}
	if len(at.Endpoint) == 0 {
		return fmt.Errorf("empty iam endpoint specified")
	}
	return nil
}
func (at CMSAuthIAM) Token() (AuthToken, error) {
	return AuthToken{}, nil
}
