package helpers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HttpReturnErrorUnauthorized(w http.ResponseWriter) {
	err := Error{
		Code:    401,
		Message: "Authorization failed",
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorInternal(w http.ResponseWriter) {
	err := Error{
		Code:    500,
		Message: "Internal error",
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorBadRequest(w http.ResponseWriter) {
	err := Error{
		Code:    400,
		Message: "Bad Request",
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorForbidden(w http.ResponseWriter) {
	err := Error{
		Code:    403,
		Message: "Forbidden",
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorNotFound(w http.ResponseWriter) {
	err := Error{
		Code:    404,
		Message: "Not Found",
	}

	HttpReturnError(w, err)
}

func HttpReturnError(w http.ResponseWriter, err Error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(err.Code)
	
	error := json.NewEncoder(w).Encode(err)
	if error != nil {
		log.Error(error)
	}
}
