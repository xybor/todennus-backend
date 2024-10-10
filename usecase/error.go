package usecase

import (
	"errors"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/pkg/xerror"
)

var (
	ErrUsernameExisted         = errors.New("username is existed")
	ErrUsernameNotFound        = errors.New("username is not found")
	ErrUsernamePasswordInvalid = errors.New("username or password is invalid")

	ErrUserNotFound = errors.New("user not found")

	ErrRefreshTokenInvalid = errors.New("token is invalid")
	ErrRefreshTokenStolen  = errors.New("IMPORTANT: refresh token was stolen, we will remove it")

	ErrGrantTypeInvalid = errors.New("grant type is not supported by provider")

	ErrUnauthorized = errors.New("unauthorized")

	ErrClientInvalid  = errors.New("client is invalid")
	ErrClientNotFound = errors.New("client is not found")

	ErrScopeInvalid = errors.New("scope is invalid")
)

func wrapDomainError(err error) xerror.ServiceError {
	switch {
	case errors.Is(err, domain.ErrUnknownRecoverable):
		return wrapNonDomainError(xerror.ServerityWarn, err)
	case errors.Is(err, domain.ErrUnknownCritical):
		return wrapNonDomainError(xerror.ServerityCritical, err)
	case errors.Is(err, domain.ErrKnown):
		return xerror.WrapDebug(err)
	default:
		return xerror.WrapCritical(err).WithMessage("[invalid] internal server error")
	}
}

func wrapNonDomainError(serverity xerror.Serverity, err error) xerror.ServiceError {
	return xerror.Wrap(err, serverity).WithMessage("internal server error")
}
