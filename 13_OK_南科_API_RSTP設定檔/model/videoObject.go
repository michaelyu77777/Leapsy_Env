package model

type VideoObject struct {

	//注意:struct名稱開頭必須要大寫...否則無法寫入mongoDB!
	H264_Profile   string `json:"H264_Profile"`
	H264_AVC_Level string `json:"H264_AVC_Level"`
	Resolution     string `json:"resolution"`
	Frame_rate     string `json:"frame_rate"`
	Video_bitrate  string `json:"video_bitrate"`
	I_Frame        string `json:"I-Frame"`
}
