package domain

import (
	"errors"
	"fmt"
)

var (
	ErrKnown              = errors.New("")
	ErrUnknownCritical    = errors.New("internal error")
	ErrUnknownRecoverable = errors.New("internal error")
	ErrInvalidUsername    = fmt.Errorf("%w%s", ErrKnown, "invalid username")
	ErrInvalidDisplayName = fmt.Errorf("%w%s", ErrKnown, "invalid display name")
	ErrInvalidPassword    = fmt.Errorf("%w%s", ErrKnown, "invalid password")
)

func Wrap(err error, format string, a ...any) error {
	return fmt.Errorf("%w: %s", err, fmt.Sprintf(format, a...))
}
