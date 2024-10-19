package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang/gddo/httputil/header"
)

type malformedRequest struct {
	Status int
	Msg    map[string]string
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) *malformedRequest {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := map[string]string{"message": "Content-Type header is not application/json"}
			return &malformedRequest{Status: http.StatusUnsupportedMediaType, Msg: msg}
		}
	}

	// Request body can be read only once.
	// Here we create a copy and assign it to r.Body
	body := http.MaxBytesReader(w, r.Body, 1048576)
	bodyBytes, _ := io.ReadAll(body)

	body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := map[string]string{"message": "Badly-formed JSON"}
			return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := map[string]string{"message": "Badly-formed JSON"}
			return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := map[string]string{"message": "Invalid fields", "fields": unmarshalTypeError.Field}
			return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := map[string]string{"message": fmt.Sprintf("Unknown field %s", fieldName)}
			return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.EOF):
			msg := map[string]string{"message": "Request body must not be empty"}
			return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case err.Error() == "http: request body too large":
			msg := map[string]string{"message": "Request body must not be larger than 1MB"}
			return &malformedRequest{Status: http.StatusRequestEntityTooLarge, Msg: msg}

		default:
			return nil
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := map[string]string{"message": "Request body must only contain a single JSON object"}
		return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	return runGoValidatorValidations(dst)
}

func runGoValidatorValidations(dst interface{}) *malformedRequest {
	validatableStruct, ok := dst.(validation.Validatable)
	if !ok {
		return nil
	}
	err := validatableStruct.Validate()
	if err != nil {
		invalidFields := make([]string, 0)

		for fieldName := range err.(validation.Errors) {
			invalidFields = append(invalidFields, fieldName)
		}

		msg := map[string]string{"message": "Invalid fields", "fields": strings.Join(invalidFields, ",")}
		return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	return nil
}
