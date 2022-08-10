package xhttp

type XResponse struct {
	Code       int32  `json:"code"`
	Msg        string `json:"msg"`
	Data       any    `json:"data,omitempty"`
	ServerTime int64  `json:"serverTime"`
}
