package model

import (
	"time"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/x/enum"
)

type UserModel struct {
	ID          int64     `gorm:"id"`
	DisplayName string    `gorm:"display_name"`
	Username    string    `gorm:"username"`
	HashedPass  string    `gorm:"hashed_pass"`
	Role        string    `gorm:"role"`
	UpdatedAt   time.Time `gorm:"updated_at"`
}

func (UserModel) TableName() string {
	return "users"
}

func NewUser(d *domain.User) *UserModel {
	return &UserModel{
		ID:          d.ID.Int64(),
		DisplayName: d.DisplayName,
		Username:    d.Username,
		HashedPass:  d.HashedPass,
		UpdatedAt:   d.UpdatedAt,
		Role:        d.Role.String(),
	}
}

func (u UserModel) To() (*domain.User, error) {
	return &domain.User{
		ID:          snowflake.ID(u.ID),
		DisplayName: u.DisplayName,
		Username:    u.Username,
		HashedPass:  u.HashedPass,
		Role:        enum.FromStr[domain.UserRole](u.Role),
		UpdatedAt:   u.UpdatedAt,
	}, nil
}
