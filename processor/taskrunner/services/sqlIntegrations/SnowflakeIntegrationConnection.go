package sqlIntegrations

import (
	"context"
	"fmt"
	"database/sql"

	_ "github.com/snowflakedb/gosnowflake"
)

type SnowflakeSQLIntegrationConnection struct {
	AccountId string
	User      string
	Password  string
	Database  string
	DB        *sql.DB
}

func (s *SnowflakeSQLIntegrationConnection) Open() error {
	var err error
	s.DB, err = sql.Open("snowflake", fmt.Sprintf("%s:%s@%s/%s", s.User, s.Password, s.AccountId, s.Database))
	if err != nil {
		return err
	}

	return nil
}

func (s *SnowflakeSQLIntegrationConnection) Close() error {
	return s.DB.Close()
}

func (s *SnowflakeSQLIntegrationConnection) Test() error {
	return s.DB.PingContext(context.Background())
}

func (s *SnowflakeSQLIntegrationConnection) Prepare(query string) (*sql.Stmt, error) {
	return s.DB.PrepareContext(context.Background(), query)
}

func (s *SnowflakeSQLIntegrationConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.ExecContext(context.Background(), query, args...)
}

func (s *SnowflakeSQLIntegrationConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.QueryContext(context.Background(), query, args...)
}
