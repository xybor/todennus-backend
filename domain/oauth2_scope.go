package domain

import (
	"github.com/xybor/todennus-backend/domain/definition"
	"github.com/xybor/todennus-backend/pkg/scope"
)

var Actions, actionMap = scope.DefineAction[definition.Actions]()
var Resources, resourceMap = scope.DefineResource[definition.Resource]()
var ScopeEngine = scope.NewEngine(actionMap, resourceMap)
