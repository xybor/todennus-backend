package database

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrRecordDuplicate = errors.New("duplicated record")
)

func convertGormError(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrRecordNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrRecordDuplicate
	default:
		return err
	}
}
