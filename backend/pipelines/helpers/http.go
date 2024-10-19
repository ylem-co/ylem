package helpers

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HttpResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if body == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(body)
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

func HttpReturnErrorForbiddenQuotaExceeded(w http.ResponseWriter) {
	err := Error{
		Code:    403,
		Message: "Your subscription doesn't allow creating more resources. Please upgrade your subscription plan or contact our support",
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

func HttpReturnErrorConflict(w http.ResponseWriter) {
	err := Error{
		Code:    http.StatusConflict,
		Message: "Conflict",
	}

	HttpReturnError(w, err)
}

func HttpReturnError(w http.ResponseWriter, err Error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(err.Code)
	
	_ = json.NewEncoder(w).Encode(err)
}
