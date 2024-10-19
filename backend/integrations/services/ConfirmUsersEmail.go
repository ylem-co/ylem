package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"ylem_integrations/config"
)

func ConfirmUsersEmail(userUuid string) error {
	var config config.Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Println(err.Error())

		return err
	}

	url := config.NetworkConfig.YlemUsersBaseUrl + "private/user/" + userUuid + "/confirm-email";

	rp, _ := json.Marshal(map[string]string{})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rp))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")

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

	return errors.New("Failed to confirm user's Email. Error: " + resp.Status)
}
