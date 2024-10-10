package scope

import (
	"fmt"
	"reflect"
	"strings"
)

type Resourcer interface {
	String() string
	IsSubset(Resourcer) bool
}

type BaseResource struct {
	full string
}

func newBaseResource(parent, current string) BaseResource {
	if parent == "" {
		return BaseResource{full: current}
	}

	return BaseResource{full: fmt.Sprintf("%s.%s", parent, current)}
}

func (resource BaseResource) String() string {
	return resource.full
}

func (resource BaseResource) IsSubset(another Resourcer) bool {
	if resource.String() == another.String() {
		return true
	}

	if another.String() == "" {
		return true
	}

	return strings.Contains(resource.String(), another.String()+".")
}

func DefineResource[T any]() (T, map[string]Resourcer) {
	var resource T
	m := map[string]Resourcer{}

	fulfill(
		m,
		"",                             // default as empty
		"",                             // root value
		reflect.ValueOf(&resource),     // rvalue
		reflect.TypeOf(BaseResource{}), // basetype
		newBaseResource,                // baseFunc
		"resource",                     // tag
	)

	return resource, m
}
