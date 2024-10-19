package services

import (
	"fmt"
	"net/http"
	"strings"
	"ylem_integrations/config"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

const PermissionActionCreate = "create"
const PermissionActionRead = "read"
const PermissionActionReadList = "read_list"
const PermissionActionUpdate = "update"
const PermissionActionDelete = "delete"

const PermissionResourceTypeIntegration = "integration"

type YlemUsers struct {
	Client *resty.Client
}

var ylemUsers YlemUsers

type AuthenticationData struct {
	Email   string
	Uuid    string
	DataKey []byte
}

func CollectAuthenticationData(token string) *AuthenticationData {
	log.Trace("Checking if token is authenticated. token: ", token)
	url := config.Cfg().NetworkConfig.AuthorizationCheckUrl

	var authData AuthenticationData
	resp, err := ylemUsers.Client.
		R().
		SetHeader("Accept", "application/json").
		SetAuthToken(token).
		SetResult(&authData).
		Post(url)

	if err != nil {
		log.Error(err.Error())

		return nil
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Expected %d code, got %s", http.StatusOK, resp.Status())

		return nil
	}

	return &authData
}

func CollectAuthenticationDataByHeader(AuthHeader string) *AuthenticationData {
	slices := strings.Split(AuthHeader, " ")

	if len(slices) != 2 {
		log.Info("Expected Authorization Bearer header, got " + AuthHeader)

		return nil
	}

	return CollectAuthenticationData(slices[1])
}

func ValidateBilledPermissions(userUuid string, organizationUuid string, action string, resourceType string, resourceUuid string, currentValue int64) (bool, error) {
	log.Tracef("Checking permissions for user uuid %s", userUuid)
	url := config.Cfg().NetworkConfig.PermissionCheckUrl

	payload := map[string]interface{}{
		"user_uuid":         userUuid,
		"organization_uuid": organizationUuid,
		"action":            action,
		"resource_type":     resourceType,
		"resource_uuid":     resourceUuid,
		"current_value":     currentValue,
	}

	resp, err := ylemUsers.Client.
		R().
		SetHeader("Accept", "application/json").
		SetBody(payload).
		Post(url)

	if err != nil {
		log.Error(err.Error())

		return false, err
	}

	if resp.StatusCode() == http.StatusForbidden {
		return false, nil
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("expected %d status code, got %s; body %s", http.StatusOK, resp.Status(), string(resp.Body()))
		log.Error(err.Error())

		return false, err
	}

	return true, nil
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

func FetchOrganizationDataKey(uuid string) ([]byte, error) {
	log.Tracef("Fetching data key for organization %s", uuid)
	urlPattern := config.Cfg().NetworkConfig.RetrieveOrganizationDataKeyUrl
	url := strings.ReplaceAll(urlPattern, "{uuid}", uuid)

	resp, err := ylemUsers.Client.
		R().
		SetHeader("Accept", "application/octet-stream").
		Get(url)

	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("http Ylem_Users call: %s", resp.Status())
		log.Error(err.Error())

		return nil, err
	}

	return resp.Body(), nil
}

func init() {
	ylemUsers = YlemUsers{
		Client: resty.New(),
	}
}
