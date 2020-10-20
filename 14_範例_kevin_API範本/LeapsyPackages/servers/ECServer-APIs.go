package servers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"

	"../configurations"
	"../databases"
	"../jsons"
	"../logings"
	"../network"
	"../times"
)

// ECAPIServer - 環控API伺服器
type ECAPIServer struct {
	server *http.Server // 伺服器指標
}

// GetConfigValueOrPanic - 取得設定值否則結束程式
/**
 * @param  string key  關鍵字
 * @return string 設定資料區塊下關鍵字對應的值
 */
func (eCAPIServer *ECAPIServer) GetConfigValueOrPanic(key string) string {
	return configurations.GetConfigValueOrPanic(reflect.TypeOf(*eCAPIServer).String(), key) // 回傳取得的設定檔區塊下關鍵字對應的值
}

// GetConfigPositiveIntValueOrPanic - 取得設定整數值否則結束程式
/**
 * @param  string key  關鍵字
 * @return int 設定資料區塊下關鍵字對應的整數值
 */
func (eCAPIServer *ECAPIServer) GetConfigPositiveIntValueOrPanic(key string) int {
	return configurations.GetConfigPositiveIntValueOrPanic(reflect.TypeOf(*eCAPIServer).String(), key) // 回傳取得的設定檔區塊下關鍵字對應的值
}

// start - 啟動環控API伺服器
func (eCAPIServer *ECAPIServer) start() {
	signalChannel := make(chan os.Signal)                         // 結束信號通道
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM) // 轉傳結束信號
	defer close(signalChannel)

	port := eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`)             // 取得預設埠
	address := fmt.Sprintf(`%s:%d`, network.GetFirstLocalIPV4String(), port) // 預設主機

	network.SetAddressAlias(address, `環控API伺服器`) // 設定預設主機別名

	router := mux.NewRouter() // 新路由
	router.HandleFunc(
		`/`,
		func(httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {
			indexHandler(eCAPIServer, httpResponseWriter, httpRequestPointer)
		}) // 處理首頁
	router.HandleFunc(
		`/{year:\d{4}}/{month:\d{2}}/{day:\d{2}}`,
		func(httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {
			dailyAPIHandler(eCAPIServer, httpResponseWriter, httpRequestPointer)
		}) // 處理顯示一日內小時資料網頁
	router.HandleFunc(
		`/{year:\d{4}}/{month:\d{2}}`,
		func(httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {
			monthlyAPIHandler(eCAPIServer, httpResponseWriter, httpRequestPointer)
		}) // 處理顯示一月內日資料網頁
	router.HandleFunc(
		`/alerts`,
		func(httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {
			alertsAPIHandler(eCAPIServer, httpResponseWriter, httpRequestPointer)
		}) // 處理顯示一月內日資料網頁
	router.HandleFunc(
		`/setAlertRead/{alertEventID:\d+}`,
		func(httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {
			setAlertIsReadTrueAPIHandler(eCAPIServer, httpResponseWriter, httpRequestPointer)
		}) // 處理設定警報資料為已讀網頁
	router.HandleFunc(
		`/setAlertUnread/{alertEventID:\d+}`,
		func(httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {
			setAlertIsReadFalseAPIHandler(eCAPIServer, httpResponseWriter, httpRequestPointer)
		}) // 處理設定警報資料為未讀網頁
	router.HandleFunc(
		`/setAlertHidden/{alertEventID:\d+}`,
		func(httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {
			setAlertIsHiddenTrueAPIHandler(eCAPIServer, httpResponseWriter, httpRequestPointer)
		}) // 處理設定警報資料為已讀網頁
	router.HandleFunc(
		`/setAlertUnhidden/{alertEventID:\d+}`,
		func(httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {
			setAlertIsHiddenFalseAPIHandler(eCAPIServer, httpResponseWriter, httpRequestPointer)
		}) // 處理設定警報資料為未讀網頁

	apiServerPointer := &http.Server{Addr: address, Handler: router} // 設定伺服器
	eCAPIServer.server = apiServerPointer                            // 儲存伺服器指標

	var apiServerPtrListenAndServeError error // 伺服器啟動錯誤

	go func() {
		apiServerPtrListenAndServeError = apiServerPointer.ListenAndServe() // 啟動伺服器或回傳伺服器啟動錯誤
	}()

	<-time.After(time.Second * 3) // 等待伺服器啟動結果

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 啟動`},
		network.GetAliasAddressPair(address),
		apiServerPtrListenAndServeError,
	)

	if nil != apiServerPtrListenAndServeError { // 若伺服器啟動錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若伺服器啟動成功
		go logger.Infof(formatString, args...) // 記錄資訊
	}

	<-signalChannel
}

// stop - 結束環控API伺服器
func (eCAPIServer *ECAPIServer) stop() {

	port := eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`)             // 取得預設埠
	address := fmt.Sprintf(`%s:%d`, network.GetFirstLocalIPV4String(), port) // 預設主機

	eCAPIServerServerShutdownError := eCAPIServer.server.Shutdown(context.TODO()) // 結束伺服器

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%s %s 結束`},
		network.GetAliasAddressPair(address),
		eCAPIServerServerShutdownError,
	)

	if nil != eCAPIServerServerShutdownError { // 若伺服器結束錯誤
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若伺服器結束成功
		go logger.Infof(formatString, args...) // 記錄資訊
	}
}

// indexHandler - 處理首頁
/**
 * @param  *ECAPIServer eCAPIServer 環控API伺服器指標
 * @param  http.ResponseWriter httpResponseWriter  回應寫入器
 * @param  *http.Request httpRequestPointer        HTTP請求指標
 */
func indexHandler(eCAPIServer *ECAPIServer, httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {

	var (
		ecsDB databases.ECSDB // 來源資料庫
	)

	go func() {
		port := eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`)             // 取得預設埠
		address := fmt.Sprintf(`%s:%d`, network.GetFirstLocalIPV4String(), port) // 預設主機

		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			[]string{`%s %s 接受 %s 請求 %s `},
			append(network.GetAliasAddressPair(address), httpRequestPointer.RemoteAddr, httpRequestPointer.URL),
			nil,
		)

		logger.Infof(formatString, args...) // 記錄資訊
	}()

	fmt.Fprintf(httpResponseWriter, `%s`, `[`+jsons.JSONString(ecsDB.Read().SecondRecord())+`]`) // 寫入回應

}

// dailyAPIHandler - 處理顯示一日內小時資料網頁
/**
 * @param  *ECAPIServer eCAPIServer 環控API伺服器指標
 * @param  http.ResponseWriter httpResponseWriter  回應寫入器
 * @param  *http.Request httpRequestPointer        HTTP請求指標
 */
func dailyAPIHandler(eCAPIServer *ECAPIServer, httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {

	var (
		mongoDB databases.MongoDB // 資料庫
	)

	go func() {
		port := eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`)             // 取得預設埠
		address := fmt.Sprintf(`%s:%d`, network.GetFirstLocalIPV4String(), port) // 預設主機

		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			[]string{`%s %s 接受 %s 請求 %s `},
			append(network.GetAliasAddressPair(address), httpRequestPointer.RemoteAddr, httpRequestPointer.URL),
			nil,
		)

		logger.Infof(formatString, args...) // 記錄資訊
	}()

	variables := mux.Vars(httpRequestPointer)    // 取得網址參數
	year, _ := strconv.Atoi(variables[`year`])   // 取得年
	month, _ := strconv.Atoi(variables[`month`]) // 取得月
	day, _ := strconv.Atoi(variables[`day`])     // 取得日

	low, upper := times.GetDailyBounds(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)) // 取得上下限日

	if 24 != mongoDB.CountHourlyRecordsBetweenTimes(low, false, upper, true) { // 若缺資料

		duration := time.Hour         // 定義期間
		currentDateTime := time.Now() // 現在時間
		dateTime := low.Add(duration) // 從下限時間下一小時開始

		for ; 0 != mongoDB.CountHourlyRecordByTime(dateTime); dateTime = dateTime.Add(duration) { // 針對每一小時
		}

		if upper.After(currentDateTime) { // 若上限時間在現在時間之後
			mongoDB.RepsertHourlyRecordsBetweenTimes(dateTime, currentDateTime) // 代添紀錄到現在時間
		} else {
			mongoDB.RepsertHourlyRecordsBetweenTimes(dateTime, upper) // 代添紀錄到上限時間
		}

	}

	fmt.Fprintf(httpResponseWriter, "%s", jsons.JSONString(mongoDB.FindHourlyRecordsBetweenTimes(low, false, upper, true))) // 寫入回應
}

// monthlyAPIHandler - 處理顯示一月內日資料網頁
/**
 * @param  *ECAPIServer eCAPIServer 環控API伺服器指標
 * @param  http.ResponseWriter httpResponseWriter  回應寫入器
 * @param  *http.Request httpRequestPointer        HTTP請求指標
 */
func monthlyAPIHandler(eCAPIServer *ECAPIServer, httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {

	var (
		mongoDB databases.MongoDB // 資料庫
	)

	go func() {
		port := eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`)             // 取得預設埠
		address := fmt.Sprintf(`%s:%d`, network.GetFirstLocalIPV4String(), port) // 預設主機

		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			[]string{`%s %s 接受 %s 請求 %s `},
			append(network.GetAliasAddressPair(address), httpRequestPointer.RemoteAddr, httpRequestPointer.URL),
			nil,
		)

		logger.Infof(formatString, args...) // 記錄資訊
	}()

	variables := mux.Vars(httpRequestPointer)    // 取得網址參數
	year, _ := strconv.Atoi(variables[`year`])   // 取得年
	month, _ := strconv.Atoi(variables[`month`]) // 取得月

	low, upper := times.GetMonthlyBounds(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)) // 取得上下限月

	daysCount := int(upper.Sub(low).Hours()) / 24

	if daysCount != mongoDB.CountDailyRecordsBetweenTimes(low, false, upper, true) { // 若缺資料

		currentDateTime := time.Now()    // 現在時間
		dateTime := low.AddDate(0, 0, 1) // 從下限時間隔天開始

		for ; 0 != mongoDB.CountDailyRecordByTime(dateTime); dateTime = dateTime.AddDate(0, 0, 1) { // 針對每一日
		}

		if upper.After(currentDateTime) { // 若上限時間在現在時間之後
			mongoDB.RepsertDailyRecordsBetweenTimes(dateTime, currentDateTime) // 代添紀錄到現在時間
		} else {
			mongoDB.RepsertDailyRecordsBetweenTimes(dateTime, upper) // 代添紀錄到上限時間
		}

	}

	fmt.Fprintf(httpResponseWriter, "%s", jsons.JSONString(mongoDB.FindDailyRecordsBetweenTimes(low, false, upper, true))) //寫入回應
}

// alertsAPIHandler - 處理顯示警報資料網頁
/**
 * @param  *ECAPIServer eCAPIServer 環控API伺服器指標
 * @param  http.ResponseWriter httpResponseWriter  回應寫入器
 * @param  *http.Request httpRequestPointer        HTTP請求指標
 */
func alertsAPIHandler(eCAPIServer *ECAPIServer, httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {

	var (
		mongoDB databases.MongoDB // 資料庫
	)

	go func() {
		port := eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`)             // 取得預設埠
		address := fmt.Sprintf(`%s:%d`, network.GetFirstLocalIPV4String(), port) // 預設主機

		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			[]string{`%s %s 接受 %s 請求 %s `},
			append(network.GetAliasAddressPair(address), httpRequestPointer.RemoteAddr, httpRequestPointer.URL),
			nil,
		)

		logger.Infof(formatString, args...) // 記錄資訊
	}()

	fmt.Fprintf(httpResponseWriter, "%s", jsons.JSONString(mongoDB.FindAllAlertRecords())) //寫入回應
}

// setAlertIsReadTrueAPIHandler - 處理設定警報資料為已讀網頁
/**
 * @param  *ECAPIServer eCAPIServer 環控API伺服器指標
 * @param  http.ResponseWriter httpResponseWriter  回應寫入器
 * @param  *http.Request httpRequestPointer        HTTP請求指標
 */
func setAlertIsReadTrueAPIHandler(eCAPIServer *ECAPIServer, httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {

	var (
		mongoDB databases.MongoDB // 資料庫
	)

	go func() {
		port := eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`)             // 取得預設埠
		address := fmt.Sprintf(`%s:%d`, network.GetFirstLocalIPV4String(), port) // 預設主機

		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			[]string{`%s %s 接受 %s 請求 %s `},
			append(network.GetAliasAddressPair(address), httpRequestPointer.RemoteAddr, httpRequestPointer.URL),
			nil,
		)

		logger.Infof(formatString, args...) // 記錄資訊
	}()

	variables := mux.Vars(httpRequestPointer)                  // 取得網址參數
	alertEventID, _ := strconv.Atoi(variables[`alertEventID`]) // 取得警報編號

	fmt.Fprintf(
		httpResponseWriter, "%s", jsons.JSONString(
			mongoDB.UpdateOneAlertRecord(
				bson.M{
					`alerteventid`: alertEventID,
				},
				bson.M{
					`$set`: bson.M{
						`isread`: true,
					},
				},
			),
		),
	) //寫入回應
}

// setAlertIsReadFalseAPIHandler - 處理設定警報資料為未讀網頁
/**
 * @param  *ECAPIServer eCAPIServer 環控API伺服器指標
 * @param  http.ResponseWriter httpResponseWriter  回應寫入器
 * @param  *http.Request httpRequestPointer        HTTP請求指標
 */
func setAlertIsReadFalseAPIHandler(eCAPIServer *ECAPIServer, httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {

	var (
		mongoDB databases.MongoDB // 資料庫
	)

	go func() {
		port := eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`)             // 取得預設埠
		address := fmt.Sprintf(`%s:%d`, network.GetFirstLocalIPV4String(), port) // 預設主機

		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			[]string{`%s %s 接受 %s 請求 %s `},
			append(network.GetAliasAddressPair(address), httpRequestPointer.RemoteAddr, httpRequestPointer.URL),
			nil,
		)

		logger.Infof(formatString, args...) // 記錄資訊
	}()

	variables := mux.Vars(httpRequestPointer)                  // 取得網址參數
	alertEventID, _ := strconv.Atoi(variables[`alertEventID`]) // 取得警報編號

	fmt.Fprintf(
		httpResponseWriter, "%s", jsons.JSONString(
			mongoDB.UpdateOneAlertRecord(
				bson.M{
					`alerteventid`: alertEventID,
				},
				bson.M{
					`$set`: bson.M{
						`isread`: false,
					},
				},
			),
		),
	) //寫入回應
}

// setAlertIsHiddenTrueAPIHandler - 處理設定警報資料為隱藏網頁
/**
 * @param  *ECAPIServer eCAPIServer 環控API伺服器指標
 * @param  http.ResponseWriter httpResponseWriter  回應寫入器
 * @param  *http.Request httpRequestPointer        HTTP請求指標
 */
func setAlertIsHiddenTrueAPIHandler(eCAPIServer *ECAPIServer, httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {

	var (
		mongoDB databases.MongoDB // 資料庫
	)

	go func() {
		port := eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`)             // 取得預設埠
		address := fmt.Sprintf(`%s:%d`, network.GetFirstLocalIPV4String(), port) // 預設主機

		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			[]string{`%s %s 接受 %s 請求 %s `},
			append(network.GetAliasAddressPair(address), httpRequestPointer.RemoteAddr, httpRequestPointer.URL),
			nil,
		)

		logger.Infof(formatString, args...) // 記錄資訊
	}()

	variables := mux.Vars(httpRequestPointer)                  // 取得網址參數
	alertEventID, _ := strconv.Atoi(variables[`alertEventID`]) // 取得警報編號

	fmt.Fprintf(
		httpResponseWriter, "%s", jsons.JSONString(
			mongoDB.UpdateOneAlertRecord(
				bson.M{
					`alerteventid`: alertEventID,
				},
				bson.M{
					`$set`: bson.M{
						`ishidden`: true,
					},
				},
			),
		),
	) //寫入回應
}

// setAlertIsHiddenFalseAPIHandler - 處理設定警報資料為非隱藏網頁
/**
 * @param  *ECAPIServer eCAPIServer 環控API伺服器指標
 * @param  http.ResponseWriter httpResponseWriter  回應寫入器
 * @param  *http.Request httpRequestPointer        HTTP請求指標
 */
func setAlertIsHiddenFalseAPIHandler(eCAPIServer *ECAPIServer, httpResponseWriter http.ResponseWriter, httpRequestPointer *http.Request) {

	var (
		mongoDB databases.MongoDB // 資料庫
	)

	go func() {
		port := eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`)             // 取得預設埠
		address := fmt.Sprintf(`%s:%d`, network.GetFirstLocalIPV4String(), port) // 預設主機

		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			[]string{`%s %s 接受 %s 請求 %s `},
			append(network.GetAliasAddressPair(address), httpRequestPointer.RemoteAddr, httpRequestPointer.URL),
			nil,
		)

		logger.Infof(formatString, args...) // 記錄資訊
	}()

	variables := mux.Vars(httpRequestPointer)                  // 取得網址參數
	alertEventID, _ := strconv.Atoi(variables[`alertEventID`]) // 取得警報編號

	fmt.Fprintf(
		httpResponseWriter, "%s", jsons.JSONString(
			mongoDB.UpdateOneAlertRecord(
				bson.M{
					`alerteventid`: alertEventID,
				},
				bson.M{
					`$set`: bson.M{
						`ishidden`: false,
					},
				},
			),
		),
	) //寫入回應
}
