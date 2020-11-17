package objects

type AudioObject struct {

	//注意:struct名稱開頭必須要大寫...否則無法寫入mongoDB!
	SampleRate   int `json:"sampleRate"`
	ChannelCount int `json:"ChannelCount"`
	Bitrate      int `json:"bitrate"`
}
