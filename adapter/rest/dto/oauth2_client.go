package dto

import (
	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/adapter/rest/dto/resource"
	"github.com/xybor/todennus-backend/usecase/dto"
)

type OAuth2ClientCreateRequestDTO struct {
	Name           string `json:"name" example:"Example Client"`
	IsConfidential bool   `json:"is_confidential" example:"true"`
}

func (req OAuth2ClientCreateRequestDTO) To() *dto.OAuth2ClientCreateRequestDTO {
	return &dto.OAuth2ClientCreateRequestDTO{
		Name:           req.Name,
		IsConfidential: req.IsConfidential,
	}
}

type OAuth2ClientCreateResponseDTO struct {
	*resource.OAuth2Client
	ClientSecret string `json:"client_secret" example:"ElBacv..."`
}

func NewOauth2ClientCreateResponseDTO(resp *dto.OAuth2ClientCreateResponseDTO) *OAuth2ClientCreateResponseDTO {
	return &OAuth2ClientCreateResponseDTO{
		OAuth2Client: resource.NewOAuth2Client(resp.Client),
		ClientSecret: resp.ClientSecret,
	}
}

type OAuth2ClientCreateFirstRequestDTO struct {
	Username string `json:"username" example:"huykingsofm"`
	Password string `json:"password" example:"s3Cr3tP@ssW0rD"`
	Name     string `json:"name" example:"First Client"`
}

func (req *OAuth2ClientCreateFirstRequestDTO) To() *dto.OAuth2ClientCreateFirstRequestDTO {
	return &dto.OAuth2ClientCreateFirstRequestDTO{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
	}
}

type OAuth2ClientCreateFirstResponseDTO struct {
	*resource.OAuth2Client
	ClientSecret string `json:"client_secret" example:"ElBacv..."`
}

func NewOauth2ClientCreateFirstResponseDTO(resp *dto.OAuth2ClientCreateByAdminResponseDTO) *OAuth2ClientCreateFirstResponseDTO {
	return &OAuth2ClientCreateFirstResponseDTO{
		OAuth2Client: resource.NewOAuth2Client(resp.Client),
		ClientSecret: resp.ClientSecret,
	}
}

type OAuth2ClientGetRequestDTO struct {
	ClientID string `param:"client_id"`
}

func (req *OAuth2ClientGetRequestDTO) To() *dto.OAuth2ClientGetRequestDTO {
	clientID, err := snowflake.ParseString(req.ClientID)
	if err != nil {
		clientID = 0
	}

	return &dto.OAuth2ClientGetRequestDTO{
		ClientID: clientID,
	}
}

type OAuth2ClientGetResponseDTO struct {
	*resource.OAuth2Client
}

func NewOAuth2ClientGetResponseDTO(resp *dto.OAuth2ClientGetResponseDTO) *OAuth2ClientGetResponseDTO {
	return &OAuth2ClientGetResponseDTO{
		OAuth2Client: resource.NewOAuth2Client(resp.Client),
	}
}
