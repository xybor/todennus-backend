package dto

import (
	"context"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/usecase/dto/resource"
)

type OAuth2ClientCreateRequest struct {
	Name           string
	IsConfidential bool
}

type OAuth2ClientCreateResponse struct {
	Client       *resource.OAuth2Client
	ClientSecret string
}

func NewOAuth2ClientCreateResponse(client *domain.OAuth2Client, secret string) *OAuth2ClientCreateResponse {
	return &OAuth2ClientCreateResponse{
		Client:       resource.NewOAuth2ClientWithoutFilter(client),
		ClientSecret: secret,
	}
}

type OAuth2ClientCreateFirstRequest struct {
	Username string
	Password string

	Name string
}

type OAuth2ClientCreateByAdminResponse struct {
	Client       *resource.OAuth2Client
	ClientSecret string
}

func NewOAuth2ClientCreateFirstResponse(ctx context.Context, client *domain.OAuth2Client, secret string) *OAuth2ClientCreateByAdminResponse {
	return &OAuth2ClientCreateByAdminResponse{
		Client:       resource.NewOAuth2ClientWithoutFilter(client),
		ClientSecret: secret,
	}
}

type OAuth2ClientGetRequest struct {
	ClientID snowflake.ID
}

type OAuth2ClientGetResponse struct {
	Client *resource.OAuth2Client
}

func NewOAuth2ClientGetResponse(ctx context.Context, client *domain.OAuth2Client) *OAuth2ClientGetResponse {
	return &OAuth2ClientGetResponse{
		Client: resource.NewOAuth2Client(ctx, client),
	}
}
