package jsons

import (
	"encoding/json"

	"../logings"
)

var (
	logger = logings.GetLogger() // 記錄器
)

// JSONString - 取得JSON字串
/**
 * @param interface{} inputObject 輸入物件
 * @return string returnJSONString 取得JSON字串
 */
func JSONString(inputObject interface{}) (returnJSONString string) {

	jsonBytes, jsonMarshalError := json.Marshal(inputObject) // 轉成JSON

	if nil != jsonMarshalError { // 若轉成JSON錯誤

		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			[]string{`轉成JSON物件`},
			[]interface{}{},
			jsonMarshalError,
		)

		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else {
		returnJSONString = string(jsonBytes) // 回傳JSON字串
	}

	return // 回傳
}
