package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/globalsign/mgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type CheckInRecord struct {
	//Key string `json:"key,omitempty"`
	Id   int    //`json:"id"`
	Name string //注意:struct名稱開頭必須要大寫...否則無法寫入mongoDB!!!不知道為什麼...
	//Check_in_time time.Time
	Check_in_time string
	Pic           string
	Leave_type    string
	Date          string
	Department    string
	Position      string
}

func main() {
	//插入一年假資料A
	//insertCheckInStatisticsOneYear()

	//插入一年假資料B
	insertCheckInRecordOneYear()

	//fmt.Println(ReadFile("1.txt"))
}

// ReadFile 讀檔 return byte[]
func ReadFile(fileName string) string {
	f, err := ioutil.ReadFile("base64/" + fileName)
	//fmt.Println("測試ReadFile ERROR:", err)

	if err != nil {
		fmt.Println("read fail", err)
	}
	return string(f)
}

// 插入一整年假資料B:寫法二
func insertCheckInRecordOneYear() {
	// Declare host and port options to pass to the Connect() method
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// fmt.Println("clientOptions TYPE:", reflect.TypeOf(clientOptions), "\n")

	// Connect to the MongoDB and return Client instance
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// Access a MongoDB collection through a database
	col := client.Database("leapsy_env").Collection("check_in_record")
	//fmt.Println("Collection type:", reflect.TypeOf(col), "\n")

	index := 0 //id序號

	// 插入一整年的統計假資料
	for myTime := time.Date(2020, 1, 1, 9, 0, 0, 0, time.Local); myTime != time.Date(2021, 1, 1, 9, 0, 0, 0, time.Local); myTime = myTime.AddDate(0, 0, 1) {

		myCheckInTime := fmt.Sprintf(`%02d:%02d:%02d`, myTime.Hour(), myTime.Minute(), myTime.Second())
		myDate := fmt.Sprintf(`%04d-%02d-%02d`, myTime.Year(), myTime.Month(), myTime.Day())

		// 實體化struct
		//Person1
		fileName := "1.txt"
		picBase64Content := ReadFile(fileName)
		index++
		DocPerson := CheckInRecord{
			Id:            index,
			Name:          "micha",
			Check_in_time: myCheckInTime,
			Pic:           picBase64Content,
			Leave_type:    "",
			Date:          myDate,
			Department:    "軟體",
			Position:      "軟體工程師",
		}
		//fmt.Println("oneDoc TYPE:", reflect.TypeOf(oneDoc), "\n")
		// 插入一筆資料到資料庫 InsertOne() method 回傳 Returns mongo.InsertOneResult
		result, insertErr := col.InsertOne(ctx, DocPerson)

		//Person2
		fileName = "2.txt"
		picBase64Content = ReadFile(fileName)
		index++
		DocPerson = CheckInRecord{
			Id:            index,
			Name:          "ken",
			Check_in_time: myCheckInTime,
			Pic:           picBase64Content,
			Leave_type:    "",
			Date:          myDate,
			Department:    "軟體",
			Position:      "軟體工程師",
		}
		result, insertErr = col.InsertOne(ctx, DocPerson)

		//Person3
		fileName = "3.txt"
		picBase64Content = ReadFile(fileName)
		index++
		DocPerson = CheckInRecord{
			Id:            index,
			Name:          "p3",
			Check_in_time: myCheckInTime,
			Pic:           picBase64Content,
			Leave_type:    "",
			Date:          myDate,
			Department:    "視覺設計",
			Position:      "視覺設計師",
		}
		result, insertErr = col.InsertOne(ctx, DocPerson)

		//Person4
		fileName = "4.txt"
		picBase64Content = ReadFile(fileName)
		index++
		DocPerson = CheckInRecord{
			Id:            index,
			Name:          "p4",
			Check_in_time: myCheckInTime,
			Pic:           picBase64Content,
			Leave_type:    "",
			Date:          myDate,
			Department:    "視覺設計",
			Position:      "視覺設計師",
		}
		result, insertErr = col.InsertOne(ctx, DocPerson)

		//Person5
		fileName = "5.txt"
		picBase64Content = ReadFile(fileName)
		index++
		DocPerson = CheckInRecord{
			Id:            index,
			Name:          "p5",
			Check_in_time: myCheckInTime,
			Pic:           picBase64Content,
			Leave_type:    "",
			Date:          myDate,
			Department:    "視覺設計",
			Position:      "視覺設計師",
		}
		result, insertErr = col.InsertOne(ctx, DocPerson)

		//Person6
		fileName = "6.txt"
		picBase64Content = ReadFile(fileName)
		index++
		DocPerson = CheckInRecord{
			Id:            index,
			Name:          "p6",
			Check_in_time: "",
			Pic:           picBase64Content,
			Leave_type:    "病",
			Date:          myDate,
			Department:    "軟體",
			Position:      "軟體工程師",
		}
		result, insertErr = col.InsertOne(ctx, DocPerson)

		//Person7
		fileName = "7.txt"
		picBase64Content = ReadFile(fileName)
		index++
		DocPerson = CheckInRecord{
			Id:            index,
			Name:          "p7",
			Check_in_time: "",
			Pic:           picBase64Content,
			Leave_type:    "事",
			Date:          myDate,
			Department:    "軟體",
			Position:      "軟體工程師",
		}
		result, insertErr = col.InsertOne(ctx, DocPerson)

		if insertErr != nil {
			fmt.Println("InsertOne ERROR:", insertErr)
			os.Exit(1) // safely exit script on error
		} else {
			// fmt.Println("InsertOne() result type: ", reflect.TypeOf(result))
			//fmt.Println("InsertOne() API result:", result)

			// get the inserted ID string
			newID := result.InsertedID
			fmt.Println("InsertOne() newID:", newID)
			//fmt.Println("InsertOne() newID type:", reflect.TypeOf(newID))
		}
	}

}

/** 插入一整年假資料A:寫法一 */
func insertCheckInStatisticsOneYear() {
	/** 連資料庫 MongoDB
	 * 寫法:透過 mgo.Dial 撥號，return seesion **/

	// 撥號
	uri := "localhost:27017"
	session, err := mgo.Dial(uri)
	if err != nil {
		log.Fatal("Couldn't connect to db.", err)
	}

	// 等所有函數都執行回傳完後 最後關閉session
	defer session.Close()

	// 取得 DB collection
	collection := session.DB("leapsy_env").C("check_in_statistics")

	/*寫入json*/
	var bdoc interface{}

	jsonString := `{ "id": 1,
	"date": "2020-01-01",
	"expected": 30,
	"attendance": 27,
	"not_arrived": 3, 
	"guests": 4 }`

	// 要比對符合此形狀(`"date": "\d{4}-\d{2}-\d{2}"`)的string 來進行部份string替換
	regularExpressionForDate := regexp.MustCompile(`"date": "\d{4}-\d{2}-\d{2}"`)
	regularExpressionForID := regexp.MustCompile(`"id": \d{1}`)

	id := 1

	// 差入一整年的統計假資料
	for myTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local); myTime != time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local); myTime = myTime.AddDate(0, 0, 1) {

		// 指定新變數(新日期) 透過Sprintf來回傳想要的格式 (0代表若沒有值則用0替代,4d表示有四位數) 以年月日來取代
		newStringDate := fmt.Sprintf(`"date": "%04d-%02d-%02d"`, myTime.Year(), myTime.Month(), myTime.Day())
		newStringID := fmt.Sprintf(`"id": %1d`, id)

		// 傳入整個jsonString，進行jsonString內容的比對，若符合 regularExpressione格式的部份，則將其部份替換成 newString，最後回傳整個新的JSON字串
		newJSONString := regularExpressionForDate.ReplaceAllString(jsonString, newStringDate) // 換掉日期
		newJSONString = regularExpressionForID.ReplaceAllString(newJSONString, newStringID)   // 換掉id

		// 將新的JSONString 轉換interface{}格式放入 bdoc中
		err = bson.UnmarshalJSON([]byte(newJSONString), &bdoc)

		if err != nil {
			panic(err)
		}

		err = collection.Insert(&bdoc)
		if err != nil {
			panic(err)
		}

		id++
	}
}
