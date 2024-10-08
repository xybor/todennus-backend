package resource

import (
	"strconv"

	"github.com/xybor/todennus-backend/domain"
)

type OAuth2Client struct {
	OwnerID  string `json:"owner_id"`
	ClientID string `json:"client_id"`
	Name     string `json:"name"`
}

func NewOAuth2Client(client domain.OAuth2Client) OAuth2Client {
	return OAuth2Client{
		OwnerID:  strconv.FormatInt(client.OwnerUserID, 10),
		ClientID: strconv.FormatInt(client.ID, 10),
		Name:     client.Name,
	}
}
