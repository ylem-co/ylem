package stats

import (
	"errors"
	"io"
	"net/http"
	"ylem_api/config"
	"ylem_api/helpers"
	"ylem_api/model/entity"
	"ylem_api/service"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ModifyResponseFunc func(statusCode int, responseBody []byte, headers http.Header) (int, []byte, http.Header, error)

func ProxyRequestTo(url string, t *entity.OauthToken, w http.ResponseWriter, r *http.Request, modRespFunc ModifyResponseFunc) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	ylemStatisticsClient := service.NewProxyingClient(config.Cfg().NetworkConfig.YlemStatisticsBaseUrl, t.InternalToken, r.Context())

	statusCode, responseBody, headers, err := ylemStatisticsClient.ProxyTo(r.Method, url, body, r.Header)
	if errors.Is(err, service.ErrorNotFound) {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	log.Tracef("Call to Ylem_Statistics service, %s %s, status code: %d, response: %s", r.Method, url, statusCode, string(responseBody))

	if modRespFunc != nil {
		statusCode, responseBody, headers, err = modRespFunc(statusCode, responseBody, headers)
		if err != nil {
			log.Error(err)
			helpers.HttpReturnErrorInternal(w)
			return
		}
	}

	helpers.SetHeaders(w, headers)
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(responseBody)), 10))
	w.WriteHeader(statusCode)
	
	_, err = w.Write(responseBody)
	if err != nil {
	  log.Error(err)
	}
}

func ProxyRequest(t *entity.OauthToken, w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitAfterN(r.URL.String(), "/stats", 2)
	log.Tracef("URL parts: %v", parts)
	if len(parts) != 2 {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	ProxyRequestTo(parts[1], t, w, r, nil)
}
