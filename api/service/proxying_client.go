package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

var ErrorServiceUnavailable = errors.New("service unavailable")
var ErrorNotFound = errors.New("not found")

type ProxyingClient interface {
	ProxyTo(method, url string, body []byte, headers http.Header) (int, []byte, http.Header, error)
}

type RateLimitClient struct {
	innerClient ProxyingClient
	limiter     *rate.Limiter
	ctx         context.Context
}

func (c *RateLimitClient) LimitRate() {
	_ = c.limiter.Wait(c.ctx)
}

func (c *RateLimitClient) ProxyTo(method, url string, body []byte, headers http.Header) (int, []byte, http.Header, error) {
	c.LimitRate()
	return c.innerClient.ProxyTo(method, url, body, headers)
}

type AuthenticatedProxyingClient struct {
	client      *resty.Client
	bearerToken string
}

func (c *AuthenticatedProxyingClient) R() *resty.Request {
	return c.client.R().SetAuthToken(c.bearerToken)
}

func (c *AuthenticatedProxyingClient) ProxyTo(method, url string, body []byte, headers http.Header) (int, []byte, http.Header, error) {
	resp, err := c.
		R().
		SetBody(body).
		SetHeaderMultiValues(headers).
		Execute(method, url)

	statusCode := 0
	respBody := make([]byte, 0)
	respHeaders := http.Header{}
	if resp != nil {
		statusCode = resp.StatusCode()
		respBody = resp.Body()
		respHeaders = resp.Header()
	}

	if err != nil {
		log.Errorf("Service call failed, status code: %d, error: %s", statusCode, err.Error())
		return statusCode, respBody, respHeaders, ErrorServiceUnavailable
	}

	return statusCode, respBody, respHeaders, nil
}

func NewProxyingClient(baseUrl string, bearerToken string, ctx context.Context) ProxyingClient {
	apc := &AuthenticatedProxyingClient{
		client:      resty.New(),
		bearerToken: bearerToken,
	}
	apc.client.SetBaseURL(baseUrl)

	rlc := &RateLimitClient{
		innerClient: apc,
		limiter:     rate.NewLimiter(100, 1),
		ctx:         ctx,
	}

	return rlc
}
