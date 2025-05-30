package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// ParseQueryInt safely parses ?key= from the URL, returning a default if invalid
func ParseQueryInt(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if v, err := strconv.Atoi(val); err == nil {
		return v
	}
	return defaultVal
}

// WriteJSON writes the given status and data as JSON
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}