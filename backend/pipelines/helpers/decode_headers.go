package helpers

import "encoding/json"

func DecodeHeaders(hjson string) (map[string]string, error) {
	headers := make(map[string]string)
	if hjson == "" {
		return headers, nil
	}

	err := json.Unmarshal([]byte(hjson), &headers)

	return headers, err
}
