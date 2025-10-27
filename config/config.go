package config

type DBType string

const (
	Mongo    DBType = "mongo"
	Postgres DBType = "postgres"
)

type Config struct {
	BaseURL string
	Host    string
	Port    int
	DBName  string
	LinksDB LinksDB
	DBType  DBType
}

type LinksDB struct {
	PostgresUrl string
	MongoUrl    string
}

func NewConfig() *Config {
	return &Config{
		BaseURL: "http://localhost:8080",
		Host:    "0.0.0.0",
		Port:    8080,
		DBName:  "shortener",
		LinksDB: LinksDB{
			PostgresUrl: "postgres://user:password@localhost:54302/shortener?sslmode=disable",
			MongoUrl:    "mongodb://user:password@localhost:27017/shortener?authSource=admin",
		},
		DBType: Mongo,
	}
}
