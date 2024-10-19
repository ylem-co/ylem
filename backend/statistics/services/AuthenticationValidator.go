package services

import (
	"bytes"
	"io"
	"strings"
	"encoding/json"
	"net/http"
	"ylem_statistics/config"

	log "github.com/sirupsen/logrus"
)

type AuthenticationData struct {
	Email string
	Uuid  string
}

func CollectAuthenticationData(Token string) *AuthenticationData {
	config := config.Cfg()

	url := config.AuthorizationCheckUrl

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Error(err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Token)

	client := &http.Client{} // that looks fishy
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		log.Println(err.Error())

		return nil
	}

	var AuthData AuthenticationData
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())

		return nil
	}

	err = json.Unmarshal(bodyBytes, &AuthData)
	if err != nil {
		log.Println(err.Error())

		return nil
	}

	return &AuthData
}

func CollectAuthenticationDataByHeader(AuthHeader string) *AuthenticationData {
	slices := strings.Split(AuthHeader, " ")

	if len(slices) != 2 {
		log.Println("Expected Authorization Bearer header, got " + AuthHeader)

		return nil
	}

	return CollectAuthenticationData(slices[1])
}
