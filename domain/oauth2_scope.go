package domain

import (
	"github.com/xybor/todennus-backend/domain/definition"
	"github.com/xybor/todennus-backend/pkg/scope"
)

var Actions, actionMap = scope.DefineAction[definition.Actions]()
var Resources, resourceMap = scope.DefineResource[definition.Resource]()
var ScopeEngine = scope.NewEngine(actionMap, resourceMap)

func init() {
	// Full permission
	ScopeEngine.Define(Actions, Resources)       // *
	ScopeEngine.Define(Actions.Read, Resources)  // read
	ScopeEngine.Define(Actions.Write, Resources) // write

	// User
	ScopeEngine.Define(Actions.Read, Resources.User).Description("grant read-only access to user information")
	ScopeEngine.Define(Actions.Read, Resources.User.Role).Description("grant read-only access to user's role")
	ScopeEngine.Define(Actions.Read, Resources.User.AllowedScope).Description("grant read-only access to user's allowed scope")

	// OAuth2 Client
	ScopeEngine.Define(Actions.Read, Resources.Client).Description("grant read-only access to client's information")
	ScopeEngine.Define(Actions.Read, Resources.Client.AllowedScope).Description("grant read-only access to client's allowed scope")
	ScopeEngine.Define(Actions.Write, Resources.Client).Description("grant write permission to OAuth2 Client")
	ScopeEngine.Define(Actions.Write.Create, Resources.Client).Description("allow user to create OAuth2 Client")
}
