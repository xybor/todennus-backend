package scope

import (
	"strings"
)

type Engine struct {
	actionMap   map[string]Actioner
	resourceMap map[string]Resourcer
}

func NewEngine(actionMap map[string]Actioner, resourceMap map[string]Resourcer) Engine {
	return Engine{
		actionMap:   actionMap,
		resourceMap: resourceMap,
	}
}

func (engine Engine) ParseScope(s string) Scoper {
	s = strings.Trim(s, " ")
	if s == "" {
		return UndefinedScope(s)
	}

	actionStr, resourceStr, found := strings.Cut(s, ":")
	if !found {
		actionStr = s
		resourceStr = ""
	}

	action, ok := engine.actionMap[actionStr]
	if !ok {
		return UndefinedScope(s)
	}

	resource, ok := engine.resourceMap[resourceStr]
	if !ok {
		return UndefinedScope(s)
	}

	scope := New(action, resource)
	return scope
}

func (engine Engine) ParseScopes(s string) Scopes {
	s = strings.Trim(s, " ")
	if s == "" {
		return Scopes{}
	}

	scopesStr := strings.Split(s, " ")
	scopes := Scopes{}
	for _, str := range scopesStr {
		scopes = append(scopes, engine.ParseScope(str))
	}

	return scopes
}
