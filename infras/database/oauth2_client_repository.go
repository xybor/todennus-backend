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
	model := model.NewOAuth2Client(client)
	return ConvertError(repo.db.Create(&model).Error)
}

func (repo *OAuth2ClientRepository) GetByID(ctx context.Context, clientID int64) (domain.OAuth2Client, error) {
	model := model.OAuth2ClientModel{}
	if err := repo.db.Take(&model, "id=?", clientID).Error; err != nil {
		return domain.OAuth2Client{}, ConvertError(err)
	}

	return model.To()
}

func (repo *OAuth2ClientRepository) Count(ctx context.Context) (int64, error) {
	var n int64
	err := repo.db.Model(&model.OAuth2ClientModel{}).Count(&n).Error
	return n, ConvertError(err)
}
