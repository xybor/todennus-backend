package xerror

import "errors"

func Is(err error, target error, otherTargets ...error) bool {
	if errors.Is(err, target) {
		return true
	}

	for _, target := range otherTargets {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func Message(err error) string {
	var serviceErr ServiceError
	if errors.As(err, &serviceErr) {
		return serviceErr.Message
	}

	return err.Error()
}
