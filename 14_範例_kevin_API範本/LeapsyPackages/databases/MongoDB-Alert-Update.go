package databases

import (
	"context"
	"fmt"

	"../logings"
	"../network"
	"../records"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// findOneAndUpdateAlertRecord - 更添警報記錄
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @param ...*options.FindOneAndUpdateOptions 選項
 * @return *mongo.SingleResult returnSingleResultPointer 更添結果
 */
func (mongoDB *MongoDB) findOneAndUpdateAlertRecord(
	filter, update primitive.M,
	opts ...*options.FindOneAndUpdateOptions) (returnSingleResultPointer *mongo.SingleResult) {

	defer mongoDB.Disconnect()

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		mongoDB.GetConfigValueOrPanic(`server`),
		mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	alertRWMutex.Lock() // 寫鎖

	// 更新警報記錄
	singleResultPointer := mongoDB.Connect().
		Database(mongoDB.GetConfigValueOrPanic(`database`)).
		Collection(mongoDB.GetConfigValueOrPanic(`alert-table`)).
		FindOneAndUpdate(
			context.TODO(),
			filter,
			update,
			opts...,
		)

	alertRWMutex.Unlock() // 寫解鎖

	findOneAndUpdateError := singleResultPointer.Err() // 更添錯誤

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 更添警報記錄`},
		network.GetAliasAddressPair(address),
		findOneAndUpdateError,
	)

	if nil != findOneAndUpdateError && mongo.ErrNoDocuments != findOneAndUpdateError { // 若更添警報紀錄錯誤且非檔案不存在錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若更新警報記錄成功
		go logger.Infof(formatString, args...)          // 記錄資訊
		returnSingleResultPointer = singleResultPointer // 回傳結果指標
	}

	return
}

// UpdateOneAlertRecord - 更新警報記錄
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @return *mongo.UpdateResult returnUpdateResult 更新結果
 */
func (mongoDB *MongoDB) UpdateOneAlertRecord(filter, update primitive.M) (results []records.AlertRecord) {

	var updatedAlertRecord records.AlertRecord // 更新的紀錄

	if nil == mongoDB.findOneAndUpdateAlertRecord(filter, update).Decode(&updatedAlertRecord) { // 若更新沒錯誤
		results = append(results, updatedAlertRecord) // 回傳結果
	}

	return
}
