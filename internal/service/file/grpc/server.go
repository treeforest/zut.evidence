package grpc

import (
	"bytes"
	"context"
	pb "github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

func New(addr string, f *file.File, opts ...grpc.ServerOption) *grpc.Server {
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
	pb.RegisterFileServer(srv, &server{srv: f})

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
	pb.UnimplementedFileServer
	srv *file.File
}

var _ pb.FileServer = &server{}

func (s *server) Upload(ctx context.Context, req *pb.UploadReq) (*pb.UploadReply, error) {
	cid, err := s.srv.UploadFile(req.Uid, req.Filename, bytes.NewReader(req.Data))
	return &pb.UploadReply{Cid: cid}, err
}
