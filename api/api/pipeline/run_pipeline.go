package pipeline

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"ylem_api/config"
	"ylem_api/helpers"
	"ylem_api/model/entity"
	"ylem_api/service"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type ModifyResponseFunc func(statusCode int, responseBody []byte, headers http.Header) (int, []byte, http.Header, error)

func RunPipeline(t *entity.OauthToken, w http.ResponseWriter, r *http.Request) {
	pipelineUuid := mux.Vars(r)["pipelineUuid"]
	url := fmt.Sprintf("/pipeline/%s/run", pipelineUuid)

	var mf ModifyResponseFunc = func(statusCode int, responseBody []byte, headers http.Header) (int, []byte, http.Header, error) {
		body := struct {
			Results         []map[string]interface{} `json:"results"`
			PipelineRunUuid string                   `json:"pipeline_run_uuid"`
		}{}
		err := json.Unmarshal(responseBody, &body)
		if err != nil {
			return 0, make([]byte, 0), make(http.Header), err
		}

		if len(body.Results) == 0 {
			return statusCode, make([]byte, 0), headers, nil
		}

		newBody, err := json.Marshal(map[string]interface{}{
			"pipeline_run_uuid": body.PipelineRunUuid,
		})

		if err != nil {
			return 0, make([]byte, 0), make(http.Header), err
		}

		return statusCode, newBody, headers, nil
	}

	ProxyRequestTo(url, t, w, r, mf)
}

func ProxyRequestTo(url string, t *entity.OauthToken, w http.ResponseWriter, r *http.Request, modRespFunc ModifyResponseFunc) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	ylemPipelinesClient := service.NewProxyingClient(config.Cfg().NetworkConfig.YlemPipelinesBaseUrl, t.InternalToken, r.Context())

	statusCode, responseBody, headers, err := ylemPipelinesClient.ProxyTo(r.Method, url, body, r.Header)
	if errors.Is(err, service.ErrorNotFound) {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	log.Tracef("Call to Ylem_Pipelines service, %s %s, status code: %d, response: %s", r.Method, url, statusCode, string(responseBody))

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
	parts := strings.SplitAfter(r.URL.String(), "/pipelines")
	if len(parts) != 2 {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	url := fmt.Sprintf("/pipelines%s", parts[1])
	ProxyRequestTo(url, t, w, r, nil)
}
