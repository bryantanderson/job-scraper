package database

import (
	"github.com/bryantanderson/go-job-assessor/internal/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DummyStore struct {
	db             *Database
	collectionName string
	databaseName   string
}

func InitializeDummyStore(db *Database) *DummyStore {
	return &DummyStore{
		db:             db,
		databaseName:   "Linkd",
		collectionName: "Dummies",
	}
}

func (s *DummyStore) Create(dummy *services.Dummy) error {
	_, err := s.getCollection().InsertOne(s.db.ctx, dummy)
	return err
}

func (s *DummyStore) FindById(id string) (*services.Dummy, error) {
	var d services.Dummy
	err := s.getCollection().FindOne(s.db.ctx, bson.M{"id": id}).Decode(&d)
	return &d, err
}

func (s *DummyStore) Update(dto *services.DummyDto, id string) (*services.Dummy, error) {
	var d services.Dummy
	err := s.getCollection().FindOneAndUpdate(s.db.ctx, bson.M{"id": id}, bson.M{"name": dto.Name}).Decode(&d)
	return &d, err
}

func (s *DummyStore) DeleteById(id string) error {
	res := s.getCollection().FindOneAndDelete(s.db.ctx, bson.M{"id": id})
	return res.Err()
}

func (s *DummyStore) getCollection() *mongo.Collection {
	return s.db.client.Database(s.db.databaseName).Collection(s.collectionName)
}
