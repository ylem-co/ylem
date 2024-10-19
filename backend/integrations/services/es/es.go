package es

import (
	"context"
	essqlclient "github.com/ylem-co/es-sql-client"
	"github.com/go-resty/resty/v2"
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

func (c *Connection) Version() (uint8, error) {
	return c.client.Version(nil)
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

	return &Connection{
		ctx:    ctx,
		client: &es,
	}, nil
}
