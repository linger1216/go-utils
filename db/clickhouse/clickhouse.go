package clickhouse

import (
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Url string `json:"url"`
}

func NewClickHouse(conf *Config) *sqlx.DB {
	db, err := sqlx.Open("clickhouse", conf.Url)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}
