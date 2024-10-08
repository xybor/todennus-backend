package domain

import (
	"errors"
	"fmt"
)

var (
	ErrKnown              = errors.New("")
	ErrUnknownCritical    = errors.New("internal error")
	ErrUnknownRecoverable = errors.New("internal error")

	ErrUsernameInvalid    = fmt.Errorf("%w%s", ErrKnown, "invalid username")
	ErrDisplayNameInvalid = fmt.Errorf("%w%s", ErrKnown, "invalid display name")
	ErrPasswordInvalid    = fmt.Errorf("%w%s", ErrKnown, "invalid password")

	ErrClientInvalid      = fmt.Errorf("%w%s", ErrKnown, "invalid client")
	ErrClientUnauthorized = fmt.Errorf("%w%s", ErrKnown, "unauthorized client")
	ErrClientNameInvalid  = fmt.Errorf("%w%s", ErrKnown, "invalid client name")
)

func Wrap(err error, format string, a ...any) error {
	return fmt.Errorf("%w: %s", err, fmt.Sprintf(format, a...))
}
