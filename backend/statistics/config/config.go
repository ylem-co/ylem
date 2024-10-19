package config

import (
	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

func Cfg() config {
	return c
}

type config struct {
	Listen                string `default:"0.0.0.0:7332"`
	AuthorizationCheckUrl string `envconfig:"YLEM_AUTHORIZATION_CHECK_URL"`
	PermissionCheckUrl    string `envconfig:"YLEM_PERMISSION_CHECK_URL"`
	DB                    db     `envconfig:"DB"`
	Kafka                 kafka  `split_words:"true"`
	LogLevel              string `envconfig:"YLEM_LOG_LEVEL"`
}

type kafka struct {
	BootstrapServers     []string `split_words:"true" envconfig:"YLEM_KAFKA_BOOTSTRAP_SERVERS"`
	TaskRunResultsTopic  string   `split_words:"true" envconfig:"YLEM_KAFKA_TASK_RUN_RESULTS_TOPIC"`
	ConsumerGroupName    string   `split_words:"true" default:"task_run_result_consumer"`
}

type db struct {
	DSN        string
	StatsTable string `split_words:"true" default:"task_runs"`
}

var c config

func init() {
	c = new()
	formattedConfig, _ := json.MarshalIndent(c, "", "    ")
	log.Debug("Configuration: ", string(formattedConfig))
}

func new() config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var s config
	err = envconfig.Process("STATISTICS", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	return s
}
