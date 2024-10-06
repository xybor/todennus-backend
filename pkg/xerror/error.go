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
