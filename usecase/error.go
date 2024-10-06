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
	ErrInvalidGrantType        = errors.New("grant type is not supported by provider")
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
