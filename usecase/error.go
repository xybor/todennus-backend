package usecase

import (
	"errors"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/x/xerror"
)

var (
	ErrServer        = xerror.Enrich(errors.New("server_error"), "an unexpected error occurred")
	ErrServerTimeout = xerror.Enrich(errors.New("server_timeout"), "server timeout")

	ErrRequestInvalid = errors.New("invalid_request")
	ErrDuplicated     = errors.New("duplicated")
	ErrNotFound       = errors.New("not_found")

	ErrCredentialsInvalid = errors.New("invalid_credentials")

	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("forbidden")

	ErrClientInvalid = errors.New("invalid_client")

	ErrScopeInvalid = errors.New("invalid_scope")

	ErrAuthorizationAccessDenied = errors.New("access_denined")
	ErrTokenInvalidGrant         = errors.New("invalid_grant")
)

var domainerr = xerror.NewWrapperConfigs(ErrServer, domain.ErrKnown)
