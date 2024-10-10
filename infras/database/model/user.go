package model

import (
	"time"

	"github.com/xybor/todennus-backend/domain"
)

type UserModel struct {
	ID           int64     `gorm:"id"`
	DisplayName  string    `gorm:"display_name"`
	Username     string    `gorm:"username"`
	HashedPass   string    `gorm:"hashed_pass"`
	AllowedScope string    `gorm:"allowed_scope"`
	UpdatedAt    time.Time `gorm:"updated_at"`
}

func (UserModel) TableName() string {
	return "users"
}

func (u *UserModel) To() (domain.User, error) {
	allowedScope, err := domain.ScopeEngine.ParseScopes(u.AllowedScope)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:           u.ID,
		DisplayName:  u.DisplayName,
		Username:     u.Username,
		HashedPass:   u.HashedPass,
		UpdatedAt:    u.UpdatedAt,
		AllowedScope: allowedScope,
	}, nil
}

func (u *UserModel) From(d domain.User) {
	u.ID = d.ID
	u.DisplayName = d.DisplayName
	u.Username = d.Username
	u.HashedPass = d.HashedPass
	u.UpdatedAt = d.UpdatedAt
	u.AllowedScope = d.AllowedScope.String()
}
