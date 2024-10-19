package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"strings"
	"ylem_integrations/config"
	"ylem_integrations/repositories"
)

func NotifyServiceIntegrationsChanged(db *sql.DB, OrganizationUuid string) {
	integrations, _ := repositories.FindAllIntegrationsBelongToOrganization(db, OrganizationUuid, "all")

	if len(integrations.Items) == 0 {
		_ = updateIntegrationConnection(OrganizationUuid, false)
	}

	if len(integrations.Items) == 1 {
		_ = updateIntegrationConnection(OrganizationUuid, true)
	}
}

func updateIntegrationConnection(organizationUuid string, IsDestinationCreated bool) error {
	var config config.Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Println(err.Error())

		return err
	}

	url := strings.Replace(config.NetworkConfig.UpdateConnectionsUrl, "{uuid}", organizationUuid, -1)

	rp, _ := json.Marshal(map[string]bool{"is_destination_created": IsDestinationCreated})
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

	return errors.New("expected HTTP 200 from Ylem_users, got " + resp.Status)
}
