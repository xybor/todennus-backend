package resource

import (
	"context"

	"github.com/xybor/todennus-backend/domain"
)

type OAuth2Client struct {
	OwnerID      int64
	ClientID     int64
	Name         string
	AllowedScope string
}

func NewOAuth2Client(ctx context.Context, client domain.OAuth2Client, needFilter bool) OAuth2Client {
	usecaseClient := OAuth2Client{
		ClientID:     client.ID,
		OwnerID:      client.OwnerUserID,
		Name:         client.Name,
		AllowedScope: client.AllowedScope.String(),
	}

	if needFilter {
		Filter(ctx, &usecaseClient.OwnerID).WhenRequestUserNot(client.OwnerUserID)
		Filter(ctx, &usecaseClient.AllowedScope).
			WhenRequestUserNot(client.OwnerUserID).
			WhenNotContainsScope(domain.ScopeEngine.New(
				domain.Actions.Read,
				domain.Resources.Client.AllowedScope,
			))
	}

	return usecaseClient
}
