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

/*
 * 初始化配置
 */
func init() {
	//file, _ := os.Open("config.json")
	file, _ := os.Open("D:\\workspace-GO\\Leapsy_Env\\10_OK_讀取日打卡紀錄+寫入mongoDB(當日檔案)\\config.json")
	buf := make([]byte, 2048)

	n, _ := file.Read(buf)
	fmt.Println(string(buf))
	err := json.Unmarshal(buf[:n], &config)
	if err != nil {
		panic(err)
		fmt.Println(err)
	}
}

// ImportDailyRecord :主程式-每日打卡資料
func ImportDailyRecord() {

	// 移除當日所有舊紀錄
	deleteDailyRecordToday()

	// 建立 channel 存放 DailyRecord型態資料
	chanDailyRecord := make(chan DailyRecord)

	// 標記完成
	dones := make(chan struct{}, worker)

	// 將日打卡紀錄檔案內容讀出，並加到 chanDailyRecord 裡面
	go addDailyRecordToChannel(chanDailyRecord)

	// 將chanDailyRecord 插入mongodb資料庫
	for i := 0; i < worker; i++ {
		go insertDailyRecord(chanDailyRecord, dones)
	}
	//等待完成
	awaitForCloseResult(dones)
	fmt.Println("日打卡紀錄插入完畢")
}

/**
 * 刪除當日所有舊紀錄
 */

func deleteDailyRecordToday() {

	session, err := mgo.Dial(config.MongodbServer)
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		fmt.Println("錯誤")
		panic(err)
	}
	defer session.Close()
	c := session.DB(config.DBName).C(config.CollectionName)
	//c := session.DB("leapsy_env").C("dailyRecord_real")

	// Delete record
	currentTime := time.Now()           //取今天日
	t := currentTime.Format("20060102") //取年月日格式
	fmt.Println("移除資料日期為 date: ", t)

	info, err := c.RemoveAll(bson.M{"date": t}) //移除今天所有舊的紀錄(格式年月日)
	if err != nil {
		fmt.Printf("移除當日所有舊紀錄失敗 %v\n", err)
		os.Exit(1)
	}

	fmt.Println("ChangeInfo info: ", info)

}

/*
 * 讀取今日打卡資料 加入到channel中
 * 讀取的檔案().csv 或 .txt檔案)，編碼要為UTF-8，繁體中文才能正確被讀取
 */
func addDailyRecordToChannel(chanDailyRecord chan<- DailyRecord) {

	// 取得今日:年月日時
	currentTime := time.Now()

	//指定要抓的csv檔名
	fileName := ""

	//若現在是九點前:取昨日
	if currentTime.Hour() > 9 {
		fmt.Println("九點前:取昨日(hour=", currentTime.Hour())

		yesterday := currentTime.AddDate(0, 0, -1)
		t := yesterday.Format("20060102") //取年月日格式
		fileName = "Rec" + t + ".csv"
	} else {
		//取今日
		fmt.Println("九點後:取今日(hour=", currentTime.Hour())

		t := currentTime.Format("20060102") //取年月日格式
		fileName = "Rec" + t + ".csv"
	}

	//檔案名稱
	//fileName := "Rec" + year + month + day + ".csv"
	fmt.Println("日打卡紀錄檔名稱:", fileName)

	// 打開每日打卡紀錄檔案(不問帳號密碼?)
	// file, err := os.Open("Z:\\" + fileName)
	//file, err := os.Open("\\\\leapsy-nas3\\CheckInRecord\\" + fileName)
	file, err := os.Open(config.DailyRecordFileFolderPath + fileName)

	if err != nil {
		fmt.Println("打開文件失敗", err)
		return
	}

	// 最後回收資源
	defer file.Close()

	fmt.Println("讀取文件")

	// 讀檔
	reader := csv.NewReader(file)

	// 一行一行讀進來
	for {

		line, err := reader.Read()

		// 若讀到結束
		if err == io.EOF {

			close(chanDailyRecord)
			fmt.Println("文件讀取完成")
			break
		} else if err != nil {
			close(chanDailyRecord)
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

	session, err := mgo.Dial(config.MongodbServer)
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		fmt.Println("錯誤")
		panic(err)
		return
	}
	defer session.Close()
	c := session.DB(config.DBName).C(config.CollectionName)
	//c := session.DB("leapsy_env").C("dailyRecord_real")

	for dailyrecord := range chanDailyRecord {
		fmt.Println("插入：", dailyrecord)
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
