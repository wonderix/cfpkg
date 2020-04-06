package cfpkg

import (
	"encoding/json"

	"github.com/wonderix/shalm/pkg/shalm"
)

type domainResource struct{}

type domainSpec struct {
	metaData
	Domain string `json:"domain,omitempty"`
}

func (o *domainResource) Apply(client *cfClient, obj *shalm.Object) error {
	var spec domainSpec
	err := json.Unmarshal(obj.Additional["spec"], &spec)
	if err != nil {
		return err
	}
	org, err := client.GetOrgByName(spec.Org)
	if err != nil {
		return nil
	}
	_, err = client.CreateDomain(spec.Domain, org.Guid)

	return ignoreCF(err, 130003)
}

func (o *domainResource) Delete(client *cfClient, obj *shalm.Object) error {
	var spec domainSpec
	err := json.Unmarshal(obj.Additional["spec"], &spec)
	if err != nil {
		return err
	}
	domain, err := client.GetDomainByName(spec.Domain)
	if err != nil {
		return nil
	}
	return client.DeleteDomain(domain.Guid)
}
