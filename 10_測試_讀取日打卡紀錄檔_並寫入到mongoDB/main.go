package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"

	//"labix.org/v2/mgo"
	"gopkg.in/mgo.v2"

	//"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

// var enc = traditionalchinese.Big5

/*
初始化配置
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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ImportPhoneInfo()
	// exampleReadGBK("Rec20201013.csv")
}

// func exampleReadGBK(filename string) {
// 	// Read UTF-8 from a GBK encoded file.
// 	f, err := os.Open(filename)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	r := transform.NewReader(f, enc.NewDecoder())

// 	// Read converted UTF-8 from `r` as needed.
// 	// As an example we'll read line-by-line showing what was read:
// 	sc := bufio.NewScanner(r)
// 	for sc.Scan() {
// 		b := sc.Bytes()
// 		fmt.Printf("Read line: %s\n", b)
// 	}
// 	if err = sc.Err(); err != nil {
// 		fmt.Println(err)
// 	}

// 	if err = f.Close(); err != nil {
// 		fmt.Println(err)
// 	}
// }

var config Config = Config{}

var worker = runtime.NumCPU()

//日打卡紀錄檔
type DailyRecord struct {
	Date       string "date"
	Name       string "name"
	CardID     string "cardID"
	Time       string "time"
	Message    string "message"
	EmployeeID string "employeeID"
}

//配置
type Config struct {
	MongodbServer   string
	DailyRecordFile string
}

/*导入手机地理信息*/
func ImportPhoneInfo() {
	var chanDailyRecord = make(chan DailyRecord)
	// 标记完成
	dones := make(chan struct{}, worker)

	//读取文件信息
	go addPhoneInfo(chanDailyRecord)
	//插入mongodb
	for i := 0; i < worker; i++ {
		go doPhoneInfo(chanDailyRecord, dones)
	}
	//等待完成
	awaitForCloseResult(dones)
	fmt.Println("插入完畢")
}

/*
 * 取得每日打卡資料
 * 讀取的檔案().csv 或 .txt檔案)，編碼要為UTF-8，繁體中文才能正確被讀取
 */

func addPhoneInfo(chanDailyRecord chan<- DailyRecord) {

	file, err := os.Open(config.DailyRecordFile)

	if err != nil {
		fmt.Println("打開文件失敗", err)
		return
	}
	defer file.Close()
	fmt.Println("讀取文件")
	reader := csv.NewReader(file)

	// 指定編碼:將繁體Big5轉成UTF-8才會正確
	big5ToUTF8Decoder := traditionalchinese.Big5.NewDecoder()

	for {
		line, err := reader.Read()

		if err == io.EOF {

			close(chanDailyRecord)
			fmt.Println("文件讀取完成")
			break
		} else if err != nil {
			close(chanDailyRecord)
			fmt.Println("Error:", err)
			break
		}

		// 將名字(Name) 繁體Big5轉成UTF-8 才能正常儲存
		big5Name := line[1] // 測試的 Big5 編碼
		utf8Name, _, _ := transform.String(big5ToUTF8Decoder, big5Name)
		//fmt.Println(utf8Name) // 顯示"名字"

		dailyrecord := DailyRecord{line[0], utf8Name, line[2], line[3], line[4], line[5]}
		chanDailyRecord <- dailyrecord
	}
}

/*
插入資料到mongodb
*/
func doPhoneInfo(chanDailyRecord <-chan DailyRecord, dones chan<- struct{}) {
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

func awaitForCloseResult(dones <-chan struct{}) {
	for {
		<-dones
		worker--
		if worker <= 0 {
			return
		}
	}
}

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// )

// func main() {

// 	file, err := os.Open("Rec20201013.csv")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	scanner := bufio.NewScanner(file)
// 	// 缺省的分隔函数是bufio.ScanLines,我们这里使用ScanWords。
// 	// 也可以定制一个SplitFunc类型的分隔函数
// 	scanner.Split(bufio.ScanWords)
// 	// scan下一个token.
// 	success := scanner.Scan()
// 	if success == false {
// 		// 出现错误或者EOF是返回Error
// 		err = scanner.Err()
// 		if err == nil {
// 			log.Println("Scan completed and reached EOF")
// 		} else {
// 			log.Fatal(err)
// 		}
// 	}
// 	// 得到数据，Bytes() 或者 Text()
// 	fmt.Println("First word found:", scanner.Text())
// 	// 再次调用scanner.Scan()发现下一个token
// }

// // ExampleScanner_emptyFinalToken return nil
// func ExampleScanner_emptyFinalToken() {
// 	// Comma-separated list; last entry is empty.
// 	// const input = "1,2,3,4,"
// 	// scanner := bufio.NewScanner(strings.NewReader(input))

// 	file, err := os.Open("Rec20201013.csv")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	scanner := bufio.NewScanner(file)

// 	// Define a split function that separates on commas.
// 	onComma := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
// 		for i := 0; i < len(data); i++ {
// 			if data[i] == ',' {
// 				return i + 1, data[:i], nil
// 			}
// 		}
// 		if !atEOF {
// 			return 0, nil, nil
// 		}
// 		// There is one final token to be delivered, which may be the empty string.
// 		// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
// 		// but does not trigger an error to be returned from Scan itself.
// 		return 0, data, bufio.ErrFinalToken
// 	}
// 	scanner.Split(onComma)
// 	// Scan.
// 	for scanner.Scan() {
// 		fmt.Printf("%q ", scanner.Text())
// 	}
// 	if err := scanner.Err(); err != nil {
// 		fmt.Fprintln(os.Stderr, "reading input:", err)
// 	}
// 	// Output: "1" "2" "3" "4" ""
// }
