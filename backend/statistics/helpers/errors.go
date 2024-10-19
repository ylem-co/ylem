package helpers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

func HttpReturnErrorUnauthorized(w http.ResponseWriter) {
	err := ErrorResponse{
		Code:    401,
		Message: "Authorization failed",
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorInternal(w http.ResponseWriter) {
	err := ErrorResponse{
		Code:    500,
		Message: "Internal error",
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorBadRequest(w http.ResponseWriter, errors []error) {

	eStrings := make([]string, 0)
	for _, e := range errors {
		eStrings = append(eStrings, e.Error())
	}

	err := ErrorResponse{
		Code:    400,
		Message: "Bad Request",
		Errors:  eStrings,
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorForbidden(w http.ResponseWriter) {
	err := ErrorResponse{
		Code:    403,
		Message: "Forbidden",
	}

	HttpReturnError(w, err)
}

func HttpReturnError(w http.ResponseWriter, err ErrorResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(err.Code)

	error := json.NewEncoder(w).Encode(err)
	if error != nil {
		log.Error(error)
	}
}
