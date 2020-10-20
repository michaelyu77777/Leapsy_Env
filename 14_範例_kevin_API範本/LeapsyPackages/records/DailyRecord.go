package records

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DailyRecord - 小時紀錄
type DailyRecord struct {
	Time                       time.Time // 記錄時間
	PMKwhThisMonth, PMKwhToday int       // 這個月、今日PM累計電度
}

// PrimitiveM - 轉成primitive.M
/*
 * @return primitive.M returnPrimitiveM 回傳結果
 */
func (dailyRecord DailyRecord) PrimitiveM() (returnPrimitiveM primitive.M) {

	returnPrimitiveM = bson.M{
		`time`:           dailyRecord.Time,
		`pmkwhthismonth`: dailyRecord.PMKwhThisMonth,
		`pmkwhtoday`:     dailyRecord.PMKwhToday,
	}

	return
}
