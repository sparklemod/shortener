package config

type Config struct {
	BaseURL string
	Host    string
	Port    int
	LinksDB PostgresDataBase
}

type PostgresDataBase struct {
	URL string
}

func NewConfig() *Config {
	return &Config{
		BaseURL: "http://localhost:8080",
		Host:    "0.0.0.0",
		Port:    8080,
		LinksDB: PostgresDataBase{
			URL: "postgres://user:password@localhost:54302/shortener?sslmode=disable",
		},
	}
}
