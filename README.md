# zeroapi

go-zero通用的api 让你不需要写api服务

# 示例

[hello](examples/hello/README.md)

# http 响应拦截

> 详细代码 [handler.go](handler.go)

```go
package zeroapi

func (h *Handler) defaultResponseHandler(in proto.Message) (code int, msg string, data interface{}) {
	if in == nil {
		return 0, "", nil
	} else {
		message, ok := in.(*dynamic.Message)
		if ok {
			if cod, exist := protoMessageValue(message, "errCode", 0); exist {
				code = int(InterfaceToInt64(cod))
			}
			if failedReason, exist := protoMessageValue(message, "failedReason", ""); exist && failedReason != "" {
				msg = failedReason.(string)
				if code == 0 {
					code = -1
				}
			}
			data = in
			return
		}
		return 0, "", in
	}
}
```

# http header 转 proto参数

> 详细代码 [requestparser.go](internal/requestparser.go)

```go
_ = dm.TrySetFieldByName("ip", r.Header.Get("X-Real-IP"))
```
