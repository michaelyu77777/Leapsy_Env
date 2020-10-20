package databases

import (
	"context"
	"fmt"
	"time"

	"../logings"
	"../network"
	"../records"
)

// DailyInsert - 每日添寫日記錄
/**
 * @param  records.DailyRecord dailyRecord  日記錄
 */
func (mongoDB *MongoDB) DailyInsert(dailyRecord records.DailyRecord) {

	defer mongoDB.Disconnect()

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		mongoDB.GetConfigValueOrPanic(`server`),
		mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	dailyRWMutex.Lock() // 寫鎖

	// 添寫日記錄
	_, insertOneError := mongoDB.Connect().
		Database(mongoDB.GetConfigValueOrPanic(`database`)).
		Collection(mongoDB.GetConfigValueOrPanic(`daily-table`)).
		InsertOne(context.TODO(), dailyRecord)

	dailyRWMutex.Unlock() // 寫解鎖

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 添寫日記錄`},
		network.GetAliasAddressPair(address),
		insertOneError,
	)

	if nil != insertOneError { // 若添寫日紀錄錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若添寫日記錄成功
		go logger.Infof(formatString, args...) // 記錄資訊
	}

}

// InsertDailyByTimeIfNotExisted - 若日紀錄不存在則添寫
/**
 * @param  @param time.Time dateTime 時間
 */
func (mongoDB *MongoDB) InsertDailyByTimeIfNotExisted(dateTime time.Time) {

	dailyDateTime := convertToDailyDateTime(dateTime) // 轉成整日

	if 0 == mongoDB.CountDailyRecordByTime(dailyDateTime) { // 若紀錄數為零

		thisPMKwhToday := mongoDB.CalculatePMKwhTodayForDailyRecord(dailyDateTime) // 計算日紀錄用的今日PM累計電度

		// 新日紀錄
		newDailyRecord := records.DailyRecord{
			Time:           dailyDateTime,
			PMKwhThisMonth: mongoDB.CalculatePMKwhThisMonth(dailyDateTime) + thisPMKwhToday,
			PMKwhToday:     thisPMKwhToday,
		}

		mongoDB.DailyInsert(newDailyRecord) // 添寫日紀錄
	}

}
