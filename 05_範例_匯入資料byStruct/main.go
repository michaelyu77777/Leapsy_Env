package main

import (

	// Built-in Golang packages
	"context" // manage multiple requests
	"fmt"     // Println() function
	"os"
	"reflect" // get an object type
	"time"

	// Official 'mongo-go-driver' packages
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoFields struct {
	Field_Str  string
	Field_Int  int
	Field_Bool bool
}

func main() {

	// Declare host and port options to pass to the Connect() method
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	fmt.Println("clientOptions TYPE:", reflect.TypeOf(clientOptions), "\n")

	// Connect to the MongoDB and return Client instance
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// Access a MongoDB collection through a database
	col := client.Database("some_database").Collection("Some Collection")
	fmt.Println("Collection type:", reflect.TypeOf(col), "\n")

	// Declare a MongoDB struct instance for the document's fields and data
	oneDoc := MongoFields{
		Field_Str:  "Some Value",
		Field_Int:  12345,
		Field_Bool: true,
	}
	fmt.Println("oneDoc TYPE:", reflect.TypeOf(oneDoc), "\n")

	// InsertOne() method Returns mongo.InsertOneResult
	result, insertErr := col.InsertOne(ctx, oneDoc)
	if insertErr !=

		nil {
		fmt.Println("InsertOne ERROR:", insertErr)
		os.Exit(1) // safely exit script on error
	} else {
		fmt.Println("InsertOne() result type: ", reflect.TypeOf(result))
		fmt.Println("InsertOne() API result:", result)

		// get the inserted ID string
		newID := result.InsertedID
		fmt.Println("InsertOne() newID:", newID)
		fmt.Println("InsertOne() newID type:", reflect.TypeOf(newID))
	}
}
