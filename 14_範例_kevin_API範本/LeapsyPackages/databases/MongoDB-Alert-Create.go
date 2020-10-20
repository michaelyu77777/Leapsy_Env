package databases

import (
	"context"
	"fmt"

	"../logings"
	"../network"
	"../records"
)

// InsertAlertRecord - 添寫警報記錄
/**
 * @param  records.AlertRecord alertRecord  警報記錄
 */
func (mongoDB *MongoDB) InsertAlertRecord(alertRecord records.AlertRecord) {

	defer mongoDB.Disconnect()

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		mongoDB.GetConfigValueOrPanic(`server`),
		mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	alertRWMutex.Lock() // 寫鎖

	// 添寫警報記錄
	_, insertOneError := mongoDB.Connect().
		Database(mongoDB.GetConfigValueOrPanic(`database`)).
		Collection(mongoDB.GetConfigValueOrPanic(`alert-table`)).
		InsertOne(context.TODO(), alertRecord)

	alertRWMutex.Unlock() // 寫解鎖

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 添寫警報記錄`},
		network.GetAliasAddressPair(address),
		insertOneError,
	)

	if nil != insertOneError { // 若添寫警報紀錄錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若添寫警報記錄成功
		go logger.Infof(formatString, args...) // 記錄資訊
	}

}
