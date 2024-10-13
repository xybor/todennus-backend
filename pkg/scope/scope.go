package scope

import (
	"fmt"
)

type scopeRelationship int

const (
	scopeRelationshipNone scopeRelationship = iota
	scopeRelationshipSubset
	scopeRelationshipSuperset
)

type Scoper interface {
	Contains(another Scoper) bool

	String() string

	IsUndefined() bool
}

var _ Scoper = Scope{}

type Scope struct {
	action   Actioner
	resource Resourcer
}

func New(action Actioner, resource Resourcer) Scope {
	return Scope{
		action:   action,
		resource: resource,
	}
}

func (scope Scope) String() string {
	if scope.resource.String() == "" {
		return scope.action.String()
	}

	return fmt.Sprintf("%s:%s", scope.action.String(), scope.resource.String())
}

func (scope Scope) Contains(another Scoper) bool {
	anotherScope, ok := another.(Scope)
	if !ok {
		return false
	}

	return anotherScope.action.IsSubset(scope.action) && anotherScope.resource.IsSubset(scope.resource)
}

func (scope Scope) AsScopes() Scopes {
	return NewScopes(scope)
}

func (scope Scope) IsUndefined() bool {
	return false
}

type UndefinedScope string

func (scope UndefinedScope) String() string {
	return string(scope)
}

func (scope UndefinedScope) Contains(another Scoper) bool {
	return false
}

func (scope UndefinedScope) IsUndefined() bool {
	return true
}

func relationship(scope, another Scoper) scopeRelationship {
	switch {
	case scope.Contains(another):
		return scopeRelationshipSuperset
	case another.Contains(scope):
		return scopeRelationshipSubset
	default:
		return scopeRelationshipNone
	}
}
