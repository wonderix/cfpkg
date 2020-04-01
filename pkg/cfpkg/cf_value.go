package cfpkg

import (
	"time"

	"github.com/kramerul/shalm/pkg/shalm"
	"go.starlark.net/starlark"
)

// MakeCF -
func MakeCF(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var k8s shalm.K8sValue
	var config Config
	var timeout int = 10
	err := starlark.UnpackArgs("cf", args, kwargs, "k8s", &k8s, "api",
		&config.CF.URL, "username", &config.CF.Username, "password", &config.CF.Password, "skip_ssl_validation", &config.CF.SkipSslValidation,
		"uaa_url?", &config.UAA.URL, "uaa_client_id?", &config.UAA.ClientID, "uaa_client_secret?", &config.UAA.ClientSecret,
		"timeout?", &timeout)
	if err != nil {
		return nil, err
	}
	cf, err := NewCF(k8s, config, time.Duration(timeout)*time.Second)
	if err != nil {
		return starlark.None, err
	}
	return shalm.NewK8sValue(cf), nil
}
