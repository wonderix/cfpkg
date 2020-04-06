package cfpkg

import (
	"encoding/json"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/wonderix/shalm/pkg/shalm"
)

type spaceResource struct{}

func (o *spaceResource) Apply(client *cfClient, obj *shalm.Object) error {
	var spec metaData
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	org, err := client.GetOrgByName(spec.Org)
	if err != nil {
		return err
	}
	_, err = client.CreateSpace(cfclient.SpaceRequest{
		Name:             obj.MetaData.Name,
		OrganizationGuid: org.Guid,
	})
	return ignoreCF(err, 40002)
}

func (o *spaceResource) Delete(client *cfClient, obj *shalm.Object) error {
	var spec metaData
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	spec.Space = obj.MetaData.Name
	spaceGUID, err := client.spaceGUIDForMetaData(spec)
	if err != nil {
		return ignoreMessage(err, "Unable to find")
	}
	return client.DeleteSpace(spaceGUID, true, false)
}
