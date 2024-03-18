package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/sivistrukov/vk-assigment/internal/entrypoints/http/schemas"
)

// validateRequestBody checks if the request body is valid according to
// the given schema. If the request body is not valid, it returns an error.
func validateRequestBody(r *http.Request, validate *validator.Validate, schema interface{}) error {
	err := parseRequestBody(r, schema)
	if err != nil {
		return errors.New("invalid request body")
	}

	err = validate.Struct(schema)
	if err != nil {
		return fmt.Errorf("invalid request body: %v", err)
	}

	return nil
}

// parseRequestBody reads the request body and unmarshals
// it into the given schema pointer.
func parseRequestBody(r *http.Request, schema any) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}

	err = json.Unmarshal(data, schema)
	if err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}

	return nil
}

// internalError is alias for writeJson.
func internalError(w http.ResponseWriter) {
	_ = writeJson(
		w,
		schemas.ErrorResponse{Error: "internal server error"},
		http.StatusInternalServerError,
	)
}

// writeJson writes the given data as JSON to the response writer
// with the given status code. It sets the Content-Type header
// to "application/json" and writes the response body.
// If an error occurs while writing the response, it is returned.
func writeJson(w http.ResponseWriter, data interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	w.WriteHeader(status)
	_, err = w.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write response body: %w", err)
	}
	return nil
}
