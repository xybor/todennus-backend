package dto

import (
	"context"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/usecase/dto/resource"
)

type OAuth2ClientCreateRequestDTO struct {
	Name           string
	IsConfidential bool
}

type OAuth2ClientCreateResponseDTO struct {
	Client       *resource.OAuth2Client
	ClientSecret string
}

func NewOAuth2ClientCreateResponseDTO(ctx context.Context, client *domain.OAuth2Client, secret string) *OAuth2ClientCreateResponseDTO {
	return &OAuth2ClientCreateResponseDTO{
		Client:       resource.NewOAuth2Client(ctx, client, false),
		ClientSecret: secret,
	}
}

type OAuth2ClientCreateFirstRequestDTO struct {
	Username string
	Password string

	Name string
}

type OAuth2ClientCreateByAdminResponseDTO struct {
	Client       *resource.OAuth2Client
	ClientSecret string
}

func NewOAuth2ClientCreateFirstResponseDTO(ctx context.Context, client *domain.OAuth2Client, secret string) *OAuth2ClientCreateByAdminResponseDTO {
	return &OAuth2ClientCreateByAdminResponseDTO{
		Client:       resource.NewOAuth2Client(ctx, client, false),
		ClientSecret: secret,
	}
}

type OAuth2ClientGetRequestDTO struct {
	ClientID snowflake.ID
}

type OAuth2ClientGetResponseDTO struct {
	Client *resource.OAuth2Client
}

func NewOAuth2ClientGetResponse(ctx context.Context, client *domain.OAuth2Client) *OAuth2ClientGetResponseDTO {
	return &OAuth2ClientGetResponseDTO{
		Client: resource.NewOAuth2Client(ctx, client, true),
	}
}
