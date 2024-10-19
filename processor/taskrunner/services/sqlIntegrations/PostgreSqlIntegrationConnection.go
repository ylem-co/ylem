package sqlIntegrations

import (
	"context"
	"fmt"
	"log"
	"database/sql"
	"ylem_taskrunner/helpers/postgresql"

	_ "github.com/lib/pq"
	"github.com/sethvargo/go-password/password"
)

type PostgreSqlSQLIntegrationConnection struct {
	Host       string
	Port       uint16
	User       string
	Password   string
	Database   string
	DriverId   string
	SSLEnabled bool
	DB         *sql.DB
}

func init() {
	driver := &postgresql.Driver{}
	sql.Register("postgres+ssh", driver)

	postgresql.InitSshPool()
}

func (s *PostgreSqlSQLIntegrationConnection) Open() error {
	var err error
	sslModeString := "disable"
	if s.SSLEnabled {
		sslModeString = "require"
	}
	s.DB, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", s.User, s.Password, s.Host, s.Port, s.Database, sslModeString))
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgreSqlSQLIntegrationConnection) OpenSsh(Host string, Port uint16, User string) error {
	var err error

	uniqueCode, err := password.Generate(10, 5, 0, false, true)
	if err != nil {
		log.Println(err)

		return err
	}

	s.DriverId = Host + "." + uniqueCode
	sslModeString := "disable"
	if s.SSLEnabled {
		sslModeString = "require"
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&sshHost=%s&sshPort=%d&sshUser=%s&sshDriverId=%s",
		s.User,
		s.Password,
		s.Host,
		s.Port,
		s.Database,
		sslModeString,
		Host,
		Port,
		User,
		s.DriverId,
	)
	s.DB, err = sql.Open("postgres+ssh", dsn)

	return err
}

func (s *PostgreSqlSQLIntegrationConnection) Test() error {
	return s.DB.PingContext(context.Background())
}

func (s *PostgreSqlSQLIntegrationConnection) Prepare(query string) (*sql.Stmt, error) {
	return s.DB.PrepareContext(context.Background(), query)
}

func (s *PostgreSqlSQLIntegrationConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.ExecContext(context.Background(), query, args...)
}

func (s *PostgreSqlSQLIntegrationConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.QueryContext(context.Background(), query, args...)
}

func (s *PostgreSqlSQLIntegrationConnection) Close() error {
	if s.DriverId != "" {
		_ = postgresql.SshPool.Unregister(s.DriverId)
	}
	return s.DB.Close()
}
