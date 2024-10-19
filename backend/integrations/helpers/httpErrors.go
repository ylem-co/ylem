package helpers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type HttpErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Fields  string `json:"fields,omitempty"`
}

func HttpReturnErrorUnauthorized(w http.ResponseWriter) {
	err := HttpErrorResponse{
		Code:    401,
		Message: "Authorization Failed",
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorInternal(w http.ResponseWriter) {
	err := HttpErrorResponse{
		Code:    500,
		Message: "Internal Server Error",
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorServiceUnavailable(w http.ResponseWriter) {
	err := HttpErrorResponse{
		Code:    503,
		Message: "Service Unavailable. Please try again later",
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorNotFound(w http.ResponseWriter) {
	err := HttpErrorResponse{
		Code:    404,
		Message: "Not Found",
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorBadRequest(w http.ResponseWriter, Message string, fields *[]string) {
	err := HttpErrorResponse{
		Code:    400,
		Message: Message,
	}

	if fields != nil {
		err.Fields = strings.Join(*fields, ",")
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorBadUuidRequest(w http.ResponseWriter) {
	err := HttpErrorResponse{
		Code:    400,
		Message: "Invalid UUID",
		Fields:  "uuid",
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorForbidden(w http.ResponseWriter) {
	err := HttpErrorResponse{
		Code:    403,
		Message: "Forbidden",
	}

	HttpReturnError(w, err)
}

func HttpReturnErrorForbiddenQuotaExceeded(w http.ResponseWriter) {
	err := HttpErrorResponse{
		Code:    403,
		Message: "Your subscription doesn't allow creating more resources. Please upgrade your subscription plan or contact our support",
	}

	HttpReturnError(w, err)
}

func HttpReturnError(w http.ResponseWriter, err HttpErrorResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(err.Code)

	_ = json.NewEncoder(w).Encode(err)
}
