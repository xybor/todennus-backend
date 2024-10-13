package xhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

const (
	ContentTypeApplicationJSON    = "application/json"
	ContentTypeXWWWFormUrlEncoded = "application/x-www-form-urlencoded"
)

func ParseRequest[T any](req *http.Request) (T, error) {
	var defaultT T
	var t T

	if err := parseURLParameter(&t, req); err != nil {
		return t, err
	}

	switch req.Method {
	case http.MethodGet:
		if err := parseURLQuery(&t, req); err != nil {
			return defaultT, err
		}

		return t, nil

	case http.MethodPost, http.MethodPut, http.MethodDelete:
		contentType := req.Header.Get("content-type")
		switch contentType {
		case ContentTypeApplicationJSON:
			if err := parseJSONBody(&t, req); err != nil {
				return defaultT, err
			}
			return t, nil

		case ContentTypeXWWWFormUrlEncoded:
			if err := parseURLEncodedFormData(&t, req); err != nil {
				return defaultT, err
			}

			return t, nil

		default:
			return defaultT, fmt.Errorf("%w%s", ErrBadRequest, fmt.Sprintf("not support content type %s", contentType))
		}

	default:
		return defaultT, fmt.Errorf("%w%s", ErrBadRequest, fmt.Sprintf("not support method %s", req.Method))
	}
}

func parseURLEncodedFormData[T any](obj T, req *http.Request) error {
	return parse(obj, req, "form", func(req *http.Request, fieldName string) any {
		return req.FormValue(fieldName)
	}, false)
}

func parseURLQuery[T any](obj T, req *http.Request) error {
	return parse(obj, req, "query", func(r *http.Request, s string) any {
		return req.URL.Query()[s][0]
	}, false)
}

func parseURLParameter[T any](obj T, req *http.Request) error {
	return parse(obj, req, "param", func(r *http.Request, s string) any {
		return chi.URLParam(r, s)
	}, false)
}

func parseJSONBody[T any](obj T, req *http.Request) error {
	m := map[string]any{}
	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		return fmt.Errorf("%winvalid json", ErrBadRequest)
	}

	return parse(obj, req, "json", func(r *http.Request, s string) any {
		return m[s]
	}, true)
}

func any2int(a any, strict bool) (int64, error) {
	switch t := a.(type) {
	case int:
		return int64(t), nil
	case int8:
		return int64(t), nil
	case int16:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case int64:
		return int64(t), nil
	case float32:
		if t == float32(int64(t)) {
			return int64(t), nil
		}

		return 0, fmt.Errorf("%wexpected int, got float", ErrBadRequest)
	case float64:
		if t == float64(int64(t)) {
			return int64(t), nil
		}

		return 0, fmt.Errorf("%wexpected int, got float", ErrBadRequest)
	case string:
		if !strict {
			if n, err := strconv.ParseInt(t, 10, 64); err == nil {
				return n, nil
			}
		}

		return 0, fmt.Errorf("%wexpected int, got string", ErrBadRequest)
	case bool:
		return 0, fmt.Errorf("%wexpected int, got bool", ErrBadRequest)
	default:
		return 0, fmt.Errorf("not handle for %T", t)
	}
}

func any2float(a any, strict bool) (float64, error) {
	switch t := a.(type) {
	case int:
		return float64(t), nil
	case int8:
		return float64(t), nil
	case int16:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case float32:
		return float64(t), nil
	case float64:
		return float64(t), nil
	case string:
		if !strict {
			if n, err := strconv.ParseFloat(t, 64); err == nil {
				return n, nil
			}
		}

		return 0, fmt.Errorf("%wexpected float, got string", ErrBadRequest)
	case bool:
		return 0, fmt.Errorf("%wexpected float, got bool", ErrBadRequest)
	default:
		return 0, fmt.Errorf("not handle for %T", t)
	}
}

func any2bool(a any, strict bool) (bool, error) {
	switch t := a.(type) {
	case int:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, fmt.Errorf("%wexpected bool, got number", ErrBadRequest)
	case int8:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, fmt.Errorf("%wexpected bool, got number", ErrBadRequest)
	case int16:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, fmt.Errorf("%wexpected bool, got number", ErrBadRequest)
	case int32:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, fmt.Errorf("%wexpected bool, got number", ErrBadRequest)
	case int64:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, fmt.Errorf("%wexpected bool, got number", ErrBadRequest)
	case float32:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, fmt.Errorf("%wexpected bool, got number", ErrBadRequest)
	case float64:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, fmt.Errorf("%wexpected bool, got number", ErrBadRequest)
	case string:
		if !strict {
			if b, err := strconv.ParseBool(t); err == nil {
				return b, nil
			}
		}

		return false, fmt.Errorf("%wexpected bool, got string", ErrBadRequest)
	case bool:
		return t, nil
	default:
		return false, fmt.Errorf("not handle for %T", t)
	}
}

func any2string(a any, strict bool) (string, error) {
	switch t := a.(type) {
	case int:
		if !strict {
			return strconv.FormatInt(int64(t), 10), nil
		}
		return "", fmt.Errorf("%wexpected string, got number", ErrBadRequest)
	case int8:
		if !strict {
			return strconv.FormatInt(int64(t), 10), nil
		}
		return "", fmt.Errorf("%wexpected string, got number", ErrBadRequest)
	case int16:
		if !strict {
			return strconv.FormatInt(int64(t), 10), nil
		}
		return "", fmt.Errorf("%wexpected string, got number", ErrBadRequest)
	case int32:
		if !strict {
			return strconv.FormatInt(int64(t), 10), nil
		}
		return "", fmt.Errorf("%wexpected string, got number", ErrBadRequest)
	case int64:
		if !strict {
			return strconv.FormatInt(int64(t), 10), nil
		}
		return "", fmt.Errorf("%wexpected string, got number", ErrBadRequest)
	case float32:
		if !strict {
			return strconv.FormatFloat(float64(t), 'f', -1, 32), nil
		}
		return "", fmt.Errorf("%wexpected string, got number", ErrBadRequest)
	case float64:
		if !strict {
			return strconv.FormatFloat(float64(t), 'f', -1, 64), nil
		}
		return "", fmt.Errorf("%wexpected string, got number", ErrBadRequest)
	case string:
		return t, nil
	case bool:
		if !strict {
			return strconv.FormatBool(t), nil
		}
		return "", fmt.Errorf("%wexpected string, got bool", ErrBadRequest)
	default:
		return "", fmt.Errorf("not handle for %T", t)
	}
}

func parse[T any](obj T, req *http.Request, tagName string, fieldVal func(*http.Request, string) any, strict bool) error {
	objType := reflect.TypeOf(obj).Elem()
	objVal := reflect.ValueOf(obj).Elem()

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldName := field.Name

		tagValue := field.Tag.Get(tagName)
		tagValue, _, _ = strings.Cut(tagValue, ",")
		if tagValue == "" {
			continue
		}

		fieldValue := fieldVal(req, tagValue)
		if fieldValue == "" {
			continue
		}

		fieldVal := objVal.FieldByName(fieldName)

		switch field.Type.Kind() {
		case reflect.String:
			s, err := any2string(fieldValue, strict)
			if err != nil {
				return fmt.Errorf("%s: %w", tagValue, err)
			}

			fieldVal.SetString(s)

		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			intFormVal, err := any2int(fieldValue, strict)
			if err != nil {
				return fmt.Errorf("%s: %w", tagValue, err)
			}
			fieldVal.SetInt(intFormVal)

		case reflect.Float32, reflect.Float64:
			floatFormVal, err := any2float(fieldValue, strict)
			if err != nil {
				return fmt.Errorf("%s: %w", tagValue, err)
			}
			fieldVal.SetFloat(floatFormVal)

		case reflect.Bool:
			boolFormVal, err := any2bool(fieldValue, strict)
			if err != nil {
				return fmt.Errorf("%s: %w", tagValue, err)
			}
			fieldVal.SetBool(boolFormVal)

		default:
			return fmt.Errorf("not support type %s", field.Type.Kind())
		}
	}

	return nil
}
