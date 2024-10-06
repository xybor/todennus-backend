package model

import (
	"time"

	"github.com/xybor/todennus-backend/domain"
)

type UserModel struct {
	ID          int64     `gorm:"id"`
	DisplayName string    `gorm:"display_name"`
	Username    string    `gorm:"username"`
	HashedPass  string    `gorm:"hashed_pass"`
	CreatedAt   time.Time `gorm:"created_at"`
	UpdatedAt   time.Time `gorm:"updated_at"`
}

func (*UserModel) TableName() string {
	return "users"
}

func (u *UserModel) To() domain.User {
	return domain.User{
		ID:          u.ID,
		DisplayName: u.DisplayName,
		Username:    u.Username,
		HashedPass:  u.HashedPass,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func (u *UserModel) From(d domain.User) {
	u.ID = d.ID
	u.DisplayName = d.DisplayName
	u.Username = d.Username
	u.HashedPass = d.HashedPass
	u.CreatedAt = d.CreatedAt
	u.UpdatedAt = d.UpdatedAt
}
