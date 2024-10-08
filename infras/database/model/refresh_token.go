package model

import "time"

type RefreshTokenModel struct {
	RefreshTokenID int64     `gorm:"refresh_token_id;primaryKey"`
	AccessTokenID  int64     `gorm:"access_token_id"`
	Seq            int       `gorm:"seq"`
	UpdatedAt      time.Time `gorm:"updated_at"`
}

func (RefreshTokenModel) TableName() string {
	return "refresh_tokens"
}
