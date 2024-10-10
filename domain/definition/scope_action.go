package definition

import "github.com/xybor/todennus-backend/pkg/scope"

type Actions struct {
	scope.BaseAction

	Read  scope.BaseAction
	Write WriteAction
}

type WriteAction struct {
	scope.BaseAction

	Create scope.BaseAction
	Update scope.BaseAction
	Delete scope.BaseAction
}
