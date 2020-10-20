package times

import (
	"time"
)

// GetHourlyBounds - 取得時間上下限小時
/**
 * @param  time.Time dateTime  時間
 * @return time.Time lower 下限小時
 * @return time.Time upper 上限小時
 */
func GetHourlyBounds(dateTime time.Time) (lower time.Time, upper time.Time) {

	duration := time.Hour

	if IsHour(dateTime) {
		lower = dateTime.Add(-duration)
		upper = dateTime
	} else {
		lower = time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), dateTime.Hour(), 0, 0, 0, time.Local)
		upper = lower.Add(duration)
	}

	lower = time.Date(lower.Year(), lower.Month(), lower.Day(), lower.Hour(), 0, 0, 0, time.Local)
	upper = time.Date(upper.Year(), upper.Month(), upper.Day(), upper.Hour(), 0, 0, 0, time.Local)

	return
}

// IsHour - 判斷時間是否為整點
/**
 * @param  time.Time dateTime  時間
 * @return bool 判斷是否為整點
 */
func IsHour(dateTime time.Time) bool {
	return dateTime.Minute() == 0 && dateTime.Second() == 0 && dateTime.Nanosecond() == 0 // 分秒微秒為零
}

// GetDailyBounds - 取得時間上下限日
/**
 * @param  time.Time dateTime  時間
 * @return  time.Time lower  下限日
 * @return  time.Time upper  上限日
 */
func GetDailyBounds(dateTime time.Time) (lower time.Time, upper time.Time) {

	if IsDay(dateTime) { // 若為整日
		lower = dateTime.AddDate(0, 0, -1) // 昨日
		upper = dateTime                   // 今日
	} else {
		lower = time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, time.Local) // 今日
		upper = lower.AddDate(0, 0, 1)                                                               // 明日
	}

	return // 回傳
}

// IsDay - 判斷是否為整日
/**
 * @param  time.Time dateTime  時間
 * @return  bool 判斷是否為整日
 */
func IsDay(dateTime time.Time) bool {
	return dateTime.Hour() == 0 && dateTime.Minute() == 0 && dateTime.Second() == 0 && dateTime.Nanosecond() == 0 // 時分秒微秒為零
}

// GetMonthlyBounds - 取得時間上下限月
/**
 * @param  time.Time dateTime  時間
 * @return  time.Time lower  下限月
 * @return  time.Time upper  上限月
 */
func GetMonthlyBounds(dateTime time.Time) (low time.Time, upper time.Time) {

	if IsMonth(dateTime) { // 若為整月
		low = dateTime.AddDate(0, -1, 0) // 上個月
		upper = dateTime                 // 這個月
	} else {
		low = time.Date(dateTime.Year(), dateTime.Month(), 1, 0, 0, 0, 0, time.Local) // 這個月
		upper = low.AddDate(0, 1, 0)                                                  // 下個月
	}

	return // 回傳
}

// IsMonth - 判斷是否為整月
/**
 * @param  time.Time dateTime  時間
 * @return  bool 判斷是否為整月
 */
func IsMonth(dateTime time.Time) bool {
	// 日為一、時分秒微秒為零
	return dateTime.Day() == 1 && dateTime.Hour() == 0 && dateTime.Minute() == 0 && dateTime.Second() == 0 && dateTime.Nanosecond() == 0
}
