package tests

import (
	"ylem_users/services"
	"ylem_users/entities"
	"testing"
)

func TestIsUserActionAllowed(t *testing.T) {
	type args struct {
		check services.HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Valid action on yourself", args{check: services.HttpPermissionCheck{UserUuid: "123456", ResourceUuid: "123456"}}, true},
		{"Invalid action. User is not found", args{check: services.HttpPermissionCheck{UserUuid: "123456", ResourceUuid: "6789"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsUserActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsUserActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInvitationActionAllowed(t *testing.T) {
	type args struct {
		check services.HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Action is not allowed", args{check: services.HttpPermissionCheck{Action: entities.ACTION_DELETE}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsInvitationActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsInvitationActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsOrganizationActionAllowed(t *testing.T) {
	type args struct {
		check services.HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Delete action is not allowed", args{check: services.HttpPermissionCheck{Action: entities.ACTION_DELETE}}, false},
		{"Create action is not allowed", args{check: services.HttpPermissionCheck{Action: entities.ACTION_CREATE}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsOrganizationActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsOrganizationActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*func TestIsPipelineActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsPipelineActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsPipelineActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMetricActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsMetricActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsMetricActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTaskActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsTaskActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsTaskActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFolderActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsFolderActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsFolderActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIntegrationActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsIntegrationActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsIntegrationActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsStatActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsStatActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsStatActionAllowed() = %v, want %v", got, tt.want)
			}
			})
	}
}

func TestIsOauthClientActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsOauthClientActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsOauthClientActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEnvVariableActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsEnvVariableActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsEnvVariableActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSubscriptionPlanActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsSubscriptionPlanActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsSubscriptionPlanActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSubscriptionActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsSubscriptionActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsSubscriptionActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCheckoutSessionActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsCheckoutSessionActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsCheckoutSessionActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsApiActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsApiActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("IsApiActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isGenericActionAllowed(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.isGenericActionAllowed(tt.args.check); got != tt.want {
				t.Errorf("isGenericActionAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_allowedForAdmin(t *testing.T) {
	type args struct {
		check HttpPermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.allowedForAdmin(tt.args.check); got != tt.want {
				t.Errorf("allowedForAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}*/
