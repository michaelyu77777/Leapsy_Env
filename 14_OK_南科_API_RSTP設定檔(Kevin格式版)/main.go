package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	router := mux.NewRouter() // 新路由
	router.HandleFunc(`/SouthSience/RtspConfig/query`, dailyAPIHandler)

	apiServerPointer := &http.Server{
		Addr:           ":8007",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	} // 設定伺服器

	log.Fatal(apiServerPointer.ListenAndServe())

}

func dailyAPIHandler(w http.ResponseWriter, r *http.Request) {

	mongoClientPointer, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(`mongodb://localhost:27017`)) // 連接預設主機

	if nil != err {
		fmt.Fprintf(w, "no data") // 寫入回應
		return
	}

	cursor, err := mongoClientPointer.
		Database(`SouthernScience`).
		Collection(`rtsp_config`).
		Find(context.TODO(), bson.M{"EnableAudioStream": false}) //find all
	//Find(context.TODO(), bson.M{"time": bson.M{`$gt`: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local), `$lt`: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)}}) //時間要大於某時間 並且小於某時間

	// cursor, err := mongoClientPointer.
	// 	Database(`Leapsy-Environmental-Control-Database`).
	// 	Collection(`second-records`).
	// 	Find(context.TODO(), bson.M{"time": bson.M{`$gt`: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local), `$lt`: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)}}) //時間要大於某時間 並且小於某時間

	if nil != err {
		fmt.Fprintf(w, "err") // 寫入回應
		return
	}

	fmt.Println("我是Cursor", cursor)

	var results []RtspConfig
	//var results []Data

	for cursor.Next(context.TODO()) { // 針對每一紀錄

		var rtspConfig RtspConfig
		//var data Data

		err = cursor.Decode(&rtspConfig) // 解析紀錄

		if nil != err {
			fmt.Fprintf(w, "解析沒有資料no data") // 寫入回應
			fmt.Println("錯誤訊息:", err)
			return
		}

		results = append(results, rtspConfig) // 儲存紀錄

	}

	jsonBytes, err := json.Marshal(results) // 轉成JSON

	if nil != err {
		fmt.Fprintf(w, "no data") // 寫入回應
		return
	}

	fmt.Fprintf(w, "%s", string(jsonBytes)) // 寫入回應

}
