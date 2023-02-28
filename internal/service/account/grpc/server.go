package grpc

import (
	"context"
	pb "github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service/account"
	"github.com/treeforest/zut.evidence/internal/service/account/conf"
	"github.com/treeforest/zut.evidence/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

func NewServer(c *conf.Config, jwtMgr *jwt.JWTManager) pb.AccountServer {
	return &server{srv: account.New(c, jwtMgr)}
}

func New(addr string, ac *account.Account, opts ...grpc.ServerOption) *grpc.Server {
	if opts == nil {
		opts = make([]grpc.ServerOption, 0)
	}
	keepParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: time.Minute * 30,
		Time:              time.Minute * 10,
		Timeout:           time.Second * 20,
	})

	srv := grpc.NewServer(append(opts, keepParams)...)

	// 注册服务
	pb.RegisterAccountServer(srv, &server{srv: ac})

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	reflection.Register(srv)

	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return srv
}

type server struct {
	pb.UnimplementedAccountServer
	srv *account.Account
}

var _ pb.AccountServer = &server{}

func (s *server) Login(ctx context.Context, req *pb.LoginReq) (resp *pb.LoginReply, err error) {
	token, role, err := s.srv.Login(ctx, req.Nick, req.Password, req.Platform)
	if err != nil {
		return &pb.LoginReply{}, err
	}
	return &pb.LoginReply{Token: token, Role: int32(role), Nick: req.Nick}, nil
}

func (s *server) Register(ctx context.Context, req *pb.RegisterReq) (resp *pb.RegisterReply, err error) {
	err = s.srv.Register(ctx, req.Nick, req.Phone, req.Email, req.Password, int(req.Role))
	if err != nil {
		return &pb.RegisterReply{}, err
	}
	return &pb.RegisterReply{}, nil
}
