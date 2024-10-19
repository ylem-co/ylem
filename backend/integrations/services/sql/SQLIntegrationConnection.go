package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"bytes"
	"strings"
    "encoding/json"
    "net/http"
	log "github.com/sirupsen/logrus"
	"ylem_integrations/entities"
	"ylem_integrations/services/bigquery"
	"ylem_integrations/services/es"
	"ylem_integrations/config"
)

type DefaultSQLIntegrationConnectionConfiguration struct {
	Host        string
	Port        uint16
	User        string
	Password    string
	Database    string
	SshHost     string
	SshPort     uint16
	SshUser     string
	ProjectId   *string
	EsVersion   *uint8
	Credentials string
	SslEnabled  bool
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

type ViaSshConnection interface {
	// Open connects to the database via SSH tunnel, returning either a db or error on opening the db handle
	// Important — it's just the pool of connections. Doesn't mean the connection itself is established
	OpenSsh(host string, port uint16, user string) error
}

type TestableConnection interface {
	// Test tests if the opened pool handle can perform any queries
	Test() error
}

type DescribableConnection interface {
	// ShowTables simply returns a list of tables
	ShowTables(db string) ([]string, error)

	// ShowDatabases simply returns a list of database
	ShowDatabases() ([]string, error)

	// DescribeTable returns list of columns of the table
	DescribeTable(db string, table string) ([]string, error)
}

func CreateSQLIntegrationConnection(Type string, Config DefaultSQLIntegrationConnectionConfiguration) (SQLIntegrationConnection, error) {
	switch Type {
	case entities.SQLIntegrationTypeMySQL,
		entities.SQLIntegrationTypeGoogleCloudSQL,
		entities.SQLIntegrationTypeMicrosoftAzureSQL,
		entities.SQLIntegrationTypePlanetScale,
		entities.SQLIntegrationTypeClickhouse,
		entities.SQLIntegrationTypeAWSRDS:
		return &MySQLIntegrationConnection{
			Host:     Config.Host,
			Port:     Config.Port,
			User:     Config.User,
			Password: Config.Password,
			Database: Config.Database,
		}, nil
	case entities.SQLIntegrationTypePostgresql,
		entities.SQLIntegrationTypeImmuta:
		return &PostgreSqlIntegrationConnection{
			Host:       Config.Host,
			Port:       Config.Port,
			User:       Config.User,
			Password:   Config.Password,
			Database:   Config.Database,
			SSLEnabled: Config.SslEnabled,
		}, nil
	case entities.SQLIntegrationTypeSnowflake:
		return &SnowflakeIntegrationConnection{
			AccountId: Config.Host,
			User:      Config.User,
			Password:  Config.Password,
			Database:  Config.Database,
		}, nil
	case entities.SQLIntegrationTypeGoogleBigQuery:
		return bigquery.NewConnection(context.Background(), Config.ProjectId, Config.Credentials)
	case entities.SQLIntegrationTypeElasticSearch:
		return es.NewConnection(
			context.Background(),
			fmt.Sprintf("https://%s:%d", Config.Host, Config.Port),
			Config.User,
			Config.Password,
			Config.EsVersion,
		)
	case entities.SQLIntegrationTypeRedshift:
		return &RedshiftIntegrationConnection{
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

func TestSQLIntegrationConnection(Type string, IsSshConnection bool, Config DefaultSQLIntegrationConnectionConfiguration) error {
	connection, err := CreateSQLIntegrationConnection(Type, Config)

	if err != nil {
		log.Error(err.Error())

		return err
	}

	if IsSshConnection {
		sshConn, ok := connection.(ViaSshConnection)
		if !ok {
			return fmt.Errorf("%s connection doesn't support SSH", Type)
		}
		err = sshConn.OpenSsh(Config.SshHost, Config.SshPort, Config.SshUser)
	} else {
		err = connection.Open()
	}

	if err != nil {
		return err
	}

	defer connection.Close()

	testableConn, ok := connection.(TestableConnection)
	if !ok {
		return fmt.Errorf("%s connection doesn't support testing", Type)
	}
	err2 := testableConn.Test()
	if err2 != nil {
		return err2
	}

	return nil
}

func UpdateSQLIntegrationConnection(organizationUuid string, isDataSourceCreated bool) bool {
    config := config.Cfg()

    url := strings.Replace(config.NetworkConfig.UpdateConnectionsUrl, "{uuid}", organizationUuid, -1);

    rp, _ := json.Marshal(map[string]bool{"is_data_source_created": isDataSourceCreated})
    var jsonStr = []byte(rp)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    if err != nil {
        return false
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        return true
    } else {
        return false
    }
}

