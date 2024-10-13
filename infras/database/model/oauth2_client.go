package model

import (
	"time"

	"github.com/xybor/todennus-backend/domain"
)

type OAuth2ClientModel struct {
	ID             int64     `gorm:"id;primaryKey"`
	UserID         int64     `gorm:"user_id"`
	Name           string    `gorm:"name"`
	HashedSecret   string    `gorm:"hashed_secret"`
	IsConfidential bool      `gorm:"is_confidential"`
	AllowedScope   string    `gorm:"allowed_scope"`
	UpdatedAt      time.Time `gorm:"updated_at"`
}

func (OAuth2ClientModel) TableName() string {
	return "oauth2_clients"
}

func (client *OAuth2ClientModel) To() (domain.OAuth2Client, error) {
	return domain.OAuth2Client{
		ID:             client.ID,
		OwnerUserID:    client.UserID,
		Name:           client.Name,
		HashedSecret:   client.HashedSecret,
		IsConfidential: client.IsConfidential,
		AllowedScope:   domain.ScopeEngine.ParseScopes(client.AllowedScope),
		UpdatedAt:      client.UpdatedAt,
	}, nil
}

func NewOAuth2Client(domain domain.OAuth2Client) OAuth2ClientModel {
	return OAuth2ClientModel{
		ID:             domain.ID,
		UserID:         domain.OwnerUserID,
		Name:           domain.Name,
		HashedSecret:   domain.HashedSecret,
		IsConfidential: domain.IsConfidential,
		UpdatedAt:      domain.UpdatedAt,
		AllowedScope:   domain.AllowedScope.String(),
	}
}
