package database

import "go.mongodb.org/mongo-driver/mongo"

func getCollection(db *Database, dn, cn string) *mongo.Collection {
	return db.client.Database(dn).Collection(cn)
}
