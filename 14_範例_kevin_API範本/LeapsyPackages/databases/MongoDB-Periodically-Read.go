package databases

import (
	"context"
	"fmt"
	"time"

	"../logings"
	"../network"
	"../records"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CountSecondRecords - 計算秒記錄個數
/**
 * @param primitive.M filter 過濾器
 * @retrun int returnCount 秒記錄個數
 */
func (mongoDB *MongoDB) CountSecondRecords(filter primitive.M) (returnCount int) {

	defer mongoDB.Disconnect()

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		mongoDB.GetConfigValueOrPanic(`server`),
		mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	periodicallyRWMutex.RLock() // 讀鎖

	// 取得秒記錄個數
	count, countError := mongoDB.Connect().
		Database(mongoDB.GetConfigValueOrPanic(`database`)).
		Collection(mongoDB.GetConfigValueOrPanic(`second-table`)).
		CountDocuments(context.TODO(), filter)

	periodicallyRWMutex.RUnlock() // 讀解鎖

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 取得秒記錄個數`},
		network.GetAliasAddressPair(address),
		countError,
	)

	if nil != countError && mongo.ErrNilDocument != countError { // 若取得秒記錄個數錯誤，且不為空資料表錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若取得秒記錄個數成功
		go logger.Infof(formatString, args...) // 記錄資訊
		returnCount = int(count)
	}

	return // 回傳
}

// CountSecondRecordByTime - 依據時間計算秒記錄數
/**
 * @param time.Time dateTime 時間
 * @return int result 取得結果
 */
func (mongoDB *MongoDB) CountSecondRecordByTime(dateTime time.Time) (result int) {

	if !dateTime.IsZero() {
		result = mongoDB.CountSecondRecords(
			bson.M{
				`time`: bson.M{
					equalToConstString: dateTime,
				},
			},
		)
	}

	return // 回傳
}

// FindSecondRecords - 取得秒紀錄
/**
 * @param bson.M filter 過濾器
 * @param ...*options.FindOptions opts 選項
 * @return []records.SecondRecord results 取得結果
 */
func (mongoDB *MongoDB) FindSecondRecords(filter primitive.M, opts ...*options.FindOptions) (results []records.SecondRecord) {

	periodicallyRWMutex.RLock() // 讀鎖

	// 查找紀錄
	cursor, findError := mongoDB.Connect().
		Database(mongoDB.GetConfigValueOrPanic(`database`)).
		Collection(mongoDB.GetConfigValueOrPanic(`second-table`)).
		Find(
			context.TODO(),
			filter,
			opts...,
		)

	periodicallyRWMutex.RUnlock() // 讀解鎖

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`查找秒記錄`},
		[]interface{}{},
		findError,
	)

	if nil != findError { // 若查找秒紀錄錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若查找秒記錄成功

		for cursor.Next(context.TODO()) { // 針對每一紀錄

			var secondRecord records.SecondRecord

			cursorDecodeError := cursor.Decode(&secondRecord) // 解析紀錄

			if nil != cursorDecodeError { // 若解析記錄錯誤

				// 取得記錄器格式和參數
				formatString, args = logings.GetLogFuncFormatAndArguments(
					[]string{`解析秒記錄`},
					[]interface{}{},
					cursorDecodeError,
				)

				logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
			}

			secondRecord.Time = secondRecord.Time.Local() // 儲存為本地時間格式

			results = append(results, secondRecord) // 儲存紀錄
		}

		if cursorErrError := cursor.Err(); nil != cursorErrError { // 若遊標錯誤

			// 取得記錄器格式和參數
			formatString, args := logings.GetLogFuncFormatAndArguments(
				[]string{`查找秒記錄遊標運作`},
				[]interface{}{},
				cursorErrError,
			)

			logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
		}

	}

	return // 回傳
}

// FindSecondRecordsBetweenTimes - 依據時間區間取得秒紀錄
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return []records.SecondRecord results 取得結果
 */
func (mongoDB *MongoDB) FindSecondRecordsBetweenTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (results []records.SecondRecord) {

	if !lowerTime.IsZero() && !upperTime.IsZero() { //若上下限時間不為零時間

		var (
			greaterThanKeyword, lessThanKeyword string // 比較關鍵字
		)

		if !isLowerTimeIncluded { // 若不包含下限時間
			greaterThanKeyword = greaterThanConstString // >
		} else {
			greaterThanKeyword = greaterThanEqualToConstString // >=
		}

		if !isUpperTimeIncluded { // 若不包含上限時間
			lessThanKeyword = lessThanConstString // <
		} else {
			lessThanKeyword = lessThanEqualToConstString // <=
		}

		// 回傳結果
		results = mongoDB.FindSecondRecords(
			bson.M{
				`time`: bson.M{
					greaterThanKeyword: lowerTime,
					lessThanKeyword:    upperTime,
				},
			},
			options.
				Find().
				SetSort(
					bson.M{
						`time`: 1,
					},
				),
		)

	}

	return // 回傳
}
