package utils

import (
	"encoding/json"
	"net/http"
)

func ParseJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func DecodeJSON(r *http.Request, v any) error {
	return ParseJSON(r, v)
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	WriteJSON(w, status, errorResponse{Error: err.Error()})
}


