package scope

import (
	"fmt"
	"reflect"
	"strings"
)

type Stringer interface {
	String() string
}

func fulfill[T Stringer, R any](
	m map[string]R,
	parent string,
	current string,
	rvalue reflect.Value,
	baseType reflect.Type,
	baseFunc func(string, string) T,
	tagName string,
) {
	rtype := rvalue.Elem().Type()
	rvalue = rvalue.Elem()
	hasBase := false

	for i := range rtype.NumField() {
		rfield := rtype.Field(i)
		fieldValue := rfield.Tag.Get(tagName)
		if fieldValue == "" {
			fieldValue = strings.ToLower(rfield.Name)
		}

		currentParent := current
		if parent != "" {
			currentParent = strings.Join([]string{parent, current}, ".")
		}

		if rfield.Type == baseType {
			var base T
			if rfield.Anonymous {
				hasBase = true
				base = baseFunc(parent, current)
			} else {
				base = baseFunc(currentParent, fieldValue)
			}

			rvalue.FieldByName(rfield.Name).Set(reflect.ValueOf(base))

			if rfield.Anonymous {
				m[base.String()] = rvalue.Interface().(R)
			} else {
				m[base.String()] = rvalue.FieldByName(rfield.Name).Interface().(R)
			}
		} else {
			fulfill(
				m,
				currentParent, fieldValue,
				rvalue.FieldByName(rfield.Name).Addr(),
				baseType, baseFunc,
				tagName,
			)
		}
	}

	if !hasBase {
		panic(fmt.Sprintf("Please add base for %s", current))
	}
}
