package util

import (
	"encoding/json"
	"net/http"
)

func ResponseOk(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}

func ResponseError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(
		map[string]string{
			"error": message,
		})
}
