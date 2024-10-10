package domain

import (
	"github.com/xybor/todennus-backend/domain/definition"
	"github.com/xybor/todennus-backend/pkg/scope"
)

var Actions, actionMap = scope.DefineAction[definition.Actions]()
var Resources, resourceMap = scope.DefineResource[definition.Resource]()
var ScopeEngine = scope.NewEngine(actionMap, resourceMap)

// Full permission
var _ = ScopeEngine.Define(Actions, Resources)

// User
var _ = ScopeEngine.Define(Actions.Read, Resources.User.AllowedScope).
	Description("grant read-only access to user's allowed scope")

// OAuth2 Client
var _ = ScopeEngine.Define(Actions.Read, Resources.Client.AllowedScope).
	Description("grant read-only access to client's allowed scope")
