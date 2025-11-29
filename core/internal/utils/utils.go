package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

//ParseJSON expects a pointer to a struct and decodes the JSON body of the request into it.
func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	
	return json.NewDecoder(r.Body).Decode(payload)	
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}