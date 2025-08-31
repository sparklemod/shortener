package config

type Config struct {
	LinksDB PostgresDataBase
}

type PostgresDataBase struct {
	URL string
}

func NewConfig() *Config {
	return &Config{
		LinksDB: PostgresDataBase{
			URL: "postgres://postgres:postgres@localhost:54302/postgres?sslmode=disable",
		},
	}
}
