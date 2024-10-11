package scope

import (
	"fmt"
	"strings"
)

type scopeRelationship int

const (
	scopeRelationshipNone scopeRelationship = iota
	scopeRelationshipSubset
	scopeRelationshipSuperset
)

var NoScope = Scope{}

type Scope struct {
	action   Actioner
	resource Resourcer
}

func NewONLYFORTEST(action Actioner, resource Resourcer) Scope {
	return Scope{
		action:   action,
		resource: resource,
	}
}

func newAction(action Actioner, resource Resourcer) Scope {
	return Scope{action: action, resource: resource}
}

func (scope Scope) String() string {
	if scope.resource.String() == "" {
		return scope.action.String()
	}

	return fmt.Sprintf("%s:%s", scope.action.String(), scope.resource.String())
}

func (scope Scope) IsSubset(another Scope) bool {
	return scope.action.IsSubset(another.action) && scope.resource.IsSubset(another.resource)
}

func (scope Scope) Contains(target Scope) bool {
	return target.IsSubset(scope)
}

func (scope Scope) AsScopes() Scopes {
	return NewScopes(scope)
}

func (scope Scope) relationship(another Scope) scopeRelationship {
	switch {
	case scope.IsSubset(another):
		return scopeRelationshipSubset
	case scope.Contains(another):
		return scopeRelationshipSuperset
	default:
		return scopeRelationshipNone
	}
}

type Scopes []Scope

func NewScopes(scopes ...Scope) Scopes {
	return scopes
}

func (scopes Scopes) String() string {
	scopeValues := []string{}
	for i := range scopes {
		scopeValues = append(scopeValues, scopes[i].String())
	}

	return strings.Join(scopeValues, " ")
}

func (scopes Scopes) Contains(target Scope) bool {
	for _, scope := range scopes {
		if scope.Contains(target) {
			return true
		}
	}

	return false
}

func (scopes Scopes) Intersect(another Scopes) Scopes {
	result := Scopes{}
	for _, scope := range scopes {
		for _, targetScope := range another {
			if relationship := scope.relationship(targetScope); relationship != scopeRelationshipNone {
				if relationship == scopeRelationshipSubset {
					result = append(result, scope)
				} else {
					result = append(result, targetScope)
				}

				break
			}
		}
	}

	return result
}
