package config

import (
	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type config struct {
	LogLevel      string        `split_words:"true" default:"info"`
	Listen        string        `default:"0.0.0.0:7339"`
	DB            db            `envconfig:"DB"`
	NetworkConfig networkConfig `envconfig:"NETWORK"`
}

type networkConfig struct {
	AuthorizationCheckUrl string `envconfig:"YLEM_AUTHORIZATION_CHECK_URL"`
	PermissionCheckUrl    string `envconfig:"YLEM_PERMISSION_CHECK_URL"`
	YlemUsersBaseUrl      string `envconfig:"YLEM_USERS_BASE_URL"`
	YlemPipelinesBaseUrl  string `envconfig:"YLEM_PIPELINES_BASE_URL"`
	YlemStatisticsBaseUrl string `envconfig:"YLEM_STATISTICS_BASE_URL"`
}

type db struct {
	Name     string `envconfig:"API_DATABASE_NAME"`
	Host     string `envconfig:"YLEM_DATABASE_HOST"`
	Port     string `envconfig:"YLEM_DATABASE_PORT"`
	User     string `envconfig:"YLEM_DATABASE_USER"`
	Password string `envconfig:"YLEM_DATABASE_PASSWORD"`
}

func Cfg() config {
	return c
}

var c config

func init() {
	c = new()
	lvl, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(lvl)

	formattedConfig, _ := json.MarshalIndent(c, "", "    ")
	log.Trace("Configuration: ", string(formattedConfig))
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

	var c config
	err = envconfig.Process("API", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	return c
}
