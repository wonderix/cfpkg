module github.com/wonderix/cfpkg

go 1.13

require (
	github.com/Masterminds/semver/v3 v3.0.3
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20201123235753-4f46d6348a05
	github.com/cloudfoundry-community/go-uaa v0.3.1
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/k14s/starlark-go v0.0.0-20200720175618-3a5c849cc368
	github.com/k14s/ytt v0.26.1-0.20200402233022-1aaca8db2e6a
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5
	github.com/wonderix/shalm v0.7.1
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.17.4
)

replace github.com/k14s/ytt => github.com/wonderix/ytt v0.28.1-0.20200908051131-36914082e903
