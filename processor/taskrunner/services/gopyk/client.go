package gopyk

import (
	"fmt"
	"net/http"
	"ylem_taskrunner/config"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

var srv *Gopyk

type Gopyk struct {
	client *resty.Client
}

type Request struct {
	Code  string `json:"code"`
	Type  string `json:"type"`
	Input string `json:"input"`
}

func (g *Gopyk) Evaluate(r Request) ([]byte, error) {
	log.Tracef("gopyk: evaluation")
	response, err := g.client.
		R().
		SetBody(r).
		Post("")

	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		log.Debug(string(response.Body()))

		return nil, fmt.Errorf("python execution failed: %s", string(response.Body()))
	}

	return response.Body(), nil
}

func Instance() *Gopyk {
	if srv == nil {
		srv = &Gopyk{
			client: resty.New().SetBaseURL(config.Cfg().Gopyk.BaseUrl),
		}
	}

	return srv
}
