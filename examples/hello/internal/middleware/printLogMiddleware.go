package middleware

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

func PrintLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 打印请求头
		logx.WithContext(r.Context()).Info("request header:", r.Header)
		next(w, r)
		// 打印响应头
		logx.WithContext(r.Context()).Info("response header:", w.Header())
	}
}
