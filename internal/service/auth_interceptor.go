package service

import (
	"context"
	log "github.com/treeforest/logger"
	"github.com/treeforest/zut.evidence/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

// AuthInterceptor 认证拦截器
type AuthInterceptor struct {
	jwtManager      *jwt.JWTManager
	accessibleRoles map[string][]int // method -> role
}

func NewAuthInterceptor(jwtManager *jwt.JWTManager, accessibleRoles map[string][]int) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, accessibleRoles}
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// log.Debug("--> unary interceptor: ", info.FullMethod)

		defer func() {
			if err := recover(); err != nil {
				log.Errorf("[Recovery] %+v", err)
			}
		}()

		uid, role, platform, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		if uid != 0 {
			ctx = SetUid(ctx, uid)
			ctx = SetRole(ctx, role)
			ctx = SetPlatform(ctx, platform)
		}

		start := time.Now()

		resp, err = handler(ctx, req)

		end := time.Now()
		latency := end.Sub(start)

		log.Infof("METHOD:%s | REQUEST:{%v} | REPLY:{%v} | ERR:{%v} | TIME:%d",
			info.FullMethod, req, resp, err, latency/time.Millisecond)

		return resp, err
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) (int64, int, string, error) {
	accessibleRoles, ok := interceptor.accessibleRoles[method]
	if !ok {
		// everyone can access
		return 0, 0, "", nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, 0, "", status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return 0, 0, "", status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	// log.Debugf("--> access token: %s", accessToken)

	uid, role, platform, _, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return 0, 0, "", status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, accessibleRole := range accessibleRoles {
		if accessibleRole == role {
			return uid, role, platform, nil
		}
	}

	return 0, 0, "", status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
