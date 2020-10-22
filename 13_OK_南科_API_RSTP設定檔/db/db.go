//package main

package db

import (
	"context"
	"log"
	"my-rest-api/settings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//GetMongoDbConnection 取得 mongodb連線
func GetMongoDbConnection() (*mongo.Client, error) {

	//連線mongodb
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:"+settings.PortOfMongoDB))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

// 取得 Collection
func GetMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := GetMongoDbConnection()

	if err != nil {
		return nil, err
	}

	collection := client.Database(DbName).Collection(CollectionName)

	return collection, nil
}
