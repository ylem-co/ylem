package api

import (
	"encoding/json"
	"net/http"
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/services"
	"strings"

	log "github.com/sirupsen/logrus"
)

func PermissionCheck(w http.ResponseWriter, r *http.Request) {
	var check services.HttpPermissionCheck

	w.Header().Set("Content-Type", "application/json")

	err := helpers.DecodeJSONBody(w, r, &check)
	if err != nil {
		rp, _ := json.Marshal(err.Msg)
		w.WriteHeader(err.Status)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	errorFields := ValidateHttpPermissionCheck(check, w)
	if len(errorFields) > 0 {
		rp, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": strings.Join(errorFields, ",")})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	ok := DoesUserHavePermissions(check)

	if ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func DoesUserHavePermissions(check services.HttpPermissionCheck) bool {
	switch check.ResourceType {
	case
		entities.RESOURCE_USER:
		return services.IsUserActionAllowed(check)
	case
		entities.RESOURCE_ORGANIZATION:
		return services.IsOrganizationActionAllowed(check)
	case
		entities.RESOURCE_PIPELINE:
		return services.IsPipelineActionAllowed(check)
	case
		entities.RESOURCE_METRICS:
		return services.IsMetricActionAllowed(check)
	case
		entities.RESOURCE_FOLDER:
		return services.IsFolderActionAllowed(check)
	case
		entities.RESOURCE_TASK:
		return services.IsTaskActionAllowed(check)
	case
		entities.RESOURCE_INTEGRATION:
		return services.IsIntegrationActionAllowed(check)
	case
		entities.RESOURCE_STAT:
		return services.IsStatActionAllowed(check)
	case
		entities.RESOURCE_ENVVARIABLE:
		return services.IsEnvVariableActionAllowed(check)
	case
		entities.RESOURCE_OAUTH_CLIENT:
		return services.IsOauthClientActionAllowed(check)

	case
		entities.RESOURCE_SUBSCRIPTION_PLAN:
		return services.IsSubscriptionPlanActionAllowed(check)

	case
		entities.RESOURCE_SUBSCRIPTION:
		return services.IsSubscriptionActionAllowed(check)

	case
		entities.RESOURCE_CHECKOUT_SESSION:
		return services.IsCheckoutSessionActionAllowed(check)

	case
		entities.RESOURCE_API_CALL:
		return services.IsApiActionAllowed(check)
	}

	return false
}

func ValidateHttpPermissionCheck(check services.HttpPermissionCheck, w http.ResponseWriter) []string {
	var errorFields []string

	if check.UserUuid == "" {
		errorFields = append(errorFields, "user_uuid")
	}

	if check.OrganizationUuid == "" {
		errorFields = append(errorFields, "organization_uuid")
	}

	if check.ResourceType == "" {
		errorFields = append(errorFields, "resource_type")
	}

	// Only "create", "read", "update" and "delete" actions are allowed
	if !entities.IsActionValid(check.Action) {
		errorFields = append(errorFields, "action")
	}

	return errorFields
}
