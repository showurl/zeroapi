package zeroapi

import (
	"encoding/json"
	"github.com/showurl/zeroapi/xhttp"
	"github.com/zeromicro/go-zero/zrpc"
)

type (
	// Config is the configuration for gateway.
	Config struct {
		PbGroup     string // 比如用户服务是 pb.userService
		RpcListenOn string // 比如用户服务是 :8080
		// ProtoSet is the file of proto set, like hello.pb
		//ProtoSets             [][]byte `json:",optional"` // 不要填写
		CallRpcTimeoutSeconds int64 `json:",default=10"`
		Port                  int   `json:",default=9000"`
	}
	// RouteMapping is a mapping between a gateway route and an upstream rpc method.
	RouteMapping struct {
		// Method is the HTTP method, like GET, POST, PUT, DELETE.
		Method string
		// Path is the HTTP path.
		Path string
		// RpcPath is the gRPC rpc method, with format of package.service/method
		RpcPath  string
		optionFs []OptionFunc
		//ResponseHandler ResponseHandler
	}

	// upstream is the configuration for an upstream.
	upstream struct {
		// Grpc is the target of the upstream.
		Grpc zrpc.RpcClientConf
		// ProtoSet is the file of proto set, like hello.pb
		ProtoSets [][]byte `json:",optional"`
		// Mapping is the mapping between gateway routes and upstream rpc methods.
		// Keep it blank if annotations are added in rpc methods.
		Mapping []RouteMapping `json:",optional"`
	}
)

var (
	defaultErrBuf []byte
)

func init() {
	defaultErrBuf, _ = json.Marshal(xhttp.XResponse{
		Code: 0,
		Msg:  "网络繁忙，请稍后再试",
	})
}
