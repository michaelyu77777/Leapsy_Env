package model

type RtspConfig struct {
	Video             VideoObject
	Audio             AudioObject //注意:struct名稱開頭必須要大寫否則無法寫入mongoDB
	EnableAudioStream string
}
