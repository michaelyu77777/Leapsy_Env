package databases

import (
	"fmt"

	"../jsons"
	"../logings"
	"../network"
	"../records"
)

//  getECSRecord - 取得環控紀錄
/**
 * @param  string sqlCommand SQL指令
 * @return records.ECSRecord ecsRecord 紀錄
 */
func (eCSDB *ECSDB) getECSRecord(sqlCommand string) (ecsRecord records.ECSRecord) {

	defer eCSDB.Disconnect()

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		eCSDB.GetConfigValueOrPanic(`server`),
		eCSDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	// 審視紀錄
	scanError := eCSDB.Connect().
		QueryRow(sqlCommand).
		Scan(
			&ecsRecord.RTEXPDTIME,
			&ecsRecord.CNC1, &ecsRecord.CNC2, &ecsRecord.CNC3, &ecsRecord.CNC4, &ecsRecord.CNC5,
			&ecsRecord.CNC6, &ecsRecord.CNC7, &ecsRecord.CNC8, &ecsRecord.CNC9, &ecsRecord.CNC10,
			&ecsRecord.CNC11, &ecsRecord.CNC12, &ecsRecord.CNC13, &ecsRecord.CNC14, &ecsRecord.CNC15,
			&ecsRecord.CNC16, &ecsRecord.CNC17, &ecsRecord.CNC18, &ecsRecord.CNC19, &ecsRecord.CNC20,
			&ecsRecord.CNC21, &ecsRecord.CNC22, &ecsRecord.CNC23, &ecsRecord.CNC24, &ecsRecord.CNC25,
			&ecsRecord.CNC26, &ecsRecord.CNC27, &ecsRecord.CNC28, &ecsRecord.CNC29, &ecsRecord.CNC30,
			&ecsRecord.CNC31, &ecsRecord.CNC32, &ecsRecord.CNC33, &ecsRecord.CNC34, &ecsRecord.CNC35,
			&ecsRecord.CNC36, &ecsRecord.CNC37, &ecsRecord.CNC38, &ecsRecord.CNC39, &ecsRecord.CNC40,
			&ecsRecord.CNC41, &ecsRecord.M42, &ecsRecord.M43, &ecsRecord.M44, &ecsRecord.M45,
			&ecsRecord.M46, &ecsRecord.M47, &ecsRecord.M48, &ecsRecord.M49, &ecsRecord.M50,
			&ecsRecord.M51, &ecsRecord.M52, &ecsRecord.M53, &ecsRecord.M54, &ecsRecord.M55,
			&ecsRecord.M56, &ecsRecord.M57, &ecsRecord.M58, &ecsRecord.M59, &ecsRecord.M60,
			&ecsRecord.M61, &ecsRecord.M62, &ecsRecord.M63, &ecsRecord.M64, &ecsRecord.M65,
			&ecsRecord.M66, &ecsRecord.M67, &ecsRecord.M68, &ecsRecord.M69, &ecsRecord.M70,
			&ecsRecord.M71, &ecsRecord.M72, &ecsRecord.M73, &ecsRecord.M74, &ecsRecord.M75,
			&ecsRecord.M76, &ecsRecord.M77, &ecsRecord.M78, &ecsRecord.M79, &ecsRecord.M80,
			&ecsRecord.M81, &ecsRecord.M82, &ecsRecord.M83, &ecsRecord.M84, &ecsRecord.M85,
			&ecsRecord.M86, &ecsRecord.M87, &ecsRecord.M88, &ecsRecord.M89, &ecsRecord.M90,
			&ecsRecord.M91, &ecsRecord.M92, &ecsRecord.M93,
		)

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 審視環控紀錄`},
		network.GetAliasAddressPair(address),
		scanError,
	)

	if nil != scanError { // 若審視紀錄錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若審視紀錄成功
		go logger.Infof(formatString, args...) // 記錄資訊
	}

	return // 回傳
}

// Read - 讀一筆紀錄
/**
 * @return records.ECSRecord ecsRecord 紀錄
 */
func (eCSDB *ECSDB) Read() (ecsRecord records.ECSRecord) {

	ecsRecord = eCSDB.getECSRecord(
		fmt.Sprintf(
			`select * from %s`,
			eCSDB.GetConfigValueOrPanic(`table`),
		),
	)

	return // 回傳
}

// ReadJSONString - 讀一筆紀錄JSON字串
/**
 * @return string jsonString JSON字串
 */
func (eCSDB *ECSDB) ReadJSONString() (jsonString string) {

	jsonString = jsons.JSONString(eCSDB.Read()) // 回傳JSON字串

	return // 回傳
}
