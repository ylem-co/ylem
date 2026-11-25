package sqlIntegrations

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/url"
	"strconv"
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

func sslModeString(enabled bool) string {
	if enabled {
		return "require"
	}

	return "disable"
}

// buildPostgresDSN percent-encodes the user and password. Supabase database
// passwords routinely include characters that break an unescaped URL.
func buildPostgresDSN(user, password, host string, port uint16, database, sslMode string, extra map[string]string) string {
	query := url.Values{}
	query.Set("sslmode", sslMode)
	for key, value := range extra {
		query.Set(key, value)
	}

	dsn := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, password),
		Host:     net.JoinHostPort(host, strconv.Itoa(int(port))),
		Path:     "/" + database,
		RawQuery: query.Encode(),
	}

	return dsn.String()
}

func init() {
	driver := &postgresql.Driver{}
	sql.Register("postgres+ssh", driver)

	postgresql.InitSshPool()
}

func (s *PostgreSqlSQLIntegrationConnection) Open() error {
	var err error
	s.DB, err = sql.Open("postgres", buildPostgresDSN(
		s.User,
		s.Password,
		s.Host,
		s.Port,
		s.Database,
		sslModeString(s.SSLEnabled),
		nil,
	))
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

	dsn := buildPostgresDSN(
		s.User,
		s.Password,
		s.Host,
		s.Port,
		s.Database,
		sslModeString(s.SSLEnabled),
		map[string]string{
			"sshHost":     Host,
			"sshPort":     strconv.Itoa(int(Port)),
			"sshUser":     User,
			"sshDriverId": s.DriverId,
		},
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
