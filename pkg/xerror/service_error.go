package xerror

import (
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"
)

type Serverity int

const (
	ServerityDebug Serverity = iota
	ServerityInfo
	ServerityWarn
	ServerityCritical
)

type ServiceError struct {
	Serverity Serverity
	Err       error
	Message   string
}

func (err ServiceError) Error() string {
	return err.Err.Error()
}

func (err ServiceError) Unwrap() error {
	return err.Err
}

func Wrap(err error, serverity Serverity) ServiceError {
	return ServiceError{Err: err, Serverity: serverity}.WithMessage(err.Error())
}

func New(serverity Serverity, msg string, a ...any) ServiceError {
	msg = fmt.Sprintf(msg, a...)
	return ServiceError{Serverity: serverity, Err: errors.New(msg)}.WithMessage(msg)
}

func Debug(msg string, a ...any) ServiceError {
	return New(ServerityDebug, msg, a...)
}

func Info(msg string, a ...any) ServiceError {
	return New(ServerityInfo, msg, a...)
}

func Warn(msg string, a ...any) ServiceError {
	return New(ServerityWarn, msg, a...)
}

func Critical(msg string, a ...any) ServiceError {
	return New(ServerityCritical, msg, a...)
}

func WrapDebug(err error) ServiceError {
	return Wrap(err, ServerityDebug)
}

func WrapInfo(err error) ServiceError {
	return Wrap(err, ServerityInfo)
}

func WrapWarn(err error) ServiceError {
	return Wrap(err, ServerityWarn)
}

func WrapCritical(err error) ServiceError {
	return Wrap(err, ServerityCritical)
}

func (err ServiceError) WithMessage(msg string, a ...any) ServiceError {
	text := fmt.Sprintf(msg, a...)
	var size int
	r, size := utf8.DecodeRuneInString(text)

	if r != utf8.RuneError {
		msg = string(unicode.ToUpper(r)) + text[size:]
	}

	return ServiceError{
		Err:       err.Err,
		Serverity: err.Serverity,
		Message:   msg,
	}
}
