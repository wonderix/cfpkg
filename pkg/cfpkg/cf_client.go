package cfpkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/cloudfoundry-community/go-cfclient"
)

type cfClient struct {
	*cfclient.Client
	spaceGUIDCache map[metaData]string
}

func newCfClient(config Config, timeout time.Duration) (*cfClient, error) {
	var client *cfclient.Client
	var err error
	waitErr := wait.PollImmediate(10*time.Second, timeout, func() (bool, error) {
		client, err = cfclient.NewClient(&cfclient.Config{
			Username:          config.CF.Username,
			Password:          config.CF.Password,
			ApiAddress:        config.CF.URL,
			SkipSslValidation: config.CF.SkipSslValidation,
		})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	if waitErr != nil {
		return nil, errors.Wrapf(err, "unable to connect to cloud foundry at %s", config.CF.URL)
	}
	return &cfClient{Client: client, spaceGUIDCache: make(map[metaData]string)}, nil
}

type featureFlag struct {
	Enabled bool `json:"enabled"`
}

func (c *cfClient) clearCache() {
	c.spaceGUIDCache = make(map[metaData]string)
}

func (c *cfClient) spaceGUIDForMetaData(data metaData) (string, error) {
	if len(data.Org) == 0 && len(data.Space) == 0 {
		return "", nil
	}
	cached, ok := c.spaceGUIDCache[data]
	if ok {
		return cached, nil
	}
	org, err := c.GetOrgByName(data.Org)
	if err != nil {
		return "", err
	}
	space, err := c.GetSpaceByName(data.Space, org.Guid)
	if err != nil {
		return "", err
	}
	c.spaceGUIDCache[data] = space.Guid
	return space.Guid, nil
}

func (c *cfClient) getAppByName(name string, spaceGUID string) (*cfclient.App, error) {
	apps, err := c.ListAppsByQuery(url.Values{"names": []string{name}, "space_guids": []string{spaceGUID}})
	if err != nil {
		return nil, err
	}
	if len(apps) == 1 {
		return &apps[0], nil
	}
	return nil, errors.New("Unable to find")
}

func (c *cfClient) getServiceInstanceByName(name string, spaceGUID string) (*cfclient.ServiceInstance, error) {
	services, err := c.ListServiceInstancesByQuery(url.Values{"names": []string{name}, "space_guids": []string{spaceGUID}})
	if err != nil {
		return nil, err
	}
	if len(services) == 1 {
		return &services[0], nil
	}
	return nil, errors.New("Unable to find")
}

func (c *cfClient) getServicePlanByName(name string, spaceGUID string) (*cfclient.ServicePlan, error) {
	plans, err := c.ListServicePlansByQuery(url.Values{"names": []string{name}, "space_guids": []string{spaceGUID}})
	if err != nil {
		return nil, err
	}
	if len(plans) == 1 {
		return &plans[0], nil
	}
	return nil, errors.New("Unable to find")
}

func (c *cfClient) getBuildpackByName(name string) (*cfclient.Buildpack, error) {
	buildPacks, err := c.ListBuildpacks()
	if err != nil {
		return nil, err
	}
	for _, buildPack := range buildPacks {
		if buildPack.Name == name {
			return &buildPack, nil
		}
	}
	return nil, nil
}

func (c *cfClient) updateFeatureFlag(name string, enabled bool) error {
	data, err := json.Marshal(&featureFlag{Enabled: enabled})
	if err != nil {
		return err
	}
	request := c.NewRequestWithBody(http.MethodPut, fmt.Sprintf("/v2/config/feature_flags/%s", name), bytes.NewBuffer(data))
	response, err := c.DoRequest(request)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		msg, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return errors.New(string(msg))
	}
	return nil

}

func (c *cfClient) updateServicePlan(servicePlanGUID string, serviceRequest map[string]interface{}) error {
	data, err := json.Marshal(serviceRequest)
	if err != nil {
		return err
	}
	request := c.NewRequestWithBody(http.MethodPut, fmt.Sprintf("/v2/service_plans/%s", servicePlanGUID), bytes.NewBuffer(data))
	response, err := c.DoRequest(request)
	if err != nil {
		return err
	}
	if response.StatusCode/100 != 2 {
		msg, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return errors.New(string(msg))
	}
	return nil

}

func (c *cfClient) listServicePlans(service string, broker string) ([]cfclient.ServicePlan, error) {
	serviceQuery := url.Values{"label": []string{service}}
	if len(broker) != 0 {
		brk, err := c.GetServiceBrokerByName(broker)
		if err != nil {
			return nil, errors.Wrapf(err, "service broker %s", broker)
		}
		serviceQuery.Set("service_broker_guid", brk.Guid)
	}
	services, err := c.ListServicesByQuery(serviceQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "service %s", service)
	}
	if len(services) == 0 {
		return nil, fmt.Errorf("service %s not found", service)
	}
	return c.ListServicePlansByQuery(url.Values{"service_guid": []string{services[0].Guid}})

}
