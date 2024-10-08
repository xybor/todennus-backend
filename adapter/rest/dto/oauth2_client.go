package dto

import (
	"strconv"

	"github.com/xybor/todennus-backend/adapter/rest/dto/resource"
	"github.com/xybor/todennus-backend/usecase/dto"
)

type OAuth2ClientCreateRequestDTO struct {
	Name           string `json:"name"`
	IsConfidential bool   `json:"is_confidential"`
}

func (req *OAuth2ClientCreateRequestDTO) To() dto.OAuth2ClientCreateRequestDTO {
	return dto.OAuth2ClientCreateRequestDTO{
		Name:           req.Name,
		IsConfidential: req.IsConfidential,
	}
}

type OAuth2ClientCreateResponseDTO struct {
	resource.OAuth2Client
	ClientSecret string `json:"client_secret"`
}

func NewOauth2ClientCreateResponseDTO(resp dto.OAuth2ClientCreateResponseDTO) OAuth2ClientCreateResponseDTO {
	return OAuth2ClientCreateResponseDTO{
		OAuth2Client: resource.NewOAuth2Client(resp.Client),
		ClientSecret: resp.ClientSecret,
	}
}

type OAuth2ClientGetRequestDTO struct {
	ClientID string `param:"client_id"`
}

func (req *OAuth2ClientGetRequestDTO) To() dto.OAuth2ClientGetRequestDTO {
	clientID, err := strconv.ParseInt(req.ClientID, 10, 64)
	if err != nil {
		clientID = 0
	}

	return dto.OAuth2ClientGetRequestDTO{
		ClientID: clientID,
	}
}

type OAuth2ClientGetResponseDTO struct {
	resource.OAuth2Client
}

func NewOAuth2ClientGetResponseDTO(resp dto.OAuth2ClientGetResponseDTO) OAuth2ClientGetResponseDTO {
	return OAuth2ClientGetResponseDTO{
		OAuth2Client: resource.NewOAuth2Client(resp.Client),
	}
}
