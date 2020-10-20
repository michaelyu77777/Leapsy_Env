package databases

import (
	"fmt"

	"../jsons"
	"../logings"
	"../network"
	"../records"
)

// CountAll - 計算所有紀錄個數
/**
 * @return int returnCount 紀錄個數
 */
func (ecsAlertDB *ECSAlertDB) CountAll() (returnCount int) {

	defer ecsAlertDB.Disconnect()

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		ecsAlertDB.GetConfigValueOrPanic(`server`),
		ecsAlertDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	// 查詢紀錄數
	row := ecsAlertDB.Connect().QueryRow(
		fmt.Sprintf(
			`select count(*) from %s`,
			ecsAlertDB.GetConfigValueOrPanic(`table`),
		),
	)

	// 審視環控警報紀錄個數
	scanError := row.Scan(&returnCount)

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 審視環控警報紀錄個數`},
		network.GetAliasAddressPair(address),
		scanError,
	)

	if nil != scanError { // 若審視環控警報紀錄個數錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若審視環控警報紀錄個數成功
		go logger.Infof(formatString, args...) // 記錄資訊
	}

	return // 回傳
}

// Read - 讀一筆紀錄
/**
 * @return []records.ECSAlertRecord ecsAlertRecord 紀錄
 */
func (ecsAlertDB *ECSAlertDB) Read() (ecsAlertRecords []records.ECSAlertRecord) {

	ecsAlertRecords = ecsAlertDB.getRecords(
		fmt.Sprintf(
			`select * from %s`,
			ecsAlertDB.GetConfigValueOrPanic(`table`),
		),
	)

	return // 回傳
}

// ReadLast - 讀末N筆紀錄
/**
 * @param int n 紀錄個數
 * @return []records.ECSAlertRecord ecsAlertRecord 紀錄
 */
func (ecsAlertDB *ECSAlertDB) ReadLast(n int) (ecsAlertRecords []records.ECSAlertRecord) {

	ecsAlertRecords = ecsAlertDB.getRecords(
		fmt.Sprintf(
			`select top %d * from %s order by ALERTEVENTID desc`,
			n,
			ecsAlertDB.GetConfigValueOrPanic(`table`),
		),
	)

	return // 回傳
}

//  getRecords - 取得紀錄
/**
 * @param  string sqlCommand SQL指令
 * @return []records.ECSAlertRecord record 紀錄
 */
func (ecsAlertDB *ECSAlertDB) getRecords(sqlCommand string) (ecsAlertRecords []records.ECSAlertRecord) {

	defer ecsAlertDB.Disconnect()

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		ecsAlertDB.GetConfigValueOrPanic(`server`),
		ecsAlertDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	// 查詢紀錄
	rows, queryError := ecsAlertDB.Connect().Query(sqlCommand)
	defer rows.Close()

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 查詢紀錄`},
		network.GetAliasAddressPair(address),
		queryError,
	)

	if nil != queryError { // 若查詢紀錄錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若查詢紀錄成功
		go logger.Infof(formatString, args...) // 記錄資訊
	}

	for rows.Next() {

		var ecsAlertRecord records.ECSAlertRecord

		// 審視環控警報紀錄
		scanError := rows.Scan(
			&ecsAlertRecord.ALERTEVENTID,
			&ecsAlertRecord.ALERTEVENTTIME,
			&ecsAlertRecord.VARTAG,
			&ecsAlertRecord.COMMENT,
			&ecsAlertRecord.ALERTTYPE,
			&ecsAlertRecord.LINETEXT,
		)

		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			[]string{`%s %s 審視環控警報紀錄`},
			network.GetAliasAddressPair(address),
			scanError,
		)

		if nil != scanError { // 若審視環控警報紀錄錯誤
			logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
		} else { // 若審視環控紀錄成功
			go logger.Infof(formatString, args...) // 記錄資訊
			ecsAlertRecords = append(ecsAlertRecords, ecsAlertRecord)
		}
	}

	return // 回傳
}

// ReadJSONString - 讀警報紀錄JSON字串
/**
 * @return string jsonString JSON字串
 */
func (ecsAlertDB *ECSAlertDB) ReadJSONString() (jsonString string) {

	jsonString = jsons.JSONString(ecsAlertDB.Read()) // 回傳JSON字串

	return // 回傳
}
