package cfpkg

import (
	"encoding/json"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/wonderix/shalm/pkg/shalm"
)

type serviceBrokerSpec struct {
	metaData
	BrokerURL string `json:"broker_url"`
	Username  string `json:"auth_username"`
	Password  string `json:"auth_password"`
}
type serviceBrokerResource struct{}

func (o *serviceBrokerResource) Apply(client *cfClient, obj *shalm.Object) error {
	var spec serviceBrokerSpec
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	serviceBroker, err := client.GetServiceBrokerByName(obj.MetaData.Name)
	if err != nil {
		if ignoreMessage(err, "Unable to find") != nil {
			return err
		}
		request := cfclient.CreateServiceBrokerRequest{
			Name:      obj.MetaData.Name,
			BrokerURL: spec.BrokerURL,
			Username:  spec.Username,
			Password:  spec.Password,
		}
		request.SpaceGUID, err = client.spaceGUIDForMetaData(spec.metaData)
		if err != nil {
			return err
		}
		_, err = client.CreateServiceBroker(request)
		return err
	}
	_, err = client.UpdateServiceBroker(serviceBroker.Guid, cfclient.UpdateServiceBrokerRequest{
		Name:      obj.MetaData.Name,
		BrokerURL: spec.BrokerURL,
		Username:  spec.Username,
		Password:  spec.Password,
	})
	return err
}

func (o *serviceBrokerResource) Delete(client *cfClient, obj *shalm.Object) error {
	serviceBroker, err := client.GetServiceBrokerByName(obj.MetaData.Name)
	if err != nil {
		return ignoreMessage(err, "Unable to find")
	}
	return client.DeleteServiceBroker(serviceBroker.Guid)
}
