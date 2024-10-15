package usecase

import (
	"errors"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/x/errorx"
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
	ErrForbidden    = errors.New("forbidden")

	ErrClientInvalid  = errors.New("client is invalid")
	ErrClientNotFound = errors.New("client is not found")

	ErrScopeInvalid = errors.New("scope is invalid")

	ErrRequestInvalid = errors.New("request is invalid")
)

func wrapDomainError(err error) errorx.ServiceError {
	switch {
	case errors.Is(err, domain.ErrUnknownRecoverable):
		return wrapNonDomainError(errorx.ServerityWarn, err)
	case errors.Is(err, domain.ErrUnknownCritical):
		return wrapNonDomainError(errorx.ServerityCritical, err)
	case errors.Is(err, domain.ErrKnown):
		return errorx.WrapDebug(err)
	default:
		return errorx.WrapCritical(err).WithMessage("[invalid] internal server error")
	}
}

func wrapNonDomainError(serverity errorx.Serverity, err error) errorx.ServiceError {
	return errorx.Wrap(err, serverity).WithMessage("internal server error")
}
