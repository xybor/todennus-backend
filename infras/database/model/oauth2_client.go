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
	UpdatedAt      time.Time `gorm:"updated_at"`
}

func (OAuth2ClientModel) TableName() string {
	return "oauth2_clients"
}

func (client *OAuth2ClientModel) To() domain.OAuth2Client {
	return domain.OAuth2Client{
		ID:             client.ID,
		OwnerUserID:    client.UserID,
		Name:           client.Name,
		HasedSecret:    client.HashedSecret,
		IsConfidential: client.IsConfidential,
		UpdatedAt:      client.UpdatedAt,
	}
}

func (client *OAuth2ClientModel) From(domain domain.OAuth2Client) {
	client.ID = domain.ID
	client.UserID = domain.OwnerUserID
	client.Name = domain.Name
	client.HashedSecret = domain.HasedSecret
	client.IsConfidential = domain.IsConfidential
	client.UpdatedAt = domain.UpdatedAt
}
