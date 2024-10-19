package config

import (
	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel string `split_words:"true" default:"info"`
	Listen   string `default:"0.0.0.0:7336"`
	DBConfig struct {
		Name     string `envconfig:"PIPELINES_DATABASE_NAME"`
		Host     string `envconfig:"YLEM_DATABASE_HOST"`
		Port     string `envconfig:"YLEM_DATABASE_PORT"`
		User     string `envconfig:"YLEM_DATABASE_USER"`
		Password string `envconfig:"YLEM_DATABASE_PASSWORD"`
	}

	NetworkConfig struct {
		AuthorizationCheckUrl string `envconfig:"YLEM_AUTHORIZATION_CHECK_URL"`
		PermissionCheckUrl    string `envconfig:"YLEM_PERMISSION_CHECK_URL"`
		UpdateConnectionsUrl  string `envconfig:"YLEM_UPDATE_CONNECTIONS_URL"`
	}

	Kafka kafka

	YlemIntegrations ylemIntegrations

	SystemOrganizationUuid string `split_words:"true" envconfig:"PIPELINES_SYSTEM_ORGANIZATION_UUID"`
}

type kafka struct {
	BootstrapServers                 []string `envconfig:"YLEM_KAFKA_BOOTSTRAP_SERVERS"`
	TaskRunsTopic                    string   `envconfig:"YLEM_KAFKA_TASK_RUNS_TOPIC"`
	TaskRunResultsTopic              string   `envconfig:"YLEM_KAFKA_TASK_RUN_RESULTS_TOPIC"`
	TriggerListenerConsumerGroupName string   `split_words:"true" default:"trigger_listener"`
}

type ylemIntegrations struct {
	BaseURL string `envconfig:"YLEM_INTEGRATIONS_BASE_URL"`
}

func Cfg() Config {
	return c
}

var c Config

func init() {
	c = new()
	lvl, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(lvl)
	log.SetReportCaller(true)

	formattedConfig, _ := json.MarshalIndent(c, "", "    ")
	log.Debug("Configuration: ", string(formattedConfig))
	log.Infof("Kafka bootstrap servers: %s", c.Kafka.BootstrapServers)
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
