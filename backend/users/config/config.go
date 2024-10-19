package config

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

var c config

func Cfg() config {
	return c
}

type config struct {
	Listen   string `default:"0.0.0.0:7333"`
	LogLevel string `envconfig:"YLEM_LOG_LEVEL"`
	DBConfig struct {
		Name     string `envconfig:"USERS_DATABASE_NAME"`
		Host     string `envconfig:"YLEM_DATABASE_HOST"`
		Port     string `envconfig:"YLEM_DATABASE_PORT"`
		User     string `envconfig:"YLEM_DATABASE_USER"`
		Password string `envconfig:"YLEM_DATABASE_PASSWORD"`
	}
	RedisDBConfig struct {
		Host     string `envconfig:"YLEM_REDIS_HOST"`
		Port     string `envconfig:"YLEM_REDIS_PORT"`
		Password string `envconfig:"YLEM_REDIS_PASSWORD"`
	}
	NetworkConfig struct {
		YlemIntegrationsBaseUrl string `envconfig:"YLEM_INTEGRATIONS_BASE_URL"`
		YlemPipelinesBaseUrl    string `envconfig:"YLEM_PIPELINES_BASE_URL"`
	}
	Aws struct {
		KmsKeyId string `split_words:"true" envconfig:"AWS_KMS_KEY_ID"`
	}
	Google struct {
		ClientId     string `split_words:"true" envconfig:"USERS_GOOGLE_CLIENT_ID"`
		ClientSecret string `split_words:"true" envconfig:"USERS_GOOGLE_CLIENT_SECRET"`
		CallbackUrl  string `split_words:"true" envconfig:"USERS_GOOGLE_CALLBACK_URL"`
	}
}

func init() {
	c = new()

	formattedConfig, _ := json.MarshalIndent(c, "", "    ")
	log.Debug("Configuration: ", string(formattedConfig))

	goth.UseProviders(
		google.New(c.Google.ClientId, c.Google.ClientSecret, c.Google.CallbackUrl, "email", "profile"),
	)

	maxAge := 86400 * 30  // 30 days
	store := sessions.NewCookieStore([]byte(""))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.Domain = ".ylem.co"
	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteNoneMode
	store.Options.Secure = true

	gothic.Store = store
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
	err = envconfig.Process("", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	return s
}
