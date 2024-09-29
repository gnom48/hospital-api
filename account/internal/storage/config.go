package storage

type Config struct {
	DatabaseUrl string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{
		DatabaseUrl: "host=localhost port=5432 dbname=hospital sslmode=disable user=postgres password=630572",
	}
}
