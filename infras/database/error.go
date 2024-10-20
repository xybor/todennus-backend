package database

import (
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/xybor/x/xerror"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrRecordDuplicate = errors.New("duplicated record")
)

func ConvertError(err error) error {
	switch {
	case xerror.Is(err, gorm.ErrRecordNotFound, redis.Nil):
		return ErrRecordNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrRecordDuplicate
	default:
		return err
	}
}
