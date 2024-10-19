package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"ylem_integrations/helpers/postgresql"

	_ "github.com/lib/pq"
	"github.com/sethvargo/go-password/password"
)

type PostgreSqlIntegrationConnection struct {
	Host       string
	Port       uint16
	User       string
	Password   string
	Database   string
	DriverId   string
	SSLEnabled bool
	DB         *sql.DB
}

func (s *PostgreSqlIntegrationConnection) Open() error {
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

func (s *PostgreSqlIntegrationConnection) OpenSsh(Host string, Port uint16, User string) error {
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

func (s *PostgreSqlIntegrationConnection) Test() error {
	return s.DB.PingContext(context.Background())
}

func (s *PostgreSqlIntegrationConnection) Prepare(query string) (*sql.Stmt, error) {
	return s.DB.PrepareContext(context.Background(), query)
}

func (s *PostgreSqlIntegrationConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.ExecContext(context.Background(), query, args...)
}

func (s *PostgreSqlIntegrationConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.QueryContext(context.Background(), query, args...)
}

func (s *PostgreSqlIntegrationConnection) Close() error {
	if s.DriverId != "" {
		_ = postgresql.SshPool.Unregister(s.DriverId)
	}
	return s.DB.Close()
}

func (s *PostgreSqlIntegrationConnection) ShowDatabases() ([]string, error) {
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

func (s *PostgreSqlIntegrationConnection) ShowTables(db string) ([]string, error) {
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

func (s *PostgreSqlIntegrationConnection) DescribeTable(db string, table string) ([]string, error) {
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
