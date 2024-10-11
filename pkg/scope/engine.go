package scope

import (
	"fmt"
	"strings"
)

type Engine struct {
	actionMap   map[string]Actioner
	resourceMap map[string]Resourcer

	definedScopes map[string]string
}

type definitionWrapper struct {
	engine *Engine
	scope  string
}

func (w definitionWrapper) Description(s string) bool {
	if _, ok := w.engine.definedScopes[w.scope]; !ok {
		panic(fmt.Sprintf("add description to a non-existed scope %s", w.scope))
	}

	w.engine.definedScopes[w.scope] = s
	return true
}

func NewEngine(actionMap map[string]Actioner, resourceMap map[string]Resourcer) Engine {
	return Engine{
		actionMap:     actionMap,
		resourceMap:   resourceMap,
		definedScopes: make(map[string]string),
	}
}

func (engine *Engine) Define(action Actioner, resource Resourcer) definitionWrapper {
	str := newAction(action, resource).String()
	if _, ok := engine.definedScopes[str]; ok {
		panic(fmt.Sprintf("add an existed scope %s", str))
	}

	engine.definedScopes[str] = ""
	return definitionWrapper{engine: engine, scope: str}
}

func (engine *Engine) New(action Actioner, resource Resourcer) Scope {
	scope := newAction(action, resource)
	if !engine.IsDefined(scope) {
		panic(fmt.Errorf("%w: %s", ErrScopeNotDefined, scope))
	}

	return scope
}

func (engine *Engine) IsDefined(scope Scope) bool {
	_, ok := engine.definedScopes[scope.String()]
	return ok
}

func (engine Engine) ParseScope(s string) (Scope, error) {
	if s == "" {
		return Scope{}, fmt.Errorf("%w: require a non-empty string", ErrScopeInvalid)
	}

	actionStr, resourceStr, found := strings.Cut(s, ":")
	if !found {
		actionStr = s
		resourceStr = ""
	}

	action, ok := engine.actionMap[actionStr]
	if !ok {
		return Scope{}, fmt.Errorf("%w: not found action %s", ErrScopeInvalid, actionStr)
	}

	resource, ok := engine.resourceMap[resourceStr]
	if !ok {
		return Scope{}, fmt.Errorf("%w: not found resource %s", ErrScopeInvalid, resourceStr)
	}

	scope := newAction(action, resource)
	if !engine.IsDefined(scope) {
		return scope, fmt.Errorf("%w: %s", ErrScopeNotDefined, scope)
	}

	return scope, nil
}

func (engine Engine) ParseScopes(s string) (Scopes, error) {
	s = strings.Trim(s, " ")
	if s == "" {
		return Scopes{}, nil
	}

	scopesStr := strings.Split(s, " ")
	scopes := Scopes{}
	for _, str := range scopesStr {
		scope, err := engine.ParseScope(str)
		if err != nil {
			return nil, err
		}

		scopes = append(scopes, scope)
	}

	return scopes, nil
}
