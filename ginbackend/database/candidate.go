package database

import (
	"sincidium/linkd/api/services"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CandidateStore struct {
	db             *Database
	collectionName string
}

func NewCandidateStore(db *Database) *CandidateStore {
	return &CandidateStore{
		db: db,
		collectionName: "candidates",
	}
}

func (s *CandidateStore) Create(c *services.Candidate) error {
	res, err := s.getCollection().InsertOne(s.db.ctx, c)
	if err != nil {
		return err
	}
	log.Infof("Candidate successfully inserted with ID: %s", res.InsertedID)
	return nil
}

func (s *CandidateStore) Get(id string) (*services.Candidate, error) {
	var candidate services.Candidate
	err := s.getCollection().FindOne(s.db.ctx, bson.M{"_id": id}).Decode(&candidate)
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}
	return &candidate, nil
}

func (s *CandidateStore) Delete(id string) error {
	res := s.getCollection().FindOneAndDelete(s.db.ctx, bson.M{"_id": id})
	return res.Err()
}

func (s *CandidateStore) getCollection() *mongo.Collection {
	return s.db.client.Database(s.db.databaseName).Collection(s.collectionName)
}
