package scope

import (
	"fmt"
	"reflect"
	"strings"
)

type Actioner interface {
	fullAction() string
	String() string
	IsSubset(Actioner) bool
}

type BaseAction struct {
	full    string
	current string
}

func newBaseAction(parent, current string) BaseAction {
	full := current
	if parent != "" {
		full = fmt.Sprintf("%s.%s", parent, current)
	}

	return BaseAction{current: current, full: full}
}

func (action BaseAction) fullAction() string {
	return action.full
}

func (action BaseAction) String() string {
	return action.current
}

func (action BaseAction) IsSubset(another Actioner) bool {
	if action.fullAction() == another.fullAction() {
		return true
	}

	return strings.Contains(action.fullAction(), another.fullAction()+".")
}

func DefineAction[A any]() (A, map[string]Actioner) {
	var action A
	m := map[string]Actioner{}

	fulfill(
		m,
		"",                           // default as empty
		"*",                          // default as empty
		reflect.ValueOf(&action),     // rvalue
		reflect.TypeOf(BaseAction{}), // basetype
		newBaseAction,                // baseFunc
		"action",                     // tag
	)
	return action, m
}
