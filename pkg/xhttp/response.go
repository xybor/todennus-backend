package xhttp

import (
	"encoding/json"
	"net/http"
)

func WriteResponseJSON(w http.ResponseWriter, code int, obj any) error {
	jsonString, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", ContentTypeApplicationJSON)
	w.WriteHeader(code)
	_, err = w.Write(jsonString)
	return err
}
