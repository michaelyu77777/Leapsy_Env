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

// CountHourlyRecords - 計算小時記錄個數
/**
 * @param primitive.M filter 過濾器
 * @retrun int returnCount 小時記錄個數
 */
func (mongoDB *MongoDB) CountHourlyRecords(filter primitive.M) (returnCount int) {

	defer mongoDB.Disconnect()

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		mongoDB.GetConfigValueOrPanic(`server`),
		mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	hourlyRWMutex.RLock() // 讀鎖

	// 取得小時記錄個數
	count, countError := mongoDB.Connect().
		Database(mongoDB.GetConfigValueOrPanic(`database`)).
		Collection(mongoDB.GetConfigValueOrPanic(`hourly-table`)).
		CountDocuments(context.TODO(), filter)

	hourlyRWMutex.RUnlock() // 讀解鎖

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 取得小時記錄個數`},
		network.GetAliasAddressPair(address),
		countError,
	)

	if nil != countError && mongo.ErrNilDocument != countError { // 若取得小時記錄個數錯誤，且不為空資料表錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若取得小時記錄個數成功
		go logger.Infof(formatString, args...) // 記錄資訊
		returnCount = int(count)
	}

	return // 回傳
}

// CountHourlyRecordByTime - 依據時間計算小時記錄數
/**
 * @param time.Time dateTime 時間
 * @return int result 取得結果
 */
func (mongoDB *MongoDB) CountHourlyRecordByTime(dateTime time.Time) (result int) {

	if !dateTime.IsZero() {
		result = mongoDB.CountHourlyRecords(
			bson.M{
				`time`: bson.M{
					equalToConstString: dateTime,
				},
			},
		)
	}

	return // 回傳
}

// CountHourlyRecordsBetweenTimes - 依據時間區間計算小時記錄數
/**
 * * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return int result 取得結果
 */
func (mongoDB *MongoDB) CountHourlyRecordsBetweenTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (result int) {

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
		result = mongoDB.CountHourlyRecords(
			bson.M{
				`time`: bson.M{
					greaterThanKeyword: lowerTime,
					lessThanKeyword:    upperTime,
				},
			},
		)

	}

	return // 回傳
}

// FindHourlyRecords - 取得小時紀錄
/**
 * @param bson.M filter 過濾器
 * @param ...*options.FindOptions opts 選項
 * @return []records.HourlyRecord results 取得結果
 */
func (mongoDB *MongoDB) FindHourlyRecords(filter primitive.M, opts ...*options.FindOptions) (results []records.HourlyRecord) {

	hourlyRWMutex.RLock() // 讀鎖

	// 查找紀錄
	cursor, findError := mongoDB.Connect().
		Database(mongoDB.GetConfigValueOrPanic(`database`)).
		Collection(mongoDB.GetConfigValueOrPanic(`hourly-table`)).
		Find(
			context.TODO(),
			filter,
			opts...,
		)

	hourlyRWMutex.RUnlock() // 讀解鎖

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`查找小時記錄`},
		[]interface{}{},
		findError,
	)

	if nil != findError { // 若查找小時紀錄錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若查找小時記錄成功

		for cursor.Next(context.TODO()) { // 針對每一紀錄

			var hourlyRecord records.HourlyRecord

			cursorDecodeError := cursor.Decode(&hourlyRecord) // 解析紀錄

			if nil != cursorDecodeError { // 若解析記錄錯誤

				// 取得記錄器格式和參數
				formatString, args = logings.GetLogFuncFormatAndArguments(
					[]string{`解析小時記錄`},
					[]interface{}{},
					cursorDecodeError,
				)

				logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
			}

			hourlyRecord.Time = hourlyRecord.Time.Local() // 儲存為本地時間格式

			results = append(results, hourlyRecord) // 儲存紀錄
		}

		if cursorErrError := cursor.Err(); nil != cursorErrError { // 若遊標錯誤

			// 取得記錄器格式和參數
			formatString, args := logings.GetLogFuncFormatAndArguments(
				[]string{`查找小時記錄遊標運作`},
				[]interface{}{},
				cursorErrError,
			)

			logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
		}

	}

	return // 回傳
}

// FindHourlyRecordsBetweenTimes - 依據時間區間取得小時紀錄
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return []records.HourlyRecord results 取得結果
 */
func (mongoDB *MongoDB) FindHourlyRecordsBetweenTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (results []records.HourlyRecord) {

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
		results = mongoDB.FindHourlyRecords(
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
