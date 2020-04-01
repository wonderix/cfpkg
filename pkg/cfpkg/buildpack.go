package cfpkg

import (
	"encoding/json"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/kramerul/shalm/pkg/shalm"
)

type buildPackResource struct{}

type buildPackSpec struct {
	cfclient.BuildpackRequest
	URL string `json:"url,omitempty"`
}

func (o *buildPackResource) Apply(client *cfClient, obj *shalm.Object) error {
	var spec buildPackSpec
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	spec.Name = &obj.MetaData.Name
	buildPack, err := client.getBuildpackByName(obj.MetaData.Name)
	if err != nil {
		return err
	}
	if buildPack == nil {
		buildPack, err = client.CreateBuildpack(&spec.BuildpackRequest)
		if err != nil {
			return err
		}
	} else {
		if err = buildPack.Update(&spec.BuildpackRequest); err != nil {
			return err
		}
	}
	reader, err := get(spec.URL)
	if err != nil {
		return err
	}
	defer reader.Close()
	return nil
	// Upload is not working (CF-BuildpackBitsUploadInvalid|290002): The buildpack upload is invalid: a filename must be specified
	// return buildPack.Upload(resp.Body, path.Base(spec.URL))
}

func (o *buildPackResource) Delete(client *cfClient, obj *shalm.Object) error {
	buildPack, err := client.getBuildpackByName(obj.MetaData.Name)
	if buildPack == nil || err != nil {
		return err
	}
	return client.DeleteBuildpack(buildPack.Guid, false)
}
