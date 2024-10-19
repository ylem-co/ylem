package helpers

import (
	"fmt"
	"strings"
	"database/sql"
	"ylem_pipelines/app/pipeline/run"
	"ylem_pipelines/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
)

const DB_TIME_TIMESTAMP = "2006-01-02 15:04:05"

func DbConn() *sql.DB {
	var config config.Config
	err := envconfig.Process("", &config)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
			config.DBConfig.User,
			config.DBConfig.Password,
			config.DBConfig.Host,
			config.DBConfig.Port,
			config.DBConfig.Name))

	if err != nil {
		panic(err)
	}

	return db
}

func NumRows(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		CheckDbErr(err)
	}
	return count
}

func CheckDbErr(err error) {
	if err != nil {
		panic(err)
	}
}

func IdListClause(field string, idList run.IdList) (string, []interface{}) {
	listType := idList.Type
	ids := idList.Ids

	clause := ``
	params := make([]interface{}, 0)
	if idList.Type == "" {
		return clause, params
	}

	clause = `AND ` + field
	if listType == run.IdListTypeEnabled {
		clause += ` IN `
	} else {
		clause += ` NOT IN `
	}

	if len(ids) == 0 {
		ids = []string{"-1"}
	}

	clause += `(?` + strings.Repeat(",?", len(ids)-1) + `)`

	for _, uuid := range ids {
		params = append(params, uuid)
	}

	return clause, params
}
