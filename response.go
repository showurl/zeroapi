package zeroapi

type XResponse struct {
	Code       int32       `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data,omitempty"`
	ServerTime int64       `json:"server_time"`
}
