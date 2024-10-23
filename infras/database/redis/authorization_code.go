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

func oauth2AuthorizationCodeKey(code string) string {
	return fmt.Sprintf("oauth2_code:%s", code)
}

func oauth2AuthorizationStoreKey(code string) string {
	return fmt.Sprintf("oauth2_store:%s", code)
}

func oauth2AuthenticationResultKey(code string) string {
	return fmt.Sprintf("oauth2_auth:%s", code)
}

type OAuth2AuthorizationCodeRepository struct {
	client *redis.Client
}

func NewOAuth2AuthorizationCodeRepository(client *redis.Client) *OAuth2AuthorizationCodeRepository {
	return &OAuth2AuthorizationCodeRepository{
		client: client,
	}
}

func (repo *OAuth2AuthorizationCodeRepository) SaveAuthorizationCode(
	ctx context.Context,
	code *domain.OAuth2AuthorizationCode,
) error {
	model := model.NewOAuth2AuthorizationCode(code)

	modelJSON, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return database.ConvertError(repo.client.SetEx(ctx,
		oauth2AuthorizationCodeKey(model.Code), modelJSON, time.Until(code.ExpiresAt)).Err())
}

func (repo *OAuth2AuthorizationCodeRepository) LoadAuthorizationCode(
	ctx context.Context,
	code string,
) (*domain.OAuth2AuthorizationCode, error) {
	result, err := repo.client.Get(ctx, oauth2AuthorizationCodeKey(code)).Result()
	if err != nil {
		return nil, database.ConvertError(err)
	}

	model := model.OAuth2AuthorizationCodeModel{Code: code}
	if err := json.Unmarshal([]byte(result), &model); err != nil {
		return nil, err
	}

	return model.To(), nil
}

func (repo *OAuth2AuthorizationCodeRepository) DeleteAuthorizationCode(
	ctx context.Context,
	code string,
) error {
	return database.ConvertError(repo.client.Del(ctx, oauth2AuthorizationCodeKey(code)).Err())
}

func (repo *OAuth2AuthorizationCodeRepository) SaveAuthorizationStore(
	ctx context.Context,
	store *domain.OAuth2AuthorizationStore,
) error {
	model := model.NewOAuth2AuthorizationStore(store)

	modelJSON, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return database.ConvertError(repo.client.SetEx(ctx,
		oauth2AuthorizationStoreKey(model.ID), modelJSON, time.Until(store.ExpiresAt)).Err())
}

func (repo *OAuth2AuthorizationCodeRepository) LoadAuthorizationStore(
	ctx context.Context,
	id string,
) (*domain.OAuth2AuthorizationStore, error) {
	result, err := repo.client.Get(ctx, oauth2AuthorizationStoreKey(id)).Result()
	if err != nil {
		return nil, database.ConvertError(err)
	}

	model := model.OAuth2AuthorizationStoreModel{ID: id}
	if err := json.Unmarshal([]byte(result), &model); err != nil {
		return nil, err
	}

	return model.To(), nil
}

func (repo *OAuth2AuthorizationCodeRepository) DeleteAuthorizationStore(
	ctx context.Context,
	id string,
) error {
	return database.ConvertError(repo.client.Del(ctx, oauth2AuthorizationStoreKey(id)).Err())
}

func (repo *OAuth2AuthorizationCodeRepository) SaveAuthenticationResult(
	ctx context.Context,
	result *domain.OAuth2AuthenticationResult,
) error {
	model := model.NewOAuth2LoginResult(result)

	modelJSON, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return database.ConvertError(repo.client.SetEx(ctx,
		oauth2AuthenticationResultKey(model.ID), modelJSON, time.Until(result.ExpiresAt)).Err())
}

func (repo *OAuth2AuthorizationCodeRepository) LoadAuthenticationResult(
	ctx context.Context,
	id string,
) (*domain.OAuth2AuthenticationResult, error) {
	result, err := repo.client.Get(ctx, oauth2AuthenticationResultKey(id)).Result()
	if err != nil {
		return nil, database.ConvertError(err)
	}

	model := model.OAuth2LoginResultModel{ID: id}
	if err := json.Unmarshal([]byte(result), &model); err != nil {
		return nil, err
	}

	return model.To(), nil
}

func (repo *OAuth2AuthorizationCodeRepository) DeleteAuthenticationResult(ctx context.Context, id string) error {
	return database.ConvertError(repo.client.Del(ctx, oauth2AuthenticationResultKey(id)).Err())
}
