package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/kramerul/cfpkg/pkg/cfpkg"
	"github.com/kramerul/shalm/cmd"
	"github.com/kramerul/shalm/pkg/shalm"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"go.starlark.net/starlark"
	"gopkg.in/yaml.v2"
)

var cfConfig string

func cfExtension(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	switch module {
	case "@cf:cf":
		return starlark.StringDict{
			"cf": starlark.NewBuiltin("cf", cfpkg.MakeCF),
		}, nil
	}
	return nil, fmt.Errorf("Unknown module '%s'", module)
}

func flags(flagsSet *pflag.FlagSet) {
	flagsSet.StringVarP(&cfConfig, "cfconfig", "c", os.Getenv("CFCONFIG"), "Set cfconfig variable")
}

func testK8s(configs ...shalm.K8sConfig) (shalm.K8s, error) {
	return createK8s(true, configs...)
}
func k8s(configs ...shalm.K8sConfig) (shalm.K8s, error) {
	return createK8s(false, configs...)
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
