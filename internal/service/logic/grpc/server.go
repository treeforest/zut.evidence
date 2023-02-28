package grpc

import (
	"context"
	"google.golang.org/grpc/reflection"
	"net"

	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service/logic"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func New(addr string, l *logic.Logic, opts ...grpc.ServerOption) *grpc.Server {
	srv := grpc.NewServer(opts...)
	pb.RegisterLogicServer(srv, &server{srv: l}) // 注册服务
	reflection.Register(srv)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return srv
}

type server struct {
	srv *logic.Logic
	pb.UnimplementedLogicServer
}

var _ pb.LogicServer = &server{}

func (s *server) ApplyKYC(ctx context.Context, req *pb.ApplyKYCReq) (*emptypb.Empty, error) {
	err := s.srv.ApplyKYC(ctx, int(req.Type), req.Name, req.IdCard, req.Cids)
	return &emptypb.Empty{}, err
}

/* 发布凭证/获取凭证列表 */

func (s *server) Issue(ctx context.Context, req *pb.IssueReq) (*emptypb.Empty, error) {
	err := s.srv.Issue(ctx, req.Issuer, req.ShortDescription, req.LongDescription, req.Endpoint, req.Website, int(req.Type))
	return &emptypb.Empty{}, err
}

func (s *server) GetIssued(ctx context.Context, req *emptypb.Empty) (*pb.GetIssuedReply, error) {
	list, err := s.srv.GetIssued(ctx)
	return &pb.GetIssuedReply{Issuers: list, Count: int32(len(list))}, err
}

func (s *server) RevokeIssued(ctx context.Context, req *pb.RevokeIssuedReq) (*emptypb.Empty, error) {
	err := s.srv.RevokeIssued(ctx, req.Id)
	return &emptypb.Empty{}, err
}

func (s *server) GetIssuerList(ctx context.Context, req *emptypb.Empty) (*pb.GetIssuerListReply, error) {
	list, err := s.srv.GetIssuerList(ctx)
	return &pb.GetIssuerListReply{Issuers: list}, err
}

/* challenge 接口*/

func (s *server) ChallengeSend(ctx context.Context, req *pb.ChallengeSendReq) (*emptypb.Empty, error) {
	err := s.srv.ChallengeSend(ctx, req.SenderDid, req.RecipientDid, req.RecipientPubkey, req.Ciphertext, req.Plaintext)
	return &emptypb.Empty{}, err
}
func (s *server) ChallengeReply(ctx context.Context, req *pb.ChallengeReplyReq) (*emptypb.Empty, error) {
	err := s.srv.ChallengeReply(ctx, uint(req.Id), req.Plaintext)
	return &emptypb.Empty{}, err
}
func (s *server) ChallengeSent(ctx context.Context, req *emptypb.Empty) (*pb.ChallengeSentReply, error) {
	c, err := s.srv.ChallengeSent(ctx)
	return &pb.ChallengeSentReply{Challenges: c}, err
}
func (s *server) ChallengeDoing(ctx context.Context, req *emptypb.Empty) (*pb.ChallengeDoingReply, error) {
	c, err := s.srv.ChallengeDoing(ctx)
	return &pb.ChallengeDoingReply{Challenges: c}, err
}
func (s *server) ChallengeDone(ctx context.Context, req *emptypb.Empty) (*pb.ChallengeDoneReply, error) {
	c, err := s.srv.ChallengeDone(ctx)
	return &pb.ChallengeDoneReply{Challenges: c}, err
}

/* proof chaim 接口 */

func (s *server) Apply(ctx context.Context, req *pb.ApplyReq) (*pb.ApplyReply, error) {
	err := s.srv.Apply(ctx, int(req.Type), req.Did, req.Issuer, req.Reason, req.Cids)
	return &pb.ApplyReply{}, err
}

func (s *server) ApplyDoing(ctx context.Context, req *emptypb.Empty) (*pb.ApplyDoingReply, error) {
	items, err := s.srv.ApplyDoing(ctx)
	if err != nil {
		return &pb.ApplyDoingReply{}, err
	}
	return &pb.ApplyDoingReply{Data: items}, nil
}

func (s *server) ApplyDone(ctx context.Context, req *emptypb.Empty) (*pb.ApplyDoneReply, error) {
	items, err := s.srv.ApplyDone(ctx)
	if err != nil {
		return &pb.ApplyDoneReply{}, err
	}
	return &pb.ApplyDoneReply{Data: items}, nil
}

func (s *server) ApplyFailed(ctx context.Context, req *emptypb.Empty) (*pb.ApplyFailedReply, error) {
	items, err := s.srv.ApplyFailed(ctx)
	if err != nil {
		return &pb.ApplyFailedReply{}, err
	}
	return &pb.ApplyFailedReply{Data: items}, nil
}

func (s *server) ApplyCount(ctx context.Context, req *pb.ApplyCountReq) (*pb.ApplyCountReply, error) {
	count, err := s.srv.ApplyCount(ctx, int(req.Type))
	return &pb.ApplyCountReply{Count: count}, err
}

func (s *server) Audit(ctx context.Context, req *pb.AuditReq) (*pb.AuditReply, error) {
	err := s.srv.Audit(ctx, req)
	return &pb.AuditReply{}, err
}

func (s *server) AuditCount(ctx context.Context, req *pb.AuditCountReq) (*pb.AuditCountReply, error) {
	count, err := s.srv.AuditCount(ctx, int(req.Type))
	return &pb.AuditCountReply{Count: count}, err
}

func (s *server) AuditDoing(ctx context.Context, req *emptypb.Empty) (*pb.AuditDoingReply, error) {
	items, err := s.srv.AuditDoing(ctx)
	return &pb.AuditDoingReply{Data: items}, err
}

func (s *server) AuditDone(ctx context.Context, req *emptypb.Empty) (*pb.AuditDoneReply, error) {
	items, err := s.srv.AuditDone(ctx)
	if err != nil {
		return &pb.AuditDoneReply{}, err
	}
	return &pb.AuditDoneReply{Data: items}, nil
}

func (s *server) AuditFailed(ctx context.Context, req *emptypb.Empty) (*pb.AuditFailedReply, error) {
	items, err := s.srv.AuditFailed(ctx)
	if err != nil {
		return &pb.AuditFailedReply{}, err
	}
	return &pb.AuditFailedReply{Data: items}, nil
}
