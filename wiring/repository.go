package wiring

import (
	"context"

	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/usecase/abstraction"
)

type Repositories struct {
	abstraction.UserRepository
}

func InitializeRepositories(ctx context.Context, db Databases) (Repositories, error) {
	r := Repositories{}

	r.UserRepository = database.NewUserRepository(db.GormPostgres)

	return r, nil
}
