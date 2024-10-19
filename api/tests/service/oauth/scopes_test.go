package tests

import (
	"reflect"
	"testing"

	"ylem_api/service/oauth"
)

func TestIsScopeGranted(t *testing.T) {
	type args struct {
		scope         string
		grantedScopes string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Granted scope", args{scope: oauth.ScopePipelinesRun, grantedScopes: oauth.ScopePipelinesRun + ", " + oauth.ScopeStatsRead}, true},
		{"Not granted scope", args{scope: oauth.ScopeStatsRead, grantedScopes: oauth.ScopePipelinesRun}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := oauth.IsScopeGranted(tt.args.scope, tt.args.grantedScopes); got != tt.want {
				t.Errorf("IsScopeGranted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeScopes(t *testing.T) {
	type args struct {
		scopes string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"Normalized scopes", args{scopes: oauth.ScopePipelinesRun + ", " + oauth.ScopeStatsRead}, []string{oauth.ScopePipelinesRun, oauth.ScopeStatsRead}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := oauth.NormalizeScopes(tt.args.scopes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NormalizeScopes() = %v, want %v", got, tt.want)
			}
		})
	}
}
