package databases

import (
	"context"
	"fmt"

	"../logings"
	"../network"
	"../records"
)

// PeriodicallyInsert - 週期性添寫秒記錄
/**
 * @param  records.SecondRecord secondRecord  秒記錄
 */
func (mongoDB *MongoDB) PeriodicallyInsert(secondRecord records.SecondRecord) {

	defer mongoDB.Disconnect()

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		mongoDB.GetConfigValueOrPanic(`server`),
		mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	periodicallyRWMutex.Lock() // 寫鎖

	// 添寫秒記錄
	_, insertOneError := mongoDB.Connect().
		Database(mongoDB.GetConfigValueOrPanic(`database`)).
		Collection(mongoDB.GetConfigValueOrPanic(`second-table`)).
		InsertOne(context.TODO(), secondRecord)

	periodicallyRWMutex.Unlock() // 寫解鎖

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 添寫秒記錄`},
		network.GetAliasAddressPair(address),
		insertOneError,
	)

	if nil != insertOneError { // 若添寫秒紀錄錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若添寫秒記錄成功
		go logger.Infof(formatString, args...) // 記錄資訊
	}

}
