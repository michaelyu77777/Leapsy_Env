package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uri        = "mongodb://localhost:27017"
	database   = "leapsy_env"
	collection = "check_in_statistics"
)

var (
	err      error
	client   *mongo.Client
	item     Item
	duration time.Duration
	times    int64 = 10
	amount   int   = 100
	method   int
)

// CheckInStatistics struct
type CheckInStatistics struct {
	Value string
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	opts := options.Client().ApplyURI(uri)
	if client, err = mongo.Connect(ctx, opts); err != nil {
		log.Fatalln(err.Error())
	}

	c := client.Database(database).Collection(collection)

	fmt.Scan(&method)

	var i int64
	for i = 0; i < times; i++ {
		if method == 1 {
			upsert(ctx, c, amount)
			continue
		}
		if method == 2 {
			bulkUpsert(ctx, c, amount)
			continue
		}
		break
	}

	log.Printf("Average time: %s", duration/(time.Duration(times)*time.Millisecond)*time.Millisecond)
}

func upsert(ctx context.Context, c *mongo.Collection, amount int) {
	defer measure(time.Now())

	for i := 0; i <= amount; i++ {
		query := bson.M{"id": i}
		update := bson.M{"$set": Item{Value: "New Item " + strconv.Itoa(i)}}

		opts := options.Update().SetUpsert(true)
		_, err := c.UpdateOne(ctx, query, update, opts)

		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	// 刪除資料庫
	// if err := c.Drop(ctx); err != nil {
	// 	log.Fatalln(err.Error())
	// }
}

func bulkUpsert(ctx context.Context, c *mongo.Collection, amount int) {
	defer measure(time.Now())

	models := []mongo.WriteModel{}

	for i := 0; i <= amount; i++ {
		query := bson.M{"id": i}
		update := bson.M{"$set": Item{Value: "New Item " + strconv.Itoa(i)}}
		model := mongo.NewUpdateOneModel()
		models = append(models, model.SetFilter(query).SetUpdate(update).SetUpsert(true))
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := c.BulkWrite(ctx, models, opts)

	if err != nil {
		log.Fatalln(err.Error())
	}

	// 刪除資料庫
	// if err = c.Drop(ctx); err != nil {
	// 	log.Fatalln(err.Error())
	// }
}

func measure(start time.Time) {
	duration += time.Since(start)
	log.Printf("Execution time: %s", time.Since(start))
	log.Printf("Elapsed time: %s ", duration)
}
