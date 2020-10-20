package databases

import (
	"reflect"
	"time"

	"../records"
	"../times"
	"go.mongodb.org/mongo-driver/bson"
)

// SumIntHourlyRecordFields - 計算兩時間內的整數型欄位和
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @param []string fieldNames 欄位名
 * @return map[string]int results 計算結果
 */
func (mongoDB *MongoDB) SumIntHourlyRecordFields(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
	fieldNames ...string,
) (results map[string]int) {

	results = make(map[string]int) // 指派空間

	for _, hourlyRecord := range mongoDB.FindHourlyRecordsBetweenTimes(lowerTime, isLowerTimeIncluded, upperTime, isUpperTimeIncluded) { // 每一紀錄

		hourlyRecordReflect := reflect.ValueOf(hourlyRecord) // 紀錄

		for _, fieldName := range fieldNames { // 每一欄位

			if hourlyRecordReflect.FieldByName(fieldName).Kind() == reflect.Int { // 若為整數型欄位
				results[fieldName] += int(hourlyRecordReflect.FieldByName(fieldName).Int()) // 加上欄位值
			} else { // 若非整數型欄位
				results[fieldName] = 0 // 結果為零
			}

		}

	}

	return // 回傳
}

// SubtractIntHourlyRecordFields - 計算兩時間內的整數型欄位相減
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @param []string fieldNames 欄位名
 * @return map[string]int results 計算結果
 */
func (mongoDB *MongoDB) SubtractIntHourlyRecordFields(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
	fieldNames ...string,
) (results map[string]int) {

	results = make(map[string]int) // 指派空間

	hourlyRecords := mongoDB.FindHourlyRecordsBetweenTimes(lowerTime, isLowerTimeIncluded, upperTime, isUpperTimeIncluded) // 查找紀錄
	lengthOfHourlyRecords := len(hourlyRecords)                                                                            // 紀錄長度

	if lengthOfHourlyRecords > 1 { // 若查找有紀錄

		firstHourlyRecordReflect := reflect.ValueOf(hourlyRecords[0])                      // 第一筆紀錄
		lastHourlyRecordReflect := reflect.ValueOf(hourlyRecords[lengthOfHourlyRecords-1]) // 最末筆紀錄

		for _, fieldName := range fieldNames { // 針對每一欄位

			if firstHourlyRecordReflect.FieldByName(fieldName).Kind() == reflect.Int { // 若為整數型欄位
				// 末紀錄欄位值 - 首紀錄欄位值
				results[fieldName] = int(
					lastHourlyRecordReflect.FieldByName(fieldName).Int() -
						firstHourlyRecordReflect.FieldByName(fieldName).Int(),
				)
			} else { // 若非整數型欄位
				results[fieldName] = 0 // 結果為零
			}

		}

	}

	return // 回傳
}

// CalculatePMKwhThisHour - 計算該小時PM累計電度
/**
 * @param time.Time dateTime 時間
 * @return int 計算結果
 */
func (mongoDB *MongoDB) CalculatePMKwhThisHour(dateTime time.Time) int {
	lower, upper := times.GetHourlyBounds(dateTime) // 取得上下限小時

	return mongoDB.SubtractIntSecondRecordFields(lower, true, upper, true, `PMKwh`)[`PMKwh`] // 計算欄位差
}

// CalculatePMKwhTodayForHourlyRecord - 計算今日的PM累計電度
/**
 * @param time.Time dateTime 時間
 * @return int 計算結果
 */
func (mongoDB *MongoDB) CalculatePMKwhTodayForHourlyRecord(dateTime time.Time) int {

	lower, upper := times.GetDailyBounds(dateTime) // 取得上下限日

	return mongoDB.SumIntHourlyRecordFields(lower, false, upper, true, `PMKwhThisHour`)[`PMKwhThisHour`] //計算欄位和
}

// CreateHourlyRecord - 產生小時記錄
/**
 * @param time.Time dateTime 時間
 * @return records.HourlyRecord returnHourlyRecord 回傳小時紀錄
 */
func (mongoDB *MongoDB) CreateHourlyRecord(dateTime time.Time) (returnHourlyRecord records.HourlyRecord) {

	if !dateTime.IsZero() { // 若非零時間

		hourlyDateTime := convertToHourlyDateTime(dateTime) // 轉成整點

		thisPMKwhThisHour := mongoDB.CalculatePMKwhThisHour(hourlyDateTime) //計算該小時PM累計電度

		// 回傳小時記錄
		returnHourlyRecord = records.HourlyRecord{
			Time:          hourlyDateTime,
			PMKwhThisHour: thisPMKwhThisHour,
			PMKwhToday:    mongoDB.CalculatePMKwhTodayForHourlyRecord(hourlyDateTime) + thisPMKwhThisHour,
		}

	}

	return
}

//  convertToHourlyDateTime - 轉成整點時間
/**
 * @param time.Time dateTime 時間
 * @return time.Time returnHourlyDateTime 回傳小時時間
 */
func convertToHourlyDateTime(dateTime time.Time) (returnHourlyDateTime time.Time) {

	if !isHourlyDateTime(dateTime) { // 若非整點時間
		// 修改時間
		returnHourlyDateTime = time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), dateTime.Hour(), 0, 0, 0, time.Local)
	} else {
		returnHourlyDateTime = dateTime // 回傳
	}

	return
}

//  isHourlyDateTime - 是否整點時間
/**
 * @param time.Time dateTime 時間
 * @return bool result 回傳結果
 */
func isHourlyDateTime(dateTime time.Time) (result bool) {

	result = 0 == dateTime.Minute() && isSecondDateTime(dateTime) // 整分整秒

	return
}

// RepsertHourlyRecordsBetweenTimes - 代添時間之間小時紀錄
/**
 * @param time.Time low 下限時間
 * @param time.Time upper 上限時間
 */
func (mongoDB *MongoDB) RepsertHourlyRecordsBetweenTimes(low, upper time.Time) {

	dateTime := convertToHourlyDateTime(low) // 時間
	duration := time.Hour                    // 期間

	if dateTime.Before(low) { // 若時間 < 下限時間
		dateTime = dateTime.Add(duration) // 多一小時
	}

	for ; !upper.Before(dateTime); dateTime = dateTime.Add(duration) { // 每小時
		mongoDB.repsertOneHourlyRecord(bson.M{`time`: dateTime}, mongoDB.CreateHourlyRecord(dateTime).PrimitiveM()) // 代添紀錄
	}

}
