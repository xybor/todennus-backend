package gorm

import (
	"context"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/infras/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OAuth2ConsentRepository struct {
	db *gorm.DB
}

func NewOAuth2ConsentRepository(db *gorm.DB) *OAuth2ConsentRepository {
	return &OAuth2ConsentRepository{db: db}
}

func (repo *OAuth2ConsentRepository) Upsert(ctx context.Context, consent *domain.OAuth2Consent) error {
	model := model.NewOAuth2Consent(consent)
	return database.ConvertError(
		repo.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "client_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"scope", "expires_at", "updated_at"}),
		}).Create(&model).Error,
	)
}

func (repo *OAuth2ConsentRepository) Get(ctx context.Context, userID, clientID int64) (*domain.OAuth2Consent, error) {
	model := model.OAuth2ConsentModel{}
	if err := repo.db.WithContext(ctx).Take(&model, "user_id=? AND client_id=?", userID, clientID).Error; err != nil {
		return nil, database.ConvertError(err)
	}

	return model.To(), nil
}
