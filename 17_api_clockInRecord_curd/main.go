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
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

//打卡紀錄
type Record struct {

	//試試把json:""去掉
	CardID     string    `json:"cardID"`
	Date       string    `json:"date"`
	Time       string    `json:"time"`
	EmployeeID string    `json:"employeeID"`
	Name       string    `json:"name"`
	Msg        string    `json:"msg"`
	DateTime   time.Time `json:"dateTime"`
}

//設定檔
type Config struct {
	ServerIP          string
	ServerPort        string
	MongodbServerIP   string
	MongodbServerPort string
	DBName            string
	Collection        string
}

var log_info *logrus.Logger //Log
var log_err *logrus.Logger  //Log
var config Config           //config設定檔

/** 初始化配置 */
func init() {

	fmt.Println("執行init()初始化")

	//設定Log
	setLog()

	//將設定config
	readConfig()
}

func setLog() {

	fmt.Println("setLog()初始化")

	/**設定LOG檔層級與輸出格式*/
	// 使用Info層級
	path := "./log/info/info"
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H",                            // 檔名格式
		rotatelogs.WithLinkName(path),               // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(10080*time.Minute),    // 文件最大保存時間(保留七天)
		rotatelogs.WithRotationTime(60*time.Minute), // 日誌切割時間間隔(一小時存一個檔案)
	)

	// 設定LOG等級
	pathMap := lfshook.WriterMap{
		logrus.InfoLevel: writer,
		// logrus.PanicLevel: writer, //若執行發生錯誤則會停止不進行下去
	}

	log_info = logrus.New()                                               // 初始化log_info
	log_info.Hooks.Add(lfshook.NewHook(pathMap, &logrus.JSONFormatter{})) // Log檔綁訂相關設定

	fmt.Println("結束Info等級設定")

	// Error層級
	path = "./log/err/err"
	writer, _ = rotatelogs.New(
		path+".%Y%m%d%H",                            // 檔名格式
		rotatelogs.WithLinkName(path),               // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(10080*time.Minute),    // 文件最大保存時間(保留七天)
		rotatelogs.WithRotationTime(60*time.Minute), // 日誌切割時間間隔(一小時存一個檔案)
	)

	// 設定LOG等級
	pathMap = lfshook.WriterMap{
		//logrus.InfoLevel: writer,
		logrus.ErrorLevel: writer,
		//logrus.PanicLevel: writer, //若執行發生錯誤則會停止不進行下去
	}

	log_err = logrus.New()
	log_err.Hooks.Add(lfshook.NewHook(pathMap, &logrus.JSONFormatter{})) //Log檔綁訂相關設定

	fmt.Println("結束Error等級設定")
	log_info.Info("結束Error等級設定")
}

// 讀設定檔
func readConfig() {

	File, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("成功開啟 config.json")
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

	fmt.Println("ServerIP: " + config.ServerIP)
	fmt.Println("ServerPort: " + config.ServerPort)
	fmt.Println("MongodbServerPort: " + config.MongodbServerPort)
	fmt.Println("DB: " + config.DBName)
	fmt.Println("Collection: " + config.Collection)

}

func main() {

	router := mux.NewRouter() // 新路由
	router.HandleFunc(`/ClockInRecord/daily/all`, ClockInRecordDilyAllAPIHandler)

	apiServerPointer := &http.Server{
		//若需要再看怎麼加入config.serverIP使用
		Addr:           ":" + config.ServerPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	} // 設定伺服器

	log.Fatal(apiServerPointer.ListenAndServe())

}

// API
func ClockInRecordDilyAllAPIHandler(w http.ResponseWriter, r *http.Request) {

	mongoClientPointer, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(`mongodb://localhost:`+config.MongodbServerPort)) // 連接預設主機

	if nil != err {
		fmt.Fprintf(w, "連 MongoDB Server錯誤") // 寫入回應
		return
	}

	// map:接收前端參數
	m := bson.M{}

	//若沒有輸入此變數，則map不加入搜尋欄位
	if `` != r.FormValue("cardID") {
		m[`cardID`] = r.FormValue("cardID")
	}

	if `` != r.FormValue("employeeID") {
		m[`employeeID`] = r.FormValue("employeeID")
	}

	if `` != r.FormValue("name") {
		m[`name`] = r.FormValue("name")
	}

	if `` != r.FormValue("msg") {
		m[`msg`] = r.FormValue("msg")
	}

	//dateTime
	if `` != r.FormValue("dateTime") {

		// 輸入String格式為 2006-01-02T15:04:05+08:00 轉成time.Time (前端可自行指定時區)
		dateTime, err := time.Parse(time.RFC3339, r.FormValue("dateTime"))

		//處理格式錯誤
		if nil != err {

			fmt.Println(err)
			log_err.WithFields(logrus.Fields{
				"輸入dateTime值": r.FormValue("dateTime"),
				"err":         err,
			}).Error("輸入日期格式:轉換錯誤")

			//date= 0001-01-01 00:00:00 +0000 UTC
			fmt.Println("輸入日期格式錯誤:date=", dateTime)
		}

		m[`dateTime`] = dateTime
	}

	// 查詢
	cursor, err := mongoClientPointer.
		Database(config.DBName).
		Collection(config.Collection).
		//Find(context.TODO(), bson.M{"dateTime": bson.M{`$lt`: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}}) //時間要大於某時間 並且小於某時間
		//Find(context.TODO(), bson.M{"time": bson.M{`$gt`: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local), `$lt`: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)}}) //時間要大於某時間 並且小於某時間
		Find(context.TODO(), m)

	// 查詢錯誤 return
	if nil != err {
		fmt.Fprintf(w, "查詢出錯", err, ",w=", w)
		if nil != err {
			log_err.WithFields(logrus.Fields{
				"err": err,
			}).Error("查詢出錯")
			fmt.Println("查詢出錯", err)
		}
		return
	}

	// 所有 record 結果
	var results []Record

	// 一筆筆拿出結果
	for cursor.Next(context.TODO()) { // 針對每一紀錄

		// 單筆 record
		var record Record

		// decode
		err = cursor.Decode(&record)

		// 錯誤處理
		if nil != err {
			//fmt.Printf("解析Record錯誤", err, ",w=", w) // 寫入回應
			fmt.Fprintf(w, "w", err, "cursor.Decode發生錯誤")
			log_err.WithFields(logrus.Fields{
				"err":      err,
				"data":     record,
				"dateTime": record.DateTime,
			}).Error("cursor.Decode發生錯誤")

			return
		}

		// 查詢結果:dateTime時區 轉成Local時區
		record.DateTime = record.DateTime.Local()
		// fmt.Println(record.DateTime)

		// decode 無法檢查到的部分
		// 檢查dateTime年(應該介於 0~9999之間)
		// 針對負年而做的檢查　dateTime:-0001-11-29T16:00:00.000+00:00
		if record.DateTime.Year() > 0 && record.DateTime.Year() < 9999 {
			// 合法的年

			// 儲入此筆record
			results = append(results, record)

		} else {
			// 非法年份

			// 不存此筆record
			fmt.Println("年份不在0~9999之間 若轉json會有問題")

			log_info.WithFields(logrus.Fields{
				"name": record.Name,
			}).Info("年份不在0~9999之間 若轉json會有問題")
		}

	}

	// 所有record結果 轉JSON
	jsonBytes, err := json.Marshal(results)

	if nil != err {

		fmt.Fprintf(w, "轉成JSON格式時發生錯誤 e=", err, ",w=", w) // 寫入回應

		log_err.WithFields(logrus.Fields{
			"err": err,
		}).Error("轉成JSON格式時發生錯誤:")

		return
	}

	fmt.Fprintf(w, "%s", string(jsonBytes)) // 寫入回應

}
