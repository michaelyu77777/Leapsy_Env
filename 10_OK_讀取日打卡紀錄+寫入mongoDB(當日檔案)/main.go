package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	//"labix.org/v2/mgo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs" //Log寫入設定
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"                    //寫log檔
	"golang.org/x/text/encoding/traditionalchinese" // 繁體中文編碼
	"golang.org/x/text/transform"
)

//設定檔
var config Config = Config{}
var worker = runtime.NumCPU()

// 指定編碼:將繁體Big5轉成UTF-8才會正確
var big5ToUTF8Decoder = traditionalchinese.Big5.NewDecoder()

// 日打卡紀錄檔
type DailyRecord struct {
	Date       string "date"
	Name       string "name"
	CardID     string "cardID"
	Time       string "time"
	Message    string "message"
	EmployeeID string "employeeID"
}

// 配置
type Config struct {
	MongodbServer             string
	DBName                    string
	CollectionName            string
	DailyRecordFileFolderPath string
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	ImportDailyRecord()

}

// init 制定LOG層級(自動呼叫?)
// func init() {
// 	//log輸出為json格式
// 	logrus.SetFormatter(&logrus.JSONFormatter{})
// 	//輸出設定為標準輸出(預設為stderr)
// 	logrus.SetOutput(os.Stdout)
// 	//設定要輸出的log等級
// 	logrus.SetLevel(logrus.DebugLevel)
// }

//Log檔
var log_info *logrus.Logger
var log_err *logrus.Logger

//var writer *rotatelogs.RotateLogs

/*
 * 初始化配置
 */
func init() {

	fmt.Println("執行init()初始化")

	/**設定LOG檔層級與輸出格式*/
	//使用Info層級
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
		//logrus.PanicLevel: writer, //若執行發生錯誤則會停止不進行下去
	}

	log_info = logrus.New()
	log_info.Hooks.Add(lfshook.NewHook(pathMap, &logrus.JSONFormatter{})) //Log檔綁訂相關設定

	fmt.Println("結束Info等級設定")

	//Error層級
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

	/**讀設定檔(config.json)*/
	//file, _ := os.Open("config.json")
	log_info.Info("打開config設定檔")
	file, err := os.Open("D:\\workspace-GO\\Leapsy_Env\\10_OK_讀取日打卡紀錄+寫入mongoDB(當日檔案)\\config.json")
	buf := make([]byte, 2048)
	if err != nil {
		log_err.WithFields(logrus.Fields{
			"trace": "trace-0001",
			"err":   err,
		}).Error("打開config錯誤")
	}

	n, err := file.Read(buf)
	fmt.Println(string(buf))
	if err != nil {
		log_err.WithFields(logrus.Fields{
			"trace": "trace-0002",
			"err":   err,
		}).Error("讀取config錯誤")
		panic(err)
		fmt.Println(err)
	}

	log_info.Info("轉換config成json")
	err = json.Unmarshal(buf[:n], &config)
	if err != nil {
		log_err.WithFields(logrus.Fields{
			"trace": "trace-0003",
			"err":   err,
		}).Error("轉換config成json發生錯誤")
		panic(err)
		fmt.Println(err)
	}
}

// ImportDailyRecord :主程式-每日打卡資料
func ImportDailyRecord() {

	//先算出要抓今日或昨日:年月日時
	currentTime := time.Now()

	//指定年月日
	date := ""

	//若現在是九點前:取昨日
	if currentTime.Hour() < 9 {
		log_info.Info("九點前:取昨日(hour=", currentTime.Hour())

		yesterday := currentTime.AddDate(0, 0, -1)
		date = yesterday.Format("20060102") //取年月日
	} else {
		//取今日
		log_info.Info("九點後:取今日(hour=", currentTime.Hour())

		date = currentTime.Format("20060102") //取年月日
	}

	//檔案名稱
	//fileName := "Rec" + year + month + day + ".csv"
	log_info.Info("取年月日:", date)

	// 移除當日所有舊紀錄
	deleteDailyRecordToday(date)

	// 建立 channel 存放 DailyRecord型態資料
	chanDailyRecord := make(chan DailyRecord)

	// 標記完成
	dones := make(chan struct{}, worker)

	// 將日打卡紀錄檔案內容讀出，並加到 chanDailyRecord 裡面
	go addDailyRecordToChannel(chanDailyRecord, date)

	// 將chanDailyRecord 插入mongodb資料庫
	for i := 0; i < worker; i++ {
		go insertDailyRecord(chanDailyRecord, dones)
	}
	//等待完成
	awaitForCloseResult(dones)
	log_info.Info("日打卡紀錄插入完畢")
}

/**
 * 刪除當日所有舊紀錄
 */

func deleteDailyRecordToday(date string) {

	log_info.Info("連接MongoDB")
	session, err := mgo.Dial(config.MongodbServer)
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log_err.WithFields(logrus.Fields{
			"trace": "trace-0004",
			"err":   err,
		}).Error("連接MongoDB發生錯誤(要刪除日打卡記錄時)")

		panic(err)
	}

	defer session.Close()
	c := session.DB(config.DBName).C(config.CollectionName)
	//c := session.DB("leapsy_env").C("dailyRecord_real")

	log_info.Info("移除當日所有舊紀錄,日期為 date: ", date)
	info, err := c.RemoveAll(bson.M{"date": date}) //移除今天所有舊的紀錄(格式年月日)
	if err != nil {
		log_err.WithFields(logrus.Fields{
			"trace": "trace-0005",
			"err":   err,
			"date":  date,
		}).Error("移除當日所有舊紀錄失敗")

		os.Exit(1)
	}

	log_info.Info("發生改變的info: ", info)

}

/*
 * 讀取今日打卡資料 加入到channel中
 * 讀取的檔案().csv 或 .txt檔案)，編碼要為UTF-8，繁體中文才能正確被讀取
 */
func addDailyRecordToChannel(chanDailyRecord chan<- DailyRecord, date string) {

	//指定要抓的csv檔名
	fileName := "Rec" + date + ".csv"

	log_info.Info("打開.csv文件", fileName)

	// 打開每日打卡紀錄檔案(windows上面登入過目的資料夾，才能運行)
	// file, err := os.Open("Z:\\" + fileName)
	// file, err := os.Open("\\\\leapsy-nas3\\CheckInRecord\\" + fileName)
	file, err := os.Open(config.DailyRecordFileFolderPath + fileName)

	if err != nil {
		log_err.WithFields(logrus.Fields{
			"trace":    "trace-0006",
			"err":      err,
			"date":     date,
			"fileName": fileName,
		}).Error("打開.csv文件失敗")

		return
	}

	// 最後回收資源
	defer file.Close()

	log_info.Info("讀取文件")

	// 讀檔
	reader := csv.NewReader(file)

	// 一行一行讀進來
	for {

		line, err := reader.Read()

		// 若讀到結束
		if err == io.EOF {

			close(chanDailyRecord)
			log_info.Info("csv文件讀取完成")
			break

		} else if err != nil {

			close(chanDailyRecord)

			log_err.WithFields(logrus.Fields{
				"trace":    "trace-0007",
				"err":      err,
				"date":     date,
				"fileName": fileName,
			}).Error("讀取csv文件失敗")

			fmt.Println("Error:", err)
			break
		}

		// 處理Name編碼問題: 將繁體(Big5)轉成 UTF-8，儲存進去才正常
		big5Name := line[1]                                             // Name(Big5)
		utf8Name, _, _ := transform.String(big5ToUTF8Decoder, big5Name) // 轉成 UTF-8
		//fmt.Println(utf8Name) // 顯示"名字"

		dailyrecord := DailyRecord{line[0], utf8Name, line[2], line[3], line[4], line[5]} // 建立每筆DailyRecord物件
		chanDailyRecord <- dailyrecord                                                    // 存到channel裡面
	}
}

/*
 * 將所有日打卡紀錄，全部插入到 mongodb
 */
func insertDailyRecord(chanDailyRecord <-chan DailyRecord, dones chan<- struct{}) {
	//开启loop个协程

	log_info.Info("連接MongoDB(插入mongodb時)")
	session, err := mgo.Dial(config.MongodbServer)
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log_err.WithFields(logrus.Fields{
			"trace": "trace-0008",
			"err":   err,
		}).Error("連接MongoDB失敗(插入mongodb時)")

		panic(err)
		return
	}

	defer session.Close()
	c := session.DB(config.DBName).C(config.CollectionName)
	//c := session.DB("leapsy_env").C("dailyRecord_real")

	for dailyrecord := range chanDailyRecord {
		log_info.Info("插入：", dailyrecord)
		c.Insert(&dailyrecord)
	}

	dones <- struct{}{}
}

// 等待結束
func awaitForCloseResult(dones <-chan struct{}) {
	for {
		<-dones
		worker--
		if worker <= 0 {
			return
		}
	}
}
