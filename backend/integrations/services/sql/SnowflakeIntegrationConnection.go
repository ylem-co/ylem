package sql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/snowflakedb/gosnowflake"
)

type SnowflakeIntegrationConnection struct {
	AccountId string
	User      string
	Password  string
	Database  string
	DB        *sql.DB
}

func (s *SnowflakeIntegrationConnection) Open() error {
	var err error
	s.DB, err = sql.Open("snowflake", fmt.Sprintf("%s:%s@%s/%s", s.User, s.Password, s.AccountId, s.Database))
	if err != nil {
		return err
	}

	return nil
}

func (s *SnowflakeIntegrationConnection) Close() error {
	return s.DB.Close()
}

func (s *SnowflakeIntegrationConnection) Test() error {
	return s.DB.PingContext(context.Background())
}

func (s *SnowflakeIntegrationConnection) Prepare(query string) (*sql.Stmt, error) {
	return s.DB.PrepareContext(context.Background(), query)
}

func (s *SnowflakeIntegrationConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.ExecContext(context.Background(), query, args...)
}

func (s *SnowflakeIntegrationConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.QueryContext(context.Background(), query, args...)
}

func (s *SnowflakeIntegrationConnection) ShowDatabases() ([]string, error) {
	rows, err := s.DB.Query("SELECT DATABASE_NAME FROM information_schema.DATABASES")

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

func (s *SnowflakeIntegrationConnection) ShowTables(db string) ([]string, error) {
	rows, err := s.DB.Query("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_CATALOG = ?", db)

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

func (s *SnowflakeIntegrationConnection) DescribeTable(table string) ([]string, error) {
	stmnt, err := s.DB.Prepare(
		"SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_NAME = ? AND TABLE_CATALOG = ?",
	)
	if err != nil {
		return nil, err
	}

	rows, err := stmnt.Query(table, s.Database)
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
