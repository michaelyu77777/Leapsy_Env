package settings

const (
	// DbName :資料庫名
	DbName = "leapsy_env" //DB

	// CollectionNameOfCheckInRecord :Collection名稱:打卡紀錄
	CollectionNameOfCheckInRecord = "check_in_record" //Collection

	// CollectionNameOfCheckInStatistics :Collection名:打卡統計
	CollectionNameOfCheckInStatistics = "check_in_statistics" //Collection

	// const CollectionName = "persion"                                //Collection //範例程式

	// PortOfAPI :開API的Port
	//const PortOfAPI = 8081 //API port
	//PortOfAPI = 8000 //API port
	PortOfAPI = 8082 //API port

	// PortOfMongoDB :MongoDB的Port
	PortOfMongoDB string = "27017"
)
