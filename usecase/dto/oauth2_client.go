package dto

import "github.com/xybor/todennus-backend/domain"

type OAuth2ClientCreateRequestDTO struct {
	Name           string
	IsConfidential bool
}

type OAuth2ClientCreateResponseDTO struct {
	Client       domain.OAuth2Client
	ClientSecret string
}

func NewOAuth2ClientCreateResponseDTO(client domain.OAuth2Client, secret string) OAuth2ClientCreateResponseDTO {
	return OAuth2ClientCreateResponseDTO{
		Client:       client,
		ClientSecret: secret,
	}
}

type OAuth2ClientGetRequestDTO struct {
	ClientID int64
}

type OAuth2ClientGetResponseDTO struct {
	Client domain.OAuth2Client
}

func NewOAuth2ClientGetResponse(client domain.OAuth2Client) OAuth2ClientGetResponseDTO {
	return OAuth2ClientGetResponseDTO{
		Client: client,
	}
}
