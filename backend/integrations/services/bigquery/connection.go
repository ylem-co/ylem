package bigquery

import (
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"
)

type BigQueryConnection struct {
	projectId   *string
	credentials []byte
	ctx         context.Context
	client      *bigquery.Client
}

func (c *BigQueryConnection) Open() error {
	var pId string
	if c.projectId != nil {
		pId = *c.projectId
	}

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

func (c *BigQueryConnection) Close() error {
	return c.client.Close()
}

func (s *BigQueryConnection) ShowDatabases() ([]string, error) {
	datasets := make([]string, 0)
	it := s.client.Datasets(s.ctx)
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		datasets = append(datasets, dataset.DatasetID)
	}

	return datasets, nil
}

func (s *BigQueryConnection) ShowTables(db string) ([]string, error) {
	tables := make([]string, 0)
	it := s.client.Dataset(db).Tables(s.ctx)
	for {
		table, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		tables = append(tables, table.TableID)
	}

	return tables, nil
}

func (s *BigQueryConnection) DescribeTable(db string, table string) ([]string, error) {
	columns := make([]string, 0)
	md, err := s.client.Dataset(db).Table(table).Metadata(s.ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range md.Schema {
		columns = append(columns, v.Name)
	}

	return columns, nil
}

func NewConnection(ctx context.Context, projectId *string, credentials string) (*BigQueryConnection, error) {

	return &BigQueryConnection{
		projectId:   projectId,
		credentials: []byte(credentials),
		ctx:         ctx,
	}, nil
}
