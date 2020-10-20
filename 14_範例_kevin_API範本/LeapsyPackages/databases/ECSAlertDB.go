package databases

import (
	"database/sql"
	"fmt"
	"reflect"

	"../configurations"
	"../logings"
	"../network"
)

// ECSAlertDB - 環控系統警報資料庫
type ECSAlertDB struct {
	db *sql.DB
}

// GetConfigValueOrPanic - 取得設定值否則結束程式
/**
 * @param  string key  關鍵字
 * @return string 設定資料區塊下關鍵字對應的值
 */
func (ecsAlertDB ECSAlertDB) GetConfigValueOrPanic(key string) string {
	return configurations.GetConfigValueOrPanic(reflect.TypeOf(ecsAlertDB).String(), key)
}

// GetConfigPositiveIntValueOrPanic - 取得正整數設定值否則結束程式
/**
 * @param  string key  關鍵字
 * @return string 設定資料區塊下關鍵字對應的正整數值
 */
func (ecsAlertDB ECSAlertDB) GetConfigPositiveIntValueOrPanic(key string) int {
	return configurations.GetConfigPositiveIntValueOrPanic(reflect.TypeOf(ecsAlertDB).String(), key)
}

// Connect - 連接資料庫
/**
 * @return *sql.DB returnDB 資料庫指標
 */
func (ecsAlertDB *ECSAlertDB) Connect() (returnDB *sql.DB) {

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		ecsAlertDB.GetConfigValueOrPanic(`server`),
		ecsAlertDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	network.SetAddressAlias(address, `環控警報資料庫`) // 設定預設主機別名

	// 連接預設主機
	db, sqlOpenError := sql.Open(
		"mssql",
		fmt.Sprintf(
			`server=%s;user id=%s;password=%s;database=%s`,
			ecsAlertDB.GetConfigValueOrPanic(`server`),
			ecsAlertDB.GetConfigValueOrPanic(`userid`),
			ecsAlertDB.GetConfigValueOrPanic(`password`),
			ecsAlertDB.GetConfigValueOrPanic(`database`),
		),
	)

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 連接`},
		network.GetAliasAddressPair(address),
		sqlOpenError,
	)

	if nil != sqlOpenError { // 若連接預設主機錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若連接預設主機成功
		go logger.Infof(formatString, args...) // 記錄資訊
	}

	dbPingError := db.Ping() // 確認主機可連接

	// 取得記錄器格式和參數
	formatString, args = logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 連接`},
		network.GetAliasAddressPair(address),
		dbPingError,
	)

	if nil != dbPingError { // 若確認主機可連接錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若確認主機可連接成功
		go logger.Infof(formatString, args...) // 記錄資訊
		ecsAlertDB.db = db                     // 儲存資料庫指標
		returnDB = db                          // 回傳資料庫指標
	}

	return // 回傳
}

// Disconnect - 中斷與資料庫的連線
/**
 * @return *sql.DB returnDB 資料庫指標
 */
func (ecsAlertDB *ECSAlertDB) Disconnect() (returnDB *sql.DB) {

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		ecsAlertDB.GetConfigValueOrPanic(`server`),
		ecsAlertDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	dbCloseError := ecsAlertDB.db.Close() // 斷接主機

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 斷接`},
		network.GetAliasAddressPair(address),
		dbCloseError,
	)

	if nil != dbCloseError { // 若斷接主機錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若斷接主機成功
		go logger.Infof(formatString, args...) // 記錄資訊
	}

	returnDB = ecsAlertDB.db //回傳資料庫指標

	return // 回傳
}
