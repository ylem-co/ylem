package services

import (
	"bytes"
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"ylem_users/config"
)

func CreateDefaultEmailDestination(organizationUuid string, email string, firstName string, token string) (string, bool) {
	url := config.Cfg().NetworkConfig.YlemIntegrationsBaseUrl + "email"

	rp, _ := json.Marshal(map[string]string{"email": email, "organization_uuid": organizationUuid, "name": firstName + "'s Email"})
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

	fmt.Println(resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		buf := new(bytes.Buffer)
    	_, err = buf.ReadFrom(resp.Body)
    	if err != nil {
    		fmt.Println("Read error")
    		log.Println(err.Error())
    		return "", false
    	}
    	body := buf.String()

    	x := map[string]map[string]string{}
    	_ = json.Unmarshal([]byte(body), &x)

    	fmt.Println(x["integration"]["uuid"])

		return x["integration"]["uuid"], true
	}

	return "", false
}
