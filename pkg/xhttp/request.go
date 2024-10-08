package xhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

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
			if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
				return defaultT, fmt.Errorf("%w%s", ErrBadRequest, err.Error())
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
	return parse(obj, req, "form", func(req *http.Request, fieldName string) string {
		return req.FormValue(fieldName)
	})
}

func parseURLQuery[T any](obj T, req *http.Request) error {
	return parse(obj, req, "query", func(r *http.Request, s string) string {
		return req.URL.Query()[s][0]
	})
}

func parseURLParameter[T any](obj T, req *http.Request) error {
	return parse(obj, req, "param", func(r *http.Request, s string) string {
		return chi.URLParam(r, s)
	})
}

func parse[T any](obj T, req *http.Request, tagName string, fieldVal func(*http.Request, string) string) error {
	objType := reflect.TypeOf(obj).Elem()
	objVal := reflect.ValueOf(obj).Elem()

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldName := field.Name

		tagValue := field.Tag.Get(tagName)
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
			fieldVal.SetString(fieldValue)

		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			intFormVal, err := strconv.ParseInt(fieldValue, 10, 64)
			if err != nil {
				return fmt.Errorf("%w%s", ErrBadRequest, err.Error())
			}
			fieldVal.SetInt(intFormVal)

		case reflect.Float32, reflect.Float64:
			floatFormVal, err := strconv.ParseFloat(fieldValue, 64)
			if err != nil {
				return fmt.Errorf("%w%s", ErrBadRequest, err.Error())
			}
			fieldVal.SetFloat(floatFormVal)

		case reflect.Bool:
			boolFormVal, err := strconv.ParseBool(fieldValue)
			if err != nil {
				return fmt.Errorf("%w:%s", ErrBadRequest, err.Error())
			}
			fieldVal.SetBool(boolFormVal)

		default:
			return fmt.Errorf("not support type %s", field.Type.Kind())
		}
	}

	return nil
}
