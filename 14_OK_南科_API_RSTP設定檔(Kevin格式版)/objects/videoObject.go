package objects

type VideoObject struct {

	//注意:struct名稱開頭必須要大寫...否則無法寫入mongoDB!
	H264_Profile   int    `json:"H264_Profile"`
	H264_AVC_Level int    `json:"H264_AVC_Level"`
	Resolution     string `json:"resolution"`
	Frame_rate     int    `json:"frame_rate"`
	Video_bitrate  int    `json:"video_bitrate"`
	I_Frame        int    `json:"I_Frame"`
}
