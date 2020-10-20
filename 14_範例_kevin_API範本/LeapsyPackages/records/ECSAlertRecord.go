package records

// ECSAlertRecord - 環控資料庫警報紀錄
type ECSAlertRecord struct {
	ALERTEVENTID, // 警報編號 int
	ALERTEVENTTIME, // 日期時間	datetime
	VARTAG, // 點名稱	nvarchar(50)
	COMMENT, // 說明	nvarchar(max)
	ALERTTYPE, // 警報群組	int
	LINETEXT string // 行文字	nvarchar(max)
}
