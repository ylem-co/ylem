package services

// ToDo replace this code with the proper permission system and matrix in the future
// Update functions accordingly when the new cases will appear

import (
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/repositories"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type HttpPermissionCheck struct {
	UserUuid         string `json:"user_uuid"`
	OrganizationUuid string `json:"organization_uuid"`
	ResourceUuid     string `json:"resource_uuid"`
	ResourceType     string `json:"resource_type"`
	Action           string `json:"action"`
	CurrentValue     int64  `json:"current_value"`
}

func IsUserActionAllowed(check HttpPermissionCheck) bool {
	if check.UserUuid == check.ResourceUuid {
		return true
	} else {
		db := helpers.DbConn()
		defer db.Close()

		user, ok := repositories.GetUserByUuid(db, check.UserUuid)
		if !ok {
			return false
		}

		if check.Action == entities.ACTION_CREATE {
			org, ok2 := repositories.GetOrganizationByUuid(db, check.OrganizationUuid)
			if !ok2 {
				return false
			}

			if entities.IsUserOrganizationAdmin(db, user) && entities.DoesUserBelongToOrganization(db, user, org) {
				return true
			}
		}

		if check.Action == entities.ACTION_READ || check.Action == entities.ACTION_DELETE {
			org, ok2 := repositories.GetOrganizationByUuid(db, check.OrganizationUuid)
			if !ok2 {
				return false
			}

			resource, ok3 := repositories.GetUserByUuid(db, check.ResourceUuid)
			if !ok3 {
				return false
			}

			if entities.IsUserOrganizationAdmin(db, user) && entities.DoesUserBelongToOrganization(db, user, org) && entities.DoesUserBelongToOrganization(db, resource, org) {
				return true
			}
		}

		if check.Action == entities.ACTION_READ_LIST {
			org, ok2 := repositories.GetOrganizationByUuid(db, check.OrganizationUuid)
			if !ok2 {
				return false
			}

			if entities.IsUserOrganizationAdmin(db, user) && entities.DoesUserBelongToOrganization(db, user, org) {
				return true
			}
		}
	}

	return false
}

func IsInvitationActionAllowed(check HttpPermissionCheck) bool {
	if check.Action == entities.ACTION_CREATE || check.Action == entities.ACTION_READ_LIST {
		db := helpers.DbConn()
		defer db.Close()

		user, ok := repositories.GetUserByUuid(db, check.UserUuid)
		if !ok {
			return false
		}

		org, ok2 := repositories.GetOrganizationByUuid(db, check.OrganizationUuid)
		if !ok2 {
			return false
		}

		if entities.IsUserOrganizationAdmin(db, user) && entities.DoesUserBelongToOrganization(db, user, org) {
			return true
		}
	}

	return false
}

func IsOrganizationActionAllowed(check HttpPermissionCheck) bool {
	// Organization is created only once during the registration and cannot be deleted
	if check.Action == entities.ACTION_CREATE || check.Action == entities.ACTION_DELETE {
		return false
	}

	db := helpers.DbConn()
	defer db.Close()

	user, ok := repositories.GetUserByUuid(db, check.UserUuid)
	if !ok {
		return false
	}

	org, ok2 := repositories.GetOrganizationByUuid(db, check.OrganizationUuid)
	if !ok2 {
		return false
	}

	if check.Action == entities.ACTION_UPDATE {
		if entities.IsUserOrganizationAdmin(db, user) && entities.DoesUserBelongToOrganization(db, user, org) {
			return true
		}
	} else {
		if entities.DoesUserBelongToOrganization(db, user, org) {
			return true
		}
	}

	return false
}

func IsPipelineActionAllowed(check HttpPermissionCheck) bool {
	db := helpers.DbConn()
	defer db.Close()

	user, ok := repositories.GetUserByUuid(db, check.UserUuid)
	if !ok {
		return false
	}

	org, ok2 := repositories.GetOrganizationByUuid(db, check.OrganizationUuid)
	if !ok2 {
		return false
	}

	if entities.DoesUserBelongToOrganization(db, user, org) {
		return true
	}

	return false
}

func IsMetricActionAllowed(check HttpPermissionCheck) bool {
	db := helpers.DbConn()
	defer db.Close()

	user, ok := repositories.GetUserByUuid(db, check.UserUuid)
	if !ok {
		return false
	}

	org, ok2 := repositories.GetOrganizationByUuid(db, check.OrganizationUuid)
	if !ok2 {
		return false
	}

	if entities.DoesUserBelongToOrganization(db, user, org) {
		return true
	}

	return false
}

func IsTaskActionAllowed(check HttpPermissionCheck) bool {
	return IsPipelineActionAllowed(check)
}

func IsFolderActionAllowed(check HttpPermissionCheck) bool {
	return IsPipelineActionAllowed(check)
}

func IsIntegrationActionAllowed(check HttpPermissionCheck) bool {
	db := helpers.DbConn()
	defer db.Close()

	user, ok := repositories.GetUserByUuid(db, check.UserUuid)
	if !ok {
		return false
	}

	org, ok2 := repositories.GetOrganizationByUuid(db, check.OrganizationUuid)
	if !ok2 {
		return false
	}

	if entities.DoesUserBelongToOrganization(db, user, org) {
		return true
	}

	return false
}

func IsStatActionAllowed(check HttpPermissionCheck) bool {
	return isGenericActionAllowed(check)
}

func IsOauthClientActionAllowed(check HttpPermissionCheck) bool {
	return isGenericActionAllowed(check)
}

func IsEnvVariableActionAllowed(check HttpPermissionCheck) bool {
	return IsPipelineActionAllowed(check)
}

func IsSubscriptionPlanActionAllowed(check HttpPermissionCheck) bool {
	return allowedForAdmin(check)
}

func IsSubscriptionActionAllowed(check HttpPermissionCheck) bool {
	return allowedForAdmin(check)
}

func IsCheckoutSessionActionAllowed(check HttpPermissionCheck) bool {
	return allowedForAdmin(check)
}

func IsApiActionAllowed(check HttpPermissionCheck) bool {
	if !isGenericActionAllowed(check) {
		return false
	}

	_, err := uuid.Parse(check.OrganizationUuid)
	if err != nil {
		log.Error(err)
		return false
	}

	return true
}

func isGenericActionAllowed(check HttpPermissionCheck) bool {
	db := helpers.DbConn()
	defer db.Close()

	user, ok := repositories.GetUserByUuid(db, check.UserUuid)
	if !ok {
		return false
	}

	org, ok2 := repositories.GetOrganizationByUuid(db, check.OrganizationUuid)
	if !ok2 {
		return false
	}

	if entities.DoesUserBelongToOrganization(db, user, org) {
		return true
	}

	return false
}

func allowedForAdmin(check HttpPermissionCheck) bool {
	db := helpers.DbConn()
	defer db.Close()

	user, ok := repositories.GetUserByUuid(db, check.UserUuid)
	if !ok {
		return false
	}

	org, ok2 := repositories.GetOrganizationByUuid(db, check.OrganizationUuid)
	if !ok2 {
		return false
	}

	return entities.DoesUserBelongToOrganization(db, user, org) && entities.IsUserOrganizationAdmin(db, user)
}
