package definition

import (
	"github.com/xybor/todennus-backend/pkg/scope"
)

type Resource struct {
	scope.BaseResource

	User   UserResource
	Client OAuth2ClientResource
}

type UserResource struct {
	scope.BaseResource

	DisplayName  scope.BaseResource
	AllowedScope scope.BaseResource `resource:"allowed_scope"`
}

type OAuth2ClientResource struct {
	scope.BaseResource

	Owner        scope.BaseResource
	AllowedScope scope.BaseResource `resource:"allowed_scope"`
}
