package ylem_integrations

import (
	"context"
	"fmt"
	"sync"
	"time"
	"net/http"
	"ylem_pipelines/config"

	"golang.org/x/time/rate"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	lru "github.com/hashicorp/golang-lru"
	log "github.com/sirupsen/logrus"
	messaging "github.com/ylem-co/shared-messaging"
)

const IntegrationStatusOnline = "online"

type LoadIntegrationFunc func(uuid.UUID) (interface{}, error)

type Client interface {
	GetApiIntegration(uid uuid.UUID) (*Api, error)
	GetEmailIntegration(uid uuid.UUID) (*Email, error)
	GetSlackIntegration(uid uuid.UUID) (*Slack, error)
	GetSmsIntegration(uid uuid.UUID) (*Sms, error)
	GetJiraIntegration(uid uuid.UUID) (*Jira, error)
	GetIncidentIoIntegration(uid uuid.UUID) (*IncidentIo, error)
	GetOpsgenieIntegration(uid uuid.UUID) (*Opsgenie, error)
	GetTableauIntegration(uid uuid.UUID) (*Tableau, error)
	GetHubspotIntegration(uid uuid.UUID) (*Hubspot, error)
	GetGoogleSheetsIntegration(uid uuid.UUID) (*GoogleSheets, error)
	GetSalesforceIntegration(uid uuid.UUID) (*Salesforce, error)
	GetJenkinsIntegration(uid uuid.UUID) (*Jenkins, error)
	GetSQLIntegration(uid uuid.UUID) (*messaging.SQLIntegration, error)
}

type cachingClient struct {
	ctx              context.Context
	innerClient      Client
	integrationCache *lru.Cache
	mu               *sync.RWMutex
}

func (c *cachingClient) init() error {
	c.mu = &sync.RWMutex{}
	err := c.initIntegrationCache()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			// reset Integration cache every 30 seconds
			case <-time.After(time.Second * 30):
				err := c.initIntegrationCache()
				if err != nil {
					panic(err)
				}

			case <-c.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *cachingClient) initIntegrationCache() error {
	var err error
	c.integrationCache, err = lru.New(10000)
	log.Trace("Integration cache reset")

	return err
}

func (c *cachingClient) GetApiIntegration(uid uuid.UUID) (*Api, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetApiIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	api, ok := r.(*Api)
	if !ok {
		return nil, nil
	}

	return api, nil
}

func (c *cachingClient) GetEmailIntegration(uid uuid.UUID) (*Email, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetEmailIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	email, ok := r.(*Email)
	if !ok {
		return nil, nil
	}

	return email, nil

}

func (c *cachingClient) GetSlackIntegration(uid uuid.UUID) (*Slack, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetSlackIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	slack, ok := r.(*Slack)
	if !ok {
		return nil, nil
	}

	return slack, nil

}

func (c *cachingClient) GetSmsIntegration(uid uuid.UUID) (*Sms, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetSmsIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	sms, ok := r.(*Sms)
	if !ok {
		return nil, nil
	}

	return sms, nil

}

func (c *cachingClient) GetJiraIntegration(uid uuid.UUID) (*Jira, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetJiraIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	jira, ok := r.(*Jira)
	if !ok {
		return nil, nil
	}

	return jira, nil

}

func (c *cachingClient) GetIncidentIoIntegration(uid uuid.UUID) (*IncidentIo, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetIncidentIoIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	io, ok := r.(*IncidentIo)
	if !ok {
		return nil, nil
	}

	return io, nil
}

func (c *cachingClient) GetOpsgenieIntegration(uid uuid.UUID) (*Opsgenie, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetOpsgenieIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	o, ok := r.(*Opsgenie)
	if !ok {
		return nil, nil
	}

	return o, nil
}

func (c *cachingClient) GetTableauIntegration(uid uuid.UUID) (*Tableau, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetTableauIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	dest, ok := r.(*Tableau)
	if !ok {
		return nil, nil
	}

	return dest, nil
}

func (c *cachingClient) GetHubspotIntegration(uid uuid.UUID) (*Hubspot, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetHubspotIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	dest, ok := r.(*Hubspot)
	if !ok {
		return nil, nil
	}

	return dest, nil
}

func (c *cachingClient) GetGoogleSheetsIntegration(uid uuid.UUID) (*GoogleSheets, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetGoogleSheetsIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	dest, ok := r.(*GoogleSheets)
	if !ok {
		return nil, nil
	}

	return dest, nil
}

func (c *cachingClient) GetSalesforceIntegration(uid uuid.UUID) (*Salesforce, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetSalesforceIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	dest, ok := r.(*Salesforce)
	if !ok {
		return nil, nil
	}

	return dest, nil
}

func (c *cachingClient) GetJenkinsIntegration(uid uuid.UUID) (*Jenkins, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetJenkinsIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	dest, ok := r.(*Jenkins)
	if !ok {
		return nil, nil
	}

	return dest, nil
}

func (c *cachingClient) GetSQLIntegration(uid uuid.UUID) (*messaging.SQLIntegration, error) {
	r, err := c.getCached(uid, func(uuid.UUID) (interface{}, error) {
		r, err := c.innerClient.GetSQLIntegration(uid)
		return r, err
	})

	if err != nil {
		return nil, err
	}

	dest, ok := r.(*messaging.SQLIntegration)
	if !ok {
		return nil, nil
	}

	return dest, nil
}

func (c *cachingClient) getCached(uid uuid.UUID, ldf LoadIntegrationFunc) (interface{}, error) {
	uidStr := uid.String()
	if d, ok := c.integrationCache.Get(uidStr); ok {
		return d, nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if d, ok := c.integrationCache.Get(uidStr); ok {
		return d, nil
	}

	d, err := ldf(uid)
	if err != nil {
		return d, err
	}

	if d == nil {
		return nil, nil
	}

	c.integrationCache.Add(uidStr, d)

	return d, nil
}

type DefaultClient struct {
	client  *resty.Client
	limiter *rate.Limiter
	ctx     context.Context
}

func (c *DefaultClient) GetApiIntegration(uid uuid.UUID) (*Api, error) {
	r := &Api{}

	err := c.getIntegration(
		fmt.Sprintf("/private/api/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) GetEmailIntegration(uid uuid.UUID) (*Email, error) {
	r := &Email{}

	err := c.getIntegration(
		fmt.Sprintf("/private/email/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) GetSlackIntegration(uid uuid.UUID) (*Slack, error) {
	r := &Slack{}

	err := c.getIntegration(
		fmt.Sprintf("/private/slack/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) GetSmsIntegration(uid uuid.UUID) (*Sms, error) {
	r := &Sms{}

	err := c.getIntegration(
		fmt.Sprintf("/private/sms/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) GetJiraIntegration(uid uuid.UUID) (*Jira, error) {
	r := &Jira{}

	err := c.getIntegration(
		fmt.Sprintf("/private/jira/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) GetIncidentIoIntegration(uid uuid.UUID) (*IncidentIo, error) {
	r := &IncidentIo{}

	err := c.getIntegration(
		fmt.Sprintf("/private/incidentio/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) GetOpsgenieIntegration(uid uuid.UUID) (*Opsgenie, error) {
	r := &Opsgenie{}

	err := c.getIntegration(
		fmt.Sprintf("/private/opsgenie/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) GetTableauIntegration(uid uuid.UUID) (*Tableau, error) {
	r := &Tableau{}

	err := c.getIntegration(
		fmt.Sprintf("/private/tableau/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) GetHubspotIntegration(uid uuid.UUID) (*Hubspot, error) {
	r := &Hubspot{}

	err := c.getIntegration(
		fmt.Sprintf("/private/hubspot/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) GetGoogleSheetsIntegration(uid uuid.UUID) (*GoogleSheets, error) {
	r := &GoogleSheets{}

	err := c.getIntegration(
		fmt.Sprintf("/private/google-sheets/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) GetSalesforceIntegration(uid uuid.UUID) (*Salesforce, error) {
	r := &Salesforce{}

	err := c.getIntegration(
		fmt.Sprintf("/private/salesforce/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) GetJenkinsIntegration(uid uuid.UUID) (*Jenkins, error) {
	r := &Jenkins{}

	err := c.getIntegration(
		fmt.Sprintf("/private/jenkins/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) GetSQLIntegration(uid uuid.UUID) (*messaging.SQLIntegration, error) {
	r := &messaging.SQLIntegration{}

	err := c.getIntegration(
		fmt.Sprintf("/private/sql/%s", uid),
		r,
	)

	return r, err
}

func (c *DefaultClient) getIntegration(url string, result interface{}) error {
	_ = c.limiter.Wait(c.ctx)

	log.Tracef("Calling Ylem_integrations service endpoint %s", url)

	resp, err := c.client.
		R().
		SetResult(result).
		Get(url)

	if err != nil {
		log.Error("Ylem_integrations service call error: ", err)
		return ErrorServiceUnavailable{}
	}

	if resp.StatusCode() == http.StatusNotFound {
		log.Tracef("Ylem_integrations: integration not found")
		return nil
	}

	if resp.StatusCode() != http.StatusOK {
		log.Error("Ylem_integrations: service is unavailable")
		return ErrorServiceUnavailable{}
	}

	return nil
}

type ErrorServiceUnavailable struct {
}

type Integration struct {
	Id               int64  `json:"-"`
	Uuid             string `json:"uuid"`
	CreatorUuid      string `json:"creator_uuid"`
	OrganizationUuid string `json:"organization_uuid"`
	Status           string `json:"status"`
	Type             string `json:"type"`
	Name             string `json:"name"`
	Value            string `json:"value"`
	UserUpdatedAt    string `json:"user_updated_at"`
}

type Api struct {
	Integration           Integration `json:"integration"`
	Method                string      `json:"method"`
	AuthType              string      `json:"auth_type"`
	AuthBearerToken       string      `json:"auth_bearer_token"`
	AuthBasicUserName     string      `json:"auth_basic_user_name"`
	AuthBasicUserPassword string      `json:"auth_basic_user_password"`
	AuthHeaderName        string      `json:"auth_header_name"`
	AuthHeaderValue       string      `json:"auth_header_value"`
}

type Email struct {
	Id          int64       `json:"-"`
	Integration Integration `json:"integration"`
	Code        string      `json:"-"`
	IsConfirmed bool        `json:"is_confirmed"`
	RequestedAt time.Time   `json:"requested_at"`
}

type Slack struct {
	Id                 int64              `json:"-"`
	Integration        Integration        `json:"integration"`
	SlackAuthorization SlackAuthorization `json:"authorization"`
	SlackChannelId     string             `json:"slack_channel_id"`
}

type SlackAuthorization struct {
	Id               int64  `json:"-"`
	Name             string `json:"name"`
	Uuid             string `json:"uuid"`
	CreatorUuid      string `json:"-"`
	OrganizationUuid string `json:"-"`
	State            string `json:"-"`
	AccessToken      string `json:"access_token"`
	Scopes           string `json:"-"`
	BotUserId        string `json:"-"`
	IsActive         bool   `json:"is_active"`
}

type Sms struct {
	Id          int64       `json:"-"`
	Integration Integration `json:"integration"`
	Code        string      `json:"-"`
	IsConfirmed bool        `json:"is_confirmed"`
	RequestedAt time.Time   `json:"requested_at"`
}

type Jira struct {
	Integration Integration `json:"integration"`
	IssueType   string      `json:"issue_type"`
	DataKey     []byte      `json:"data_key"`
	AccessToken []byte      `json:"access_token"`
	CloudId     string      `json:"cloudid"`
}

type IncidentIo struct {
	Integration Integration `json:"integration"`
	DataKey     []byte      `json:"data_key"`
	ApiKey      []byte      `json:"api_key"`
	Mode        string      `json:"mode"`
	Visibility  string      `json:"visibility"`
}

type Opsgenie struct {
	Integration Integration `json:"integration"`
	DataKey     []byte      `json:"data_key"`
	ApiKey      []byte      `json:"api_key"`
}

type Tableau struct {
	Integration    Integration `json:"integration"`
	Server         string      `json:"server"`
	DataKey        []byte      `json:"data_key"`
	Username       []byte      `json:"username"`
	Password       []byte      `json:"password"`
	Sitename       string      `json:"site_name"`
	ProjectName    string      `json:"project_name"`
	DatasourceName string      `json:"datasource_name"`
	Mode           string      `json:"mode"`
}

type Hubspot struct {
	Integration       Integration `json:"integration"`
	DataKey           []byte      `json:"data_key"`
	AccessToken       []byte      `json:"access_token"`
	PipelineStageCode string      `json:"pipeline_stage_code"`
	OwnerCode         string      `json:"owner_code"`
}

type GoogleSheets struct {
	Integration   Integration `json:"integration"`
	DataKey       []byte      `json:"data_key"`
	Credentials   []byte      `json:"credentials"`
	Mode          string      `json:"mode"`
	SpreadsheetId string      `json:"spreadsheet_id"`
	SheetId       int64       `json:"sheet_id"`
	WriteHeader   bool        `json:"write_header"`
}

type Salesforce struct {
	Integration Integration `json:"integration"`
	DataKey     []byte      `json:"data_key"`
	AccessToken []byte      `json:"access_token"`
	Domain      string      `json:"domain"`
}

type Jenkins struct {
	Integration Integration `json:"integration"`
	DataKey     []byte      `json:"data_key"`
	Token       []byte      `json:"token"`
	BaseUrl     string      `json:"base_url"`
}

func (e ErrorServiceUnavailable) Error() string {
	return "Ylem_integrations service is unavailable"
}

func NewClient(ctx context.Context) (Client, error) {
	c := &DefaultClient{
		client:  resty.New(),
		limiter: rate.NewLimiter(20, 1),
		ctx:     ctx,
	}

	c.client.SetBaseURL(config.Cfg().YlemIntegrations.BaseURL)

	cc := &cachingClient{
		ctx:         ctx,
		innerClient: c,
	}
	err := cc.init()

	if err != nil {
		return nil, err
	}

	log.Trace("Ylem_integrations client initialized")

	return cc, nil
}
