package helpers

import (
	"database/sql"
	"fmt"
	"ylem_integrations/config"
)

const DB_TIME_TIMESTAMP = "2006-01-02 15:04:05"

func DbConn() *sql.DB {
	config := config.Cfg()

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
