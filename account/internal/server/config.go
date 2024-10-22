package server

import (
	"account/internal/storage"
	"log"

	"github.com/olivere/elastic/v7"
)

type Config struct {
	BindAddress    string `toml:"bind_address"`
	LogLevel       string `toml:"log_level"`
	LogHeaders     bool   `toml:"log_headers"`
	LogBody        bool   `toml:"log_body"`
	LogQueryParams bool   `toml:"log_query_params"`
	StorageConfig  *storage.Config
	ElasticClient  *elastic.Client
}

func NewConfig() *Config {
	client, err := elastic.NewClient(
		elastic.SetURL("http://elasticsearch:9200"),
		elastic.SetBasicAuth("elastic", "password"),
	)
	if err != nil {
		log.Fatal("Elasticsearch ElasticClient was not created")
	}

	return &Config{
		BindAddress:    ":8081",
		LogLevel:       "debug",
		LogHeaders:     false,
		LogBody:        true,
		LogQueryParams: true,
		StorageConfig:  storage.NewConfig(),
		ElasticClient:  client,
	}
}
