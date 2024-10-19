package services

import (
	"bytes"
	"strings"
	"io"
	"encoding/json"
	"net/http"
	"ylem_pipelines/config"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type AuthenticationData struct {
	Email            string `json:"email"`
	Uuid             string `json:"uuid"`
	OrganizationUuid string `json:"organization_uuid"`
}

func InitialAuthorization(AuthHeader string) *AuthenticationData {
	slices := strings.Split(AuthHeader, " ")

	if len(slices) != 2 {
		log.Println("Expected Authorization Bearer header, got " + AuthHeader)
		return nil
	}

	token := slices[1]

	var config config.Config
	err := envconfig.Process("", &config)

	if err != nil {
		log.Error(err)
		return nil
	}

	url := config.NetworkConfig.AuthorizationCheckUrl

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Error(err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		log.Error(err)
		return nil
	}

	var AuthData AuthenticationData
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil
	}

	err = json.Unmarshal(bodyBytes, &AuthData)
	if err != nil {
		log.Error(err)
		return nil
	}

	return &AuthData
}
