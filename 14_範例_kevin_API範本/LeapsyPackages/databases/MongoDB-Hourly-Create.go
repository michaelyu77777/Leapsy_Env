package databases

import (
	"context"
	"fmt"
	"time"

	"../logings"
	"../network"
	"../records"
)

// HourlyInsert - 每小時添寫小時記錄
/**
 * @param  records.HourlyRecord hourlyRecord  小時記錄
 */
func (mongoDB *MongoDB) HourlyInsert(hourlyRecord records.HourlyRecord) {

	defer mongoDB.Disconnect()

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		mongoDB.GetConfigValueOrPanic(`server`),
		mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	hourlyRWMutex.Lock() // 寫鎖

	// 添寫小時記錄
	_, insertOneError := mongoDB.Connect().
		Database(mongoDB.GetConfigValueOrPanic(`database`)).
		Collection(mongoDB.GetConfigValueOrPanic(`hourly-table`)).
		InsertOne(context.TODO(), hourlyRecord)

	hourlyRWMutex.Unlock() // 寫解鎖

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 添寫小時記錄`},
		network.GetAliasAddressPair(address),
		insertOneError,
	)

	if nil != insertOneError { // 若添寫小時紀錄錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若添寫小時記錄成功
		go logger.Infof(formatString, args...) // 記錄資訊
	}

}

// InsertHourlyByTimeIfNotExisted - 若小時紀錄不存在則添寫
/**
 * @param time.Time dateTime 時間
 */
func (mongoDB *MongoDB) InsertHourlyByTimeIfNotExisted(dateTime time.Time) {

	hourlyDateTime := convertToHourlyDateTime(dateTime) // 轉成整點時間

	if 0 == mongoDB.CountHourlyRecordByTime(hourlyDateTime) { // 若紀錄數為零
		thisPMKwhThisHour := mongoDB.CalculatePMKwhThisHour(hourlyDateTime) //計算該小時PM累計電度

		// 新小時記錄
		newHourlyRecord := records.HourlyRecord{
			Time:          hourlyDateTime,
			PMKwhThisHour: thisPMKwhThisHour,
			PMKwhToday:    mongoDB.CalculatePMKwhTodayForHourlyRecord(hourlyDateTime) + thisPMKwhThisHour,
		}

		mongoDB.HourlyInsert(newHourlyRecord) // 添寫小時紀錄
	}

}
