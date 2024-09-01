package database

import (
	"encoding/json"
	"github.com/bryantanderson/go-job-assessor/internal/services"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AssessStore struct {
	db                     *Database
	elastic                *ElasticDatabase
	collectionName         string
	criteriaCollectionName string
}

func InitializeAssessStore(db *Database, elastic *ElasticDatabase) *AssessStore {
	return &AssessStore{
		db:                     db,
		elastic:                elastic,
		collectionName:         "assessments",
		criteriaCollectionName: "criteria",
	}
}

func (s *AssessStore) Create(a *services.Assessment) error {
	log.Infoln("Creating candidate assessment in database")
	elasticRes, err := s.elastic.CreateDocument(s.collectionName, a)

	if err != nil {
		log.Errorln(err)
		return err
	}

	a.ElasticId = elasticRes.Id_
	_, err = s.getCollection(s.collectionName).InsertOne(s.db.ctx, a)
	return err
}

func (s *AssessStore) CreateInternalJobCriteria(jc *services.Rubric) error {
	log.Infoln("Creating Linkd job criteria in database")
	_, err := s.getCollection(s.collectionName).InsertOne(s.db.ctx, jc)
	return err
}

func (s *AssessStore) QueryInternalJobCriteria(id string) (*services.Rubric, error) {
	var rubric services.Rubric
	err := s.getCollection(s.criteriaCollectionName).FindOne(s.db.ctx, bson.M{"id": id}).Decode(&rubric)
	return &rubric, err
}

func (s *AssessStore) FindById(assessmentId string) (*services.Assessment, error) {
	var a services.Assessment
	err := s.getCollection(s.collectionName).FindOne(s.db.ctx, bson.M{"_id": assessmentId}).Decode(&a)
	return &a, err
}

func (s *AssessStore) Query(params map[string]string) ([]*services.Assessment, error) {
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
	foundAssessments := make([]*services.Assessment, 0)

	for _, hit := range res.Hits.Hits {
		var ass services.Assessment
		json.Unmarshal(hit.Source_, &ass)
		foundAssessments = append(foundAssessments, &ass)
	}

	return foundAssessments, nil
}

func (s *AssessStore) Delete(userId string) error {
	// Get assessment from MongoDB as it contains the elastic ID
	assessment, err := s.FindById(userId)

	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	// Delete assessment from MongoDB
	mongoRes, err := s.getCollection(s.collectionName).DeleteOne(s.db.ctx, bson.M{"_id": assessment.Id})
	log.Infoln(mongoRes)

	if err != nil {
		log.Errorf("Failed to delete assessment from MongoDB: %s", err.Error())
		return err
	}

	// Delete assessment from ElasticSearch
	_, err = s.elastic.DeleteDocument(s.db.ctx, s.collectionName, assessment.ElasticId)
	return nil
}

func (s *AssessStore) getCollection(c string) *mongo.Collection {
	return getCollection(s.db, s.db.databaseName, c)
}
