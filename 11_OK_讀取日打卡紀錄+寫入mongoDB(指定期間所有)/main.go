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

// 日打卡紀錄(.csv檔 有工號的)
type DailyRecord struct {
	Date       string `date`
	Name       string `name`
	CardID     string `cardID`
	Time       string `time`
	EmployeeID string `employeeID`
	//Message    string `msg`
	//Message    string `message`
}

// 打卡紀錄(.ts檔)
// type DailyRecordByTsFile struct {
// 	Date       string `date`       //日期
// 	Name       string `name`       //姓名
// 	CardID     string `cardID`     //卡號
// 	Time       string `time`       //時間
// 	EmployeeID string `employeeID` //員工編號
// 	Message    string `msg`        //進出訊息
// }

// 打卡紀錄(.ts檔) 按照ts檔順序
type DailyRecordByTsFile struct {
	CardID     string `cardID`     //卡號
	Date       string `date`       //日期
	Time       string `time`       //時間
	EmployeeID string `employeeID` //員工編號
	Name       string `name`       //姓名
	Message    string `msg`        //進出訊息
}

/** 初始化配置 */
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
	MongodbServerIP string //IP
	DBName          string
	Collection      string
	StartDate       string // 讀檔開始日
	EndDate         string // 讀檔結束日
}

//Log檔
var log_info *logrus.Logger
var log_err *logrus.Logger

//檔案的開始與結束日期(轉Time格式)
var dateStart time.Time
var dateEnd time.Time

func main() {

	/**轉換逗號+有員工編號+csv檔*/
	// countDateStartAndEnd()
	//runtime.GOMAXPROCS(runtime.NumCPU())
	//ImportDailyRecord()

	/**轉換空白區隔+ts檔*/
	//計算開始結束日期
	countDateStartAndEnd()
	runtime.GOMAXPROCS(runtime.NumCPU())
	//轉換開始結束日期
	ImportDailyRecordBy_TsFile()

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

/*導入每日打卡資料(.ts)*/
func ImportDailyRecordBy_TsFile() {

	// 建立 channel 存放 DailyRecord型態資料
	chanDailyRecordByTsFile := make(chan DailyRecordByTsFile)

	// 標記完成
	dones := make(chan struct{}, worker)

	// 將日打卡紀錄檔案內容讀出，並加到 chanDailyRecord 裡面
	go addDailyRecordForManyDays_TsFile(chanDailyRecordByTsFile)

	log_info.Info("抓出chanDailyRecordByTsFile: ", chanDailyRecordByTsFile)

	log_info.WithFields(logrus.Fields{
		"trace":                        "trace-00xx-.ts",
		"len(chanDailyRecordByTsFile)": len(chanDailyRecordByTsFile),
	}).Info("確認初始資料量")

	// 將chanDailyRecord 插入mongodb資料庫
	for i := 0; i < worker; i++ {
		go insertDailyRecord_TsFile(chanDailyRecordByTsFile, dones)
	}
	//等待完成
	awaitForCloseResult(dones)
	fmt.Println("日打卡紀錄(.ts)插入完畢")
}

/*
 * 讀取有逗號+員工編號的打卡資料(多日)
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

			//dailyrecord := DailyRecord{line[0], utf8Name, line[2], line[3], line[4], line[5]}
			dailyrecord := DailyRecord{line[0], utf8Name, line[2], line[3], line[4]} //拿掉message

			// 存到channel裡面
			chanDailyRecord <- dailyrecord
		}
	}

	close(chanDailyRecord) // 關閉儲存的channel
}

func addDailyRecordForManyDays_TsFile(chanDailyRecordByTsFile chan<- DailyRecordByTsFile) {

	//readTsFile("20170626.st")
	//readTsFile("20170630.st")

	// 會包含dateEnd最後一天
	for myTime := dateStart; myTime != dateEnd.AddDate(0, 0, 1); myTime = myTime.AddDate(0, 0, 1) {

		// 檔名(年月日).ts
		fileName := myTime.Format("20060102") + ".st"
		log_info.Info("讀檔: ", fileName)
		fmt.Println("讀檔:", fileName)

		// 年月資料夾路徑
		folderNameByYearMonth := myTime.Format("200601")

		// 判斷檔案是否存在
		_, err := os.Lstat("\\\\leapsy-nas3\\CheckInRecord\\20170605-20201011(st)\\" + folderNameByYearMonth + "\\" + fileName)
		//_, err := os.Lstat("\\\\leapsy-nas3\\CheckInRecord\\20170605-20201011(st)\\201706\\" + fileName)

		// 檔案不存在
		if err != nil {
			log_info.WithFields(logrus.Fields{
				"trace": "trace-00xx",
				"err":   err,
				//"fileName": "\\\\leapsy-nas3\\CheckInRecord\\20170605-20201011(st)\\201706\\" + fileName,
				"fileName": "\\\\leapsy-nas3\\CheckInRecord\\20170605-20201011(st)\\" + folderNameByYearMonth + "\\" + fileName,
			}).Info("檔案不存在")

		} else {
			//檔案若存在

			//file, err := os.Open("Z:\\" + fileName) 打開每日打卡紀錄檔案(本機要先登入過目的地磁碟機才能正常運作)
			//file, err := os.Open("\\\\leapsy-nas3\\CheckInRecord\\20170605-20201011(st)\\201706\\" + fileName)
			file, err := os.Open("\\\\leapsy-nas3\\CheckInRecord\\20170605-20201011(st)\\" + folderNameByYearMonth + "\\" + fileName)

			log_info.WithFields(logrus.Fields{
				"trace":    "trace-00xx",
				"err":      err,
				"fileName": "\\\\leapsy-nas3\\CheckInRecord\\20170605-20201011(st)\\201706\\" + fileName,
			}).Info("打開檔案")

			// 讀檔
			if err != nil {

				log_err.WithFields(logrus.Fields{
					"trace":    "trace-0005",
					"err":      err,
					"fileName": fileName,
				}).Error("打開檔案失敗")

				//log.Fatal(err)
			}

			// 最後回收資源
			defer func() {
				if err = file.Close(); err != nil {
					log_err.WithFields(logrus.Fields{
						"trace":    "trace-0006",
						"err":      err,
						"fileName": fileName,
					}).Error("關閉檔案失敗")

					//log.Fatal(err)
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

				// fmt.Println("fileName=", fileName, "行號=", counter, "big5Name=", big5Name)
				// log_info.Info("fileName=", fileName, "行號=", counter, "big5Name="+big5Name)
				// log_info.Info("big5Name[140:145]=", big5Name[140:145])
				// log_info.Info("big5Name[140:146]=", big5Name[140:146])
				// log_info.Info("big5Name[140:147]=", big5Name[140:147])
				// log_info.Info("big5Name[140:148]=", big5Name[140:148])
				// log_info.Info("big5Name[140:149]=", big5Name[140:149])
				// log_info.Info("big5Name[140:150]=", big5Name[140:150])

				// fmt.Println("fileName=", fileName, "行號=", counter, "utf8Name=", utf8Name)
				// log_info.Info("fileName=", fileName, "行號=", counter, "utf8Name="+utf8Name)
				// log_info.Info("utf8Name[139:140]=", utf8Name[139:140])
				// log_info.Info("utf8Name[139:141]=", utf8Name[139:141])
				// log_info.Info("utf8Name[139:142]=", utf8Name[139:142])
				// log_info.Info("utf8Name[139:143]=", utf8Name[139:143])
				// log_info.Info("utf8Name[139:144]=", utf8Name[139:144])
				// log_info.Info("utf8Name[139:145]=", utf8Name[139:145])
				// log_info.Info("utf8Name[139:146]=", utf8Name[139:146])
				// log_info.Info("utf8Name[139:147]=", utf8Name[139:147])
				// log_info.Info("utf8Name[139:148]=", utf8Name[139:148])
				// log_info.Info("utf8Name[139:149]=", utf8Name[139:149])
				// log_info.Info("utf8Name[139:150]=", utf8Name[139:150])
				// log_info.Info("utf8Name[139:151]=", utf8Name[139:151])

				// log_info.Info("檢查點[144:145]=", utf8Name[144:145])
				// if strings.Compare(" ", utf8Name[144:145]) == 0 {
				// 	log_info.Info("有檢測到空白:", utf8Name[144:145])

				// } else {
				// 	log_info.Info("沒有檢查到空白:", utf8Name[144:145])
				// }

				//fmt.Println(utf8Name[140:145]) //判斷ADMIN
				//fmt.Println(utf8Name[144:145]) //判斷空白
				//fmt.Println(utf8Name[144:153]) //判斷中文名

				// ADMIN(不入資料庫)
				if strings.Compare("ADMIN", utf8Name[140:145]) == 0 {
					fmt.Println("找到ADMIN:", utf8Name[140:145])

					log_info.WithFields(logrus.Fields{
						"fileName": fileName,
						"trace":    "trace-0008",
						"行號":       counter,
						"值":        utf8Name[140:145],
					}).Info("找到ADMIN")

				} else if strings.Compare(" ", utf8Name[144:145]) == 0 {
					// 空白(不入資料庫)
					fmt.Println("找到空白:", utf8Name[144:145])

					log_info.WithFields(logrus.Fields{
						"fileName": fileName,
						"trace":    "trace-0009",
						"行號":       counter,
						"值":        utf8Name[144:145],
					}).Info("找到空白")

				} else if strings.Compare("按密碼", utf8Name[58:67]) == 0 {
					// 人名(正常進出(按密碼))(入資料庫)
					fmt.Println("找到(按密碼):", utf8Name[15:27], utf8Name[27:37], utf8Name[37:45], utf8Name[139:144], utf8Name[58:67], utf8Name[45:68])

					log_info.WithFields(logrus.Fields{
						"fileName": fileName,
						"trace":    "trace-000x",
						"行號":       counter,
						"卡號":       utf8Name[15:27],
						"日期":       utf8Name[27:37],
						"時間":       utf8Name[37:45],
						"員工編號":     utf8Name[142:147],
						"姓名":       utf8Name[147:156],
						"進出訊息":     utf8Name[45:68],
					}).Info("找到(按密碼):")

					/**順序:
					日期
					姓名
					卡號
					時間
					員工逼號
					進出訊息*/
					dailyrecordbytsfile := DailyRecordByTsFile{
						utf8Name[15:27],
						utf8Name[27:37],
						utf8Name[37:45],
						utf8Name[147:156],
						utf8Name[147:156],
						utf8Name[45:68]}

					// 存到channel裡面
					chanDailyRecordByTsFile <- dailyrecordbytsfile

				} else {
					// 人名(正常進出 / 密碼錯誤)(入資料庫)
					fmt.Println("找到人名", utf8Name[15:27], utf8Name[27:37], utf8Name[37:45], utf8Name[139:144], utf8Name[144:153], utf8Name[45:57])

					log_info.WithFields(logrus.Fields{
						"fileName": fileName,
						"trace":    "trace-0010",
						"行號":       counter,
						"卡號":       utf8Name[15:27],
						"日期":       utf8Name[27:37],
						"時間":       utf8Name[37:45],
						"員工編號":     utf8Name[139:144],
						"姓名":       utf8Name[144:153],
						"進出訊息":     utf8Name[45:57],
					}).Info("找到人名")

					//順序:日期 姓名 卡號 時間 員工逼號 進出訊息
					dailyrecordbytsfile := DailyRecordByTsFile{
						utf8Name[15:27],
						utf8Name[27:37],
						utf8Name[37:45],
						utf8Name[139:144],
						utf8Name[144:153],
						utf8Name[45:57]}

					// 存到channel裡面
					chanDailyRecordByTsFile <- dailyrecordbytsfile

				}
			}
		}
	}

	close(chanDailyRecordByTsFile) // 關閉儲存的channel

	log_info.WithFields(logrus.Fields{
		"trace": "trace-00xx",
	}).Info("所有檔案讀取完成，已關閉儲存的channel")

}

//讀取單檔 TsFile
func readTsFiles(fileName string) {

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

		//fmt.Println(utf8Name[140:145]) //判斷ADMIN
		//fmt.Println(utf8Name[140:141]) //判斷空白
		//fmt.Println(utf8Name[144:153]) //判斷中文名

		// 排除 ADMIN / 排除空白/ 才取名字
		if strings.Compare("ADMIN", utf8Name[140:145]) == 0 {
			fmt.Println("找到ADMIN:", utf8Name[140:145])

			log_info.WithFields(logrus.Fields{
				"fileName": fileName,
				"trace":    "trace-0008",
				"行號":       counter,
				"值":        utf8Name[140:145],
			}).Info("找到ADMIN")

		} else if strings.Compare(" ", utf8Name[140:141]) == 0 {
			fmt.Println("找到空白:", utf8Name[140:141])

			log_info.WithFields(logrus.Fields{
				"fileName": fileName,
				"trace":    "trace-0009",
				"行號":       counter,
				"值":        utf8Name[140:141],
			}).Info("找到空白")

		} else {

			fmt.Println("找到人名", utf8Name[15:27], utf8Name[27:37], utf8Name[37:45], utf8Name[139:144], utf8Name[144:153])

			log_info.WithFields(logrus.Fields{
				"fileName": fileName,
				"trace":    "trace-0010",
				"行號":       counter,
				"卡號":       utf8Name[15:27],
				"日期":       utf8Name[27:37],
				"時間":       utf8Name[37:45],
				"員工編號":     utf8Name[139:144],
				"姓名":       utf8Name[144:153],
			}).Info("找到人名")

		}

	}

}

/**轉換開始結束日期格式 變成time.Time格式*/
func countDateStartAndEnd() {
	/** 取得開始日期 **/
	stringStartDate := config.StartDate

	// 取出年月日
	startYear, _ := strconv.Atoi(stringStartDate[0:4])
	startMonth, _ := strconv.Atoi(stringStartDate[4:6])
	startDay, _ := strconv.Atoi(stringStartDate[6:8])

	// 轉成time格式
	dateStart = time.Date(startYear, time.Month(startMonth), startDay, 0, 0, 0, 0, time.Local)
	fmt.Println("檔案開始日期:", dateStart)
	log_info.Info("檔案開始日期: ", dateStart)

	/** 取得結束日期 **/
	stringEndDate := config.EndDate

	// 取出年月日
	endYear, _ := strconv.Atoi(stringEndDate[0:4])
	endMonth, _ := strconv.Atoi(stringEndDate[4:6])
	endDay, _ := strconv.Atoi(stringEndDate[6:8])

	// 轉成time格式
	dateEnd = time.Date(endYear, time.Month(endMonth), endDay, 0, 0, 0, 0, time.Local)
	fmt.Println("檔案結束日期:", dateEnd)
	log_info.Info("檔案結束日期: ", dateEnd)
}

/*
 * 將所有日打卡紀錄，全部插入到 mongodb
 */
func insertDailyRecord(chanDailyRecord <-chan DailyRecord, dones chan<- struct{}) {
	//开启loop个协程

	session, err := mgo.Dial(config.MongodbServerIP)
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

	c := session.DB("leapsy_env").C("dailyRecord_real")

	for dailyrecord := range chanDailyRecord {
		fmt.Println("插入一筆打卡資料：", dailyrecord)
		log_info.Info("插入一筆打卡資料:", dailyrecord)

		c.Insert(&dailyrecord)
	}

	dones <- struct{}{}
}

/*
 * 將所有日打卡紀錄，全部插入到 mongodb
 */
func insertDailyRecord_TsFile(chanDailyRecordByTsFile <-chan DailyRecordByTsFile, dones chan<- struct{}) {
	//开启loop个协程

	log_info.Info("開始插入MONGODB")

	session, err := mgo.Dial(config.MongodbServerIP)
	if err != nil {
		fmt.Println("打卡紀錄插入錯誤(insertDailyRecord_TsFile)")

		log_err.WithFields(logrus.Fields{
			"trace": "trace-00xx-.ts",
			"err":   err,
		}).Error("打卡紀錄插入錯誤(insertDailyRecord_TsFile)")

		panic(err)
		return
	}

	defer session.Close()

	c := session.DB(config.DBName).C(config.Collection)
	log_info.Info("連上DBName:", config.DBName, "Collection", config.Collection)

	//確認資料筆數
	// ch := make(chan int, 100)
	// for i := 0; i < 34; i++ {
	// 	ch <- 0
	// }
	// fmt.Println("資料量:", len(ch))

	log_info.WithFields(logrus.Fields{
		"trace":                        "trace-00xx-.ts",
		"len(chanDailyRecordByTsFile)": len(chanDailyRecordByTsFile),
	}).Info("確認資料量")

	for dailyrecord := range chanDailyRecordByTsFile {
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
