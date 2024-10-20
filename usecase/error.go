package usecase

import (
	"errors"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/x/xerror"
)

var (
	ErrRequestInvalid = errors.New("request is invalid")

	ErrUsernameExisted  = errors.New("username is existed")
	ErrUsernameNotFound = errors.New("username is not found")

	ErrUserNotFound = errors.New("user not found")

	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	ErrClientInvalid  = errors.New("client is invalid")
	ErrClientNotFound = errors.New("client is not found")

	ErrIdPInvalid = errors.New("idp is invalid")

	ErrScopeInvalid = errors.New("invalid scope")

	ErrAuthorizationAccessDenied        = errors.New("access denined")
	ErrAuthorizationResponseTypeInvalid = errors.New("invalid response type")

	ErrTokenInvalidGrant     = errors.New("invalid grant")
	ErrTokenInvalidGrantType = errors.New("invalid grant type")
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
