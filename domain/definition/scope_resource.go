package definition

import (
	"github.com/xybor/x/scope"
)

type Resource struct {
	*scope.BaseResource

	User   *UserResource
	Client *OAuth2ClientResource
}

type UserResource struct {
	*scope.BaseResource

	AllowedScope *scope.BaseResource `resource:"allowed_scope"`
	Role         *scope.BaseResource `resource:"role"`
}

type OAuth2ClientResource struct {
	*scope.BaseResource

	Owner        *scope.BaseResource `resource:"owner"`
	AllowedScope *scope.BaseResource `resource:"allowed_scope"`
}
