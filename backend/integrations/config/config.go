package config

import (
	"os"
	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Listen   string `default:"0.0.0.0:7337"`
	LogLevel string `envconfig:"YLEM_LOG_LEVEL"`
	DBConfig struct {
		Name     string `envconfig:"INTEGRATIONS_DATABASE_NAME"`
		Host     string `envconfig:"YLEM_DATABASE_HOST"`
		Port     string `envconfig:"YLEM_DATABASE_PORT"`
		User     string `envconfig:"YLEM_DATABASE_USER"`
		Password string `envconfig:"YLEM_DATABASE_PASSWORD"`
	}
	NetworkConfig struct {
		AuthorizationCheckUrl                   string `envconfig:"YLEM_AUTHORIZATION_CHECK_URL"`
		PermissionCheckUrl                      string `envconfig:"YLEM_PERMISSION_CHECK_URL"`
		UpdateConnectionsUrl                    string `envconfig:"YLEM_UPDATE_CONNECTIONS_URL"`
		YlemUsersBaseUrl                        string `envconfig:"YLEM_USERS_BASE_URL"`
		RetrieveOrganizationDataKeyUrl          string `envconfig:"INTEGRATIONS_RETRIEVE_ORGANIZATION_DATA_KEY_URL"`
		EmailConfirmationUrl                    string `envconfig:"INTEGRATIONS_EMAIL_CONFIRMATION_LINK"`
		EmailAfterConfirmationRedirectUrl       string `envconfig:"INTEGRATIONS_EMAIL_AFTER_CONFIRMATION_REDIRECT_URL"`
		SlackAfterAuthorizationRedirectUrl      string `envconfig:"INTEGRATIONS_SLACK_AFTER_AUTHORIZATION_REDIRECT_URL"`
		JiraAfterAuthorizationRedirectUrl       string `envconfig:"INTEGRATIONS_JIRA_AFTER_AUTHORIZATION_REDIRECT_URL"`
		HubspotAfterAuthorizationRedirectUrl    string `envconfig:"INTEGRATIONS_HUBSPOT_AFTER_AUTHORIZATION_REDIRECT_URL"`
		SalesforceAfterAuthorizationRedirectUrl string `envconfig:"INTEGRATIONS_SALESFORCE_AFTER_AUTHORIZATION_REDIRECT_URL"`
	}
	Slack struct {
		ClientId     string `envconfig:"INTEGRATIONS_SLACK_CLIENT_ID"`
		ClientSecret string `envconfig:"INTEGRATIONS_SLACK_CLIENT_SECRET"`
	}
	Twilio struct {
		AccountSid string `envconfig:"YLEM_INTEGRATIONS_TWILIO_ACCOUNT_SID"`
		AuthToken  string `envconfig:"YLEM_INTEGRATIONS_TWILIO_AUTH_TOKEN"`
		NumberFrom string `envconfig:"YLEM_INTEGRATIONS_TWILIO_NUMBER_FROM"`
	}
	Kafka      kafka
	Jira       jira       `split_words:"true"`
	Aws        aws        `split_words:"true"`
	Hubspot    hubspot    `split_words:"true"`
	Salesforce salesforce `split_words:"true"`
	Ssh   ssh
}

type kafka struct {
	BootstrapServers                   []string `envconfig:"YLEM_KAFKA_BOOTSTRAP_SERVERS"`
	TaskRunsTopic                      string   `envconfig:"YLEM_KAFKA_TASK_RUNS_TOPIC"`
	TaskRunResultsTopic                string   `envconfig:"YLEM_KAFKA_TASK_RUN_RESULTS_TOPIC"`
	TaskRunsLoadBalancedTopic          string   `envconfig:"YLEM_KAFKA_TASK_RUNS_LOAD_BALANCED_TOPIC"`
	QueryTaskRunResultsTopic           string   `envconfig:"YLEM_KAFKA_QUERY_TASK_RUN_RESULTS_TOPIC"`
	NotificationTaskRunResultsTopic    string   `envconfig:"YLEM_KAFKA_NOTIFICATION_TASK_RUN_RESULTS_TOPIC"`
	YlemIntegrationsConsumerGroupName  string   `split_words:"true" default:"ylem_integrations_consumer"`
	LoadBalancerConsumerGroupName      string   `split_words:"true" default:"task_load_balancer"`
}

type jira struct {
	OauthClientId     string `split_words:"true"`
	OauthClientSecret string `split_words:"true"`
	OauthRedirectUri  string `split_words:"true"`
}

type hubspot struct {
	OauthClientId     string `split_words:"true"`
	OauthClientSecret string `split_words:"true"`
	OauthRedirectUri  string `split_words:"true"`
}

type salesforce struct {
	OauthClientId     string `split_words:"true"`
	OauthClientSecret string `split_words:"true"`
	OauthRedirectUri  string `split_words:"true"`
}

type aws struct {
	KmsKeyId    string `split_words:"true" envconfig:"AWS_KMS_KEY_ID"`
}

type ssh struct {
	PrivateKey []byte
}

func Cfg() Config {
	return c
}

var c Config

func init() {
	c = new()

	readSshPrivateKey()

	formattedConfig, _ := json.MarshalIndent(c, "", "    ")
	log.Debug("Configuration: ", string(formattedConfig))
}

func readSshPrivateKey() {
	var err error

	pwd, _ := os.Getwd()
	path := pwd + "/config/keys/id_rsa"
	c.Ssh.PrivateKey, err = os.ReadFile(path)

	if err != nil {
		log.Errorf("could not read a private ssh key from %s: %s", path, err.Error())
	}
}

func new() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var c Config
	err = envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	return c
}
