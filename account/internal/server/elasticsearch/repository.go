package elasticsearch

import (
	"account/internal/storage"
	"context"
	"errors"
	"reflect"

	models "github.com/gnom48/hospital-api-lib"
	"github.com/olivere/elastic/v7"
)

type Repository struct {
	ElasticClient *elastic.Client
}

func (r *Repository) CreateElasticsearchIndexes() error {
	exists, err := r.ElasticClient.IndexExists(INDEX_USERS).Do(context.Background())
	if err != nil || !exists {
		_, err = r.ElasticClient.CreateIndex(INDEX_USERS).Do(context.Background())
	}
	return err
}

func (r *Repository) AddIndex(user *models.User) (*elastic.IndexResponse, error) {
	return r.ElasticClient.Index().
		Index(INDEX_USERS).
		Id(user.Id).
		BodyJson(CachedUser{
			Id:       user.Id,
			Username: user.Username,
			Password: user.Password,
		}).
		Do(context.Background())
}

func (r *Repository) GetUserInfoByLoginPasswordElasticsearch(username string, password string) (*CachedUser, error) {
	searchResult, err := r.ElasticClient.Search().
		Index(INDEX_USERS).
		Query(elastic.NewMatchQuery("username", username)).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	if searchResult.TotalHits() > 0 {
		for _, item := range searchResult.Each(reflect.TypeOf(CachedUser{})) {
			if user, ok := item.(CachedUser); ok {
				if user.Password == storage.EncryptString(password) {
					return &user, nil
				}
			}
		}
	}
	return nil, errors.New("Not found")
}

func (r *Repository) GetUserInfoByIdElasticsearch(userId string) (*CachedUser, error) {
	searchResult, err := r.ElasticClient.Search().
		Index(INDEX_USERS).
		Query(elastic.NewMatchQuery("user_id", userId)).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	if searchResult.TotalHits() > 0 {
		for _, item := range searchResult.Each(reflect.TypeOf(CachedUser{})) {
			if user, ok := item.(CachedUser); ok {
				return &user, nil
			}
		}
	}
	return nil, errors.New("Not found")
}
