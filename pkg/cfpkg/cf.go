package cfpkg

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/wonderix/shalm/pkg/shalm"
)

type resource interface {
	Apply(client *cfClient, obj *shalm.Object) error
	Delete(client *cfClient, obj *shalm.Object) error
}

type uaaResource interface {
	Apply(client *uaa.API, obj *shalm.Object) error
	Delete(client *uaa.API, obj *shalm.Object) error
}

type metaData struct {
	Org   string `json:"org,omitempty"`
	Space string `json:"space,omitempty"`
}

// UaaConfig -
type UaaConfig struct {
	URL          string `yaml:"url,omitempty"`
	ClientID     string `yaml:"client_id,omitempty"`
	ClientSecret string `yaml:"client_secret,omitempty"`
}

// CFConfig -
type CFConfig struct {
	URL               string `yaml:"url"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	SkipSslValidation bool   `yaml:"skip_ssl_validation"`
}

// Config -
type Config struct {
	UAA UaaConfig `yaml:"uaa,omitempty"`
	CF  CFConfig  `yaml:"cf,omitempty"`
}

// CF -
type CF struct {
	client       *cfClient
	uaaClient    *uaa.API
	k8s          shalm.K8s
	resources    map[string]resource
	uaaResources map[string]uaaResource
}

var _ shalm.K8s = (*CF)(nil)

// NewCF -
func NewCF(k8s shalm.K8s, config Config, timeout time.Duration) (*CF, error) {
	client, err := newCfClient(config, timeout)
	if err != nil {
		return nil, err
	}
	var uaaClient *uaa.API
	if config.UAA.URL != "" {
		uaaClient, err = uaa.New(config.UAA.URL, uaa.WithClientCredentials(config.UAA.ClientID, config.UAA.ClientSecret, uaa.JSONWebToken),
			uaa.WithSkipSSLValidation(config.CF.SkipSslValidation))
		if err != nil {
			return nil, err
		}
	}
	return &CF{k8s: k8s, client: client,
		resources: map[string]resource{
			"OrgQuota":        &orgQuotaResource{},
			"Org":             &orgResource{},
			"Space":           &spaceResource{},
			"Domain":          &domainResource{},
			"SharedDomain":    &sharedDomainResource{},
			"Buildpack":       &buildPackResource{},
			"ServiceBroker":   &serviceBrokerResource{},
			"FeatureFlag":     &featureFlagResource{},
			"ServiceAccess":   &serviceAccessResource{},
			"ServiceBinding":  &serviceBindingResource{},
			"ServiceInstance": &serviceInstanceResource{},
			"ServiceKey":      &serviceKeyResource{},
			"Route":           &routeResource{},
			"App":             &appResource{},
		},
		uaaResources: map[string]uaaResource{
			"User":   &userResource{},
			"Member": &memberResource{},
		},
		uaaClient: uaaClient}, nil
}

// Get -
func (c *CF) Get(kind string, name string, options *shalm.K8sOptions) (*shalm.Object, error) {
	return c.k8s.Get(kind, name, options)
}

// IsNotExist -
func (c *CF) IsNotExist(err error) bool {
	return c.k8s.IsNotExist(err)
}

// ForSubChart -
func (c *CF) ForSubChart(namespace string, app string, version semver.Version) shalm.K8s {
	return &CF{
		client:       c.client,
		uaaClient:    c.uaaClient,
		k8s:          c.k8s.ForSubChart(namespace, app, version),
		resources:    c.resources,
		uaaResources: c.uaaResources,
	}
}

// Inspect -
func (c *CF) Inspect() string {
	return c.k8s.Inspect()
}

// Host -
func (c *CF) Host() string {
	return c.k8s.Host()
}

// SetTool -
func (c *CF) SetTool(tool shalm.Tool) {
	c.k8s.SetTool(tool)
}

// Watch -
func (c *CF) Watch(kind string, name string, options *shalm.K8sOptions) shalm.ObjectStream {
	return c.k8s.Watch(kind, name, options)
}

// RolloutStatus -
func (c *CF) RolloutStatus(kind string, name string, options *shalm.K8sOptions) error {
	return c.k8s.RolloutStatus(kind, name, options)
}

// Wait -
func (c *CF) Wait(kind string, name string, condition string, options *shalm.K8sOptions) error {
	return c.k8s.Wait(kind, name, condition, options)
}

// DeleteObject -
func (c *CF) DeleteObject(kind string, name string, options *shalm.K8sOptions) error {
	return c.k8s.DeleteObject(kind, name, options)
}

// Delete -
func (c *CF) Delete(output shalm.ObjectStream, options *shalm.K8sOptions) error {
	grouped := output.Sort(compare, true).GroupBy(grouping)
	err := grouped("cf")(func(obj *shalm.Object) error {
		resource, ok := c.resources[obj.Kind]
		if ok {
			fmt.Printf("deleting %s.%s/%s\n", obj.Kind, obj.APIVersion, obj.MetaData.Name)
			return resource.Delete(c.client, obj)
		}
		return fmt.Errorf("Invalid object kind %s", obj.Kind)
	})
	c.client.clearCache()
	if err != nil {
		return err
	}
	err = grouped("uaa")(func(obj *shalm.Object) error {
		if c.uaaClient == nil {
			return errors.New("Uaa client not configured")
		}
		resource, ok := c.uaaResources[obj.Kind]
		if ok {
			fmt.Printf("deleting %s.%s/%s\n", obj.Kind, obj.APIVersion, obj.MetaData.Name)
			return resource.Delete(c.uaaClient, obj)
		}
		return fmt.Errorf("Invalid object kind %s", obj.Kind)
	})
	if err != nil {
		return err
	}
	return c.k8s.Delete(grouped("k8s"), options)
}

// ConfigContent -
func (c *CF) ConfigContent() *string {
	return c.k8s.ConfigContent()
}

// ForConfig -
func (c *CF) ForConfig(config string) (shalm.K8s, error) {
	return c.k8s.ForConfig(config)
}

// Progress -
func (c *CF) Progress(progress int) {
	c.k8s.Progress(progress)
}

// Tool -
func (c *CF) Tool() shalm.Tool {
	return c.k8s.Tool()
}

// Apply -
func (c *CF) Apply(output shalm.ObjectStream, options *shalm.K8sOptions) error {
	grouped := output.Sort(compare, false).GroupBy(grouping)
	err := grouped("uaa")(func(obj *shalm.Object) error {
		if c.uaaClient == nil {
			return errors.New("Uaa client not configured")
		}
		resource, ok := c.uaaResources[obj.Kind]
		if ok {
			fmt.Printf("applying %s.%s/%s\n", obj.Kind, obj.APIVersion, obj.MetaData.Name)
			return resource.Apply(c.uaaClient, obj)
		}
		return fmt.Errorf("Invalid object kind %s", obj.Kind)
	})
	if err != nil {
		return err
	}
	err = grouped("cf")(func(obj *shalm.Object) error {
		resource, ok := c.resources[obj.Kind]
		if ok {
			fmt.Printf("applying %s.%s/%s\n", obj.Kind, obj.APIVersion, obj.MetaData.Name)
			return resource.Apply(c.client, obj)
		}
		return fmt.Errorf("Invalid object kind %s", obj.Kind)
	})
	c.client.clearCache()
	if err != nil {
		return err
	}
	return c.k8s.Apply(grouped("k8s"), options)
}

func grouping(o *shalm.Object) string {
	switch o.APIVersion {
	case "cloudfoundry.io/v1beta1":
		return "cf"
	case "uaa.io/v1beta1":
		return "uaa"
	default:
		return "k8s"
	}
}

func kindOrdinal(o *shalm.Object) int {
	switch o.Kind {
	case "User":
		return 1
	case "Member":
		return 2
	case "OrgQuota":
		return 3
	case "FeatureFlag":
		return 4
	case "Buildpack":
		return 5
	case "Org":
		return 6
	case "Space":
		return 7
	case "SharedDomain":
		return 8
	case "Domain":
		return 9
	case "Route":
		return 10
	case "ServiceBroker":
		return 11
	case "ServiceAccess":
		return 12
	case "ServiceInstance":
		return 13
	case "CustomerServiceInstance":
		return 14
	case "ServiceKey":
		return 15
	case "App":
		return 16
	case "ServiceBinding":
		return 17
	default:
		return 1000
	}
}

func compare(o1 *shalm.Object, o2 *shalm.Object) int {
	diff := kindOrdinal(o1) - kindOrdinal(o2)
	if diff == 0 {
		diff = strings.Compare(o1.MetaData.Name, o2.MetaData.Name)
	}
	return diff
}

func ignoreCF(err error, code int) error {
	if err == nil {
		return nil
	}
	cfError, ok := err.(cfclient.CloudFoundryError)
	if ok {
		if cfError.Code == code {
			return nil
		}
	}
	return err
}

func convertRequestError(err error) error {
	if err == nil {
		return nil
	}
	requestErr, ok := err.(uaa.RequestError)
	if ok {
		return fmt.Errorf("%s: %s", requestErr.Error(), string(requestErr.ErrorResponse))
	}
	return err
}
func ignoreMessage(err error, msg string) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), msg) {
		return nil
	}
	return err
}

func isNotExist(err error) bool {
	return strings.Contains(err.Error(), "not found")
}
