package handlers

import (
	"encoding/json"
	"net/http"
)

func JSONError(w http.ResponseWriter, error string, code int) {
	type JSONError struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(JSONError{Error: error})
}
