package storage

type Config struct {
	DatabaseUrl string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{
		// DatabaseUrl: "host=localhost port=7432 dbname=hospital sslmode=disable user=postgres password=630572",
		DatabaseUrl: "host=db port=7432 dbname=postgres sslmode=disable user=postgres password=postgres",
	}
}
