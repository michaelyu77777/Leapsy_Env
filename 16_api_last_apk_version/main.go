package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Apk_Version_And_File_Path struct {
	Last_version  string `json:"last_version"`
	Download_path string `json:"download_path"`
}

func main() {
	router := mux.NewRouter() // 新路由
	router.HandleFunc(`/lastVersion/apk`, lastVersionApkAPIHandler)

	apiServerPointer := &http.Server{
		Addr:           ":5002",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	} // 設定伺服器

	log.Fatal(apiServerPointer.ListenAndServe())

}

func lastVersionApkAPIHandler(w http.ResponseWriter, r *http.Request) {

	mongoClientPointer, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(`mongodb://localhost:27017`)) // 連接預設主機

	if nil != err {
		fmt.Fprintf(w, "連 MongoDB Server錯誤") // 寫入回應
		return
	}

	cursor, err := mongoClientPointer.
		Database(`Leapsy_Launcher`).
		Collection(`Apk_Version_And_File_Path`).
		//Find(context.TODO(), bson.M{"time": bson.M{`$gt`: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local), `$lt`: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)}}) //時間要大於某時間 並且小於某時間
		Find(context.TODO(), bson.M{}) //時間要大於某時間 並且小於某時間

	if nil != err {
		fmt.Fprintf(w, "連 MongoDB.Collection 錯誤") // 寫入回應
		return
	}

	// 取得DB結果
	var results []Apk_Version_And_File_Path

	for cursor.Next(context.TODO()) { // 針對每一紀錄

		var data Apk_Version_And_File_Path

		err = cursor.Decode(&data) // 解析紀錄

		if nil != err {
			fmt.Fprintf(w, "no data") // 寫入回應
			return
		}

		results = append(results, data) // 儲存紀錄

	}

	jsonBytes, err := json.Marshal(results) // 轉成JSON

	if nil != err {
		fmt.Fprintf(w, "轉成JSON格式錯誤") // 寫入回應
		return
	}

	fmt.Fprintf(w, "%s", string(jsonBytes)) // 寫入回應

}
