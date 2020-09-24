package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
)

type MongoFields struct {
	//Key string `json:"key,omitempty"`
	ID        int    `json:"_id"`
	Fiel_dStr string `json:"Field_Str,omitempty"` //json欄位要長這樣 若field非
	FieldInt  int    `json:"Field_Int,omitempty"`
	FieldBool bool   `json:"Field_Bool,omitempty"`
}

func main() {

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

	// 插入一年的統計假資料
	// for i := 0; i <= 360; i++ {

	// 	// 統計數據一
	// 	err = bson.UnmarshalJSON([]byte(`{	"id": 1,
	// 										"date": "2020-01-01",
	// 										"expected": 30,
	// 										"attendance": 27,
	// 										"not_arrived": 3,
	// 										"guests": 4 }`), &bdoc)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	err = collection.Insert(&bdoc)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
}
