package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"./objects"
	"github.com/sirupsen/logrus"
)

var log_info logrus.Logger
var log_err logrus.Logger
var config objects.Config

func main() {

	// 建立新路由
	router := mux.NewRouter()

	// 建立路由路徑與Handler
	router.HandleFunc(`/SouthSience/RtspConfig/query`, dailyAPIHandler)

	log_info.WithFields(logrus.Fields{}).Info("建立路由路徑")

	//讀入config
	readConfig()

	// 建立Server
	apiServerPointer := &http.Server{
		//Addr:           ":8007",
		Addr:           config.portofAPI,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log_info.WithFields(logrus.Fields{}).Info("建立Server")

	log.Fatal(apiServerPointer.ListenAndServe())

}

func dailyAPIHandler(w http.ResponseWriter, r *http.Request) {

	mongoClientPointer, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(`mongodb://localhost:27017`)) // 連接預設主機

	log_info.WithFields(logrus.Fields{}).Info("建立資料庫連線")

	if nil != err {
		fmt.Fprintf(w, "建立資料庫連線") // 寫入回應
		log_err.WithFields(logrus.Fields{"err": err}).Error("建立資料庫連線錯誤")
		return
	}

	content, err := mongoClientPointer.
		Database(`SouthernScience`).
		Collection(`rtsp_config`).
		Find(context.TODO(), bson.M{"EnableAudioStream": false}) //find all

	log_info.WithFields(logrus.Fields{}).Info("找colleciton全部")

	if nil != err {
		fmt.Fprintf(w, "查Collection錯誤") // 寫入回應
		log_err.WithFields(logrus.Fields{"err": err}).Error("查Collection錯誤")
		return
	}

	var results []objects.RtspConfig

	for content.Next(context.TODO()) { // 針對每一紀錄

		var rtspConfig objects.RtspConfig

		err = content.Decode(&rtspConfig) // 解析紀錄

		if nil != err {
			fmt.Fprintf(w, "解析錯誤") // 寫入回應
			log_err.WithFields(logrus.Fields{"err": err}).Error("解析錯誤")

			return
		}

		results = append(results, rtspConfig) // 儲存紀錄

	}

	jsonBytes, err := json.Marshal(results) // 轉成JSON

	if nil != err {
		fmt.Fprintf(w, "轉成JSON錯誤") // 寫入回應
		log_err.WithFields(logrus.Fields{"err": err}).Error("轉成JSON錯誤")
		return
	}

	fmt.Fprintf(w, "%s", string(jsonBytes)) // 寫入回應

}

// 讀設定檔
func readConfig() {

	File, err := os.Open("config.ini")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("成功開啟 config.ini")
	defer File.Close()

	byteValue, _ := ioutil.ReadAll(File)

	//將json設定轉存入變數
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		panic(err)

		log_err.WithFields(logrus.Fields{
			"trace": "trace-0001",
			"err":   err,
		}).Error("將設定讀到config變數中失敗")

		fmt.Println(err)
	}
}
