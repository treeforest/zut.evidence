package service

import (
	"context"
	"google.golang.org/grpc/metadata"
)

type mdUidKey struct{}
type mdRoleKey struct{}
type mdPlatformKey struct{}

const (
	authorizationKey = "authorization"
)

// NewOutgoingCtxFromInComingCtx 通过服务端接收的 IncomingContext 来创建GRPC客户端请求的 OutgoingContext，
// 主要是将 IncomingContext 中的元数据透传给 OutgoingContext。
func NewOutgoingCtxFromInComingCtx(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return context.TODO()
	}

	return metadata.NewOutgoingContext(context.TODO(), md)
}

func SetUid(ctx context.Context, uid int64) context.Context {
	return context.WithValue(ctx, mdUidKey{}, uid)
}

func GetUid(ctx context.Context) int64 {
	val := ctx.Value(mdUidKey{})
	if val == nil {
		return 0
	}
	return val.(int64)
}

func SetRole(ctx context.Context, role int) context.Context {
	return context.WithValue(ctx, mdRoleKey{}, role)
}

func GetRole(ctx context.Context) int {
	val := ctx.Value(mdRoleKey{})
	if val == nil {
		panic("not found role from context")
	}
	return val.(int)
}

func SetPlatform(ctx context.Context, platform string) context.Context {
	return context.WithValue(ctx, mdPlatformKey{}, platform)
}

func GetPlatform(ctx context.Context) string {
	val := ctx.Value(mdRoleKey{})
	if val == nil {
		panic("not found platform from context")
	}
	return val.(string)
}
