package cfpkg

import (
	"encoding/json"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/kramerul/shalm/pkg/shalm"
)

type appSpec struct {
	metaData
	cfclient.AppCreateRequest
	Ports []int  `json:"ports,omitempty"`
	URL   string `json:"url,omitempty"`
}

func (a *appSpec) setEtag(etag string) {
	if a.Environment == nil {
		a.Environment = make(map[string]interface{})
	}
	a.Environment["CFPKG_ETAG"] = etag
}

type appResource struct{}

func (o *appResource) Apply(client *cfClient, obj *shalm.Object) error {
	var spec appSpec
	var err error
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}

	spec.SpaceGuid, err = client.spaceGUIDForMetaData(spec.metaData)
	if err != nil {
		return err
	}
	state := spec.State
	app, err := client.getAppByName(obj.MetaData.Name, spec.SpaceGuid)
	etag := ""
	if err != nil {
		if ignoreMessage(err, "Unable to find") != nil {
			return err
		}
		if len(spec.URL) != 0 {
			spec.State = cfclient.APP_STOPPED
		}
		spec.Name = obj.MetaData.Name
		_, err = client.CreateApp(spec.AppCreateRequest)
		if err != nil {
			return err
		}
	} else {
		if app.Environment != nil {
			x, ok := app.Environment["CFPKG_ETAG"]
			if ok {
				etag = x.(string)
			}
		}
		spec.setEtag(etag)
		_, err = client.UpdateApp(app.Guid, cfclient.AppUpdateResource{
			Name:                    obj.MetaData.Name,
			Memory:                  spec.Memory,
			Instances:               spec.Instances,
			DiskQuota:               spec.DiskQuota,
			SpaceGuid:               spec.SpaceGuid,
			StackGuid:               spec.StackGuid,
			State:                   spec.State,
			Command:                 spec.Command,
			Buildpack:               spec.Buildpack,
			HealthCheckHttpEndpoint: spec.HealthCheckHttpEndpoint,
			HealthCheckType:         string(spec.HealthCheckType),
			Diego:                   spec.Diego,
			EnableSSH:               spec.EnableSSH,
			DockerImage:             spec.DockerImage,
			DockerCredentials:       map[string]interface{}{"username": spec.DockerCredentials.Username, "password": spec.DockerCredentials.Password},
			Environment:             spec.Environment,
			Ports:                   spec.Ports,
		})
		if err != nil {
			return err
		}
	}
	if len(spec.URL) != 0 {
		reader, newEtag, err := getWithEtag(spec.URL, etag)
		if err != nil {
			return err
		}
		if newEtag != etag {
			defer reader.Close()
			err = client.UploadAppBits(reader, app.Guid)
			if err != nil {
				return err
			}
			spec.setEtag(newEtag)
			_, err = client.UpdateApp(app.Guid, cfclient.AppUpdateResource{State: state, Environment: spec.Environment})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *appResource) Delete(client *cfClient, obj *shalm.Object) error {
	var spec appSpec
	var err error
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	spec.SpaceGuid, err = client.spaceGUIDForMetaData(spec.metaData)
	if err != nil {
		return ignoreMessage(err, "Unable to find")
	}
	app, err := client.getAppByName(obj.MetaData.Name, spec.SpaceGuid)
	if err != nil {
		return ignoreMessage(err, "Unable to find")
	}
	return client.DeleteApp(app.Guid)
}
