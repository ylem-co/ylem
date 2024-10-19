package sql

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"ylem_integrations/helpers"
	"ylem_integrations/helpers/postgresql"
)

func init() {
	mysql.RegisterDialContext(helpers.MySqlSSHTCPNetName, helpers.SSHDialContextFunc)

	driver := &postgresql.Driver{}
	sql.Register("postgres+ssh", driver)

	postgresql.InitSshPool()
}
