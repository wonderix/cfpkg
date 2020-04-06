package cfpkg

import (
	"encoding/json"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/wonderix/shalm/pkg/shalm"
)

type orgQuotaResource struct{}

func (o *orgQuotaResource) Apply(client *cfClient, obj *shalm.Object) error {
	var spec cfclient.OrgQuotaRequest
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	spec.Name = obj.MetaData.Name
	quota, err := client.GetOrgQuotaByName(obj.MetaData.Name)
	if err != nil {
		if ignoreMessage(err, "Unable to find") == nil {
			_, err := client.CreateOrgQuota(spec)
			return err
		}
		return err
	}
	_, err = client.UpdateOrgQuota(quota.Guid, spec)
	return err
}

func (o *orgQuotaResource) Delete(client *cfClient, obj *shalm.Object) error {
	quota, err := client.GetOrgQuotaByName(obj.MetaData.Name)
	if err != nil {
		return ignoreMessage(err, "Unable to find")
	}
	if err = client.DeleteOrgQuota(quota.Guid, false); err != nil {
		return err
	}
	return nil
}
