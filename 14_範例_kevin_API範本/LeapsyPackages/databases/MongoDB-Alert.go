package databases

import (
	"reflect"
	"time"
)

// SumIntAlertRecordFields - 計算兩時間內的整數型欄位和
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @param []string fieldNames 欄位名
 * @return map[string]int results 計算結果
 */
func (mongoDB *MongoDB) SumIntAlertRecordFields(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
	fieldNames ...string,
) (results map[string]int) {

	results = make(map[string]int) // 指派空間

	for _, alertRecord := range mongoDB.FindAlertRecordsBetweenTimes(lowerTime, isLowerTimeIncluded, upperTime, isUpperTimeIncluded) { // 每一紀錄

		alertRecordReflect := reflect.ValueOf(alertRecord) // 紀錄

		for _, fieldName := range fieldNames { // 每一欄位

			if alertRecordReflect.FieldByName(fieldName).Kind() == reflect.Int { // 若為整數型欄位
				results[fieldName] += int(alertRecordReflect.FieldByName(fieldName).Int()) // 加上欄位值
			} else { // 若非整數型欄位
				results[fieldName] = 0 // 結果為零
			}

		}

	}

	return // 回傳
}

// SubtractIntAlertRecordFields - 計算兩時間內的整數型欄位相減
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @param []string fieldNames 欄位名
 * @return map[string]int results 計算結果
 */
func (mongoDB *MongoDB) SubtractIntAlertRecordFields(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
	fieldNames ...string,
) (results map[string]int) {

	results = make(map[string]int) // 指派空間

	alertRecords := mongoDB.FindAlertRecordsBetweenTimes(lowerTime, isLowerTimeIncluded, upperTime, isUpperTimeIncluded) // 查找紀錄
	lengthOfAlertRecords := len(alertRecords)                                                                            // 紀錄長度

	if lengthOfAlertRecords > 1 { // 若查找有紀錄

		firstAlertRecordReflect := reflect.ValueOf(alertRecords[0])                     // 第一筆紀錄
		lastAlertRecordReflect := reflect.ValueOf(alertRecords[lengthOfAlertRecords-1]) // 最末筆紀錄

		for _, fieldName := range fieldNames { // 針對每一欄位

			if firstAlertRecordReflect.FieldByName(fieldName).Kind() == reflect.Int { // 若為整數型欄位
				// 末紀錄欄位值 - 首紀錄欄位值
				results[fieldName] = int(
					lastAlertRecordReflect.FieldByName(fieldName).Int() -
						firstAlertRecordReflect.FieldByName(fieldName).Int(),
				)
			} else { // 若非整數型欄位
				results[fieldName] = 0 // 結果為零
			}

		}

	}

	return // 回傳
}
