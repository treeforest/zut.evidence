package grpc

import (
	"context"
	"net"
	"time"

	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service/comt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

func New(addr string, s *comt.Server, opts ...grpc.ServerOption) *grpc.Server {
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
	pb.RegisterCometServer(srv, &server{srv: s})

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
	pb.UnimplementedCometServer
	srv *comt.Server
}

func (s *server) Ping(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *server) PushMsg(ctx context.Context, req *pb.PushMsgReq) (*pb.PushMsgReply, error) {
	s.srv.PushMsg(req.Uid, int(req.Type), req.Msg)
	return &pb.PushMsgReply{}, nil
}
