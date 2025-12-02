package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"log"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

//ParseJSON expects a pointer to a struct and decodes the JSON body of the request into it.
func ParseJSON(r *http.Request, payload any) error {

	if err := Validate.Struct(payload); err != nil {
		log.Printf("invalid request body: %v", err)
		return errors.New("invalid request body")
	}

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

// MakeJSONRequest creates an HTTP request with JSON body
func MakeJSONRequest(ctx context.Context, method, url string, body interface{}) (*http.Request, error) {
	var reqBody io.Reader
	
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}
	
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	return req, nil
}

//executes an HTTP request and returns the response body and status code
func ExecuteRequest(client *http.Client, req *http.Request) ([]byte, int, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
	}
	
	return body, resp.StatusCode, nil
}

//unmarshals a JSON response body with proper error handling
func UnmarshalJSONResponse[T any](body []byte, statusCode int, serviceName string) (*T, int, error) {
	var resp T
	if err := json.Unmarshal(body, &resp); err != nil {
		if statusCode != http.StatusOK {
			return nil, statusCode, fmt.Errorf("%s returned status %d: %s", serviceName, statusCode, string(body))
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &resp, statusCode, nil
}