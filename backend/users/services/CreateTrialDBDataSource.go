package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"ylem_users/config"
)

type dataSource struct {
	OrganizationUuid string  `json:"organization_uuid"`
	Type             string  `json:"type"`
	Name             string  `json:"name"`
	Host             string  `json:"host"`
	Port             int     `json:"port"`
	User             string  `json:"user"`
	Password         string  `json:"password"`
	Database         string  `json:"database"`
	ConnectionType   string  `json:"connection_type"`
	SslEnabled       bool    `json:"ssl_enabled"`
	SshHost          string  `json:"ssh_host"`
	SshPort          int     `json:"ssh_port"`
	SshUser          string  `json:"ssh_user"`
}

func newTrialDataSource(organizationUuid string) *dataSource {
	d := dataSource{OrganizationUuid: organizationUuid}
	d.Type = "mysql"
	d.Name = "Public Rfam Database for testing"
	d.Host = "mysql-rfam-public.ebi.ac.uk"
	d.Port = 4497
	d.User = "rfamro"
	d.Password = ""
	d.Database = "Rfam"
	d.ConnectionType = "direct"

	return &d
}

func CreateTrialDBDataSource(organizationUuid string, token string) (string, bool) {
	url := config.Cfg().NetworkConfig.YlemIntegrationsBaseUrl + "sql/mysql";

	trialDBData := newTrialDataSource(organizationUuid);

	rp, _ := json.Marshal(trialDBData)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rp))
	if err != nil {
		log.Println(err.Error())
		return "", false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return "", false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		buf := new(bytes.Buffer)
    	_, err = buf.ReadFrom(resp.Body)
    	if err != nil {
    		log.Println(err.Error())
    		return "", false
    	}

    	body := buf.String()

    	x := map[string]map[string]string{}
    	_ = json.Unmarshal([]byte(body), &x)

		return x["integration"]["uuid"], true
	}

	return "", false
}
