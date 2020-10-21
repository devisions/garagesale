package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Register the postgres database driver.
)

// Config is what we require to open a database connection.
type Config struct {
	Host       string
	DBName     string
	Username   string
	Password   string
	DisableTLS bool
}

// Open knows how to open a database connection.
func Open(cfg Config) (*sqlx.DB, error) {

	q := url.Values{}
	q.Set("sslmode", "require")
	if cfg.DisableTLS {
		q.Set("sslmode", "disable")
	}
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.Username, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.DBName,
		RawQuery: q.Encode(),
	}
	return sqlx.Open("postgres", u.String())
}
