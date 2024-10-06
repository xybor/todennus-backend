package database

import (
	"context"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/infras/database/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Save(ctx context.Context, user domain.User) error {
	model := &model.UserModel{}
	model.From(user)

	return convertGormError(repo.db.Save(&model).Error)
}

func (repo *UserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	model := model.UserModel{}
	if err := repo.db.Take(&model, "username=?", username).Error; err != nil {
		return domain.User{}, convertGormError(err)
	}

	return model.To(), nil
}