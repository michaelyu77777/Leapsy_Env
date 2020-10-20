package network

import (
	"net"
	"regexp"
	"strings"
	"sync"

	"../logings"
)

var (
	logger = logings.GetLogger() // 記錄器

	addressToAliasMap = make(map[string]string) // 網址對應別名
	readWriteLock     = new(sync.RWMutex)       // 讀寫鎖
)

// GetFirstLocalIPV4String - 取得本機第一個IPV4字串
/**
 * @return string firstLocalIPV4String 第一個本地IPV4字串
 */
func GetFirstLocalIPV4String() (firstLocalIPV4String string) {

	localIPV4Strings := getLocalIPV4Strings() // 取得本地IPV4字串

	if len(localIPV4Strings) > 0 { // 若取得本地IPV4字串成功，則回傳第一個本地IP字串
		firstLocalIPV4String = localIPV4Strings[0]
	} else { // 若取得本地IPV4字串成功，則回傳迴環位址
		firstLocalIPV4String = `localhost`
	}

	return // 回傳
}

// getLocalIPV4Strings - 取得本機IPV4字串
/**
 * @return []string localIPV4Strings 本地IPV4字串陣列
 */
func getLocalIPV4Strings() (localIPV4Strings []string) {

	addrs, netInterfaceAddrsError := net.InterfaceAddrs() // 取得單播介面位址

	formatStringSlices := []string{`取得單播介面位址`} // 記錄器格式片段
	defaultArgs := []interface{}{}             // 記錄器預設參數

	// 取得記錄器格式與參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		formatStringSlices,
		defaultArgs,
		netInterfaceAddrsError,
	)

	if nil != netInterfaceAddrsError { // 若取得單播介面位址錯誤，則記錄錯誤並逐層結束程式
		logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
	} else { // 若取得單播介面位址錯誤，則回傳加入IPV4字串

		for _, addr := range addrs { // 針對每一個單播介面位址

			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && nil != ipNet.IP.To4() { // 若單播介面位址轉型別成功、非迴環且為IPV4
				localIPV4Strings = append(localIPV4Strings, ipNet.IP.String()) // 回傳加入此單播介面位址的IPV4字串
			}

		}

	} // end else nil != netInterfaceAddrsError

	return // 回傳
}

// LookupHostString - 查找主機
/**
 * @param  string hostString  主機
 * @return []string results 查找主機結果
 */
func LookupHostString(hostString string) (results []string) {

	// 主機正規表示式
	hostStringRegularExpression :=
		`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]):\d{1,5}$`

	if regexp.MustCompile(hostStringRegularExpression).MatchString(hostString) { // 若主機符合格式

		hostSlices := strings.Split(hostString, `:`) // 將主機名稱切開

		formatStringSlices := []string{`查找主機 %s `} // 記錄器格式片段
		defaultArgs := []interface{}{hostString}   // 記錄器預設參數

		addrs, netLookupHostError := net.LookupHost(hostSlices[0]) // 查找主機

		// 取得記錄器格式和參數
		formatString, args := logings.GetLogFuncFormatAndArguments(
			formatStringSlices,
			defaultArgs,
			netLookupHostError,
		)

		go logger.Infof(formatString, args...) // 記錄資訊

		if nil == netLookupHostError { // 若查找主機成功

			for _, addr := range addrs { // 針對每一結果

				updatedHostString := addr + `:` + hostSlices[1] // 結果加埠

				if _, ok := addressToAliasMap[updatedHostString]; ok { // 若結果加埠存在別名
					results = append(results, updatedHostString) // 將結果加埠加入回傳
				}

			}

		}
	} else { // 若主機不符合格式，則記錄資訊
		go logger.Infof(`主機不符合格式"[hostname]:[port]": %s`, hostString)
	}

	return // 回傳
}

// GetAddressAlias - 取得位址別名
/**
 * @param  string addressString 位址字串
 * @return string 位址別名
 */
func GetAddressAlias(addressString string) string {
	readWriteLock.RLock()                   // 讀鎖
	defer readWriteLock.RUnlock()           // 記得解開讀鎖
	return addressToAliasMap[addressString] // 回傳位址別名
}

// SetAddressAlias - 設定位址別名
/**
 * @param  string addressString 位址字串
 * @param  string aliasString 別名字串
 */
func SetAddressAlias(addressString, aliasString string) {
	readWriteLock.Lock()                           // 寫鎖
	defer readWriteLock.Unlock()                   // 記得解開寫鎖
	addressToAliasMap[addressString] = aliasString // 設定位址別名
}

// GetAliasAddressPair - 取得(別名,位址)對
/**
 * @param  string addressString 位址字串
 * @return  []interface{} (別名,位址)對
 */
func GetAliasAddressPair(addressString string) []interface{} {
	return []interface{}{GetAddressAlias(addressString), addressString} // 回傳(別名,位址)對
}
