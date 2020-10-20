package records

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HourlyRecord - 小時紀錄
type HourlyRecord struct {
	Time                      time.Time // 記錄時間
	PMKwhThisHour, PMKwhToday int       // 這小時、今日PM累計電度
}

// PrimitiveM - 轉成primitive.M
/*
 * @return primitive.M returnPrimitiveM 回傳結果
 */
func (hourlyRecord HourlyRecord) PrimitiveM() (returnPrimitiveM primitive.M) {

	returnPrimitiveM = bson.M{
		`time`:          hourlyRecord.Time,
		`pmkwhthishour`: hourlyRecord.PMKwhThisHour,
		`pmkwhtoday`:    hourlyRecord.PMKwhToday,
	}

	return
}
