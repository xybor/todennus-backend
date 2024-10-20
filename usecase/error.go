package usecase

import (
	"errors"

	"github.com/xybor/x/xerror"
)

var (
	ErrServer = xerror.Wrap(errors.New("server_error"), "an unexpected error occurred")

	ErrRequestInvalid = errors.New("invalid_request")

	ErrUsernameExisted  = errors.New("username_exists")
	ErrUsernameNotFound = errors.New("username_not_found")

	ErrUserNotFound = errors.New("user_not_found")

	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("forbidden")

	ErrClientInvalid = errors.New("invalid_client")

	ErrIdPInvalid = errors.New("invalid_idp")

	ErrScopeInvalid = errors.New("invalid_scope")

	ErrAuthorizationAccessDenied = errors.New("access_denined")
	ErrTokenInvalidGrant         = errors.New("invalid_grant")
)
