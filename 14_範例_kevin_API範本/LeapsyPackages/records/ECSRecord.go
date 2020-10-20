package records

// ECSRecord - 環控資料庫紀錄
type ECSRecord struct {
	RTEXPDTIME, // 紀錄時間 varchar(30)
	CNC1, // 樓層 varchar(20)
	CNC2, // 廠區 varchar(20)
	CNC3, // 設備編號 varchar(20)
	CNC4, // 參數更新時間 varchar(20)
	CNC5, // CNC名稱 varchar(20)
	CNC6, // CNC狀態 varchar(20)
	CNC7, // 當前程式 varchar(20)
	CNC8, // 當前段落號 varchar(20)
	CNC9, // 主軸負載 varchar(20)
	CNC10, // 主軸溫度 varchar(20)
	CNC11, // 主軸電流 varchar(20)
	CNC12, // X軸負載 varchar(20)
	CNC13, // Y軸負載 varchar(20)
	CNC14, // Z軸負載 varchar(20)
	CNC15, // X軸溫度 varchar(20)
	CNC16, // Y軸溫度 varchar(20)
	CNC17, // Z軸溫度 varchar(20)
	CNC18, // X軸電流 varchar(20)
	CNC19, // Y軸電流 varchar(20)
	CNC20, // Z軸電流 varchar(20)
	CNC21, // 主軸轉速 varchar(20)
	CNC22, // 主軸進給率 varchar(20)
	CNC23, // 主軸扭矩 varchar(20)
	CNC24, // X軸扭矩 varchar(20)
	CNC25, // Y軸扭矩 varchar(20)
	CNC26, // Z軸扭矩 varchar(20)
	CNC27, // 報警發生時間 varchar(20)
	CNC28, // 報警復位時間 varchar(20)
	CNC29, // 記錄時間 varchar(20)
	CNC30, // 報警后加工時間 varchar(20)
	CNC31, // 廠區 varchar(50)
	CNC32, // 樓層 varchar(50)
	CNC33, // 設備編號 varchar(50)
	CNC34, // 設備名稱 varchar(50)
	CNC35, // 報警代碼 varchar(50)
	CNC36, // 通電時長 varchar(50)
	CNC37, // 加工時長 varchar(50)
	CNC38, // 產量計數 varchar(50)
	CNC39, // 機台復位耗時 varchar(50)
	CNC40, // 機台復原耗時 varchar(50)
	CNC41, // 報警訊息 varchar(300)
	M42, // PM2.5 懸浮微粒_1 int
	M43, // CO2 二氧化碳_1 int
	M44, // HCHO 甲醛_1 int
	M45, // TVOC 總揮發性有機化合物_1 int
	M46, // Temperature 溫度_1 int
	M47, // Humidity 濕度_1 int
	M48, // PM10 懸浮微粒_1 int
	M49, // PM2.5 懸浮微粒_2 int
	M50, // CO2 二氧化碳_2 int
	M51, // HCHO 甲醛_2 int
	M52, // TVOC 總揮發性有機化合物_2 int
	M53, // Temperature 溫度_2 int
	M54, // Humidity 濕度_2 int
	M55, // PM10 懸浮微粒_2 int
	M56, // PM2.5 懸浮微粒_3 int
	M57, // CO2 二氧化碳_3 int
	M58, // HCHO 甲醛_3 int
	M59, // TVOC 總揮發性有機化合物_3 int
	M60, // Temperature 溫度_3 int
	M61, // Humidity 濕度_3 int
	M62, // PM10 懸浮微粒_3 int
	M63, // PM2.5 懸浮微粒_4 int
	M64, // CO2 二氧化碳_4 int
	M65, // HCHO 甲醛_4 int
	M66, // TVOC 總揮發性有機化合物_4 int
	M67, // Temperature 溫度_4 int
	M68, // Humidity 濕度_4 int
	M69, // PM10 懸浮微粒_4 int
	M70, // PM2.5 懸浮微粒_5 int
	M71, // CO2 二氧化碳_5 int
	M72, // HCHO 甲醛_5 int
	M73, // TVOC 總揮發性有機化合物_5 int
	M74, // Temperature 溫度_5 int
	M75, // Humidity 濕度_5 int
	M76, // PM10 懸浮微粒_5 int
	M77, // PM2.5 懸浮微粒_6 int
	M78, // CO2 二氧化碳_6 int
	M79, // HCHO 甲醛_6 int
	M80, // TVOC 總揮發性有機化合物_6 int
	M81, // Temperature 溫度_6 int
	M82, // Humidity 濕度_6 int
	M83, // PM10 懸浮微粒_6 int
	M84, // PM_RS線電壓 int
	M85, // PM_ST線電壓 int
	M86, // PM_TR線電壓 int
	M87, // PM_R相電流 int
	M88, // PM_S相電流 int
	M89, // PM_T相電流 int
	M90, // PM_頻率 int
	M91, // PM_有效功率 int
	M92, // PM_功率因數 int
	M93 string // PM_累計電度 int
}
