package incidentio

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var srv *IncidentIo

type IncidentIo struct {
	client *resty.Client
}

type Severity struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"rank"`
}

type ioseverities struct {
	Severities []Severity `json:"severities"`
}

func (i *IncidentIo) GetSeverities(ApiKey string) ([]Severity, error){
	log.Tracef("incident.io: getting severities")
	var result ioseverities
	response, err := i.client.
		R().
		SetAuthToken(ApiKey).
		SetResult(&result).
		Get("v1/severities")

	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		log.Debug(string(response.Body()))

		return nil, fmt.Errorf("incident.io: getting severities, expected http 200, got %s", response.Status())
	}

	return result.Severities, nil
}

func Instance() *IncidentIo {
	if srv == nil {
		srv = &IncidentIo{
			client: resty.New().SetBaseURL("https://api.incident.io/"),
		}
	}

	return srv
}
