package sqlIntegrations

import (
	"context"
	"fmt"
	"database/sql"
	_ "github.com/alexbrainman/odbc"
)

type RedshiftSQLIntegrationConnection struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Database string
	DB       *sql.DB
	Context  context.Context
}

func (s *RedshiftSQLIntegrationConnection) Open() error {
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

func (s *RedshiftSQLIntegrationConnection) Close() error {
	return s.DB.Close()
}

func (s *RedshiftSQLIntegrationConnection) Test() error {
	return s.DB.PingContext(s.Context)
}

func (s *RedshiftSQLIntegrationConnection) Prepare(query string) (*sql.Stmt, error) {
	return s.DB.PrepareContext(s.Context, query)
}

func (s *RedshiftSQLIntegrationConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.ExecContext(s.Context, query, args...)
}

func (s *RedshiftSQLIntegrationConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.QueryContext(s.Context, query, args...)
}
