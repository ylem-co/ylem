package api

import (
	"bytes"
	"strings"

	"github.com/go-resty/resty/v2"
)

type File struct {
	Content []byte
	Name    string
}

type Config struct {
	Headers               map[string]string
	File                  *File
	Method                string
	AuthType              string
	AuthBearerToken       string
	AuthBasicUserName     string
	AuthBasicUserPassword string
	AuthHeaderName        string
	AuthHeaderValue       string
	Severity              string
}

const AuthTypeBasic = "Basic"
const AuthTypeBearer = "Bearer"
const AuthTypeHeader = "Header"

func Call(Url string, Payload string, Config Config) ([]byte, error) {
	client := resty.New()
	client.SetRetryCount(3)

	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(Payload)

	switch Config.AuthType {
	case AuthTypeBasic:
		request.SetBasicAuth(Config.AuthBasicUserName, Config.AuthBasicUserPassword)
	case AuthTypeBearer:
		request.SetAuthToken(Config.AuthBearerToken)
	case AuthTypeHeader:
		request.SetHeader(Config.AuthHeaderName, Config.AuthHeaderValue)
	}

	if Config.File != nil {
		request.SetFileReader("file", Config.File.Name, bytes.NewReader(Config.File.Content))
	}

	if Config.Headers != nil {
		for k, v := range Config.Headers {
			request.SetHeader(k, v)
		}
	}

	resp, err := request.Execute(strings.ToUpper(Config.Method), Url)

	return resp.Body(), err
}
