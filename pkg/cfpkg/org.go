package cfpkg

import (
	"encoding/json"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/wonderix/shalm/pkg/shalm"
)

type orgResource struct{}

type orgSpec struct {
	Status string `json:"status,omitempty"`
	Quota  string `json:"quota,omitempty"`
}

func (o *orgResource) Apply(client *cfClient, obj *shalm.Object) error {
	request := cfclient.OrgRequest{
		Name: obj.MetaData.Name,
	}
	specJSON, ok := obj.Additional["spec"]
	if ok {
		var spec orgSpec
		err := json.Unmarshal(specJSON, &spec)
		if err != nil {
			return err
		}
		request.Status = spec.Status
		if spec.Quota != "" {
			quota, err := client.GetOrgQuotaByName(spec.Quota)
			if err != nil {
				return err
			}
			request.QuotaDefinitionGuid = quota.Guid
		}
	}
	_, err := client.CreateOrg(request)
	return ignoreCF(err, 30002)
}

func (o *orgResource) Delete(client *cfClient, obj *shalm.Object) error {
	org, err := client.GetOrgByName(obj.MetaData.Name)
	if err != nil {
		return nil
	}
	return client.DeleteOrg(org.Guid, true, false)
}
