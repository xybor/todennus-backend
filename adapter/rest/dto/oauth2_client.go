package dto

import (
	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/adapter/rest/dto/resource"
	"github.com/xybor/todennus-backend/usecase/dto"
)

type OAuth2ClientCreateRequest struct {
	Name           string `json:"name" example:"Example Client"`
	IsConfidential bool   `json:"is_confidential" example:"true"`
}

func (req OAuth2ClientCreateRequest) To() *dto.OAuth2ClientCreateRequest {
	return &dto.OAuth2ClientCreateRequest{
		Name:           req.Name,
		IsConfidential: req.IsConfidential,
	}
}

type OAuth2ClientCreateResponse struct {
	*resource.OAuth2Client
	ClientSecret string `json:"client_secret" example:"ElBacv..."`
}

func NewOauth2ClientCreateResponse(resp *dto.OAuth2ClientCreateResponse) *OAuth2ClientCreateResponse {
	if resp == nil {
		return nil
	}

	return &OAuth2ClientCreateResponse{
		OAuth2Client: resource.NewOAuth2Client(resp.Client),
		ClientSecret: resp.ClientSecret,
	}
}

type OAuth2ClientCreateFirstRequest struct {
	Username string `json:"username" example:"huykingsofm"`
	Password string `json:"password" example:"s3Cr3tP@ssW0rD"`
	Name     string `json:"name" example:"First Client"`
}

func (req *OAuth2ClientCreateFirstRequest) To() *dto.OAuth2ClientCreateFirstRequest {
	return &dto.OAuth2ClientCreateFirstRequest{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
	}
}

type OAuth2ClientCreateFirstResponse struct {
	*resource.OAuth2Client
	ClientSecret string `json:"client_secret" example:"ElBacv..."`
}

func NewOauth2ClientCreateFirstResponse(resp *dto.OAuth2ClientCreateByAdminResponse) *OAuth2ClientCreateFirstResponse {
	if resp == nil {
		return nil
	}

	return &OAuth2ClientCreateFirstResponse{
		OAuth2Client: resource.NewOAuth2Client(resp.Client),
		ClientSecret: resp.ClientSecret,
	}
}

type OAuth2ClientGetRequest struct {
	ClientID string `param:"client_id"`
}

func (req *OAuth2ClientGetRequest) To() *dto.OAuth2ClientGetRequest {
	clientID, err := snowflake.ParseString(req.ClientID)
	if err != nil {
		clientID = 0
	}

	return &dto.OAuth2ClientGetRequest{
		ClientID: clientID,
	}
}

type OAuth2ClientGetResponse struct {
	*resource.OAuth2Client
}

func NewOAuth2ClientGetResponse(resp *dto.OAuth2ClientGetResponse) *OAuth2ClientGetResponse {
	if resp == nil {
		return nil
	}

	return &OAuth2ClientGetResponse{
		OAuth2Client: resource.NewOAuth2Client(resp.Client),
	}
}
