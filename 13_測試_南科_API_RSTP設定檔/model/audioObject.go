package model

type AudioObject struct {

	//注意:struct名稱開頭必須要大寫...否則無法寫入mongoDB!
	SampleRate   string `json:"sampleRate"`
	ChannelCount string `json:"ChannelCount"`
	Bitrate      string `json:"bitrate"`
}
