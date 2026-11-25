package anthropic

import (
	"ylem_taskrunner/config"

	"github.com/go-resty/resty/v2"
)

const apiVersion = "2026-10-21"

var srv *Anthropic

type Anthropic struct {
	client *resty.Client
}

func Instance() *Anthropic {
	if srv == nil {
		srv = &Anthropic{
			client: resty.New().
				SetBaseURL("https://api.anthropic.com/v1").
				SetAuthToken(config.Cfg().Anthropic.Key).
				SetHeader("anthropic-version", apiVersion),
		}
	}

	return srv
}
