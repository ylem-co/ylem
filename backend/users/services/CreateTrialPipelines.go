package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"ylem_users/config"
)

func CreateTrialPipelines(organizationUuid string, destinationUuid string, sourceUuid string, token string) error {
	url := config.Cfg().NetworkConfig.YlemPipelinesBaseUrl + "pipeline/trials"

	rp, _ := json.Marshal(map[string]string{"organization_uuid": organizationUuid, "destination_uuid": destinationUuid, "source_uuid": sourceUuid})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rp))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())

		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return errors.New("Failed to create trial pipelines. Error: " + resp.Status)
}
