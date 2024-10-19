package services

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
	"ylem_integrations/config"
)

type Jira struct {
	Client *resty.Client
}

var jira Jira

type JiraAuthorizedGrant struct {
	Scopes      string `json:"scope,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int
}

type JiraAvailableResource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func CreateJiraCloudGrantLink(State string) string {
	return fmt.Sprintf(
		"https://auth.atlassian.com/authorize?audience=api.atlassian.com&client_id=%s&scope=read%%3Ajira-user%%20write%%3Ajira-work&&state=%s&redirect_uri=%s&response_type=code&prompt=consent",
		config.Cfg().Jira.OauthClientId,
		State,
		url.QueryEscape(config.Cfg().Jira.OauthRedirectUri),
	)
}

func JiraGrantAuthorization(Token string) (*JiraAuthorizedGrant, error) {
	cfg := config.Cfg().Jira

	payload := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     cfg.OauthClientId,
		"client_secret": cfg.OauthClientSecret,
		"code":          Token,
		"redirect_uri":  cfg.OauthRedirectUri,
	}

	var grant JiraAuthorizedGrant
	resp, err := jira.Client.
		R().
		SetHeader("Accept", "application/json").
		SetResult(&grant).
		SetBody(payload).
		Post("https://auth.atlassian.com/oauth/token")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("jira api call: %s", resp.Status())
	}

	return &grant, nil
}

func JiraListAvailableResources(AccessToken string) ([]JiraAvailableResource, error) {
	resp, err := jira.Client.
		R().
		SetAuthToken(AccessToken).
		SetHeader("Accept", "application/json").
		SetResult([]JiraAvailableResource{}).
		Get("https://api.atlassian.com/oauth/token/accessible-resources")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("jira api call: %s", resp.Status())
	}

	return *resp.Result().(*[]JiraAvailableResource), nil
}

func init() {
	jira = Jira{
		Client: resty.New(),
	}
}
