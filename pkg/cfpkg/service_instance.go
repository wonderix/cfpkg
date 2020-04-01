package cfpkg

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/kramerul/shalm/pkg/shalm"
)

type serviceInstanceSpec struct {
	metaData
	Plan       string                 `json:"plan"`
	Broker     string                 `json:"broker"`
	Service    string                 `json:"service"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
}
type serviceInstanceResource struct{}

func (o *serviceInstanceResource) Apply(client *cfClient, obj *shalm.Object) error {
	var request cfclient.ServiceInstanceRequest
	var spec serviceInstanceSpec
	var err error
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	request.Tags = spec.Tags
	request.Parameters = spec.Parameters
	request.SpaceGuid, err = client.spaceGUIDForMetaData(spec.metaData)
	if err != nil {
		return err
	}
	servicePlans, err := client.listServicePlans(spec.Service, spec.Broker)
	if err != nil {
		return err
	}
	for _, servicePlan := range servicePlans {
		if spec.Plan == servicePlan.Name {
			request.ServicePlanGuid = servicePlan.Guid
		}
	}
	if len(request.ServicePlanGuid) == 0 {
		return fmt.Errorf("service plan %s not found", spec.Plan)
	}
	request.Name = obj.MetaData.Name
	serviceInstance, err := client.getServiceInstanceByName(spec.Service, request.SpaceGuid)
	if err != nil {
		if ignoreMessage(err, "Unable to find") != nil {
			return err
		}
		_, err = client.CreateServiceInstance(request)
		if err != nil {
			return err
		}
	} else {
		data, err := json.Marshal(request)
		if err != nil {
			return err
		}
		err = client.UpdateServiceInstance(serviceInstance.Guid, bytes.NewBuffer(data), false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *serviceInstanceResource) Delete(client *cfClient, obj *shalm.Object) error {
	var spec serviceInstanceSpec
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	spaceGUID, err := client.spaceGUIDForMetaData(spec.metaData)
	if err != nil {
		return ignoreMessage(err, "Unable to find")
	}
	serviceInstance, err := client.getServiceInstanceByName(obj.MetaData.Name, spaceGUID)
	if err != nil {
		return ignoreMessage(err, "Unable to find")
	}
	return client.DeleteServiceInstance(serviceInstance.Guid, false, false)
}
