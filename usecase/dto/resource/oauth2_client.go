package resource

import (
	"context"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/x/scope"
)

type OAuth2Client struct {
	OwnerID      snowflake.ID
	ClientID     snowflake.ID
	Name         string
	AllowedScope string
}

func NewOAuth2Client(ctx context.Context, client *domain.OAuth2Client, needFilter bool) *OAuth2Client {
	usecaseClient := &OAuth2Client{
		ClientID:     client.ID,
		OwnerID:      client.OwnerUserID,
		Name:         client.Name,
		AllowedScope: client.AllowedScope.String(),
	}

	if needFilter {
		Filter(ctx, &usecaseClient.OwnerID).WhenRequestUserNot(client.OwnerUserID)
		Filter(ctx, &usecaseClient.AllowedScope).
			WhenRequestUserNot(client.OwnerUserID).
			WhenNotContainsScope(scope.New(domain.Actions.Read, domain.Resources.Client.AllowedScope))
	}

	return usecaseClient
}
