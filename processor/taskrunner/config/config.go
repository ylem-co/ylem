package config

import (
	"os"
	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

func Cfg() config {
	return c
}

type config struct {
	Listen           string     `default:"0.0.0.0:7335"`
	Twilio           twilio     `split_words:"true"`
	Aws              aws        `split_words:"true"`
	LogLevel         string     `split_words:"true" default:"debug"`
	TaskRunner       taskRunner `split_words:"true" envconfig:"TASK_RUNNER"`
	Kafka            kafka
	Ssh              ssh
	Redis            redis
	Tableau          tableau
	Gopyk            gopyk
	YlemStatistics   ylemStatistics
	Openai           openai
}

type taskRunner struct {
	Instances uint `split_words:"true" default:"1"`
}

type kafka struct {
	BootstrapServers                   []string `envconfig:"YLEM_KAFKA_BOOTSTRAP_SERVERS"`
	TaskRunsLoadBalancedTopic          string   `envconfig:"YLEM_KAFKA_TASK_RUNS_LOAD_BALANCED_TOPIC"`
	TaskRunsTopic                      string   `envconfig:"YLEM_KAFKA_TASK_RUNS_TOPIC"`
	TaskRunResultsTopic                string   `envconfig:"YLEM_KAFKA_TASK_RUN_RESULTS_TOPIC"`
	QueryTaskRunResultsTopic           string   `envconfig:"YLEM_KAFKA_QUERY_TASK_RUN_RESULTS_TOPIC"`
	NotificationTaskRunResultsTopic    string   `envconfig:"YLEM_KAFKA_NOTIFICATION_TASK_RUN_RESULTS_TOPIC"`
	LoadBalancerConsumerGroupName      string   `split_words:"true" default:"task_load_balancer"`
	TaskRunnerConsumerGroupName        string   `split_words:"true" default:"task_runner"`
}

type aws struct {
	Region      string `split_words:"true" default:"eu-central-1" envconfig:"AWS_REGION"`
	KmsKeyId    string `split_words:"true" envconfig:"AWS_KMS_KEY_ID"`
}

type twilio struct {
	AccountSid string `envconfig:"YLEM_INTEGRATIONS_TWILIO_ACCOUNT_SID"`
	AuthToken  string `envconfig:"YLEM_INTEGRATIONS_TWILIO_AUTH_TOKEN"`
	NumberFrom string `envconfig:"YLEM_INTEGRATIONS_TWILIO_NUMBER_FROM"`
}

type ssh struct {
	PrivateKey []byte
}

type redis struct {
	Host      string `envconfig:"YLEM_REDIS_HOST"`
	Port      string `envconfig:"YLEM_REDIS_PORT"`
	Password  string `envconfig:"YLEM_REDIS_PASSWORD"`
	KeyPrefix string `split_words:"true" default:"taskrunner."`
}

type tableau struct {
	HttpWrapperBaseUrl string `envconfig:"TASK_RUNNER_TABLEAU_HTTP_WRAPPER_BASE_URL"`
}

type gopyk struct {
	BaseUrl string `envconfig:"TASK_RUNNER_GOPYK_BASE_URL"`
}

type ylemStatistics struct {
	BaseURL string `envconfig:"YLEM_STATISTICS_BASE_URL"`
}

type openai struct {
	GptKey string `split_words:"true"`
	Model  string `split_words:"true" default:"gpt-4o-mini"`
}

var c config

func init() {
	c = new()

	readSshPrivateKey()

	lvl, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(lvl)

	formattedConfig, _ := json.MarshalIndent(c, "", "    ")
	log.Debug("Configuration: ", string(formattedConfig))
}

func new() config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var s config
	err = envconfig.Process("TASK_RUNNER", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	return s
}

func readSshPrivateKey() {
	var err error

	pwd, _ := os.Getwd()
	path := pwd + "/config/keys/id_rsa"
	c.Ssh.PrivateKey, err = os.ReadFile(path)

	if err != nil {
		log.Infof("could not read a private ssh key from %s: %s", path, err.Error())
	}
}
