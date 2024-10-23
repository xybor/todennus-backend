package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/infras/database/model"
)

func oauth2ConsentResultKey(userID, clientID int64) string {
	return fmt.Sprintf("oauth2_consent:%d:%d", userID, clientID)
}

type OAuth2ConsentResultRepository struct {
	client *redis.Client
}

func NewOAuth2ConsentResultRepository(client *redis.Client) *OAuth2ConsentResultRepository {
	return &OAuth2ConsentResultRepository{
		client: client,
	}
}

func (repo *OAuth2ConsentResultRepository) SaveResult(
	ctx context.Context,
	result *domain.OAuth2ConsentResult,
) error {
	key := oauth2ConsentResultKey(result.UserID.Int64(), result.ClientID.Int64())
	value, err := json.Marshal(model.NewOAuth2ConsentResultModel(result))
	if err != nil {
		return err
	}

	return database.ConvertError(repo.client.SetEx(ctx, key, value, time.Until(result.ExpiresAt)).Err())
}

func (repo *OAuth2ConsentResultRepository) LoadResult(
	ctx context.Context,
	userID, clientID int64,
) (*domain.OAuth2ConsentResult, error) {
	result, err := repo.client.Get(ctx, oauth2ConsentResultKey(userID, clientID)).Result()
	if err != nil {
		return nil, database.ConvertError(err)
	}

	model := model.OAuth2ConsentResultModel{}
	if err := json.Unmarshal([]byte(result), &model); err != nil {
		return nil, err
	}

	return model.To(userID, clientID), nil
}

func (repo *OAuth2ConsentResultRepository) DeleteResult(
	ctx context.Context,
	userID, clientID int64,
) error {
	return database.ConvertError(repo.client.Del(ctx, oauth2ConsentResultKey(userID, clientID)).Err())
}
