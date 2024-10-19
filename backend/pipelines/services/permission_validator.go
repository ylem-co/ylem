package services

import (
	"bytes"
	"errors"
	"encoding/json"
	"net/http"
	"ylem_pipelines/app/pipeline/common"
	"ylem_pipelines/config"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

const PermissionActionCreate = "create"
const PermissionActionRead = "read"
const PermissionActionReadList = "read_list"
const PermissionActionUpdate = "update"
const PermissionActionDelete = "delete"
const PermissionActionRun = "run"
const PermissionResourceTypePipeline = "pipeline"
const PermissionResourceTypeMetrics = "metrics"
const PermissionResourceTypeFolder = "folder"
const PermissionResourceTypeEnvVariable = "envvariable"

func ValidateBilledPermissions(userUuid string, organizationUuid string, action string, resourceType string, resourceUuid string, currentValue int64) (bool, error) {
	config := config.Cfg()
	err := envconfig.Process("", &config)
	if err != nil {
		return false, err
	}

	url := config.NetworkConfig.PermissionCheckUrl

	rp, _ := json.Marshal(map[string]interface{}{
		"user_uuid":         userUuid,
		"organization_uuid": organizationUuid,
		"action":            action,
		"resource_type":     resourceType,
		"resource_uuid":     resourceUuid,
		"current_value":     currentValue,
	})
	var jsonStr = []byte(rp)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else {
		return false, nil
	}
}

// Example of the usage
// services.ValidatePermissions("15eeebce-849a-4003-9200-423057790e61", "10c45d4a-e909-4946-86b2-c2165241cfde", "create", "source", "")
func ValidatePermissions(userUuid string, organizationUuid string, action string, resourceType string, resourceUuid string) bool {
	allowed, err := ValidateBilledPermissions(userUuid, organizationUuid, action, resourceType, resourceUuid, 0)
	if err != nil {
		log.Error(err)
	}
	return allowed
}

func GetPipelinePermissionResourceType(wfType string) (string, error) {
	switch wfType {
	case common.PipelineTypeGeneric:
		return PermissionResourceTypePipeline, nil
	case common.PipelineTypeMetric:
		return PermissionResourceTypeMetrics, nil
	}

	return "", errors.New("unknown pipeline type")

}
