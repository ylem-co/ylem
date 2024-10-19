package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"ylem_users/config"
)

func TestTrialDBDataSource(sourceUuid string, token string) bool {
	url := config.Cfg().NetworkConfig.YlemIntegrationsBaseUrl + "integration/sql/" + sourceUuid + "/test";

	rp, _ := json.Marshal("{}")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rp))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
