package records

import (
	"reflect"
	"strconv"
	"time"

	"../logings"
	"../times"
)

// SecondRecord - 環控主機秒紀錄
type SecondRecord struct {
	Time, ModifyTime, AlarmTime, RecoverTime, RecordTime, StartTime                                                                                                                                                                                                                                                                                                                                                                                          time.Time
	Floor, Site, EquipmentNo, CNCName, Status, CurrentProgram, SequenceNumber, SpindleLoad, SpindleTemp, SpindleCurrent, AxisXLoad, AxisYLoad, AxisZLoad, AxisXTemp, AxisYTemp, AxisZTemp, AxisXCurrent, AxisYCurrent, AxisZCurrent, SpindleSpeed, SpindleFeedRate, SpindleTorque, AxisXTorque, AxisYTorque, AxisZTorque, Factory, Floor2, EquipmentID, EquipmentName, FailureReason, PowerTime, OperationTime, Num, RepairTime, IdleTime, AlarmMessages     string
	Sax01_0, Sax01_1, Sax01_2, Sax01_3, Sax01_4, Sax01_5, Sax01_6, Sax02_0, Sax02_1, Sax02_2, Sax02_3, Sax02_4, Sax02_5, Sax02_6, Sax03_0, Sax03_1, Sax03_2, Sax03_3, Sax03_4, Sax03_5, Sax03_6, Sax04_0, Sax04_1, Sax04_2, Sax04_3, Sax04_4, Sax04_5, Sax04_6, Sax05_0, Sax05_1, Sax05_2, Sax05_3, Sax05_4, Sax05_5, Sax05_6, Sax06_0, Sax06_1, Sax06_2, Sax06_3, Sax06_4, Sax06_5, Sax06_6, PMVrs, PMVst, PMVtr, PMIr, PMIs, PMIt, PMHz, PMKw, PMPf, PMKwh int
}

var (
	secondRecordToECSRecordMap = map[string]string{ // secondRecordToRecordMap - 秒紀錄與環控紀錄欄位對照
		`Time`:            `RTEXPDTIME`, // 紀錄時間 varchar(30)
		`Floor`:           `CNC1`,       // 樓層 varchar(20)
		`Site`:            `CNC2`,       // 廠區 varchar(20)
		`EquipmentNo`:     `CNC3`,       // 設備編號 varchar(20)
		`ModifyTime`:      `CNC4`,       // 參數更新時間 varchar(20)
		`CNCName`:         `CNC5`,       // CNC名稱 varchar(20)
		`Status`:          `CNC6`,       // CNC狀態 varchar(20)
		`CurrentProgram`:  `CNC7`,       // 當前程式 varchar(20)
		`SequenceNumber`:  `CNC8`,       // 當前段落號 varchar(20)
		`SpindleLoad`:     `CNC9`,       // 主軸負載 varchar(20)
		`SpindleTemp`:     `CNC10`,      // 主軸溫度 varchar(20)
		`SpindleCurrent`:  `CNC11`,      // 主軸電流 varchar(20)
		`AxisXLoad`:       `CNC12`,      // X軸負載 varchar(20)
		`AxisYLoad`:       `CNC13`,      // Y軸負載 varchar(20)
		`AxisZLoad`:       `CNC14`,      // Z軸負載 varchar(20)
		`AxisXTemp`:       `CNC15`,      // X軸溫度 varchar(20)
		`AxisYTemp`:       `CNC16`,      // Y軸溫度 varchar(20)
		`AxisZTemp`:       `CNC17`,      // Z軸溫度 varchar(20)
		`AxisXCurrent`:    `CNC18`,      // X軸電流 varchar(20)
		`AxisYCurrent`:    `CNC19`,      // Y軸電流 varchar(20)
		`AxisZCurrent`:    `CNC20`,      // Z軸電流 varchar(20)
		`SpindleSpeed`:    `CNC21`,      // 主軸轉速 varchar(20)
		`SpindleFeedRate`: `CNC22`,      // 主軸進給率 varchar(20)
		`SpindleTorque`:   `CNC23`,      // 主軸扭矩 varchar(20)
		`AxisXTorque`:     `CNC24`,      // X軸扭矩 varchar(20)
		`AxisYTorque`:     `CNC25`,      // Y軸扭矩 varchar(20)
		`AxisZTorque`:     `CNC26`,      // Z軸扭矩 varchar(20)
		`AlarmTime`:       `CNC27`,      // 報警發生時間 varchar(20)
		`RecoverTime`:     `CNC28`,      // 報警復位時間 varchar(20)
		`RecordTime`:      `CNC29`,      // 記錄時間 varchar(20)
		`StartTime`:       `CNC30`,      // 報警后加工時間 varchar(20)
		`Factory`:         `CNC31`,      // 廠區 varchar(50)
		`Floor2`:          `CNC32`,      // 樓層 varchar(50)
		`EquipmentID`:     `CNC33`,      // 設備編號 varchar(50)
		`EquipmentName`:   `CNC34`,      // 設備名稱 varchar(50)
		`FailureReason`:   `CNC35`,      // 報警代碼 varchar(50)
		`PowerTime`:       `CNC36`,      // 通電時長 varchar(50)
		`OperationTime`:   `CNC37`,      // 加工時長 varchar(50)
		`Num`:             `CNC38`,      // 產量計數 varchar(50)
		`RepairTime`:      `CNC39`,      // 機台復位耗時 varchar(50)
		`IdleTime`:        `CNC40`,      // 機台復原耗時 varchar(50)
		`AlarmMessages`:   `CNC41`,      // 報警訊息 varchar(300)
		`Sax01_0`:         `M42`,        // PM2.5 懸浮微粒_1 int
		`Sax01_1`:         `M43`,        // CO2 二氧化碳_1 int
		`Sax01_2`:         `M44`,        // HCHO 甲醛_1 int
		`Sax01_3`:         `M45`,        // TVOC 總揮發性有機化合物_1 int
		`Sax01_4`:         `M46`,        // Temperature 溫度_1 int
		`Sax01_5`:         `M47`,        // Humidity 濕度_1 int
		`Sax01_6`:         `M48`,        // PM10 懸浮微粒_1 int
		`Sax02_0`:         `M49`,        // PM2.5 懸浮微粒_2 int
		`Sax02_1`:         `M50`,        // CO2 二氧化碳_2 int
		`Sax02_2`:         `M51`,        // HCHO 甲醛_2 int
		`Sax02_3`:         `M52`,        // TVOC 總揮發性有機化合物_2 int
		`Sax02_4`:         `M53`,        // Temperature 溫度_2 int
		`Sax02_5`:         `M54`,        // Humidity 濕度_2 int
		`Sax02_6`:         `M55`,        // PM10 懸浮微粒_2 int
		`Sax03_0`:         `M56`,        // PM2.5 懸浮微粒_3 int
		`Sax03_1`:         `M57`,        // CO2 二氧化碳_3 int
		`Sax03_2`:         `M58`,        // HCHO 甲醛_3 int
		`Sax03_3`:         `M59`,        // TVOC 總揮發性有機化合物_3 int
		`Sax03_4`:         `M60`,        // Temperature 溫度_3 int
		`Sax03_5`:         `M61`,        // Humidity 濕度_3 int
		`Sax03_6`:         `M62`,        // PM10 懸浮微粒_3 int
		`Sax04_0`:         `M63`,        // PM2.5 懸浮微粒_4 int
		`Sax04_1`:         `M64`,        // CO2 二氧化碳_4 int
		`Sax04_2`:         `M65`,        // HCHO 甲醛_4 int
		`Sax04_3`:         `M66`,        // TVOC 總揮發性有機化合物_4 int
		`Sax04_4`:         `M67`,        // Temperature 溫度_4 int
		`Sax04_5`:         `M68`,        // Humidity 濕度_4 int
		`Sax04_6`:         `M69`,        // PM10 懸浮微粒_4 int
		`Sax05_0`:         `M70`,        // PM2.5 懸浮微粒_5 int
		`Sax05_1`:         `M71`,        // CO2 二氧化碳_5 int
		`Sax05_2`:         `M72`,        // HCHO 甲醛_5 int
		`Sax05_3`:         `M73`,        // TVOC 總揮發性有機化合物_5 int
		`Sax05_4`:         `M74`,        // Temperature 溫度_5 int
		`Sax05_5`:         `M75`,        // Humidity 濕度_5 int
		`Sax05_6`:         `M76`,        // PM10 懸浮微粒_5 int
		`Sax06_0`:         `M77`,        // PM2.5 懸浮微粒_6 int
		`Sax06_1`:         `M78`,        // CO2 二氧化碳_6 int
		`Sax06_2`:         `M79`,        // HCHO 甲醛_6 int
		`Sax06_3`:         `M80`,        // TVOC 總揮發性有機化合物_6 int
		`Sax06_4`:         `M81`,        // Temperature 溫度_6 int
		`Sax06_5`:         `M82`,        // Humidity 濕度_6 int
		`Sax06_6`:         `M83`,        // PM10 懸浮微粒_6 int
		`PMVrs`:           `M84`,        // PM_RS線電壓 int
		`PMVst`:           `M85`,        // PM_ST線電壓 int
		`PMVtr`:           `M86`,        // PM_TR線電壓 int
		`PMIr`:            `M87`,        // PM_R相電流 int
		`PMIs`:            `M88`,        // PM_S相電流 int
		`PMIt`:            `M89`,        // PM_T相電流 int
		`PMHz`:            `M90`,        // PM_頻率 int
		`PMKw`:            `M91`,        // PM_有效功率 int
		`PMPf`:            `M92`,        // PM_功率因數 int
		`PMKwh`:           `M93`,        // PM_累計電度 int
	}
)

// getMappedToECSRecordFieldName - 取得秒紀錄對應的環控紀錄欄位名
/**
 * @param string secondRecordFieldName 秒紀錄欄位名
 * @return string 環控紀錄欄位名
 */
func getMappedToECSRecordFieldName(secondRecordFieldName string) string {
	return secondRecordToECSRecordMap[secondRecordFieldName] // 回傳秒紀錄對應的環控紀錄欄位名
}

// SecondRecord - 將ECSRecord轉成SecondRecord
/**
 * @return SecondRecord 秒紀錄
 */
func (ecsRecord ECSRecord) SecondRecord() (secondRecord SecondRecord) {

	valueOfECSRecord := reflect.ValueOf(ecsRecord)               // 環控紀錄的值
	typeOfSecondRecord := reflect.TypeOf(secondRecord)           // 秒紀錄的資料型別
	valueOfSecondRecord := reflect.ValueOf(&secondRecord).Elem() // 秒紀錄的值

	for index := 0; index < typeOfSecondRecord.NumField(); index++ { // 針對秒紀錄每一個欄位

		secondRecordFieldName := typeOfSecondRecord.Field(index).Name              // 秒紀錄欄位名
		secondRecordFieldValue := valueOfSecondRecord.Field(index)                 // 秒紀錄欄位值
		ecsRecordFieldName := getMappedToECSRecordFieldName(secondRecordFieldName) // 環控紀錄欄位名

		if `` != ecsRecordFieldName { // 若有對應的環控紀錄欄位名

			ecsRecordFieldValue := valueOfECSRecord.FieldByName(ecsRecordFieldName) // 環控紀錄欄位值

			switch typeOfSecondRecord.Field(index).Type.String() { // 若秒紀錄欄位型別為

			case `int`: // 整數

				integer, strconvAtoiError := strconv.Atoi(ecsRecordFieldValue.String()) // 將環控紀錄欄位值字串轉為整數

				// 取得記錄器格式和參數
				formatString, args := logings.GetLogFuncFormatAndArguments(
					[]string{`環控紀錄欄位 %s 值轉成整數`},
					[]interface{}{ecsRecordFieldName},
					strconvAtoiError,
				)

				if nil != strconvAtoiError { // 若將環控紀錄欄位值字串轉為整數錯誤
					logger.Panicf(formatString, args...) // 記錄錯誤並逐層結束程式
				} else { // 若將環控紀錄欄位值字串轉為整數成功
					secondRecordFieldValue.SetInt(int64(integer)) // 設定秒紀錄欄位值為環控紀錄欄位轉化後的整數值
				}

			case `time.Time`: // 時間
				secondRecordFieldValue.Set(reflect.ValueOf(times.RTEXPDTIMEStringToTime(ecsRecordFieldValue.String()))) // 設定秒紀錄欄位值為環控紀錄欄位的時間值

			default: // 預設
				secondRecordFieldValue.SetString(ecsRecordFieldValue.String()) // 設定秒紀錄欄位值為環控紀錄欄位的字串值

			}

		}

	}

	return // 回傳
}
