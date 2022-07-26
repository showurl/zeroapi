package gateway

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"path"
	"strings"
	"time"
)

type (
	GatewayEngine struct {
		RestConf  rest.RestConf
		Config    Config
		prefix    string
		router    *router
		ProtoSets [][]byte
	}
	router struct {
		routers []RouteMapping
	}
)

func (r *router) add(s string, url string, rpcPath string, fs ...OptionFunc) {
	r.routers = append(r.routers, RouteMapping{
		Method:   s,
		Path:     url,
		RpcPath:  rpcPath,
		optionFs: fs,
	})
}

func Engine(restConf rest.RestConf, conf Config, protoSets [][]byte) *GatewayEngine {
	return &GatewayEngine{
		RestConf:  restConf,
		Config:    conf,
		router:    &router{},
		ProtoSets: protoSets,
	}
}

func (e *GatewayEngine) Server(opts ...Option) *Server {
	svr := &Server{
		Server:    rest.MustNewServer(e.RestConf),
		upstreams: e.upstreams(),
		timeout:   time.Duration(e.Config.CallRpcTimeoutSeconds) * time.Second,
	}
	for _, opt := range opts {
		opt(svr)
	}
	return svr
}

func (e *GatewayEngine) formatPrefix(prefix string) string {
	return path.Join(e.prefix, prefix)
}

func (e *GatewayEngine) Group(prefix string) (n *GatewayEngine) {
	n = &GatewayEngine{
		RestConf:  e.RestConf,
		Config:    e.Config,
		prefix:    e.formatPrefix(prefix),
		router:    e.router,
		ProtoSets: e.ProtoSets,
	}
	return
}

func (e *GatewayEngine) GET(url string, handler interface{}, optionFs ...OptionFunc) {
	rpcPath := path.Join(e.Config.PbGroup, funcName(handler))
	e.router.add("get", e.formatPrefix(url), rpcPath, optionFs...)
}
func (e *GatewayEngine) POST(url string, handler interface{}, optionFs ...OptionFunc) {
	rpcPath := path.Join(e.Config.PbGroup, funcName(handler))
	e.router.add("post", e.formatPrefix(url), rpcPath, optionFs...)
}
func (e *GatewayEngine) PUT(url string, handler interface{}, optionFs ...OptionFunc) {
	rpcPath := path.Join(e.Config.PbGroup, funcName(handler))
	e.router.add("put", e.formatPrefix(url), rpcPath, optionFs...)
}
func (e *GatewayEngine) DELETE(url string, handler interface{}, optionFs ...OptionFunc) {
	rpcPath := path.Join(e.Config.PbGroup, funcName(handler))
	e.router.add("delete", e.formatPrefix(url), rpcPath, optionFs...)
}
func (e *GatewayEngine) PATCH(url string, handler interface{}, optionFs ...OptionFunc) {
	rpcPath := path.Join(e.Config.PbGroup, funcName(handler))
	e.router.add("patch", e.formatPrefix(url), rpcPath, optionFs...)
}

func (e *GatewayEngine) upstreams() []upstream {
	endpoint := e.Config.RpcListenOn
	if strings.HasPrefix(endpoint, ":") {
		endpoint = "127.0.0.1" + endpoint
	}
	return []upstream{{
		Grpc: zrpc.RpcClientConf{
			Endpoints: []string{endpoint},
			NonBlock:  true,
		},
		ProtoSets: e.ProtoSets,
		Mapping:   e.router.routers,
	}}
}
