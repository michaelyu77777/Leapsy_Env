package databases

import (
	"context"
	"fmt"
	"reflect"

	"../configurations"
	"../logings"
	"../network"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB - 資料庫
type MongoDB struct {
	clientPointer *mongo.Client // Mongo 客戶端指標
}

// GetConfigValueOrPanic - 取得設定值否則結束程式
/**
 * @param  string key  關鍵字
 * @return string 設定資料區塊下關鍵字對應的值
 */
func (mongoDB MongoDB) GetConfigValueOrPanic(key string) string {
	return configurations.GetConfigValueOrPanic(reflect.TypeOf(mongoDB).String(), key)
}

// GetConfigPositiveIntValueOrPanic - 取得正整數設定值否則結束程式
/**
 * @param  string key  關鍵字
 * @return string 設定資料區塊下關鍵字對應的正整數值
 */
func (mongoDB MongoDB) GetConfigPositiveIntValueOrPanic(key string) int {
	return configurations.GetConfigPositiveIntValueOrPanic(reflect.TypeOf(mongoDB).String(), key)
}

// Connect - 連接資料庫
/**
 * @return *mongo.Client mongoClientPointer 資料庫客戶端指標
 */
func (mongoDB *MongoDB) Connect() (returnMongoClientPointer *mongo.Client) {

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		mongoDB.GetConfigValueOrPanic(`server`),
		mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	network.SetAddressAlias(address, `環控歷史紀錄資料庫`) // 設定預設主機別名

	clientOptions := options.Client().ApplyURI(`mongodb://` + address)                    // 連線選項
	mongoClientPointer, mongoConnectError := mongo.Connect(context.TODO(), clientOptions) // 連接預設主機

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 連接`},
		network.GetAliasAddressPair(address),
		mongoConnectError,
	)

	if nil != mongoConnectError { // 若連接預設主機錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若連接預設主機成功
		go logger.Infof(formatString, args...) // 記錄資訊
	}

	mongoClientPointerPingError := mongoClientPointer.Ping(context.TODO(), nil) // 確認主機可連接

	// 取得記錄器格式和參數
	formatString, args = logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 連接`},
		network.GetAliasAddressPair(address),
		mongoClientPointerPingError,
	)

	if nil != mongoClientPointerPingError { // 若確認主機可連接錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若確認主機可連接成功
		go logger.Infof(formatString, args...)        // 記錄資訊
		mongoDB.clientPointer = mongoClientPointer    // 儲存資料庫指標
		returnMongoClientPointer = mongoClientPointer // 回傳資料庫指標
	}

	return // 回傳
}

// Disconnect - 中斷與資料庫的連線
/**
 * @return *mongo.Client mongoClientPointer 資料庫客戶端指標
 */
func (mongoDB *MongoDB) Disconnect() (returnMongoClientPointer *mongo.Client) {

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		mongoDB.GetConfigValueOrPanic(`server`),
		mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	mongoDBClientDisconnectError := mongoDB.clientPointer.Disconnect(context.TODO()) // 斷接主機

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 斷接`},
		network.GetAliasAddressPair(address),
		mongoDBClientDisconnectError,
	)

	if nil != mongoDBClientDisconnectError { // 若斷接主機錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若斷接主機成功
		go logger.Infof(formatString, args...)           // 記錄資訊
		returnMongoClientPointer = mongoDB.clientPointer //回傳 Mongo DB 客戶端
	}

	return // 回傳
}
