package model

import (
	"time"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/domain"
)

type OAuth2AuthorizationCodeModel struct {
	Code                string `json:"-"`
	UserID              int64  `json:"uid"`
	ClientID            int64  `json:"cid"`
	Scope               string `json:"scp"`
	CodeChallenge       string `json:"chl"`
	CodeChallengeMethod string `json:"cmt"`
	ExpiresAt           int64  `json:"exp"`
}

func NewOAuth2AuthorizationCode(code *domain.OAuth2AuthorizationCode) *OAuth2AuthorizationCodeModel {
	return &OAuth2AuthorizationCodeModel{
		Code:                code.Code,
		UserID:              code.UserID.Int64(),
		ClientID:            code.ClientID.Int64(),
		Scope:               code.Scope.String(),
		CodeChallenge:       code.CodeChallenge,
		CodeChallengeMethod: code.CodeChallengeMethod,
		ExpiresAt:           code.ExpiresAt.UnixMilli(),
	}
}

func (code OAuth2AuthorizationCodeModel) To() *domain.OAuth2AuthorizationCode {
	return &domain.OAuth2AuthorizationCode{
		Code:                code.Code,
		UserID:              snowflake.ID(code.UserID),
		ClientID:            snowflake.ID(code.ClientID),
		Scope:               domain.ScopeEngine.ParseScopes(code.Scope),
		CodeChallenge:       code.CodeChallenge,
		CodeChallengeMethod: code.CodeChallengeMethod,
		ExpiresAt:           time.UnixMilli(code.ExpiresAt),
	}
}

type OAuth2AuthorizationStoreModel struct {
	ID                  string `json:"-"`
	HasAuthenticated    bool   `json:"hat"`
	ResponseType        string `json:"res"`
	ClientID            int64  `json:"cid"`
	RedirectURI         string `json:"rdr"`
	Scope               string `json:"scp"`
	State               string `json:"sta"`
	CodeChallenge       string `json:"chl"`
	CodeChallengeMethod string `json:"cmt"`
	ExpiresAt           int64  `json:"exp"`
}

func NewOAuth2AuthorizationStore(store *domain.OAuth2AuthorizationStore) *OAuth2AuthorizationStoreModel {
	return &OAuth2AuthorizationStoreModel{
		ID:                  store.ID,
		HasAuthenticated:    store.HasAuthenticated,
		ResponseType:        store.ResponseType,
		ClientID:            store.ClientID.Int64(),
		RedirectURI:         store.RedirectURI,
		Scope:               store.Scope.String(),
		State:               store.State,
		CodeChallenge:       store.CodeChallenge,
		CodeChallengeMethod: store.CodeChallengeMethod,
		ExpiresAt:           store.ExpiresAt.UnixMilli(),
	}
}

func (store OAuth2AuthorizationStoreModel) To() *domain.OAuth2AuthorizationStore {
	return &domain.OAuth2AuthorizationStore{
		ID:                  store.ID,
		HasAuthenticated:    store.HasAuthenticated,
		ResponseType:        store.ResponseType,
		ClientID:            snowflake.ID(store.ClientID),
		RedirectURI:         store.RedirectURI,
		Scope:               domain.ScopeEngine.ParseScopes(store.Scope),
		State:               store.State,
		CodeChallenge:       store.CodeChallenge,
		CodeChallengeMethod: store.CodeChallengeMethod,
		ExpiresAt:           time.UnixMilli(store.ExpiresAt),
	}
}

type OAuth2LoginResultModel struct {
	ID              string `json:"-"`
	AuthorizationID string `json:"aid"`
	Ok              bool   `json:"ok"`
	Err             string `json:"err,omitempty"`
	UserID          int64  `json:"uid,omitempty"`
	Username        string `json:"usn,omitempty"`
	ExpiresAt       int64  `json:"exp"`
}

func NewOAuth2LoginResult(result *domain.OAuth2AuthenticationResult) *OAuth2LoginResultModel {
	return &OAuth2LoginResultModel{
		ID:              result.ID,
		AuthorizationID: result.AuthorizationID,
		Ok:              result.Ok,
		Err:             result.Error,
		UserID:          result.UserID.Int64(),
		Username:        result.Username,
		ExpiresAt:       result.ExpiresAt.UnixMilli(),
	}
}

func (result OAuth2LoginResultModel) To() *domain.OAuth2AuthenticationResult {
	return &domain.OAuth2AuthenticationResult{
		ID:              result.ID,
		AuthorizationID: result.AuthorizationID,
		Ok:              result.Ok,
		Error:           result.Err,
		UserID:          snowflake.ID(result.UserID),
		Username:        result.Username,
		ExpiresAt:       time.UnixMilli(result.ExpiresAt),
	}
}
