package sql

import (
	"context"
	"database/sql"
	"fmt"
	"ylem_integrations/helpers"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLIntegrationConnection struct {
	Host       string
	Port       uint16
	User       string
	Password   string
	Database   string
	Context    context.Context
	CancelFunc context.CancelFunc
	DB         *sql.DB
}

func (s *MySQLIntegrationConnection) Open() error {
	var err error
	s.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s", s.User, s.Password, helpers.MySqlSSHTCPNetName, s.Host, s.Port, s.Database))
	if err != nil {
		return err
	}

	s.Context = context.Background()
	s.CancelFunc = nil

	return nil
}

func (s *MySQLIntegrationConnection) OpenSsh(host string, port uint16, user string) error {
	return s.OpenSshContext(context.Background(), host, port, user)
}

func (s *MySQLIntegrationConnection) OpenSshContext(ctx context.Context, host string, port uint16, user string) error {
	var err error
	s.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s", s.User, s.Password, helpers.MySqlSSHTCPNetName, s.Host, s.Port, s.Database))
	if err != nil {
		return err
	}

	s.Context, s.CancelFunc = helpers.CreateSSHConfigurationContext(
		ctx,
		host,
		port,
		user,
	)

	return nil
}

func (s *MySQLIntegrationConnection) Close() error {
	err := s.DB.Close()
	if s.CancelFunc != nil {
		s.CancelFunc()
	}

	return err
}

func (s *MySQLIntegrationConnection) Test() error {
	return s.DB.PingContext(s.Context)
}

func (s *MySQLIntegrationConnection) Prepare(query string) (*sql.Stmt, error) {
	return s.DB.PrepareContext(s.Context, query)
}

func (s *MySQLIntegrationConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.ExecContext(s.Context, query, args...)
}

func (s *MySQLIntegrationConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.QueryContext(s.Context, query, args...)
}

func (s *MySQLIntegrationConnection) ShowDatabases() ([]string, error) {
	rows, err := s.DB.QueryContext(s.Context, "SHOW DATABASES")

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

func (s *MySQLIntegrationConnection) ShowTables(db string) ([]string, error) {
	rows, err := s.DB.QueryContext(s.Context, "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?", db)

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

func (s *MySQLIntegrationConnection) DescribeTable(db string, table string) ([]string, error) {
	stmnt, err := s.DB.PrepareContext(
		s.Context,
		"SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ?",
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
