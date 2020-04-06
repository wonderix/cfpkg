package cfpkg

import (
	"encoding/json"

	"github.com/wonderix/shalm/pkg/shalm"
)

type sharedDomainResource struct{}

type sharedDomainSpec struct {
	Domain   string `json:"domain,omitempty"`
	Internal bool   `json:"internal,omitempty"`
}

func (o *sharedDomainResource) Apply(client *cfClient, obj *shalm.Object) error {
	var spec sharedDomainSpec
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	_, err := client.CreateSharedDomain(spec.Domain, spec.Internal, "")
	return ignoreCF(err, 130003)
}

func (o *sharedDomainResource) Delete(client *cfClient, obj *shalm.Object) error {
	var spec sharedDomainSpec
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	domain, err := client.GetSharedDomainByName(spec.Domain)
	if err != nil {
		return nil
	}
	return client.DeleteSharedDomain(domain.Guid, false)
}
