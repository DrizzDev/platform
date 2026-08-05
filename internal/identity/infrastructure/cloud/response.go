package cloud

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DrizzDev/platform/internal/identity/application/organization"
	tenant "github.com/DrizzDev/platform/internal/identity/domain/organization"
)

// membership is the resolved organization and its display name.
type membership struct {
	org  tenant.Organization
	name string
}

// payload accepts either a direct organization body or one wrapped in a JSend
// envelope, so the parser tolerates both cloud response shapes.
type payload struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Data *struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

// interpret maps the cloud response to the organization contract. A 403 is a
// denied membership; a 401 requires a new sign-in; anything else, or an unusable
// body, is treated as unavailable.
func (Client) interpret(response *http.Response) (membership, error) {
	switch response.StatusCode {
	case http.StatusOK:
		var body payload
		if failure := json.NewDecoder(response.Body).Decode(&body); failure != nil {
			return membership{}, organization.Unavailable{}
		}
		id, name := body.Id, body.Name
		if body.Data != nil {
			id, name = body.Data.Id, body.Data.Name
		}
		if id <= 0 {
			return membership{}, organization.Unavailable{}
		}
		org, failure := tenant.New(strconv.Itoa(id))
		if failure != nil {
			return membership{}, organization.Unavailable{}
		}
		return membership{org: org, name: name}, nil
	case http.StatusUnauthorized:
		return membership{}, organization.Required{}
	case http.StatusForbidden:
		return membership{}, organization.Forbidden{}
	default:
		return membership{}, organization.Unavailable{}
	}
}
