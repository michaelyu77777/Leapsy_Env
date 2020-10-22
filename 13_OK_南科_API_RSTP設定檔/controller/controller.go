//package main
package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber"
	"gopkg.in/mgo.v2/bson"

	"my-rest-api/db"
	"my-rest-api/settings"
	//"my-rest-api/model" //只有在create(insert)或update 才需要import model
)

// const dbName = "leapsy_env"                                     //DB
// const collectionNameOfCheckInRecord = "check_in_record"         //Collection
// const collectionNameOfCheckInStatistics = "check_in_statistics" //Collection
// const collectionName = "persion"                                //Collection
// const port = 8081                                               //API port
// const port = 8000 //API port

// 建立GET POST 路徑
func NewPersonController() {

	fmt.Println("建立API URL")
	app := fiber.New()

	/*建立 checkInRecord 路徑*/
	app.Get("/SouthSience/RtspConfig/query", getRtspConfig)
	//app.Post("/person", createPerson)
	//app.Put("/person/:id", updatePerson)
	//app.Delete("/person/:id", deletePerson)

	/*建立範例 person 路徑*/
	// app.Get("/person/:id?", getPerson)
	// app.Post("/person", createPerson)
	// app.Put("/person/:id", updatePerson)
	// app.Delete("/person/:id", deletePerson)

	fmt.Println("建立API URL結束")

	app.Listen(settings.PortOfAPI)
}

/* 以下為 CheckInRecord 相關 functions */
// 取得指定日期<應到>人員資料
func getRtspConfig(c *fiber.Ctx) {

	// 取得 collection
	collection, err := db.GetMongoDbCollection(settings.DbName, settings.CollectionNameOfRtspConfig)

	// 若連線有誤
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	//var filter bson.M = bson.M{}

	// 若有給date
	// if c.Params("date") != "" {

	// 	/*按照date取出當筆資料*/

	// 	//取出date參數
	// 	myDate := c.Params("date")
	// 	fmt.Println("查詢日期=", myDate)

	// 	filter = bson.M{"date": myDate}
	// 	fmt.Println("filter=", filter) //filter 型態 Map[date:2020-01-01]

	// }

	//cur, err := collection.Find(context.Background(), filter)
	filter := bson.M{}
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var results []bson.M
	cur.All(context.Background(), &results)

	// 若查無資料
	if results == nil {
		c.SendStatus(404)
		return
	}

	fmt.Println(results[0])

	// 	var anotherResults []model.RtspConfig

	// 	for _,value:=range results {
	// 		var rtspConfig model.RtspConfig

	// 		rtspConfig.Audio=value.
	// anotherResults=append(anotherResults,)
	// 	}

	json, _ := json.Marshal(results)
	c.Send(json)
}

// 取得指定日期<實到>人員資料
// func getAttendanceOfCheckInStatistics(c *fiber.Ctx) {

// 	// 取得 collection
// 	collection, err := db.GetMongoDbCollection(settings.DbName, settings.CollectionNameOfCheckInRecord)

// 	// 若連線有誤
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var filter bson.M = bson.M{}

// 	// 若有給date
// 	if c.Params("date") != "" {

// 		/*按照date取出當筆資料*/

// 		//取出date參數
// 		myDate := c.Params("date")
// 		fmt.Println("查詢日期=", myDate)

// 		//bson.M{} 裡面所用的欄位名稱 必須使用mongoDb欄位名稱 而非struct的欄位名稱 (與JAVA相異)
// 		filter = bson.M{"date": myDate, "leave_type": ""} //應到:leave_type is NULL
// 		fmt.Println("filter=", filter)                    //filter 型態 Map[date:2020-01-01]

// 	}

// 	var results []bson.M
// 	cur, err := collection.Find(context.Background(), filter)
// 	defer cur.Close(context.Background())

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	cur.All(context.Background(), &results)

// 	// 若查無資料
// 	if results == nil {
// 		c.SendStatus(404)
// 		return
// 	}

// 	json, _ := json.Marshal(results)
// 	c.Send(json)
// }

// // 取得指定日期<未到>人員資料
// func getNotArrivedOfCheckInStatistics(c *fiber.Ctx) {

// 	// 取得 collection
// 	collection, err := db.GetMongoDbCollection(settings.DbName, settings.CollectionNameOfCheckInRecord)

// 	// 若連線有誤
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var filter bson.M = bson.M{}

// 	// 若有給date
// 	if c.Params("date") != "" {

// 		/*按照date取出當筆資料*/

// 		//取出date參數
// 		myDate := c.Params("date")
// 		fmt.Println("查詢日期=", myDate)

// 		//bson.M{} 裡面所用的欄位名稱 必須使用mongoDb欄位名稱 而非struct的欄位名稱 (與JAVA相異)
// 		filter = bson.M{"date": myDate, "leave_type": bson.M{"$ne": ""}} //應到:leave_type is NOT Equal NULL
// 		fmt.Println("filter=", filter)

// 	}

// 	var results []bson.M
// 	cur, err := collection.Find(context.Background(), filter)
// 	defer cur.Close(context.Background())

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	cur.All(context.Background(), &results)

// 	// 若查無資料
// 	if results == nil {
// 		c.SendStatus(404)
// 		return
// 	}

// 	json, _ := json.Marshal(results)
// 	c.Send(json)
// }

// /* 以下為 CheckInStatistics 相關 functions */
// // 取得指定日期統計資料
// func getCheckInStatistics(c *fiber.Ctx) {

// 	// 取得 collection
// 	collection, err := db.GetMongoDbCollection(settings.DbName, settings.CollectionNameOfCheckInStatistics)

// 	// 若連線有誤
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var filter bson.M = bson.M{}

// 	// 若有給date
// 	if c.Params("date") != "" {

// 		/*按照date取出當筆資料*/

// 		//取出date參數
// 		myDate := c.Params("date")
// 		fmt.Println("查詢日期=", myDate)

// 		filter = bson.M{"date": myDate}
// 		fmt.Println("filter=", filter) //filter 型態 Map[date:2020-01-01]

// 	}

// 	var results []bson.M
// 	cur, err := collection.Find(context.Background(), filter)
// 	defer cur.Close(context.Background())

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	cur.All(context.Background(), &results)

// 	// 若查無資料
// 	if results == nil {
// 		c.SendStatus(404)
// 		return
// 	}

// 	json, _ := json.Marshal(results)
// 	c.Send(json)
// }

/* 以下為範例 Person 相關 functions */
// func getPerson(c *fiber.Ctx) {
// 	collection, err := db.GetMongoDbCollection(dbName, collectionName)
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var filter bson.M = bson.M{}

// 	if c.Params("id") != "" {

// 		/* 按照_id來取出當筆資料*/
// 		id := c.Params("id")
// 		objID, _ := primitive.ObjectIDFromHex(id)
// 		filter = bson.M{"_id": objID}
// 	}

// 	var results []bson.M
// 	cur, err := collection.Find(context.Background(), filter)
// 	defer cur.Close(context.Background())

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	cur.All(context.Background(), &results)

// 	if results == nil {
// 		c.SendStatus(404)
// 		return
// 	}

// 	json, _ := json.Marshal(results)
// 	c.Send(json)
// }

// func createPerson(c *fiber.Ctx) {
// 	collection, err := db.GetMongoDbCollection(dbName, collectionName)
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var person model.Person
// 	json.Unmarshal([]byte(c.Body()), &person)

// 	res, err := collection.InsertOne(context.Background(), person)
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	response, _ := json.Marshal(res)
// 	c.Send(response)
// }

// func updatePerson(c *fiber.Ctx) {
// 	collection, err := db.GetMongoDbCollection(dbName, collectionName)
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}
// 	var person model.Person
// 	json.Unmarshal([]byte(c.Body()), &person)

// 	update := bson.M{
// 		"$set": person,
// 	}

// 	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
// 	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	response, _ := json.Marshal(res)
// 	c.Send(response)
// }

// func deletePerson(c *fiber.Ctx) {
// 	collection, err := db.GetMongoDbCollection(dbName, collectionName)

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
// 	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	jsonResponse, _ := json.Marshal(res)
// 	c.Send(jsonResponse)
// }
