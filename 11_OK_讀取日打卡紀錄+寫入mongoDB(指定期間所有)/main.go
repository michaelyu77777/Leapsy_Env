package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	//"labix.org/v2/mgo"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
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

	/**主功能*/
	//runtime.GOMAXPROCS(runtime.NumCPU())
	//ImportDailyRecord()

	/**測試開文字檔*/
	readTsFile("20170630.st")
}

func readTsFile(fileName string) {

	// 讀檔
	file, err := os.Open(fileName)
	if err != nil {

		log_err.WithFields(logrus.Fields{
			"trace":    "trace-0005",
			"err":      err,
			"fileName": fileName,
		}).Error("打開檔案失敗")

		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log_err.WithFields(logrus.Fields{
				"trace":    "trace-0006",
				"err":      err,
				"fileName": fileName,
			}).Error("關閉檔案失敗")

			log.Fatal(err)
		}
	}()

	// 讀檔
	scanner := bufio.NewScanner(file)

	// 一行一行讀
	counter := 0         // 行號
	for scanner.Scan() { // internally, it advances token based on sperator

		// 讀進一行(Big5)
		big5Name := scanner.Text()
		counter++

		//轉成utf8(繁體)
		utf8Name, _, _ := transform.String(big5ToUTF8Decoder, big5Name)

		//fmt.Println(utf8Name)
		//fmt.Println(scanner.Bytes()) // token in bytes

		// 若內容等於 空白" " ==0
		if strings.Compare(" ", utf8Name[139:140]) == 0 {

			fmt.Println("姓名空白, 行號:", counter)
			log_info.WithFields(logrus.Fields{
				"fileName": fileName,
				"trace":    "trace-0006",
				"行號":       counter,
				"日期":       utf8Name[27:37],
				"時間":       utf8Name[37:45],
				"整列內容":     utf8Name,
			}).Info("姓名空白列")

		} else if strings.Compare("ADMIN", utf8Name[139:144]) == 0 {
			//若內容等於ADMIN

			fmt.Println("ADMIN管理員, 行號:", counter)
			log_info.WithFields(logrus.Fields{
				"fileName": fileName,
				"trace":    "trace-0007",
				"行號":       counter,
				"日期":       utf8Name[27:37],
				"時間":       utf8Name[37:45],
				"整列內容":     utf8Name,
			}).Info("ADMIN列")

		} else {
			fmt.Println(utf8Name[15:27], utf8Name[27:37], utf8Name[37:45], utf8Name[139:144], utf8Name[144:153])
			log_info.WithFields(logrus.Fields{
				"fileName": fileName,
				"trace":    "trace-0008",
				"行號":       counter,
				"卡號":       utf8Name[15:27],
				"日期":       utf8Name[27:37],
				"時間":       utf8Name[37:45],
				"員工編號":     utf8Name[139:144],
				"姓名":       utf8Name[144:153],
			}).Info("打卡內容")
		}

	}

}

//Log檔
var log_info *logrus.Logger
var log_err *logrus.Logger

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

	/*json設定檔*/
	file, _ := os.Open("config.json")
	buf := make([]byte, 2048)

	//將設定讀到config變數中
	n, _ := file.Read(buf)
	fmt.Println(string(buf))
	err := json.Unmarshal(buf[:n], &config)
	if err != nil {
		panic(err)

		log_err.WithFields(logrus.Fields{
			"trace": "trace-0001",
			"err":   err,
		}).Error("將設定讀到config變數中失敗")

		fmt.Println(err)
	}
}

// 設定檔
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
	fmt.Println("檔案開始日期:", dateStart)
	log_info.Info("檔案開始日期: ", dateStart)

	/** 取得結束日期 **/
	stringEndDate := config.EndDate

	// 取出年月日
	endYear, _ := strconv.Atoi(stringEndDate[0:4])
	endMonth, _ := strconv.Atoi(stringEndDate[4:6])
	endDay, _ := strconv.Atoi(stringEndDate[6:8])

	// 轉成time格式
	dateEnd := time.Date(endYear, time.Month(endMonth), endDay, 0, 0, 0, 0, time.Local)
	fmt.Println("檔案結束日期:", dateEnd)
	log_info.Info("檔案結束日期: ", dateEnd)

	// for myTime := time.Date(2020, 1, 1, 9, 0, 0, 0, time.Local); myTime != time.Date(2021, 1, 1, 9, 0, 0, 0, time.Local); myTime = myTime.AddDate(0, 0, 1) {
	// }

	// 會包含dateEnd最後一天
	for myTime := dateStart; myTime != dateEnd.AddDate(0, 0, 1); myTime = myTime.AddDate(0, 0, 1) {

		// 檔案名稱(年月日)
		fileName := "Rec" + myTime.Format("20060102") + ".csv"
		log_info.Info("讀取檔名: ", fileName)
		fmt.Println("讀取檔名:", fileName)

		// 打開每日打卡紀錄檔案(本機要先登入過目的地磁碟機才能正常運作)
		//file, err := os.Open("Z:\\" + fileName)
		file, err := os.Open("\\\\leapsy-nas3\\CheckInRecord\\" + fileName)

		if err != nil {
			fmt.Println("打開檔案失敗", err)
			log_err.WithFields(logrus.Fields{
				"trace":    "trace-0002",
				"err":      err,
				"fileName": fileName,
			}).Error("打開檔案失敗")
			return
		}

		// 最後回收資源
		defer file.Close()

		fmt.Println("開始讀取檔案")
		log_info.Info("開始讀取檔案: ", fileName)

		// 讀檔
		reader := csv.NewReader(file)

		// 一行一行讀進來
		for {

			line, err := reader.Read()

			// 若讀到結束
			if err == io.EOF {

				fmt.Println(fileName, "此份檔案讀取完成")
				log_info.Info("此份檔案讀取完成: ", fileName)

				break
			} else if err != nil {

				close(chanDailyRecord)

				fmt.Println("關閉channel")
				fmt.Println("Error:", err)

				log_err.WithFields(logrus.Fields{
					"trace":    "trace-0003",
					"err":      err,
					"fileName": fileName,
				}).Error("關閉channel")

				break
			}

			// 處理Name編碼問題: 將繁體(Big5)轉成 UTF-8，儲存進去才正常
			log_info.Info("轉成UTF8配合繁體")
			big5Name := line[1]                                             // Name(Big5)
			utf8Name, _, _ := transform.String(big5ToUTF8Decoder, big5Name) // 轉成 UTF-8
			//fmt.Println(utf8Name) // 顯示"名字"

			// 建立每筆DailyRecord物件
			log_info.WithFields(logrus.Fields{
				"line[0]":  line[0],
				"utf8Name": utf8Name,
				"line[2]":  line[2],
				"line[3]":  line[3],
				"line[4]":  line[4],
				"line[5]":  line[5],
			}).Info("dailyrecord")

			dailyrecord := DailyRecord{line[0], utf8Name, line[2], line[3], line[4], line[5]}

			// 存到channel裡面
			chanDailyRecord <- dailyrecord
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
		fmt.Println("打卡紀錄插入錯誤(insertDailyRecord)")

		log_err.WithFields(logrus.Fields{
			"trace": "trace-0004",
			"err":   err,
		}).Error("打卡紀錄插入錯誤(insertDailyRecord)")

		panic(err)
		return
	}

	defer session.Close()

	log_info.Info("DB:leapsy_env, Collection:dailyRecord_real")
	c := session.DB("leapsy_env").C("dailyRecord_real")

	for dailyrecord := range chanDailyRecord {
		fmt.Println("插入一筆打卡資料：", dailyrecord)
		log_info.Info("插入一筆打卡資料:", dailyrecord)

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
