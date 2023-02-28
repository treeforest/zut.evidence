package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/treeforest/zut.evidence/pkg/discovery"
	"github.com/treeforest/zut.evidence/pkg/discovery/example/pb"
	"github.com/treeforest/zut.evidence/pkg/graceful"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

type Server struct{}

var _ pb.GreeterServer = &Server{}

func (s *Server) Hello(ctx context.Context, req *pb.HelleReq) (*pb.HelloReply, error) {
	fmt.Println("recv: ", req.Name)
	return &pb.HelloReply{Say: fmt.Sprintf("Hello %s", req.Name)}, nil
}

func main() {
	port := flag.Int("port", 8769, "listen port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
	defer lis.Close()
	gSrv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterGreeterServer(gSrv, new(Server))
	go func() {
		if err = gSrv.Serve(lis); err != nil {
			panic(err)
		}
	}()
	fmt.Println("serve at ", lis.Addr().String())

	dis, err := discovery.New(discovery.BaseUrl)
	if err != nil {
		panic(err)
	}
	unRegister, err := dis.Register("Greeter", lis.Addr().String())
	if err != nil {
		panic(err)
	}

	graceful.Stop(func() {
		_ = unRegister(context.Background())
		gSrv.GracefulStop()
	})
}
