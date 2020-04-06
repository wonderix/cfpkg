package cfpkg

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/wonderix/shalm/pkg/shalm"
)

type serviceAccessSpec struct {
	metaData
	Service string `json:"service"`
	Broker  string `json:"broker"`
	Plan    string `json:"plan"`
	Enabled bool   `json:"enabled"`
}
type serviceAccessResource struct{}

func (o *serviceAccessResource) Apply(client *cfClient, obj *shalm.Object) error {
	return o.apply(client, obj, false)
}

func (o *serviceAccessResource) Delete(client *cfClient, obj *shalm.Object) error {
	return o.apply(client, obj, true)
}

func (o *serviceAccessResource) apply(client *cfClient, obj *shalm.Object, negate bool) error {
	var spec serviceAccessSpec
	var err error
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	servicePlans, err := client.listServicePlans(spec.Service, spec.Broker)
	if err != nil {
		return ignoreMessage(err, "Unable to find ")
	}
	for _, servicePlan := range servicePlans {
		if len(spec.Plan) == 0 || spec.Plan == servicePlan.Name {
			err = client.updateServicePlan(servicePlan.Guid, map[string]interface{}{"public": (spec.Enabled)})
			if err != nil {
				return errors.Wrapf(err, "service %s", spec.Service)
			}
		}
	}
	return nil
}
