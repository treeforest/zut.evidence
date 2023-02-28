package grpc

import (
	"context"
	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service"
	"github.com/treeforest/zut.evidence/internal/service/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

func New(addr string, w *wallet.Wallet, opts ...grpc.ServerOption) *grpc.Server {
	srv := grpc.NewServer(opts...)

	// 注册服务
	pb.RegisterWalletServer(srv, &server{srv: w})

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
	pb.UnimplementedWalletServer
	srv *wallet.Wallet
}

var _ pb.WalletServer = &server{}

func (s *server) GenerateKey(ctx context.Context, req *emptypb.Empty) (*pb.GenerateKeyReply, error) {
	publicKey, err := s.srv.GenerateKey(service.GetUid(ctx))
	if err != nil {
		return &pb.GenerateKeyReply{}, err
	}
	return &pb.GenerateKeyReply{PublicKey: publicKey}, nil
}

func (s *server) DownloadKey(ctx context.Context, req *emptypb.Empty) (*pb.DownloadKeyReply, error) {
	key, pub, err := s.srv.DownloadKey(service.GetUid(ctx))
	if err != nil {
		return &pb.DownloadKeyReply{}, err
	}
	return &pb.DownloadKeyReply{PrivateKey: key, PublicKey: pub}, nil
}

func (s *server) GetPrivateKey(ctx context.Context, req *emptypb.Empty) (*pb.GetPrivateKeyReply, error) {
	keyBytes, err := s.srv.GetPrivateKey(service.GetUid(ctx))
	if err != nil {
		return &pb.GetPrivateKeyReply{}, err
	}
	return &pb.GetPrivateKeyReply{PrivateKey: keyBytes}, nil
}

func (s *server) GetPubKey(ctx context.Context, req *emptypb.Empty) (*pb.GetPubKeyReply, error) {
	publicKey, _ := s.srv.GetPubKey(service.GetUid(ctx))
	return &pb.GetPubKeyReply{PublicKey: publicKey}, nil
	//if err != nil {
	//	return &pb.GetPubKeyReply{}, err
	//}
	//// log.Debug("pub:", publicKey)
	//return &pb.GetPubKeyReply{PublicKey: publicKey}, nil
}

func (s *server) Sign(ctx context.Context, req *pb.SignReq) (*pb.SignReply, error) {
	sig, err := s.srv.Sign(service.GetUid(ctx), req.Data)
	if err != nil {
		return &pb.SignReply{}, err
	}
	// log.Debug("sig:", sig)
	return &pb.SignReply{Signature: sig}, nil
}

func (s *server) Verify(ctx context.Context, req *pb.VerifyReq) (*pb.VerifyReply, error) {
	// log.Debug(*req)
	b := s.srv.Verify(req.Signature, req.PublicKey, req.Data)
	return &pb.VerifyReply{Status: b}, nil
}

func (s *server) Encrypt(ctx context.Context, req *pb.EncryptReq) (*pb.EncryptReply, error) {
	ciphertext, err := s.srv.Encrypt(service.GetUid(ctx), req.Data)
	return &pb.EncryptReply{Ciphertext: ciphertext}, err
}

func (s *server) EncryptByPubKey(ctx context.Context, req *pb.EncryptByPubKeyReq) (
	*pb.EncryptByPubKeyReply, error) {
	ciphertext, err := s.srv.EncryptByPubKey(req.Data, req.PublicKey)
	return &pb.EncryptByPubKeyReply{Ciphertext: ciphertext}, err
}

func (s *server) Decrypt(ctx context.Context, req *pb.DecryptReq) (*pb.DecryptReply, error) {
	data, err := s.srv.Decrypt(service.GetUid(ctx), req.Ciphertext)
	return &pb.DecryptReply{Data: data}, err
}
