package ylem_users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"ylem_api/config"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

var ErrorServiceUnavailable = errors.New("Ylem_Users service is unavailable")

type Client struct {
	client  *resty.Client
	limiter *rate.Limiter
}

func (c *Client) IssueJWT(OrganizationUuid uuid.UUID, UserUuid uuid.UUID) (string, error) {
	_ = c.limiter.Wait(context.Background())
	url := "/private/jwt-tokens/"

	log.Tracef("Calling Ylem_Users endpoint %s", url)

	resp, err := c.client.
		R().
		SetBody(map[string]interface{}{
			"organization_uuid": OrganizationUuid.String(),
			"user_uuid":         UserUuid.String(),
		}).
		Post(url)

	if err != nil {
		log.Error("Ylem_Users call error: ", err)
		return "", ErrorServiceUnavailable
	}

	if resp.StatusCode() != http.StatusCreated {
		log.Errorf("Ylem_Users: service is unavailable. Response code: %d, response: %s", resp.StatusCode(), string(resp.Body()))
		return "", errors.New("unable to issue JWT")
	}

	decodedBody := make(map[string]string)
	err = json.Unmarshal(resp.Body(), &decodedBody)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return decodedBody["token"], nil
}

func NewClient() *Client {
	c := &Client{
		client:  resty.New(),
		limiter: rate.NewLimiter(100, 1),
	}

	c.client.SetBaseURL(config.Cfg().NetworkConfig.YlemUsersBaseUrl)
	log.Trace("Ylem_Users client initialized")

	return c
}
