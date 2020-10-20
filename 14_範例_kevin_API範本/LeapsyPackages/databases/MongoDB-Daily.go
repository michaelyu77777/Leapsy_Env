package databases

import (
	"reflect"
	"time"

	"../records"
	"../times"

	"go.mongodb.org/mongo-driver/bson"
)

// SumIntDailyRecordFields - 計算兩時間內的整數型欄位和
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @param []string fieldNames 欄位名
 * @return map[string]int results 計算結果
 */
func (mongoDB *MongoDB) SumIntDailyRecordFields(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
	fieldNames ...string,
) (results map[string]int) {

	results = make(map[string]int) // 指派空間

	for _, dailyRecord := range mongoDB.FindDailyRecordsBetweenTimes(lowerTime, isLowerTimeIncluded, upperTime, isUpperTimeIncluded) { // 每一紀錄

		dailyRecordReflect := reflect.ValueOf(dailyRecord) // 紀錄

		for _, fieldName := range fieldNames { // 每一欄位

			if dailyRecordReflect.FieldByName(fieldName).Kind() == reflect.Int { // 若為整數型欄位
				results[fieldName] += int(dailyRecordReflect.FieldByName(fieldName).Int()) // 加上欄位值
			} else { // 若非整數型欄位
				results[fieldName] = 0 // 結果為零
			}

		}

	}

	return // 回傳
}

// SubtractIntDailyRecordFields - 計算兩時間內的整數型欄位相減
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @param []string fieldNames 欄位名
 * @return map[string]int results 計算結果
 */
func (mongoDB *MongoDB) SubtractIntDailyRecordFields(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
	fieldNames ...string,
) (results map[string]int) {

	results = make(map[string]int) // 指派空間

	dailyRecords := mongoDB.FindDailyRecordsBetweenTimes(lowerTime, isLowerTimeIncluded, upperTime, isUpperTimeIncluded) // 查找紀錄
	lengthOfDailyRecords := len(dailyRecords)                                                                            // 紀錄長度

	if lengthOfDailyRecords > 1 { // 若查找有紀錄

		firstDailyRecordReflect := reflect.ValueOf(dailyRecords[0])                     // 第一筆紀錄
		lastDailyRecordReflect := reflect.ValueOf(dailyRecords[lengthOfDailyRecords-1]) // 最末筆紀錄

		for _, fieldName := range fieldNames { // 針對每一欄位

			if firstDailyRecordReflect.FieldByName(fieldName).Kind() == reflect.Int { // 若為整數型欄位
				// 末紀錄欄位值 - 首紀錄欄位值
				results[fieldName] = int(
					lastDailyRecordReflect.FieldByName(fieldName).Int() -
						firstDailyRecordReflect.FieldByName(fieldName).Int(),
				)
			} else { // 若非整數型欄位
				results[fieldName] = 0 // 結果為零
			}

		}

	}

	return // 回傳
}

// CalculatePMKwhThisMonth - 計算這個月的PM累計電度
/**
 * @param time.Time dateTime 時間
 * @return int 計算結果
 */
func (mongoDB *MongoDB) CalculatePMKwhThisMonth(dateTime time.Time) int {

	lower, upper := times.GetMonthlyBounds(dateTime) // 取得上下限月

	return mongoDB.SumIntDailyRecordFields(lower, false, upper, true, `PMKwhToday`)[`PMKwhToday`] // 計算欄位差
}

// CalculatePMKwhTodayForDailyRecord - 計算日紀錄用的今日PM累計電度
/**
 * @param time.Time dateTime 時間
 * @return int 計算結果
 */
func (mongoDB *MongoDB) CalculatePMKwhTodayForDailyRecord(dateTime time.Time) int {

	lower, _ := times.GetDailyBounds(dateTime) // 取得上下限日

	results := mongoDB.FindHourlyRecordsBetweenTimes(lower, false, dateTime, true) // 查找紀錄
	lengthOfResults := len(results)                                                // 紀錄長度

	if lengthOfResults > 0 { // 若查找有紀錄
		return results[lengthOfResults-1].PMKwhToday // 回傳末紀錄欄位值
	}

	return 0 // 回傳零
}

// CreateDailyRecord - 產生小時記錄
/**
 * @param time.Time dateTime 時間
 * @return records.DailyRecord returnDailyRecord 回傳小時紀錄
 */
func (mongoDB *MongoDB) CreateDailyRecord(dateTime time.Time) (returnDailyRecord records.DailyRecord) {

	if !dateTime.IsZero() { // 若非零時間

		dailyDateTime := convertToDailyDateTime(dateTime) // 轉成整點

		thisPMKwhToday := mongoDB.CalculatePMKwhTodayForDailyRecord(dailyDateTime) // 計算日紀錄用的今日PM累計電度

		// 回傳日紀錄
		returnDailyRecord = records.DailyRecord{
			Time:           dailyDateTime,
			PMKwhThisMonth: mongoDB.CalculatePMKwhThisMonth(dailyDateTime) + thisPMKwhToday,
			PMKwhToday:     thisPMKwhToday,
		}

	}

	return
}

//  convertToDailyDateTime - 轉成整日時間
/**
 * @param time.Time dateTime 時間
 * @return time.Time returnDailyDateTime 回傳小時時間
 */
func convertToDailyDateTime(dateTime time.Time) (returnDailyDateTime time.Time) {

	if !isDailyDateTime(dateTime) { // 若非整日時間
		// 修改時間
		returnDailyDateTime = time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, time.Local)
	} else {
		returnDailyDateTime = dateTime // 回傳
	}

	return
}

//  isDailyDateTime - 是否整日時間
/**
 * @param time.Time dateTime 時間
 * @return bool result 回傳結果
 */
func isDailyDateTime(dateTime time.Time) (result bool) {

	result = 0 == dateTime.Hour() && isHourlyDateTime(dateTime) // 整日整點

	return
}

// RepsertDailyRecordsBetweenTimes - 代添時間之間小時紀錄
/**
 * @param time.Time low 下限時間
 * @param time.Time upper 上限時間
 */
func (mongoDB *MongoDB) RepsertDailyRecordsBetweenTimes(low, upper time.Time) {

	dateTime := convertToDailyDateTime(low) // 時間

	if dateTime.Before(low) { // 若時間 < 下限時間
		dateTime = dateTime.AddDate(0, 0, 1) // 多一日
	}

	for ; !upper.Before(dateTime); dateTime = dateTime.AddDate(0, 0, 1) { // 每一日
		mongoDB.repsertOneDailyRecord(bson.M{`time`: dateTime}, mongoDB.CreateDailyRecord(dateTime).PrimitiveM()) // 代添紀錄
	}

}
