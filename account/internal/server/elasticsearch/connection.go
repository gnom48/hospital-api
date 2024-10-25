package elasticsearch

import (
	"errors"
	"fmt"
	"log"

	"github.com/olivere/elastic/v7"
)

type ElasticsearchConnection struct {
	repository *Repository
}

const INDEX_USERS = "users"

type CachedUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (e *ElasticsearchConnection) configureElasticsearchClient(elasticUrl string) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(elasticUrl),
	)
	if err != nil {
		return nil, errors.New("Elasticsearch: ElasticClient was not created: " + err.Error())
	}
	version, err := client.ElasticsearchVersion(elasticUrl)
	if err != nil {
		return nil, errors.New("Elasticsearch: ElasticClient couldn't connect to server: " + err.Error())
	} else {
		fmt.Println("Elasticsearch version = " + version)
	}

	return client, nil
}

func (e *ElasticsearchConnection) Repository() *Repository {
	if e.repository == nil {
		if client, err := e.configureElasticsearchClient("http://elasticsearch:9200"); err != nil {
			log.Fatal(err)
		} else {
			e.repository = &Repository{
				ElasticClient: client,
			}
		}
	}
	return e.repository
}
