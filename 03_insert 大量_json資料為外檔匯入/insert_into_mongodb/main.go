package main

import (

	// Built-in Golang packages
	// manage multiple requests
	// Println() function
	// io.ReadFile
	"log"
	// get an object type
	// Import the JSON

	// encoding package

	// Official 'mongo-go-driver' packages
	"github.com/globalsign/mgo"

	//"go.mongodb.org/mongo-driver/bson"
	// "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoFields struct {
	//Key string `json:"key,omitempty"`
	ID        int    `json:"_id"`
	FieldStr  string `json:"Field_Str,omitempty"` //json欄位要長這樣
	FieldInt  int    `json:"Field_Int,omitempty"`
	FieldBool bool   `json:"Field_Bool,omitempty"`
}

func main() {

	/** 連資料庫 MongoDB
	 * 寫法:透過 mgo.Dial 撥號，return seesion **/

	// 撥號
	uri := "localhost:27017"
	session, err := mgo.Dial(uri)
	if err != nil {
		log.Fatal("Couldn't connect to db.", err)
	}

	// 等所有函數都執行回傳完後 最後關閉session
	defer session.Close()

	// 取得 DB collection
	collection := session.DB("leapsy_env").C("check_in_statistics")

	/*寫入json*/
	var bdoc interface{}

	// 插入一年的統計假資料
	for i := 0; i <= 360; i++ {

		// 統計數據一
		err = bson.UnmarshalJSON([]byte(`{	"id": 1,
											"date": "2020-01-01",
											"expected": 30,
											"attendance": 27,
											"not_arrived": 3, 
											"guests": 4 }`), &bdoc)
		if err != nil {
			panic(err)
		}

		err = collection.Insert(&bdoc)
		if err != nil {
			panic(err)
		}
	}
}
