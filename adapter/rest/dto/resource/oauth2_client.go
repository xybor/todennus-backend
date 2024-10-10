package resource

import (
	"github.com/xybor/todennus-backend/pkg/xstring"
	"github.com/xybor/todennus-backend/usecase/dto/resource"
)

type OAuth2Client struct {
	OwnerID      string `json:"owner_id,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	Name         string `json:"name,omitempty"`
	AllowedScope string `json:"allowed_scope,omitempty"`
}

func NewOAuth2Client(client resource.OAuth2Client) OAuth2Client {
	return OAuth2Client{
		OwnerID:      xstring.FormatID(client.OwnerID),
		ClientID:     xstring.FormatID(client.ClientID),
		Name:         client.Name,
		AllowedScope: client.AllowedScope,
	}
}
