package api

import (
	"fmt"
	"encoding/json"
	"net/http"
	"ylem_taskrunner/helpers"

	"github.com/gorilla/mux"
)

var ExampleHandler = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	response, err := json.Marshal(&struct {
		Message string
	}{
		Message: fmt.Sprintf("Hello, %s!", vars["name"]),
	})

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
