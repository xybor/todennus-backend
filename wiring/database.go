package wiring

import (
	"context"

	"github.com/redis/go-redis/v9"
	config "github.com/xybor/todennus-config"
	"github.com/xybor/todennus-migration/postgres"
	"gorm.io/gorm"
)

type Databases struct {
	GormPostgres *gorm.DB
	Redis        *redis.Client
}

func InitializeDatabases(ctx context.Context, config config.Config) (Databases, error) {
	db := Databases{}
	var err error

	db.GormPostgres, err = postgres.Initialize(ctx, config)
	if err != nil {
		return db, err
	}

	db.Redis = redis.NewClient(&redis.Options{
		Addr:     config.Variable.Redis.Addr,
		DB:       config.Variable.Redis.DB,
		Username: config.Secret.Redis.Username,
		Password: config.Secret.Redis.Password,
	})

	return db, nil
}
