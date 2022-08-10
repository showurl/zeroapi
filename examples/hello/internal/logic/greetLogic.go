package logic

import (
	"context"
	"math/rand"
	"time"

	"github.com/showurl/zeroapi/examples/hello/internal/svc"
	"github.com/showurl/zeroapi/examples/hello/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GreetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGreetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GreetLogic {
	return &GreetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GreetLogic) Greet(in *pb.StreamReq) (*pb.StreamResp, error) {
	failedReason := ""
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		failedReason = "随机失败"
	}
	return &pb.StreamResp{Greet: "your ip:" + in.Ip, FailedReason: failedReason}, nil
}
