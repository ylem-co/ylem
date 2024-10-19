package tests

import (
	"database/sql"
	"ylem_users/entities"
	"ylem_users/helpers"
	"testing"
)

func TestIsUserTeamMember(t *testing.T) {
	type args struct {
		db   *sql.DB
		user entities.User
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Valid 1", args{db: helpers.DbConn(), user: entities.User{Roles: entities.ROLE_TEAM_MEMBER + ", " + entities.ROLE_ORGANIZATION_ADMIN}}, true},
		{"Invalid 1", args{db: helpers.DbConn(), user: entities.User{Roles: entities.ROLE_ORGANIZATION_ADMIN}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entities.IsUserTeamMember(tt.args.db, tt.args.user); got != tt.want {
				t.Errorf("IsUserTeamMember() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUserAllowedToSwitch(t *testing.T) {
	type args struct {
		db   *sql.DB
		user entities.User
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Valid 1", args{db: helpers.DbConn(), user: entities.User{Roles: entities.ROLE_TEAM_MEMBER + ", " + entities.ROLE_ALLOWED_TO_SWITCH}}, true},
		{"Invalid 1", args{db: helpers.DbConn(), user: entities.User{Roles: entities.ROLE_TEAM_MEMBER + ", " + entities.ROLE_ORGANIZATION_ADMIN}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entities.IsUserAllowedToSwitch(tt.args.db, tt.args.user); got != tt.want {
				t.Errorf("IsUserAllowedToSwitch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUserOrganizationAdmin(t *testing.T) {
	type args struct {
		db   *sql.DB
		user entities.User
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Valid 1", args{db: helpers.DbConn(), user: entities.User{Roles: entities.ROLE_TEAM_MEMBER + ", " + entities.ROLE_ORGANIZATION_ADMIN}}, true},
		{"Invalid 1", args{db: helpers.DbConn(), user: entities.User{Roles: entities.ROLE_TEAM_MEMBER + ", " + entities.ROLE_ALLOWED_TO_SWITCH}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entities.IsUserOrganizationAdmin(tt.args.db, tt.args.user); got != tt.want {
				t.Errorf("IsUserOrganizationAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoesUserBelongToOrganization(t *testing.T) {
	type args struct {
		db           *sql.DB
		user         entities.User
		organization entities.Organization
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Valid 1", args{db: helpers.DbConn(), user: entities.User{OrganizationId: 123}, organization: entities.Organization{Id: 123}}, true},
		{"Invalid 1", args{db: helpers.DbConn(), user: entities.User{OrganizationId: 123}, organization: entities.Organization{Id: 321}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entities.DoesUserBelongToOrganization(tt.args.db, tt.args.user, tt.args.organization); got != tt.want {
				t.Errorf("DoesUserBelongToOrganization() = %v, want %v", got, tt.want)
			}
		})
	}
}
