package servers

import (
	"time"

	"../databases"
	"../logings"
	"../times"
)

// startPeriodicallyRecord - 開始週期性記錄
/**
 * @param  *databases.ECSDB eCSDB 來源資料庫
 * @param  *databases.MongoDB mongoDB 目的資料庫
 */
func startPeriodicallyRecord(eCSDB *databases.ECSDB, mongoDB *databases.MongoDB) {

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`啟動 週期性記錄 `},
		[]interface{}{},
		nil,
	)

	go logger.Infof(formatString, args...) // 記錄資訊

	durationString := `15s`
	duration, timeParseDurationError := time.ParseDuration(durationString) // 解析期間

	if nil != timeParseDurationError { // 若解析期間錯誤

		// 取得記錄器格式和參數
		formatString, args = logings.GetLogFuncFormatAndArguments(
			[]string{`解析期間字串 %s `},
			[]interface{}{durationString},
			timeParseDurationError,
		)

		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式

	} else {

		for upper := time.Now(); ; <-time.After(time.Until(upper)) { // 針對每一時間
			upper = time.Date(upper.Year(), upper.Month(), upper.Day(), upper.Hour(), upper.Minute(), upper.Second(), 0, time.Local) // 修改時間
			secondRecord := eCSDB.Read().SecondRecord()                                                                              // 讀取秒紀錄
			secondRecord.Time = upper                                                                                                // 儲存時間
			go mongoDB.PeriodicallyInsert(secondRecord)                                                                              // 添寫秒記錄
			upper = upper.Add(duration)                                                                                              // 設定下一時間
		}

	}

}

// stopPeriodicallyRecord - 結束週期性記錄
/**
 * @param  *databases.ECSDB eCSDB 來源資料庫
 * @param  *databases.MongoDB mongoDB 目的資料庫
 */
func stopPeriodicallyRecord(eCSDB *databases.ECSDB, mongoDB *databases.MongoDB) {

	eCSDB.Disconnect()   // 中斷與資料庫的連線
	mongoDB.Disconnect() // 中斷與資料庫的連線

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`結束 週期性記錄 `},
		[]interface{}{},
		nil,
	)

	logger.Infof(formatString, args...) // 記錄資訊

}

// startHourlyRecord - 開始每時記錄
/**
 * @param  *databases.MongoDB mongoDB 資料庫
 */
func startHourlyRecord(mongoDB *databases.MongoDB) {

	go func() {
		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			[]string{`啟動 每小時記錄 `},
			[]interface{}{},
			nil,
		)

		logger.Infof(formatString, args...) // 記錄資訊

	}()

	for { // 循環

		_, upper := times.GetHourlyBounds(time.Now()) //取得上下限小時

		<-time.After(time.Until(upper)) // 等到上限小時

		mongoDB.InsertHourlyByTimeIfNotExisted(upper) // 若小時紀錄不存在則添寫

	}

}

// stopHourlyRecord - 結束每時記錄
/**
 * @param  *databases.MongoDB mongoDB 資料庫
 */
func stopHourlyRecord(mongoDB *databases.MongoDB) {

	mongoDB.Disconnect() // 中斷與資料庫的連線

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`結束 每小時記錄 `},
		[]interface{}{},
		nil,
	)

	logger.Infof(formatString, args...) // 記錄資訊

}

// startDailyRecord - 開始每日記錄
/**
 * @param  *databases.MongoDB mongoDB 資料庫
 */
func startDailyRecord(mongoDB *databases.MongoDB) {

	go func() {
		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			[]string{`啟動 每日記錄 `},
			[]interface{}{},
			nil,
		)

		logger.Infof(formatString, args...) // 記錄資訊

	}()

	for { // 循環

		_, upper := times.GetDailyBounds(time.Now()) // 取得上下限日

		<-time.After(time.Until(upper)) // 等到上限日

		mongoDB.InsertDailyByTimeIfNotExisted(upper) // 若日紀錄不存在則添寫
	}

}

// stopDailyRecord - 結束每日記錄
/**
 * @param  *databases.MongoDB mongoDB 資料庫
 */
func stopDailyRecord(mongoDB *databases.MongoDB) {

	mongoDB.Disconnect() // 中斷與資料庫的連線

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`啟動 每日記錄 `},
		[]interface{}{},
		nil,
	)

	logger.Infof(formatString, args...) // 記錄資訊

}

// startRecordAlerts - 開始記錄警報
/**
 * @param  *databases.ECSAlertDB ecsAlertDB 來源資料庫
 * @param  *databases.MongoDB mongoDB 目的資料庫
 */
func startRecordAlerts(ecsAlertDB *databases.ECSAlertDB, mongoDB *databases.MongoDB) {

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`啟動 記錄警報 `},
		[]interface{}{},
		nil,
	)

	go logger.Infof(formatString, args...) // 記錄資訊

	durationString := `15s`
	duration, timeParseDurationError := time.ParseDuration(durationString) // 解析期間

	if nil != timeParseDurationError { // 若解析期間錯誤

		// 取得記錄器格式和參數
		formatString, args = logings.GetLogFuncFormatAndArguments(
			[]string{`解析期間字串 %s `},
			[]interface{}{durationString},
			timeParseDurationError,
		)

		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式

	} else {

		go logger.Infof(formatString, args...) // 記錄資訊

		for ; ; <-time.After(duration) { // 針對每一時間

			if ecsAlertDB.CountAll() > mongoDB.CountAllAlertRecords() { // 若環控資料庫警報數 > 警報數

				for _, ecsAlertRecord := range ecsAlertDB.ReadLast(ecsAlertDB.CountAll() - mongoDB.CountAllAlertRecords()) {
					mongoDB.InsertAlertRecord(ecsAlertRecord.AlertRecord()) // 添寫警報記錄
				}

			}
		}

	}

}

// stopRecordAlerts - 結束記錄警報
/**
 * @param  *databases.ECSAlertDB ecsAlertDB 來源資料庫
 * @param  *databases.MongoDB mongoDB 目的資料庫
 */
func stopRecordAlerts(ecsAlertDB *databases.ECSAlertDB, mongoDB *databases.MongoDB) {

	ecsAlertDB.Disconnect() // 中斷與資料庫的連線
	mongoDB.Disconnect()    // 中斷與資料庫的連線

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`結束 記錄警報 `},
		[]interface{}{},
		nil,
	)

	logger.Infof(formatString, args...) // 記錄資訊

}
