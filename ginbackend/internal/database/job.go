package database

import (
	"encoding/json"

	"github.com/bryantanderson/go-job-assessor/internal/services"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobStore struct {
	db             *Database
	elastic        *ElasticDatabase
	collectionName string
}

func InitializeJobStore(db *Database, elastic *ElasticDatabase) *JobStore {
	return &JobStore{
		db:             db,
		elastic:        elastic,
		collectionName: "jobs",
	}
}

func (s *JobStore) Create(j *services.Job) error {
	elasticRes, err := s.elastic.CreateDocument(s.collectionName, j)

	if err != nil {
		log.Errorln(err)
		return err
	}

	j.ElasticId = elasticRes.Id_
	_, err = s.getCollection().InsertOne(s.db.ctx, j)
	return err
}

func (s *JobStore) Query(params map[string]string) ([]*services.Job, error) {
	query := &search.Request{
		Query: &types.Query{
			Match: make(map[string]types.MatchQuery),
		},
	}

	// Construct elastic search query from request query parameters
	for k, v := range params {
		query.Query.Match[k] = types.MatchQuery{
			Query: v,
		}
	}
	res, err := s.elastic.QueryDocument(s.db.ctx, s.collectionName, query)

	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}

	log.Infof("Elastic Search query succeeded in %d, with %d hits", res.Took, res.Hits.Total.Value)
	foundJobs := make([]*services.Job, 0)

	for _, hit := range res.Hits.Hits {
		var job services.Job
		json.Unmarshal(hit.Source_, &job)
		foundJobs = append(foundJobs, &job)
	}

	return foundJobs, nil
}

func (s *JobStore) getCollection() *mongo.Collection {
	return getCollection(s.db, s.db.databaseName, s.collectionName)
}
