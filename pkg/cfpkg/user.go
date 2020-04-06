package cfpkg

import (
	"encoding/json"

	"github.com/cloudfoundry-community/go-uaa"
	"github.com/wonderix/shalm/pkg/shalm"
)

type userSpec struct {
	Password string      `json:"password,omitempty"`
	Emails   []uaa.Email `json:"emails,omitempty"`
}

type userResource struct{}

func (o *userResource) Apply(client *uaa.API, obj *shalm.Object) error {
	specJSON, ok := obj.Additional["spec"]
	var spec userSpec
	if ok {
		err := json.Unmarshal(specJSON, &spec)
		if err != nil {
			return err
		}
	}
	user := uaa.User{
		Username: obj.MetaData.Name,
		Password: spec.Password,
		Emails:   spec.Emails,
	}
	// _, err := client.UpdateUser(user)
	// return convertRequestError(err)

	_, err := client.CreateUser(user)
	return ignoreMessage(convertRequestError(err), "scim_resource_already_exists")
}

func (o *userResource) Delete(client *uaa.API, obj *shalm.Object) error {
	user, err := client.GetUserByUsername(obj.MetaData.Name, "", "")
	if err != nil {
		if isNotExist(err) {
			return nil
		}
		return err
	}
	_, err = client.DeleteUser(user.ID)
	return convertRequestError(err)
}
