package postgres

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func Conn(config *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		url.QueryEscape(config.User),
		url.QueryEscape(config.Password),
		config.Host,
		config.Port,
		config.Name,
	))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
