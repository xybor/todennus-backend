package resource

import (
	"github.com/xybor/todennus-backend/usecase/dto/resource"
)

type OAuth2Client struct {
	OwnerID      string `json:"owner_id,omitempty" example:"330559330522759168"`
	ClientID     string `json:"client_id,omitempty" example:"332974701238012989"`
	Name         string `json:"name,omitempty" example:"Example Client"`
	AllowedScope string `json:"allowed_scope,omitempty" example:"read:user"`
}

func NewOAuth2Client(client *resource.OAuth2Client) *OAuth2Client {
	return &OAuth2Client{
		OwnerID:      client.OwnerID.String(),
		ClientID:     client.ClientID.String(),
		Name:         client.Name,
		AllowedScope: client.AllowedScope,
	}
}
