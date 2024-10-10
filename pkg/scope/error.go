package scope

import "errors"

var (
	ErrScopeInvalid    = errors.New("scope is invalid")
	ErrScopeNotDefined = errors.New("scope has not been defined yet")
)
