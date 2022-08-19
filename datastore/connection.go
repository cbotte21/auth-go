//This package is meant to be boilerplate for every service.
//This is a mongodb wrapper

package datastore

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const uri string = "mongodb+srv://root:Asdfasdf1@cluster0.btzdz.mongodb.net/?retryWrites=true&w=majority" //TODO: Pull from env variables

func GetMongoClient() (*mongo.Client, bool) {
	//connect
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	//error check
	if err != nil {
		return nil, false
	}

	//ping
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
		return nil, false
	}

	//return
	return client, true
}
