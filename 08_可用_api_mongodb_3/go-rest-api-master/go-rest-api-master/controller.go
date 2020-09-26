package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const dbName = "leapsy_env"              //DB
const collectionName = "check_in_record" //Collection
const port = 8000                        //API port

// 建立GET POST 路徑
func newPersonController() {
	app := fiber.New()

	// app.Get("/person/:id?", getPerson)
	// app.Post("/person", createPerson)
	// app.Put("/person/:id", updatePerson)
	// app.Delete("/person/:id", deletePerson)

	app.Get("/checkInRecord/query/:date?", getCheckInRecord)
	app.Post("/person", createPerson)
	app.Put("/person/:id", updatePerson)
	app.Delete("/person/:id", deletePerson)

	app.Listen(port)
}

/* 以下為 CheckInRecord 相關 functions */
func getCheckInRecord(c *fiber.Ctx) {

	// 取得 collection
	collection, err := getMongoDbCollection(dbName, collectionName)

	// 若連線有誤
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var filter bson.M = bson.M{}

	// 若有給date
	if c.Params("date") != "" {

		/*按照date取出當筆資料*/

		//取出date參數
		myDate := c.Params("date")
		fmt.Println("查詢日期=", myDate)

		filter = bson.M{"date": myDate}
		fmt.Println("filter=", filter) //filter 型態 Map[date:2020-01-01]

	}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	cur.All(context.Background(), &results)

	// 若查無資料
	if results == nil {
		c.SendStatus(404)
		return
	}

	json, _ := json.Marshal(results)
	c.Send(json)
}

/* 以下為 Person 相關 functions */
func getPerson(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var filter bson.M = bson.M{}

	if c.Params("id") != "" {

		/* 按照_id來取出當筆資料*/
		id := c.Params("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	cur.All(context.Background(), &results)

	if results == nil {
		c.SendStatus(404)
		return
	}

	json, _ := json.Marshal(results)
	c.Send(json)
}

func createPerson(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var person Person
	json.Unmarshal([]byte(c.Body()), &person)

	res, err := collection.InsertOne(context.Background(), person)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	response, _ := json.Marshal(res)
	c.Send(response)
}

func updatePerson(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	var person Person
	json.Unmarshal([]byte(c.Body()), &person)

	update := bson.M{
		"$set": person,
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	response, _ := json.Marshal(res)
	c.Send(response)
}

func deletePerson(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(dbName, collectionName)

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	jsonResponse, _ := json.Marshal(res)
	c.Send(jsonResponse)
}
