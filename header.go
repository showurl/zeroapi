package zeroapi

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func getHeader(ctx context.Context) metadata.MD {
	incomingContext, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadata.MD{}
	}
	return incomingContext
}

func GetValueByKey(ctx context.Context, key string) (value string) {
	strings := getHeader(ctx).Get(key)
	if len(strings) > 0 {
		value = strings[0]
	}
	return
}
