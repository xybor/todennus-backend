package database

import (
	"context"

	"github.com/xybor/todennus-backend/infras/database/model"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (repo *RefreshTokenRepository) Save(
	ctx context.Context,
	refreshTokenId,
	accessTokenID int64,
	seq int,
) error {
	return convertGormError(repo.db.Save(&model.RefreshTokenModel{
		RefreshTokenID: refreshTokenId,
		AccessTokenID:  accessTokenID,
		Seq:            seq,
	}).Error)
}

func (repo *RefreshTokenRepository) UpdateByRefreshTokenID(
	ctx context.Context,
	refreshTokenID, accessTokenID int64,
	expectedCurSeq int,
) error {
	result := repo.db.Model(&model.RefreshTokenModel{}).
		Where("refresh_token_id=? AND seq=?", refreshTokenID, expectedCurSeq).
		Updates(map[string]any{
			"seq":             expectedCurSeq + 1,
			"access_token_id": accessTokenID,
		})

	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}

	return convertGormError(result.Error)
}

func (repo *RefreshTokenRepository) DeleteByRefreshTokenID(
	ctx context.Context, refreshTokenID int64,
) error {
	return convertGormError(repo.db.Delete(&model.RefreshTokenModel{}, refreshTokenID).Error)
}
