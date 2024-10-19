package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"ylem_users/config"

	"github.com/go-redis/redis/v8"
)

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

func RedisDbConn(ctx context.Context) *redis.Client {
	c := config.Cfg()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.RedisDBConfig.Host, c.RedisDBConfig.Port),
		Password: c.RedisDBConfig.Password,
		DB:       0,
	}).WithContext(ctx)

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return rdb
}
