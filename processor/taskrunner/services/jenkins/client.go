package jenkins

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

var instance *Jenkins

type Jenkins struct {
	client *resty.Client
	ctx context.Context
}

func (j Jenkins) RunBuild(baseUrl string, project string, token string) error {
	log.Tracef("jenkins: run build")
	url := fmt.Sprintf("%s/job/%s/build?token=%s", baseUrl, project, token)

	response, err := j.client.
		R().
		Post(url)

	if err != nil {
		return err
	}

	if response.StatusCode() >= http.StatusBadRequest {
		log.Debug(string(response.Body()))

		return fmt.Errorf("jenkins: run build: expected http 2xx, got %s", response.Status())
	}

	return nil
}

func Instance(ctx context.Context) *Jenkins {
	if instance == nil {
		instance = &Jenkins{
			client: resty.New(),
			ctx:    ctx,
		}

	}

	return instance
}
