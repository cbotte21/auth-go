package datastore

import (
	"context"
	"github.com/cbotte21/games-auth/schema"
	"go.mongodb.org/mongo-driver/mongo"
)

type Datastore[T schema.Schema[any]] interface {
	Find(model T) (T, bool)
	Create(model T) bool
	Update(model T) T
	Delete(model T) bool
}

func Find[T schema.Schema[any]](schema T) (T, bool) {
	client, status := GetMongoClient()
	if !status {
		return schema, false
	}
	var result T
	collection := client.Database(schema.Database()).Collection(schema.Collection())
	err := collection.FindOne(context.TODO(), schema).Decode(&result)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			panic(err)
		}
		return schema, false
	}

	return result, true
}

func Create[T schema.Schema[any]](schema T) bool {
	client, status := GetMongoClient()
	if !status {
		return false
	}

	collection := client.Database(schema.Database()).Collection(schema.Collection())
	_, err := collection.InsertOne(context.TODO(), schema)

	return err == nil
}

func Update[X, Y schema.Schema[any]](filter X, updated Y) bool {
	client, status := GetMongoClient()
	if !status {
		return false
	}

	collection := client.Database(filter.Database()).Collection(filter.Collection())
	_, err := collection.UpdateOne(context.TODO(), filter, updated)

	return err == nil
}

func Delete[T schema.Schema[any]](schema T) bool {
	client, status := GetMongoClient()
	if !status {
		return false
	}

	collection := client.Database(schema.Database()).Collection(schema.Collection())
	_, err := collection.DeleteOne(context.TODO(), schema)

	return err == nil
}
