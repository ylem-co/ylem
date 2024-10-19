package sqlIntegrations

import (
	"context"
	"errors"
	"fmt"
	"database/sql"
	"encoding/json"
	"ylem_taskrunner/helpers"
	"ylem_taskrunner/services/bigquery"
	"ylem_taskrunner/services/es"

	messaging "github.com/ylem-co/shared-messaging"
)

type DefaultSQLIntegrationConnectionConfiguration struct {
	Host        string
	Port        uint16
	User        string
	Password    string
	Database    string
	ProjectId   string
	Credentials string
	SslEnabled  bool
	EsVersion   *uint8
}

type SQLIntegrationConnection interface {
	// Open connects to the database, returning either a db or error on opening the db handle
	// Important — it's just the pool of connections. Doesn't mean the connection itself is established
	// It's unclear yet, but the idea is you should not really use the handle directly
	// Rather than this interface.
	Open() error

	// Close closes the opened handle pool
	Close() error
}

type SQLDriverConnection interface {
	// Prepare is a proxy of sql.PrepareContext (adds our own context)
	Prepare(query string) (*sql.Stmt, error)

	// Exec is a proxy of sql.ExecContext (adds our own context)
	Exec(query string, args ...interface{}) (sql.Result, error)

	// Query is a proxy of sql.ExecQuery (adds our own context)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type QueryableConnection interface {
	Query(q string, args ...interface{}) ([]map[string]interface{}, []string, error)
}

type ViaSshConnection interface {
	// Open connects to the database via SSH tunnel, returning either a db or error on opening the db handle
	// Important — it's just the pool of connections. Doesn't mean the connection itself is established
	OpenSsh(Host string, Port uint16, User string) error
}

type TestableConnection interface {
	// Test tests if the opened pool handle can perform any queries
	Test() error
}

func CreateSQLIntegrationConnection(Type string, Config DefaultSQLIntegrationConnectionConfiguration) (SQLIntegrationConnection, error) {
	switch Type {
	case messaging.SQLIntegrationTypeMySQL,
		messaging.SQLIntegrationTypeGoogleCloudSQL,
		messaging.SQLIntegrationTypeMicrosoftAzureSQL,
		messaging.SQLIntegrationTypePlanetScale,
		messaging.SQLIntegrationTypeClickhouse,
		messaging.SQLIntegrationTypeAWSRDS:
		return &MySQLSQLIntegrationConnection{
			Host:     Config.Host,
			Port:     Config.Port,
			User:     Config.User,
			Password: Config.Password,
			Database: Config.Database,
		}, nil
	case messaging.SQLIntegrationTypePostgresql,
		messaging.SQLIntegrationTypeImmuta:
		return &PostgreSqlSQLIntegrationConnection{
			Host:       Config.Host,
			Port:       Config.Port,
			User:       Config.User,
			Password:   Config.Password,
			Database:   Config.Database,
			SSLEnabled: Config.SslEnabled,
		}, nil
	case messaging.SQLIntegrationTypeSnowflake:
		return &SnowflakeSQLIntegrationConnection{
			AccountId: Config.Host,
			User:      Config.User,
			Password:  Config.Password,
			Database:  Config.Database,
		}, nil
	case messaging.SQLIntegrationTypeGoogleBigQuery:
		return bigquery.NewConnection(context.Background(), Config.ProjectId, Config.Credentials)
	case messaging.SQLIntegrationTypeElasticsearch:
		return es.NewConnection(
			context.Background(),
			fmt.Sprintf("https://%s:%d", Config.Host, Config.Port),
			Config.User,
			Config.Password,
			Config.EsVersion,
		)
	case messaging.SQLIntegrationTypeRedshift:
		return &RedshiftSQLIntegrationConnection{
			Host:     Config.Host,
			Port:     Config.Port,
			User:     Config.User,
			Password: Config.Password,
			Database: Config.Database,
		}, nil
	default:
		return nil, errors.New("The type " + Type + " is not supported")
	}
}

func CollectDataFromSQLSQLIntegrationAsJSON(Conn SQLDriverConnection, Query string) ([]byte, []string, error) {
	rows, err := Conn.Query(Query)

	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	columnTypes, _ := rows.ColumnTypes()
	columnNames := make([]string, 0)
	for _, val := range columnTypes {
		columnNames = append(columnNames, val.Name())
	}

	finalRows, err := helpers.SQLToJSON(
		helpers.SqlToJSONParams{
			Rows: rows,
		},
	)

	if err != nil {
		return nil, nil, err
	}

	bytes, err := json.Marshal(finalRows)

	return bytes, columnNames, err
}

func CollectDataFromQueryableSQLIntegrationAsJSON(Conn QueryableConnection, Query string) ([]byte, []string, error) {
	result, columns, err := Conn.Query(Query)
	if err != nil {
		return nil, nil, err
	}

	bytes, err := json.Marshal(result)

	return bytes, columns, err
}
