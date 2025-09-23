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
			URL: "postgres://user:password@localhost:54302/shortener?sslmode=disable",
		},
	}
}
