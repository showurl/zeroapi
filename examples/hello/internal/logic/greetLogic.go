package logic

import (
	"context"
	"github.com/showurl/zeroapi"

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
	return &pb.StreamResp{Greet: zeroapi.GetValueByKey(l.ctx, "User-Agent")}, nil
}
