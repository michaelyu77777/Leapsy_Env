package servers

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"../databases"
	"../logings"
)

// ECServer - 環控伺服器
type ECServer struct {
	server *http.Server // 伺服器指標
}

// StartECServer - 啟動環控伺服器
func StartECServer() {

	var (
		eCAPIServer                                                    ECAPIServer          // 環控API伺服器
		ecsDB                                                          databases.ECSDB      // 來源資料庫
		ecsAlertDB                                                     databases.ECSAlertDB // 警報來源資料庫
		periodicallyMongoDB, hourlyMongoDB, dailyMongoDB, alertMongoDB databases.MongoDB    // 記錄用資料庫
	)

	signalChannel := make(chan os.Signal) // 結束信號通道
	noticeChannel := make(chan bool)      // 通知通道

	defer func() {
		close(signalChannel)                                 // 關閉信號通道
		close(noticeChannel)                                 // 關閉通知通道
		eCAPIServer.stop()                                   // 結束環控API伺服器
		stopPeriodicallyRecord(&ecsDB, &periodicallyMongoDB) // 結束周期性記錄
		stopHourlyRecord(&hourlyMongoDB)                     // 結束每時記錄
		stopDailyRecord(&dailyMongoDB)                       // 結束每日記錄
		stopRecordAlerts(&ecsAlertDB, &alertMongoDB)         // 結束記錄警報
		StopECServer()                                       // 結束環控伺服器
	}()

	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM) // 轉傳結束信號

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`啟動 環控伺服器 `},
		[]interface{}{},
		nil,
	)

	go logger.Infof(formatString, args...) // 記錄資訊

	go eCAPIServer.start() // 啟動環控API伺服器

	go startPeriodicallyRecord(&ecsDB, &periodicallyMongoDB) // 開始週期性記錄

	go startHourlyRecord(&hourlyMongoDB) // 開始每時記錄

	go startDailyRecord(&dailyMongoDB) // 開始每日記錄

	go startRecordAlerts(&ecsAlertDB, &alertMongoDB) // 開始記錄警報

	go catchSignal(signalChannel, noticeChannel) // 抓取信號並通知

	<-noticeChannel // 抓取通知

}

// StopECServer - 結束環控伺服器
func StopECServer() {
	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`結束 環控伺服器 `},
		[]interface{}{},
		nil,
	)

	logger.Infof(formatString, args...) // 記錄資訊
}

/**
 * catchSignal - 抓取信號並通知
 *
 * @param  chan os.Signal signalChannel  信號通道
 * @param  chan bool noticeChannel       通知通道
 */
func catchSignal(signalChannel chan os.Signal, noticeChannel chan bool) {
	<-signalChannel       // 抓取信號
	noticeChannel <- true // 通知
}
