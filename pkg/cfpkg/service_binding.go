package cfpkg

import (
	"encoding/json"
	"net/url"

	"github.com/kramerul/shalm/pkg/shalm"
)

type serviceBindingSpec struct {
	metaData
	ServiceInstance string `json:"service_instance"`
	App             string `json:"app"`
}
type serviceBindingResource struct{}

func (o *serviceBindingResource) Apply(client *cfClient, obj *shalm.Object) error {
	serviceGUID, appGUID, err := getServiceAndApp(client, obj)
	if err != nil {
		return err
	}
	_, err = client.CreateServiceBinding(appGUID, serviceGUID)
	if err != nil {
		return ignoreCF(err, 90003)
	}
	return nil
}

func (o *serviceBindingResource) Delete(client *cfClient, obj *shalm.Object) error {
	serviceGUID, appGUID, err := getServiceAndApp(client, obj)
	if err != nil {
		return ignoreMessage(err, "Unable to find")
	}
	serviceBindings, err := client.ListServiceBindingsByQuery(url.Values{"app_guid": []string{appGUID}, "service_instance_guid": []string{serviceGUID}})
	if err != nil {
		return err
	}
	if len(serviceBindings) == 0 {
		return nil
	}
	return client.DeleteServiceBinding(serviceBindings[0].Guid)
}

func getServiceAndApp(client *cfClient, obj *shalm.Object) (string, string, error) {
	var spec serviceBindingSpec
	var err error
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return "", "", err
	}
	spaceGUID, err := client.spaceGUIDForMetaData(spec.metaData)
	if err != nil {
		return "", "", err
	}
	service, err := client.getServiceInstanceByName(spec.ServiceInstance, spaceGUID)
	if err != nil {
		return "", "", err
	}
	app, err := client.getAppByName(spec.App, spaceGUID)
	if err != nil {
		return "", "", err
	}
	return service.Guid, app.Guid, nil
}
