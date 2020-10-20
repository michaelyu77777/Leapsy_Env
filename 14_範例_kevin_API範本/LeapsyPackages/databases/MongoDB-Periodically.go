package databases

import (
	"reflect"
	"time"
)

// SumIntSecondRecordFields - 計算兩時間內的整數型欄位和
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @param []string fieldNames 欄位名
 * @return map[string]int results 計算結果
 */
func (mongoDB *MongoDB) SumIntSecondRecordFields(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
	fieldNames ...string,
) (results map[string]int) {

	results = make(map[string]int) // 指派空間

	for _, secondRecord := range mongoDB.FindSecondRecordsBetweenTimes(lowerTime, isLowerTimeIncluded, upperTime, isUpperTimeIncluded) { // 每一紀錄

		secondRecordReflect := reflect.ValueOf(secondRecord) // 紀錄

		for _, fieldName := range fieldNames { // 每一欄位

			if secondRecordReflect.FieldByName(fieldName).Kind() == reflect.Int { // 若為整數型欄位
				results[fieldName] += int(secondRecordReflect.FieldByName(fieldName).Int()) // 加上欄位值
			} else { // 若非整數型欄位
				results[fieldName] = 0 // 結果為零
			}

		}

	}

	return // 回傳
}

// SubtractIntSecondRecordFields - 計算兩時間內的整數型欄位相減
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @param []string fieldNames 欄位名
 * @return map[string]int results 計算結果
 */
func (mongoDB *MongoDB) SubtractIntSecondRecordFields(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
	fieldNames ...string,
) (results map[string]int) {

	results = make(map[string]int) // 指派空間

	secondRecords := mongoDB.FindSecondRecordsBetweenTimes(lowerTime, isLowerTimeIncluded, upperTime, isUpperTimeIncluded) // 查找紀錄
	lengthOfSecondRecords := len(secondRecords)                                                                            // 紀錄長度

	if lengthOfSecondRecords > 1 { // 若查找有紀錄

		firstSecondRecordReflect := reflect.ValueOf(secondRecords[0])                      // 第一筆紀錄
		lastSecondRecordReflect := reflect.ValueOf(secondRecords[lengthOfSecondRecords-1]) // 最末筆紀錄

		for _, fieldName := range fieldNames { // 針對每一欄位

			if firstSecondRecordReflect.FieldByName(fieldName).Kind() == reflect.Int { // 若為整數型欄位
				// 末紀錄欄位值 - 首紀錄欄位值
				results[fieldName] = int(
					lastSecondRecordReflect.FieldByName(fieldName).Int() -
						firstSecondRecordReflect.FieldByName(fieldName).Int(),
				)
			} else { // 若非整數型欄位
				results[fieldName] = 0 // 結果為零
			}

		}

	}

	return // 回傳
}

//  isSecondDateTime - 是否整秒時間
/**
 * @param time.Time dateTime 時間
 * @return bool result 回傳結果
 */
func isSecondDateTime(dateTime time.Time) (result bool) {

	result = 0 == dateTime.Second() && 0 == dateTime.Nanosecond() // 整秒

	return
}
