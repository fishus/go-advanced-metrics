package middleware

import (
	"encoding/json"
	"net/http"
)

func JSONError(w http.ResponseWriter, error string, code int) {
	type JSONError struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(JSONError{Error: error})
}
