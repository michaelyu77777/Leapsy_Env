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

// CountDailyRecords - 計算日記錄個數
/**
 * @param primitive.M filter 過濾器
 * @retrun int returnCount 日記錄個數
 */
func (mongoDB *MongoDB) CountDailyRecords(filter primitive.M) (returnCount int) {

	defer mongoDB.Disconnect()

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		mongoDB.GetConfigValueOrPanic(`server`),
		mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	dailyRWMutex.RLock() // 讀鎖

	// 取得日記錄個數
	count, countError := mongoDB.Connect().
		Database(mongoDB.GetConfigValueOrPanic(`database`)).
		Collection(mongoDB.GetConfigValueOrPanic(`daily-table`)).
		CountDocuments(context.TODO(), filter)

	dailyRWMutex.RUnlock() // 讀解鎖

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 取得日記錄個數`},
		network.GetAliasAddressPair(address),
		countError,
	)

	if nil != countError && mongo.ErrNilDocument != countError { // 若取得日記錄個數錯誤，且不為空資料表錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若取得日記錄個數成功
		go logger.Infof(formatString, args...) // 記錄資訊
		returnCount = int(count)
	}

	return // 回傳
}

// CountDailyRecordByTime - 依據時間計算日記錄數
/**
 * @param time.Time dateTime 時間
 * @return int result 取得結果
 */
func (mongoDB *MongoDB) CountDailyRecordByTime(dateTime time.Time) (result int) {

	if !dateTime.IsZero() {
		result = mongoDB.CountDailyRecords(
			bson.M{
				`time`: bson.M{
					equalToConstString: dateTime,
				},
			},
		)
	}

	return // 回傳
}

// CountDailyRecordsBetweenTimes - 依據時間區間計算日紀錄數
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return int result 取得結果
 */
func (mongoDB *MongoDB) CountDailyRecordsBetweenTimes(
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
		result = mongoDB.CountDailyRecords(
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

// FindDailyRecords - 取得日紀錄
/**
 * @param bson.M filter 過濾器
 * @param ...*options.FindOptions opts 選項
 * @return []records.DailyRecord results 取得結果
 */
func (mongoDB *MongoDB) FindDailyRecords(filter primitive.M, opts ...*options.FindOptions) (results []records.DailyRecord) {

	dailyRWMutex.RLock() // 讀鎖

	// 查找紀錄
	cursor, findError := mongoDB.Connect().
		Database(mongoDB.GetConfigValueOrPanic(`database`)).
		Collection(mongoDB.GetConfigValueOrPanic(`daily-table`)).
		Find(
			context.TODO(),
			filter,
			opts...,
		)

	dailyRWMutex.RUnlock() // 讀解鎖

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`查找日記錄`},
		[]interface{}{},
		findError,
	)

	if nil != findError { // 若查找日紀錄錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若查找日記錄成功

		for cursor.Next(context.TODO()) { // 針對每一紀錄

			var DailyRecord records.DailyRecord

			cursorDecodeError := cursor.Decode(&DailyRecord) // 解析紀錄

			if nil != cursorDecodeError { // 若解析記錄錯誤

				// 取得記錄器格式和參數
				formatString, args = logings.GetLogFuncFormatAndArguments(
					[]string{`解析日記錄`},
					[]interface{}{},
					cursorDecodeError,
				)

				logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
			}

			DailyRecord.Time = DailyRecord.Time.Local() // 儲存為本地時間格式

			results = append(results, DailyRecord) // 儲存紀錄
		}

		if cursorErrError := cursor.Err(); nil != cursorErrError { // 若遊標錯誤

			// 取得記錄器格式和參數
			formatString, args := logings.GetLogFuncFormatAndArguments(
				[]string{`查找日記錄遊標運作`},
				[]interface{}{},
				cursorErrError,
			)

			logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
		}

	}

	return // 回傳
}

// FindDailyRecordsBetweenTimes - 依據時間區間取得日紀錄
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return []records.DailyRecord results 取得結果
 */
func (mongoDB *MongoDB) FindDailyRecordsBetweenTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (results []records.DailyRecord) {

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
		results = mongoDB.FindDailyRecords(
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
