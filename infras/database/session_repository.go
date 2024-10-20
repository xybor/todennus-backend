package database

import (
	"context"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/infras/database/model"
	"github.com/xybor/x/session"
	"github.com/xybor/x/xcontext"
)

type SessionRepository struct {
	store session.Store[model.SessionModel]
}

func NewSessionRepository(store session.Store[model.SessionModel]) *SessionRepository {
	return &SessionRepository{store: store}
}

func (repo *SessionRepository) Load(ctx context.Context) (domain.Session, error) {
	model, err := repo.store.Load(ctx, xcontext.Session(ctx))
	if err != nil {
		return domain.Session{}, err
	}

	return model.To(), nil
}

func (repo *SessionRepository) Save(ctx context.Context, session domain.Session) error {
	model := model.NewSession(session)

	return repo.store.Save(ctx, xcontext.Session(ctx), model)
}
