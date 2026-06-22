package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("failed to marshal json:", err)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	return err
}

func checkIDParam(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	idString := chi.URLParam(r, "id")
	if idString == "" {
		return uuid.Nil, fmt.Errorf("id parameter is required")
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid id parameter:%w", err)
	}

	return id, nil

}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError
		var invalidUnmarshalErr *json.InvalidUnmarshalError
		var maxBytesErr *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("body contains badly formed JSON (at character %d)", syntaxErr.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeErr):
			if unmarshalTypeErr.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeErr.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeErr.Offset)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown field %s", fieldName)

		case errors.As(err, &maxBytesErr):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesErr.Limit)

		// programmer error: data passed into readJSON wasn't a non-nil pointer
		case errors.As(err, &invalidUnmarshalErr):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("body must only contain a single JSON value")
	}

	return nil
}
