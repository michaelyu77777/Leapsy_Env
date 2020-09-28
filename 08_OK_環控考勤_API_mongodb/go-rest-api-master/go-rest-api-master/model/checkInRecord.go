package model

type CheckInRecord struct {
	id            int    //`json:"_id"`
	Name          string //注意:struct名稱開頭必須要大寫...否則無法寫入mongoDB!!!不知道為什麼...
	Check_in_time string
	Pic           string
	Leave_type    string
	Date          string
	Department    string
	Position      string
}
