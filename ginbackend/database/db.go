package database

import (
	"context"
	"sincidium/linkd/api/setup"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	client *mongo.Client
	ctx    context.Context
	cancel func()

	dsn          string
	databaseName string
}

func InitDB(settings *setup.ApplicationSettings) *Database {
	db := &Database{}
	db.dsn = settings.DatabaseUri
	db.databaseName = "Linkd"
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *Database) Open() {
	db.client, _ = mongo.Connect(db.ctx, options.Client().ApplyURI(db.dsn))
	err := db.Ping()
	if err != nil {
		log.Errorln(err)
		panic("unable to connect to MongoDB")
	} else {
		log.Infoln("connected to MongoDB")
	}
}

func (db *Database) Ping() error {
	if err := db.client.Ping(db.ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
}
