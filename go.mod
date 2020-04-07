module github.com/wonderix/cfpkg

go 1.13

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20190808214049-35bcce23fc5f
	github.com/cloudfoundry-community/go-uaa v0.3.1
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/k14s/ytt v0.26.1-0.20200402233022-1aaca8db2e6a
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5
	github.com/wonderix/shalm v0.5.3
	go.starlark.net v0.0.0-20191021185836-28350e608555
	golang.org/x/sys v0.0.0-20200120151820-655fe14d7479 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.17.4
)

replace go.starlark.net => github.com/k14s/starlark-go v0.0.0-20200402152745-409c85f3828d // ytt branch
