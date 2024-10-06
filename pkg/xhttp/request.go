package xhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

const (
	ContentTypeApplicationJSON    = "application/json"
	ContentTypeXWWWFormUrlEncoded = "application/x-www-form-urlencoded"
)

func ParseRequest[T any](req *http.Request) (T, error) {
	var t T

	contentType := req.Header.Get("content-type")
	switch contentType {
	case ContentTypeApplicationJSON:
		err := json.NewDecoder(req.Body).Decode(&t)
		if err != nil {
			return t, fmt.Errorf("%w%s", ErrBadRequest, err.Error())
		}
		return t, nil
	case ContentTypeXWWWFormUrlEncoded:
		return parseURLEncodedFormData[T](req)
	default:
		return t, fmt.Errorf("%w%s", ErrBadRequest, fmt.Sprintf("not support content type %s", contentType))
	}
}

func parseURLEncodedFormData[T any](req *http.Request) (T, error) {
	var result T
	resultType := reflect.TypeOf(result)
	resultVal := reflect.ValueOf(&result).Elem()

	for i := 0; i < resultType.NumField(); i++ {
		field := resultType.Field(i)
		fieldName := field.Name

		fieldTag := field.Tag.Get("form")
		if fieldTag == "" {
			return result, fmt.Errorf("not found form tag of field %s", fieldName)
		}

		formValue := req.FormValue(fieldTag)
		if formValue == "" {
			continue
		}

		fieldVal := resultVal.FieldByName(fieldName)

		switch field.Type.Kind() {
		case reflect.String:
			fieldVal.SetString(formValue)

		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			intFormVal, err := strconv.ParseInt(formValue, 10, 64)
			if err != nil {
				return result, fmt.Errorf("%w%s", ErrBadRequest, err.Error())
			}
			fieldVal.SetInt(intFormVal)

		case reflect.Float32, reflect.Float64:
			floatFormVal, err := strconv.ParseFloat(formValue, 64)
			if err != nil {
				return result, fmt.Errorf("%w%s", ErrBadRequest, err.Error())
			}
			fieldVal.SetFloat(floatFormVal)

		case reflect.Bool:
			boolFormVal, err := strconv.ParseBool(formValue)
			if err != nil {
				return result, fmt.Errorf("%w:%s", ErrBadRequest, err.Error())
			}
			fieldVal.SetBool(boolFormVal)

		default:
			return result, fmt.Errorf("not support type %s", field.Type.Kind())
		}
	}

	return result, nil
}
