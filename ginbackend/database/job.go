package database

import (
	"sincidium/linkd/api/services"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobStore struct {
	db             *Database
	elastic        *ElasticDatabase
	collectionName string
}

func NewJobStore(db *Database, elastic *ElasticDatabase) *JobStore {
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

func (s *JobStore) getCollection() *mongo.Collection {
	return getCollection(s.db, s.db.databaseName, s.collectionName)
}
