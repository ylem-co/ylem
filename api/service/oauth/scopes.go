package oauth

import "strings"

const (
	ScopePipelinesRun = "pipelines:run"
	ScopeStatsRead    = "stats:read"
)

func IsScopeGranted(scope, grantedScopes string) bool {
	for _, v := range strings.Split(grantedScopes, ",") {
		if scope == strings.TrimSpace(v) {
			return true
		}
	}

	return false
}

func NormalizeScopes(scopes string) []string {
	scopeArr := strings.Split(scopes, ",")
	for k, v := range scopeArr {
		scopeArr[k] = strings.TrimSpace(v)
	}

	return scopeArr
}
