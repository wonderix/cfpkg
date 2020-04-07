package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/k14s/ytt/pkg/yttlibrary"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/wonderix/cfpkg/pkg/cfpkg"
	"github.com/wonderix/shalm/cmd"
	"github.com/wonderix/shalm/pkg/shalm"
	"go.starlark.net/starlark"
	"gopkg.in/yaml.v2"
)

var cfConfig string
var skipK8s bool

func iacModule() (starlark.StringDict, error) {
	state := starlark.NewDict(2)
	stateDir := os.Getenv("IAC_DEPLOYMENT_STATE_DIR")
	state.SetKey(starlark.String("write"), starlark.NewBuiltin("write", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var value starlark.Value
		if err := starlark.UnpackArgs("write", args, kwargs, "value", &value); err != nil {
			return starlark.None, err
		}
		os.MkdirAll(stateDir, 0755)
		return starlark.None, shalm.WriteYamlFile(path.Join(stateDir, "cfpkg.yml"), value)
	}))
	state.SetKey(starlark.String("read"), starlark.NewBuiltin("read", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		dict, err := shalm.ReadYamlFile(path.Join(stateDir, "cfpkg.yml"))
		if err != nil {
			return starlark.None, err
		}
		return shalm.UnwrapDict(dict), nil
	}))
	content, err := shalm.ReadYamlFile(os.Getenv("IAC_DEPLOYMENT_CONTEXT_FILE"))
	if err != nil {
		return nil, err
	}
	context, err := content.(starlark.HasAttrs).Attr("context")
	if err != nil {
		return nil, err
	}
	return starlark.StringDict{
		"context": context,
		"state":   shalm.WrapDict(state),
	}, nil

}
func cfExtension(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	switch module {
	case "@cf:cf":
		return starlark.StringDict{
			"cf": starlark.NewBuiltin("cf", cfpkg.MakeCF),
		}, nil
	case "@iac:iac":
		return iacModule()
	case "@ytt:yaml":
		return yttlibrary.YAMLAPI, nil
	case "@ytt:base64":
		return yttlibrary.Base64API, nil
	case "@ytt:json":
		return yttlibrary.JSONAPI, nil
	}
	return nil, fmt.Errorf("Unknown module '%s'", module)
}

func flags(flagsSet *pflag.FlagSet) {
	flagsSet.StringVar(&cfConfig, "cfconfig", os.Getenv("CFCONFIG"), "Set cfconfig variable")
	flagsSet.BoolVar(&skipK8s, "skip-k8s", false, "Skip all deployments to kubernetes")
}

func testK8s(configs ...shalm.K8sConfig) (shalm.K8s, error) {
	return createK8s(true, configs...)
}
func k8s(configs ...shalm.K8sConfig) (shalm.K8s, error) {
	return createK8s(skipK8s, configs...)
}

func createK8s(memory bool, configs ...shalm.K8sConfig) (shalm.K8s, error) {
	var k shalm.K8s
	var err error
	if memory {
		k = shalm.NewK8sInMemory("default")
	} else {
		k, err = shalm.NewK8s(configs...)
		if err != nil {
			return nil, err
		}
	}
	if len(cfConfig) == 0 {
		return k, nil
	}
	config := cfpkg.Config{}
	content, err := ioutil.ReadFile(cfConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "reading yaml file %s", cfConfig)
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal %s", cfConfig)
	}
	return cfpkg.NewCF(k, config, 10*time.Second)
}

func main() {
	cmd.Execute(cmd.WithModules(cfExtension), cmd.WithApplyFlags(flags), cmd.WithTestFlags(flags), cmd.WithK8s(k8s), cmd.WithTestK8s(testK8s))
}
