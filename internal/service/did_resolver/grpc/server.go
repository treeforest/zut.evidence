package grpc

import (
	"context"
	log "github.com/treeforest/logger"
	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service"
	"github.com/treeforest/zut.evidence/internal/service/did_resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

func New(addr string, d *did_resolver.DIDResolver, opts ...grpc.ServerOption) *grpc.Server {
	srv := grpc.NewServer(opts...)

	// 注册服务
	pb.RegisterDIDResolverServer(srv, &server{srv: d})

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
	pb.UnimplementedDIDResolverServer
	srv *did_resolver.DIDResolver
}

var _ pb.DIDResolverServer = &server{}

func (s *server) CreateDID(ctx context.Context, req *emptypb.Empty) (*pb.CreateDIDReply, error) {
	did, created, err := s.srv.CreateDID(ctx, service.GetUid(ctx))
	return &pb.CreateDIDReply{Did: did, Created: created}, err
}

func (s *server) RevokeDID(ctx context.Context, req *pb.RevokeDIDReq) (*emptypb.Empty, error) {
	err := s.srv.RevokeDID(ctx, req.Did)
	return &emptypb.Empty{}, err
}

func (s *server) GetDIDs(ctx context.Context, req *emptypb.Empty) (*pb.GetDIDsReply, error) {
	b, err := s.srv.GetDIDs(ctx, service.GetUid(ctx))
	return &pb.GetDIDsReply{Dids: b}, err
}

func (s *server) GetDIDDocument(ctx context.Context, req *pb.GetDIDDocumentReq) (*pb.GetDIDDocumentRely, error) {
	log.Debugf("GetDIDDocument | req:%v", *req)
	doc, err := s.srv.GetDIDDocument(ctx, req.Did)
	if err != nil {
		return &pb.GetDIDDocumentRely{}, err
	}
	return &pb.GetDIDDocumentRely{
		DidDocument: &pb.DIDDocument{
			Context:        doc.Context,
			Id:             doc.Id,
			Created:        doc.Created,
			Updated:        doc.Updated,
			Authentication: doc.Authentication,
			PublicKey: []*pb.DIDDocument_PublicKey{
				{
					Id:           doc.PublicKey[0].Id,
					Type:         doc.PublicKey[0].Type,
					PublicKeyHex: doc.PublicKey[0].PublicKeyHex,
				},
			},
			Proof: &pb.DIDDocument_Proof{
				Type:      doc.Proof.Type,
				Creator:   doc.Proof.Creator,
				Signature: doc.Proof.Signature,
			},
		},
	}, nil
}

func (s *server) GetPublicKeyByDID(ctx context.Context, req *pb.GetPublicKeyByDIDReq) (*pb.GetPublicKeyByDIDReply, error) {
	pub, err := s.srv.GetPublicKeyByDID(ctx, req.Did)
	return &pb.GetPublicKeyByDIDReply{PublicKey: pub}, err
}

func (s *server) ExistDID(ctx context.Context, req *pb.ExistDIDReq) (*pb.ExistDIDReply, error) {
	exists := s.srv.ExistDiD(ctx, req.Uid, req.Did)
	return &pb.ExistDIDReply{Exists: exists}, nil
}

func (s *server) GetOwnerByDID(ctx context.Context, req *pb.GetOwnerByDIDReq) (*pb.GetOwnerByDIDReply, error) {
	owner, err := s.srv.GetOwnerByDID(ctx, req.Did)
	return &pb.GetOwnerByDIDReply{Uid: owner}, err
}
