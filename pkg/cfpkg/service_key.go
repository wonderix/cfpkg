package cfpkg

import (
	"encoding/json"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/wonderix/shalm/pkg/shalm"
)

type serviceKeySpec struct {
	metaData
	ServiceInstance string      `json:"service_instance"`
	Parameters      interface{} `json:"parameters,omitempty"`
}
type serviceKeyResource struct{}

func (o *serviceKeyResource) Apply(client *cfClient, obj *shalm.Object) error {
	var spec serviceKeySpec
	var err error
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	serviceGUID, err := getServiceInstance(client, spec)
	if err != nil {
		return err
	}

	_, err = client.CreateServiceKey(cfclient.CreateServiceKeyRequest{
		Name:                obj.MetaData.Name,
		ServiceInstanceGuid: serviceGUID,
		Parameters:          spec.Parameters,
	})
	if err != nil {
		return ignoreCF(err, 360001)
	}
	return nil
}

func (o *serviceKeyResource) Delete(client *cfClient, obj *shalm.Object) error {
	var spec serviceKeySpec
	var err error
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	serviceGUID, err := getServiceInstance(client, spec)
	if err != nil {
		return ignoreMessage(err, "Unable to find")
	}
	serviceKeys, err := client.ListServiceKeysByQuery(url.Values{"name": []string{obj.MetaData.Name}, "service_instance_guid": []string{serviceGUID}})
	if err != nil {
		return err
	}
	if len(serviceKeys) == 0 {
		return nil
	}
	return client.DeleteServiceKey(serviceKeys[0].Guid)
}

func getServiceInstance(client *cfClient, spec serviceKeySpec) (string, error) {
	spaceGUID, err := client.spaceGUIDForMetaData(spec.metaData)
	if err != nil {
		return "", err
	}
	service, err := client.getServiceInstanceByName(spec.ServiceInstance, spaceGUID)
	if err != nil {
		return "", err
	}
	return service.Guid, nil
}
