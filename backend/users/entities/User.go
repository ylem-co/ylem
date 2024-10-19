package entities

import (
	"database/sql"
	"strings"
)

type User struct {
	Id                       int
	FirstName                string
	LastName                 string
	Source                   string
	Uuid                     string
	ExternalSystemId         string
	Email                    string
	Phone                    string
	HashedPassword           string
	Roles                    string
	OrganizationId           int
	IsEmailConfirmed         int
}

type UserToExpose struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Uuid             string `json:"uuid"`
	Email            string `json:"email"`
	Roles            string `json:"roles"`
	IsActive         int    `json:"is_active"`
	IsEmailConfirmed int    `json:"is_email_confirmed"`
}

const IsEmailConfirmedTrue = 1
const IsEmailConfirmedFalse = 0

const ROLE_ORGANIZATION_ADMIN = "ROLE_ORGANIZATION_ADMIN"
const ROLE_TEAM_MEMBER = "ROLE_TEAM_MEMBER"
const ROLE_ALLOWED_TO_SWITCH = "ROLE_ALLOWED_TO_SWITCH"

const SourceApp = "app"
const SourceGoogle = "google"

func IsUserTeamMember(db *sql.DB, user User) bool {
	return strings.Contains(user.Roles, ROLE_TEAM_MEMBER)
}

func IsUserAllowedToSwitch(db *sql.DB, user User) bool {
	return strings.Contains(user.Roles, ROLE_ALLOWED_TO_SWITCH)
}

func IsUserOrganizationAdmin(db *sql.DB, user User) bool {
	return strings.Contains(user.Roles, ROLE_ORGANIZATION_ADMIN)
}

func DoesUserBelongToOrganization(db *sql.DB, user User, organization Organization) bool {
	return organization.Id == user.OrganizationId
}
