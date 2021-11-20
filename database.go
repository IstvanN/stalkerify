package main

import (
	"fmt"
	"os"

	"github.com/globalsign/mgo"
)

var (
	mongoURI  = os.Getenv("MONGODB_URI")
	dbName    = os.Getenv("MONGODB_DB_NAME")
	collName  = os.Getenv("MONGODB_COLL_NAME")
	dbSession *mgo.Session
)

// StartupMongo is the init call of the mongo DB, supposed to be called in the main function
func startupMongo() (*mgo.Session, error) {
	var err error
	dbSession, err = mgo.Dial(mongoURI)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to MongoDB: %v", err)
	}

	return dbSession, nil
}

func getMongoCollection() *mgo.Collection {
	return dbSession.DB(dbName).C(collName)
}
