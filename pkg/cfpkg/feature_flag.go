package cfpkg

import (
	"encoding/json"

	"github.com/kramerul/shalm/pkg/shalm"
)

type featureFlagResource struct{}

func (o *featureFlagResource) Apply(client *cfClient, obj *shalm.Object) error {
	var spec featureFlag
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	return client.updateFeatureFlag(obj.MetaData.Name, spec.Enabled)
}

func (o *featureFlagResource) Delete(client *cfClient, obj *shalm.Object) error {
	return nil
}
