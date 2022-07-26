## proto file

### common.proto

```protobuf
syntax = "proto3";

package pb;
option go_package = "./pb";

message CommonReq {
  string id = 1;
}

```

### hello.proto

```protobuf
syntax = "proto3";

package pb;
option go_package = "./pb";
import "common.proto";

message StreamReq {
  string name = 1;
  CommonReq common = 2;
}

message StreamResp {
  string greet = 1;
}

service StreamGreeter {
  rpc greet(StreamReq) returns (StreamResp);
}
```

### ProtoSet file

```go
// protoset.go
package pb

import _ "embed"

//go:embed common.pb
var ProtoSetCommon []byte

//go:embed hello.pb
var ProtoSetHello []byte
```

## gateway route

> 参考 [gateway.go](internal/server/gateway.go)

> 参考 [middleware](internal/middleware/printLogMiddleware.go)

```go
engine := zeroapi.Engine(restConf, gatewatConf, protoSets)
engine.GET("/ping", s.Ping)
svr := engine.Server(options...)
svr.Use(middlewares...)
svr.Start()
```

## rpc logic 如何获取 header 参数

```go
// gateway.go
svr := engine.Server(zeroapi.WithHeaderProcessor(func(header http.Header) []string {
return []string{
"User-Agent:" + header.Get("User-Agent"),
"X-Forwarded-For:" + header.Get("X-Forwarded-For"),
"X-Real-IP:" + header.Get("X-Real-IP"),
"token:" + header.Get("token"),
}
}))
// logic.go
func (l *GreetLogic) Greet(in *pb.StreamReq) (*pb.StreamResp, error) {
return &pb.StreamResp{Greet: zeroapi.GetValueByKey(l.ctx, "User-Agent")}, nil
}
```

# 运行 hello.go

```shell
go run .
curl http://127.0.0.1:9999/hello/v1/greet
# 得到结果 {"code":0,"msg":"","data":{"greet":"grpc-go/1.48.0"},"server_time":1658829159831}
```

# 如何返回自定义错误

- 增加failedReason字段

```protobuf
message StreamResp {
  string greet = 1;
  string failedReason = 2;
}
```

- defaultResponseHandler 增加逻辑
```go
func (h *Handler) defaultResponseHandler(in proto.Message) (code int, msg string, data interface{}) {
	if in == nil {
		return 0, "", nil
	} else {
		message, ok := in.(*dynamic.Message)
		if ok {
			if failedReason := message.GetFieldByName("failedReason"); failedReason != nil {
				if s, ok := failedReason.(string); ok && s != "" {
					return -1, s, nil
				}
			}
		}
		return 0, "", in
	}
}
```

- logic 测试
```go
func (l *GreetLogic) Greet(in *pb.StreamReq) (*pb.StreamResp, error) {
	failedReason := ""
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		failedReason = "随机失败"
	}
	return &pb.StreamResp{Greet: zeroapi.GetValueByKey(l.ctx, "User-Agent"), FailedReason: failedReason}, nil
}
```

```shell
go run .
curl http://127.0.0.1:9999/hello/v1/greet
# 得到结果 {"code":-1,"msg":"随机失败","server_time":1658832552397}
```