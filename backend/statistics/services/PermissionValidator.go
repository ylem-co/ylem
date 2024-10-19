package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"ylem_statistics/config"

	log "github.com/sirupsen/logrus"
)

const PermissionActionCreate = "create"
const PermissionActionRead = "read"
const PermissionActionReadList = "read_list"
const PermissionActionUpdate = "update"
const PermissionActionDelete = "delete"
const PermissionResourceTypeStat = "stat"

// Example of the usage
// services.ValidatePermissions("15eeebce-849a-4003-9200-423057790e61", "10c45d4a-e909-4946-86b2-c2165241cfde", "create", "source", "")
func ValidatePermissions(userUuid string, organizationUuid string, action string, resourceType string, resourceUuid string) bool {
	config := config.Cfg()

	url := config.PermissionCheckUrl

	rp, _ := json.Marshal(map[string]string{"user_uuid": userUuid, "organization_uuid": organizationUuid, "action": action, "resource_type": resourceType, "resource_uuid": resourceUuid})
	var jsonStr = []byte(rp)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Error(err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	} else {
		return false
	}
}
