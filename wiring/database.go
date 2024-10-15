package wiring

import (
	"context"

	config "github.com/xybor/todennus-config"
	"github.com/xybor/todennus-migration/postgres"
	"gorm.io/gorm"
)

type Databases struct {
	GormPostgres *gorm.DB
}

func InitializeDatabases(ctx context.Context, config config.Config) (Databases, error) {
	db := Databases{}
	var err error

	db.GormPostgres, err = postgres.Initialize(ctx, config)
	if err != nil {
		return db, err
	}

	return db, nil
}
