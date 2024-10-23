package model

import (
	"time"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
)

type OAuth2ConsentResultModel struct {
	Accepted  bool   `json:"acp"`
	Scope     string `json:"scp"`
	ExpiresAt int64  `json:"exp"`
}

func NewOAuth2ConsentResultModel(result *domain.OAuth2ConsentResult) *OAuth2ConsentResultModel {
	return &OAuth2ConsentResultModel{
		Accepted:  result.Accepted,
		Scope:     result.Scope.String(),
		ExpiresAt: result.ExpiresAt.UnixMilli(),
	}
}

func (model *OAuth2ConsentResultModel) To(userID, clientID int64) *domain.OAuth2ConsentResult {
	return &domain.OAuth2ConsentResult{
		UserID:    snowflake.ID(userID),
		ClientID:  snowflake.ID(clientID),
		Accepted:  model.Accepted,
		Scope:     domain.ScopeEngine.ParseScopes(model.Scope),
		ExpiresAt: time.UnixMilli(model.ExpiresAt),
	}
}

type OAuth2ConsentModel struct {
	UserID    int64     `gorm:"user_id;primaryKey"`
	ClientID  int64     `gorm:"client_id;primaryKey"`
	Scope     string    `gorm:"scope"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
	ExpiresAt time.Time `gorm:"expires_at"`
}

func (OAuth2ConsentModel) TableName() string {
	return "oauth2_consents"
}

func NewOAuth2Consent(consent *domain.OAuth2Consent) *OAuth2ConsentModel {
	return &OAuth2ConsentModel{
		UserID:    consent.UserID.Int64(),
		ClientID:  consent.ClientID.Int64(),
		Scope:     consent.Scope.String(),
		CreatedAt: consent.CreatedAt,
		UpdatedAt: consent.UpdatedAt,
		ExpiresAt: consent.ExpiresAt,
	}
}

func (model OAuth2ConsentModel) To() *domain.OAuth2Consent {
	return &domain.OAuth2Consent{
		UserID:    snowflake.ID(model.UserID),
		ClientID:  snowflake.ID(model.ClientID),
		Scope:     domain.ScopeEngine.ParseScopes(model.Scope),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		ExpiresAt: model.ExpiresAt,
	}
}
