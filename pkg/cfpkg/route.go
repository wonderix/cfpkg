package cfpkg

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/kramerul/shalm/pkg/shalm"
)

type routeSpec struct {
	metaData
	Domain string `json:"domain"`
	Host   string `json:"host"`
	Path   string `json:"path"`
	Port   int    `json:"port"`
}
type routeResource struct{}

func (o *routeResource) Apply(client *cfClient, obj *shalm.Object) error {
	var spec routeSpec
	var err error
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	request := cfclient.RouteRequest{
		Host: spec.Host,
		Path: spec.Path,
		Port: spec.Port,
	}
	request.SpaceGuid, err = client.spaceGUIDForMetaData(spec.metaData)
	if err != nil {
		return err
	}
	request.DomainGuid, err = getDomainGUID(client, spec.Domain)
	if err != nil {
		return err
	}
	_, err = client.CreateRoute(request)
	return ignoreMessage(err, "CF-RouteHostTaken")
}

func (o *routeResource) Delete(client *cfClient, obj *shalm.Object) error {
	var spec routeSpec
	var err error
	if err := json.Unmarshal(obj.Additional["spec"], &spec); err != nil {
		return err
	}
	domainGUID, err := getDomainGUID(client, spec.Domain)
	if err != nil {
		return err
	}
	routes, err := client.ListRoutesByQuery(url.Values{"host": []string{spec.Host}, "domain_guid": []string{domainGUID}})
	if err != nil {
		return err
	}
	if len(routes) == 0 {
		return nil
	}
	return client.DeleteRoute(routes[0].Guid)
}

func getDomainGUID(client *cfClient, domain string) (string, error) {
	sharedDomains, err := client.ListSharedDomainsByQuery(url.Values{"name": []string{domain}})
	if err != nil {
		return "", err
	}
	if len(sharedDomains) == 0 {
		domains, err := client.ListDomainsByQuery(url.Values{"name": []string{domain}})
		if err != nil {
			return "", err
		}
		if len(domains) == 0 {
			return "", fmt.Errorf("Unable to find domain %s", domain)
		}
		return domains[0].Guid, nil
	}
	return sharedDomains[0].Guid, nil

}
