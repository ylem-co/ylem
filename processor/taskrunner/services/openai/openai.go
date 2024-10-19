package openai

import (
	"ylem_taskrunner/config"

	"github.com/go-resty/resty/v2"
)

var srv *OpenAi

type OpenAi struct {
	client *resty.Client
}

func Instance() *OpenAi {
	if srv == nil {
		srv = &OpenAi{
			client: resty.New().
				SetBaseURL("https://api.openai.com/v1").
				SetAuthToken(config.Cfg().Openai.GptKey),
		}
	}

	return srv
}
