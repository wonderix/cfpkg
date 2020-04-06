package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

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

func cfExtension(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	exports := starlark.NewDict(2)
	exports.SetKey(starlark.String("write"), starlark.NewBuiltin("write", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var value starlark.Value
		if err := starlark.UnpackArgs("write", args, kwargs, "value", &value); err != nil {
			return starlark.None, err
		}
		return starlark.None, shalm.WriteYamlFile(path.Join(os.Getenv("IAC_DEPLOYMENT_GEN_DIR"), "exports.yml"), value)
	}))
	exports.SetKey(starlark.String("read"), starlark.NewBuiltin("read", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		content, err := ioutil.ReadFile(path.Join(os.Getenv("IAC_DEPLOYMENT_GEN_DIR"), "exports.yml"))
		if err != nil {
			return starlark.None, err
		}
		return starlark.String(string(content)), nil
	}))
	switch module {
	case "@cf:cf":
		return starlark.StringDict{
			"cf": starlark.NewBuiltin("cf", cfpkg.MakeCF),
			"context": starlark.NewBuiltin("context", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
				return shalm.ReadYamlFile(os.Getenv("IAC_DEPLOYMENT_CONTEXT_FILE"))
			}),
			"exports": exports,
		}, nil
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
