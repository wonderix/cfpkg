package cfpkg

import (
	"encoding/json"

	"github.com/cloudfoundry-community/go-uaa"
	"github.com/wonderix/shalm/pkg/shalm"
)

type memberResource struct{}

type memberSpec struct {
	Member string `json:"member,omitempty"`
	Type   string `json:"type,omitempty"`
	Group  string `json:"group,omitempty"`
}

func (o *memberResource) Apply(client *uaa.API, obj *shalm.Object) error {
	groupID, userID, typ, err := extractMemberParameter(client, obj.Additional["spec"])
	if err != nil {
		return err
	}
	err = client.AddGroupMember(groupID, userID, typ, "")
	return ignoreMessage(convertRequestError(err), "member_already_exists")
}

func (o *memberResource) Delete(client *uaa.API, obj *shalm.Object) error {
	groupID, userID, typ, err := extractMemberParameter(client, obj.Additional["spec"])
	if err != nil {
		return ignoreMessage(err, "not found")
	}
	err = client.RemoveGroupMember(groupID, userID, typ, "")
	return ignoreMessage(convertRequestError(err), "member_not_found")
}

func extractMemberParameter(client *uaa.API, specJSON json.RawMessage) (string, string, string, error) {
	var spec memberSpec
	err := json.Unmarshal(specJSON, &spec)
	if err != nil {
		return "", "", "", err
	}
	group, err := client.GetGroupByName(spec.Group, "")
	if err != nil {
		return "", "", "", convertRequestError(err)
	}
	user, err := client.GetUserByUsername(spec.Member, "", "")
	if err != nil {
		return "", "", "", convertRequestError(err)
	}
	return group.ID, user.ID, spec.Type, nil

}
