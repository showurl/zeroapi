package main

import (
	"flag"
	"fmt"

	"github.com/showurl/zeroapi/examples/hello/internal/config"
	"github.com/showurl/zeroapi/examples/hello/internal/server"
	"github.com/showurl/zeroapi/examples/hello/internal/svc"
	"github.com/showurl/zeroapi/examples/hello/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/hello.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	svr := server.NewStreamGreeterServer(ctx)
	go svr.Gateway()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterStreamGreeterServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
