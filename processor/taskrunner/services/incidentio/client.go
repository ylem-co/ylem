package incidentio

import (
	"context"
	"fmt"
	"net/http"
	"ylem_taskrunner/config"
	"ylem_taskrunner/services/aws/kms"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

var srv *IncidentIo

type IncidentIo struct {
	client *resty.Client
}

type Incident struct {
	IdempotencyKey string `json:"idempotency_key"`
	Mode           string `json:"mode"`
	Name           string `json:"name"`
	SeverityId     string `json:"severity_id"`
	Status         string `json:"status"`
	Summary        string `json:"summary"`
	Visibility     string `json:"visibility"`
}

func (i *IncidentIo) CreateIncident(ApiKey string, inc Incident) error {
	log.Tracef("incident.io: creating incident")
	response, err := i.client.
		R().
		SetAuthToken(ApiKey).
		SetBody(inc).
		Post("v1/incidents")

	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusOK {
		log.Debug(string(response.Body()))

		return fmt.Errorf("incident.io: getting severities, expected http 200, got %s", response.Status())
	}

	return nil
}

func Instance() *IncidentIo {
	if srv == nil {
		srv = &IncidentIo{
			client: resty.New().SetBaseURL("https://api.incident.io/"),
		}
	}

	return srv
}

func DecryptKeyAndCreateIncident(ctx context.Context, dataKey []byte, apiKey []byte, inc Incident) error {
	decryptedDataKey, err := kms.DecryptDataKey(
		ctx,
		config.Cfg().Aws.KmsKeyId,
		dataKey,
	)

	if err != nil {
		return err
	}

	decryptedApiKey, err := kms.Decrypt(apiKey, decryptedDataKey)

	if err != nil {
		return err
	}

	return Instance().CreateIncident(string(decryptedApiKey), inc)
}
