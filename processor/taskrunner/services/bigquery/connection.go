package bigquery

import (
	"context"
	"errors"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type BigQueryConnection struct {
	projectId   string
	credentials []byte
	ctx         context.Context
	client      *bigquery.Client
}

func (c *BigQueryConnection) Open() error {
	pId := c.projectId
	if pId == "" {
		pId = bigquery.DetectProjectID
	}

	client, err := bigquery.NewClient(c.ctx, pId, option.WithCredentialsJSON(c.credentials))
	if err != nil {
		return err
	}
	c.client = client

	return nil
}

func (c *BigQueryConnection) Test() error {
	q := c.client.Query("SELECT 1")
	rows, err := q.Read(c.ctx)
	if err != nil {
		return err
	}

	values := make([]bigquery.Value, 1)
	err = rows.Next(&values)
	if err != nil {
		return err
	}

	if values[0] != int64(1) {
		return errors.New("unable to execute test query")
	}

	return nil
}

func (c *BigQueryConnection) Query(q string, args ...interface{}) ([]map[string]interface{}, []string, error) {
	it, err := c.client.Query(q).Read(c.ctx)
	if err != nil {
		return nil, nil, err
	}

	var rows []map[string]interface{}
	for {
		var bqValues map[string]bigquery.Value
		err := it.Next(&bqValues)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		values := make(map[string]interface{})
		for k, v := range bqValues {
			values[k] = v
		}

		rows = append(rows, values)
	}

	columnNames := make([]string, 0)
	for _, fieldSchema := range it.Schema {
		columnNames = append(columnNames, fieldSchema.Name)
	}

	return rows, columnNames, nil
}

func (c *BigQueryConnection) Close() error {
	return c.client.Close()
}

func NewConnection(ctx context.Context, projectId string, credentials string) (*BigQueryConnection, error) {
	return &BigQueryConnection{
		projectId:   projectId,
		credentials: []byte(credentials),
		ctx:         ctx,
	}, nil
}
