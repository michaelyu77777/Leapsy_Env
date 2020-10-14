package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"

	//"labix.org/v2/mgo"
	"gopkg.in/mgo.v2"

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

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	ImportDailyRecord()
}

/*
 * 初始化配置
 */
func init() {
	file, _ := os.Open("config.json")
	buf := make([]byte, 2048)

	n, _ := file.Read(buf)
	fmt.Println(string(buf))
	err := json.Unmarshal(buf[:n], &config)
	if err != nil {
		panic(err)
		fmt.Println(err)
	}
}

//配置
type Config struct {
	MongodbServer string
	StartDate     string // 讀檔開始日
	EndDate       string // 讀檔結束日
}

/*導入每日打卡資料*/
func ImportDailyRecord() {

	// 建立 channel 存放 DailyRecord型態資料
	chanDailyRecord := make(chan DailyRecord)

	// 標記完成
	dones := make(chan struct{}, worker)

	// 將日打卡紀錄檔案內容讀出，並加到 chanDailyRecord 裡面
	go addDailyRecordForManyDays(chanDailyRecord)

	// 將chanDailyRecord 插入mongodb資料庫
	for i := 0; i < worker; i++ {
		go insertDailyRecord(chanDailyRecord, dones)
	}
	//等待完成
	awaitForCloseResult(dones)
	fmt.Println("日打卡紀錄插入完畢")
}

/*
 * 讀取"多日"打卡資料
 * 讀取的檔案().csv 或 .txt檔案)，編碼要為UTF-8，繁體中文才能正確被讀取
 */
func addDailyRecordForManyDays(chanDailyRecord chan<- DailyRecord) {

	/** 取得開始日期 **/
	stringStartDate := config.StartDate
	// 取出年月日
	startYear, _ := strconv.Atoi(stringStartDate[0:4])
	startMonth, _ := strconv.Atoi(stringStartDate[4:6])
	startDay, _ := strconv.Atoi(stringStartDate[6:8])
	// 轉成time格式
	dateStart := time.Date(startYear, time.Month(startMonth), startDay, 0, 0, 0, 0, time.Local)

	/** 取得結束日期 **/
	stringEndDate := config.EndDate
	// 取出年月日
	endYear, _ := strconv.Atoi(stringEndDate[0:4])
	endMonth, _ := strconv.Atoi(stringEndDate[4:6])
	endDay, _ := strconv.Atoi(stringEndDate[6:8])
	// 轉成time格式
	dateEnd := time.Date(endYear, time.Month(endMonth), endDay, 0, 0, 0, 0, time.Local)

	fmt.Println("開始日期:", dateStart)
	fmt.Println("結束日期:", dateEnd)

	// for myTime := time.Date(2020, 1, 1, 9, 0, 0, 0, time.Local); myTime != time.Date(2021, 1, 1, 9, 0, 0, 0, time.Local); myTime = myTime.AddDate(0, 0, 1) {
	// }

	// 會包含dateEnd最後一天
	for myTime := dateStart; myTime != dateEnd.AddDate(0, 0, 1); myTime = myTime.AddDate(0, 0, 1) {

		// 檔案名稱(年月日)
		fileName := "Rec" + myTime.Format("20060102") + ".csv"
		fmt.Println("讀取檔名:", fileName)

		// 打開每日打卡紀錄檔案(本機要先登入過目的地磁碟機才能正常運作)
		//file, err := os.Open("Z:\\" + fileName)
		file, err := os.Open("\\\\leapsy-nas3\\CheckInRecord\\" + fileName)

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

				fmt.Println(fileName, "此份檔案讀取完成")
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

	close(chanDailyRecord) // 關閉儲存的channel
}

/*
 * 將所有日打卡紀錄，全部插入到 mongodb
 */
func insertDailyRecord(chanDailyRecord <-chan DailyRecord, dones chan<- struct{}) {
	//开启loop个协程

	session, err := mgo.Dial(config.MongodbServer)
	if err != nil {
		fmt.Println("錯誤")
		panic(err)
		return
	}
	defer session.Close()
	c := session.DB("leapsy_env").C("dailyRecord_real")

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
