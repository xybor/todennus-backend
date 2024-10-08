package database

import (
	"context"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/infras/database/model"
	"gorm.io/gorm"
)

type OAuth2ClientRepository struct {
	db *gorm.DB
}

func NewOAuth2ClientRepository(db *gorm.DB) *OAuth2ClientRepository {
	return &OAuth2ClientRepository{db: db}
}

func (repo *OAuth2ClientRepository) Create(ctx context.Context, client domain.OAuth2Client) error {
	model := model.OAuth2ClientModel{}
	model.From(client)

	return convertGormError(repo.db.Create(&model).Error)
}

func (repo *OAuth2ClientRepository) GetByID(ctx context.Context, clientID int64) (domain.OAuth2Client, error) {
	model := model.OAuth2ClientModel{}
	if err := repo.db.Take(&model, "id=?", clientID).Error; err != nil {
		return domain.OAuth2Client{}, convertGormError(err)
	}

	return model.To(), nil
}
