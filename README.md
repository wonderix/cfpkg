# cfpkg

Install Cloud Foundry artifacts the kubernetes way.

This product brings the wold of [shalm](https://github.com/kramerul/shalm) to Cloud Foundry.


## Features

* Supports everything provided by [shalm](https://github.com/kramerul/shalm)
* Allows installation of Cloud Foundry artifacts in a k8s way


## Example

The following artefact will create (or update) an organization within Cloud Foundry

```
---
apiVersion: cloudfoundry.io/v1beta1
kind: Org
metadata:
  name: test
```

## Reference

The following types are supported. Other types will follow.

| Type            | Example                                                                 |
|-----------------|-------------------------------------------------------------------------|
| User            | [user.yaml](test/resources/templates/user.yaml)                         |
| Member          | [member.yaml](test/resources/templates/member.yaml)                     |
| Domain          | [domain.yaml](test/resources/templates/domain.yaml)                     |
| SharedDomain    | [shared_domain.yaml](test/resources/templates/shared_domain.yaml)       |
| OrgQuota        | [org_quota.yaml](test/resources/templates/org_quota.yaml)               |
| FeatureFlag     | [feature_flag.yaml](test/resources/templates/feature_flag.yaml)         |
| Org             | [org.yaml](test/resources/templates/org.yaml)                           |
| Space           | [space.yaml](test/resources/templates/space.yaml)                       |
| Route           | [route.yaml](test/resources/templates/route.yaml)                       |
| Buildpack       | [buildpack.yaml](test/resources/templates/buildpack.yaml)               |
| ServiceBroker   | [service_broker.yaml](test/resources/templates/service_broker.yaml)     |
| ServiceAccess   | [service_access.yaml](test/resources/templates/service_access.yaml)     |
| ServiceInstance | [service_instance.yaml](test/resources/templates/service_instance.yaml) |
| ServiceBinding  | [service_binding.yaml](test/resources/templates/service_binding.yaml) |
| ServiceKey      | [service_key.yaml](test/resources/templates/service_key.yaml) |
| App             | [app.yaml](test/resources/templates/app.yaml)                           |