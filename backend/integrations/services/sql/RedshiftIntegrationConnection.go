package sql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/alexbrainman/odbc"
)

type RedshiftIntegrationConnection struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Database string
	DB       *sql.DB
	Context  context.Context
}

func (s *RedshiftIntegrationConnection) Open() error {
	var err error
	s.DB, err = sql.Open(
		"odbc",
		fmt.Sprintf(
			"Driver={Amazon Redshift (x64)}; Server=%s; Port=%d; Database=%s; UID=%s; PWD=%s",
			s.Host,
			s.Port,
			s.Database,
			s.User,
			s.Password,
		),
	)
	if err != nil {
		return err
	}

	s.Context = context.Background()

	return nil
}

func (s *RedshiftIntegrationConnection) Close() error {
	return s.DB.Close()
}

func (s *RedshiftIntegrationConnection) Test() error {
	return s.DB.PingContext(s.Context)
}

func (s *RedshiftIntegrationConnection) Prepare(query string) (*sql.Stmt, error) {
	return s.DB.PrepareContext(s.Context, query)
}

func (s *RedshiftIntegrationConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.ExecContext(s.Context, query, args...)
}

func (s *RedshiftIntegrationConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.QueryContext(s.Context, query, args...)
}

func (s *RedshiftIntegrationConnection) ShowDatabases() ([]string, error) {
	rows, err := s.DB.Query("SELECT datname FROM pg_database")

	if err != nil {
		return nil, err
	}

	tables := make([]string, 0)
	var table string
	for rows.Next() {
		err = rows.Scan(&table)
		if err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}

	return tables, nil
}

func (s *RedshiftIntegrationConnection) ShowTables(db string) ([]string, error) {
	stmnt, err := s.DB.Prepare(`
SELECT
     table_name
FROM
    information_schema.tables
WHERE
    table_type = 'BASE TABLE'
AND
    table_catalog = $1
    AND table_schema NOT IN ('pg_catalog', 'information_schema');
`)

	if err != nil {
		return nil, err
	}

	rows, err := stmnt.Query(db)
	if err != nil {
		return nil, err
	}

	tables := make([]string, 0)
	var table string
	for rows.Next() {
		err = rows.Scan(&table)
		if err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}

	return tables, nil
}

func (s *RedshiftIntegrationConnection) DescribeTable(db string, table string) ([]string, error) {
	stmnt, err := s.DB.Prepare(
		`
SELECT
	column_name
FROM
	information_schema.columns
WHERE
	table_name = $1
	AND table_catalog = $2`,
	)
	if err != nil {
		return nil, err
	}

	rows, err := stmnt.Query(table, db)
	if err != nil {
		return nil, err
	}

	columns := make([]string, 0)
	var column string
	for rows.Next() {
		err = rows.Scan(&column)
		if err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	return columns, err
}
