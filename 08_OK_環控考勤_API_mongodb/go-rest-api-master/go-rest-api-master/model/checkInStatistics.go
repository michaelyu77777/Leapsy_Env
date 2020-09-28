package model

type CheckInStatistics struct {
	Date        string //注意:struct名稱開頭必須要大寫...否則無法寫入mongoDB!!!不知道為什麼...
	Expected    string
	Attendance  string
	Not_arrived string
	Guests      string
}
