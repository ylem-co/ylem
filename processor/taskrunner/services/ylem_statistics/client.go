package ylem_statistics

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"ylem_taskrunner/config"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	client *resty.Client
}

func (c *Client) GetAverageMetricValue(pipelineUuid uuid.UUID, period string, periodCount int) (float64, error) {
	var result = struct {
		Value float64 `json:"value"`
	}{}

	url := fmt.Sprintf("/private/pipelines/%s/values/avg/%s/%d", pipelineUuid.String(), period, periodCount)

	log.Tracef("Calling Ylem_statistics endpoint %s", url)

	resp, err := c.client.
		R().
		SetResult(&result).
		Get(url)

	if err != nil {
		log.Errorf("Ylem_statistics call failed, status code: %d, error: %s", resp.StatusCode(), err)
		return result.Value, ErrorServiceUnavilable{}
	}

	if resp.StatusCode() == http.StatusBadRequest {
		errResponse := decodeErrorResponse(resp.Body())
		if errResponse != nil {
			return result.Value, errors.New("bad function call.\n" + strings.Join(errResponse.Errors, "\n"))
		} else {
			return result.Value, errors.New("bad function call")
		}
	}

	if resp.StatusCode() != http.StatusOK {
		log.Error("Ylem_statistics: service unavailable")
		return result.Value, ErrorServiceUnavilable{}
	}

	return result.Value, nil
}

func (c *Client) GetMetricValueQuantile(pipelineUuid uuid.UUID, level float64, period string, periodCount int) (float64, error) {
	var result = struct {
		Value float64 `json:"value"`
	}{}

	url := fmt.Sprintf("/private/pipelines/%s/values/quantile/%f/%s/%d", pipelineUuid.String(), level, period, periodCount)

	log.Tracef("Calling Ylem_statistics endpoint %s", url)

	resp, err := c.client.
		R().
		SetResult(&result).
		Get(url)

	if err != nil {
		log.Errorf("Ylem_statistics call failed, status code: %d, error: %s", resp.StatusCode(), err)
		return result.Value, ErrorServiceUnavilable{}
	}

	if resp.StatusCode() == http.StatusBadRequest {
		errResponse := decodeErrorResponse(resp.Body())
		if errResponse != nil {
			return result.Value, errors.New("bad function call.\n" + strings.Join(errResponse.Errors, "\n"))
		} else {
			return result.Value, errors.New("bad function call")
		}
	}

	if resp.StatusCode() != http.StatusOK {
		log.Error("Ylem_statistics: service unavailable")
		return result.Value, ErrorServiceUnavilable{}
	}

	return result.Value, nil
}

func (c *Client) GetApproximatePipelineExecutionTime(pipelineUuid uuid.UUID) (int, error) {
	var result = struct {
		Value int `json:"value"`
	}{}

	url := fmt.Sprintf("/private/pipelines/%s/duration-stats", pipelineUuid.String())

	log.Tracef("Calling Ylem_statistics endpoint %s", url)

	resp, err := c.client.
		R().
		SetResult(&result).
		Get(url)

	if err != nil {
		log.Errorf("Ylem_statistics call failed, status code: %d, error: %s", resp.StatusCode(), err)
		return result.Value, ErrorServiceUnavilable{}
	}

	if resp.StatusCode() == http.StatusBadRequest {
		errResponse := decodeErrorResponse(resp.Body())
		if errResponse != nil {
			return result.Value, errors.New("bad function call.\n" + strings.Join(errResponse.Errors, "\n"))
		} else {
			return result.Value, errors.New("bad function call")
		}
	}

	if resp.StatusCode() != http.StatusOK {
		log.Error("Ylem_statistics: service unavailable")
		return result.Value, ErrorServiceUnavilable{}
	}

	return result.Value, nil
}

type ErrorServiceUnavilable struct {
}

func (e ErrorServiceUnavilable) Error() string {
	return "Service is unavailable"
}

type errorResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func decodeErrorResponse(body []byte) *errorResponse {
	errResponse := errorResponse{}
	err := json.Unmarshal(body, &errResponse)
	if err != nil {
		return nil
	}

	return &errResponse
}

func NewClient() *Client {
	c := &Client{
		client: resty.New(),
	}

	c.client.SetBaseURL(config.Cfg().YlemStatistics.BaseURL)
	log.Trace("Ylem_statistics client initialized")

	return c
}
