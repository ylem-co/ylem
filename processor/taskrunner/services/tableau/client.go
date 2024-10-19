package tableau

import (
	"ylem_taskrunner/config"

	"github.com/go-resty/resty/v2"
)

type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type client struct {
	client *resty.Client
}

type TableauServerError struct {
	msg string
}

func (e TableauServerError) Error() string {
	return e.msg
}

func (c *client) Insert(server, username, password, siteName, projectName, datasourceName, mode string, columns []Column, rows [][]interface{}) error {
	reqData := map[string]interface{}{
		"mode": mode,
		"connection": map[string]string{
			"server":          server,
			"username":        username,
			"password":        password,
			"site_name":       siteName,
			"project_name":    projectName,
			"datasource_name": datasourceName,
		},
		"table": map[string]interface{}{
			"name":    datasourceName,
			"columns": columns,
			"rows":    rows,
		},
	}

	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqData).
		Post("/private/tableau/insert")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return TableauServerError{msg: resp.String()}
	}

	return nil
}

func NewClient() *client {
	return &client{
		client: resty.New().SetBaseURL(config.Cfg().Tableau.HttpWrapperBaseUrl),
	}
}
