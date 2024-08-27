package database

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	documentDelete "github.com/elastic/go-elasticsearch/v8/typedapi/core/delete"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/get"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	indexDelete "github.com/elastic/go-elasticsearch/v8/typedapi/indices/delete"
	log "github.com/sirupsen/logrus"
	"sincidium/linkd/api/setup"
	"strings"
	"time"
)

type DocumentDto struct {
	IndexName    string      `json:"indexName"`
	DocumentBody interface{} `json:"documentBody"`
}

type ElasticDatabase struct {
	client *elasticsearch.TypedClient
}

func InitElasticSearch(settings *setup.ApplicationSettings) *ElasticDatabase {
	elasticClient, err := elasticsearch.NewTypedClient(
		elasticsearch.Config{
			CloudID: settings.ElasticCloudId,
			APIKey:  settings.ElasticApiKey,
		},
	)
	if err != nil {
		panic(err)
	}
	elastic := &ElasticDatabase{
		client: elasticClient,
	}
	// Initialize the index
	indices := []string{"assessments", "jobs"}

	for _, indexName := range indices {
		res, err := elastic.createIndex(indexName)

		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				log.Infof("Index %s already exists\n", indexName)
			} else {
				log.Errorf("Response %v with error %s occurred while creating Elastic index\n", res, err.Error())
			}
		}
	}

	return elastic
}

func (s *ElasticDatabase) createIndex(name string) (*create.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	res, err := s.client.Indices.Create(name).Do(ctx)
	return res, err
}

func (s *ElasticDatabase) deleteIndex(indexName string) (*indexDelete.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	res, err := s.client.Indices.Delete(indexName).Do(ctx)
	return res, err
}

func (s *ElasticDatabase) CreateDocument(indexName string, document any) (*index.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	res, err := s.client.Index(indexName).Request(document).Do(ctx)
	return res, err
}

func (s *ElasticDatabase) GetDocument(ctx context.Context, indexName, id string) (*get.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	res, err := s.client.Get(indexName, id).Do(ctx)
	return res, err
}

func (s *ElasticDatabase) QueryDocument(ctx context.Context, indexName string, query *search.Request) (*search.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	res, err := s.client.Search().Index(indexName).Request(query).Do(ctx)
	return res, err
}

func (s *ElasticDatabase) DeleteDocument(ctx context.Context, indexName, id string) (*documentDelete.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	res, err := s.client.Delete(indexName, id).Do(ctx)
	return res, err
}
