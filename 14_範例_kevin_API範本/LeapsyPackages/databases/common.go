package databases

import (
	"sync"

	"../logings"
)

const (
	equalToConstString            = `$eq`  // =
	greaterThanConstString        = `$gt`  // >
	greaterThanEqualToConstString = `$gte` // >=
	lessThanConstString           = `$lt`  // <
	lessThanEqualToConstString    = `$lte` // <=
)

var (
	logger                                                         = logings.GetLogger() // 記錄器
	periodicallyRWMutex, hourlyRWMutex, dailyRWMutex, alertRWMutex sync.RWMutex          // 讀寫鎖
)
