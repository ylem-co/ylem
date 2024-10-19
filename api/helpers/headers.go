package helpers

import (
	"encoding/json"
	"net/http"
)

func DecodeHeaders(hjson string) (map[string]string, error) {
	headers := make(map[string]string)
	if hjson == "" {
		return headers, nil
	}

	err := json.Unmarshal([]byte(hjson), &headers)

	return headers, err
}

func SetHeaders(w http.ResponseWriter, headers http.Header) {
	for header, vals := range headers {
		for _, v := range vals {
			w.Header().Add(header, v)
		}
	}
}
