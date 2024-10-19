package es

import (
	"context"
	
	essqlclient "github.com/ylem-co/es-sql-client"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type Connection struct {
	ctx    context.Context
	client *essqlclient.ES
}

func (c *Connection) Open() error {
	return nil
}

func (c *Connection) Test() error {
	_, err := c.client.Version(nil)

	return err
}

func (c *Connection) Close() error {
	return nil
}

func (c *Connection) Query(q string, args ...interface{}) ([]map[string]interface{}, []string, error) {
	res, err := c.client.SqlQuery(q)
	if err != nil {
		return nil, nil, err
	}

	cols := make([]string, len(res.Columns))
	for _,v := range res.Columns {
		cols = append(cols, v.Name)
	}

	return res.Rows, cols, nil
}

func NewConnection(ctx context.Context, url string, user string, password string, version *uint8) (*Connection, error) {
	es := essqlclient.CreateWithBaseUrl(ctx, url, nil, func(c *resty.Client) {
		if password == "" {
			return
		}

		c.SetBasicAuth(user, password)
	})

	if version != nil {
		_, err := es.Version(version)

		if err != nil {
			return nil, err
		}
	}

	es.SetLogger(log.StandardLogger())

	return &Connection{
		ctx:    ctx,
		client: &es,
	}, nil
}
