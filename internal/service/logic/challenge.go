package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/treeforest/logger"
	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (l *Logic) ChallengeSend(ctx context.Context, senderDID, recipientDID, recipientPubkey, ciphertext, plaintext string) error {
	senderUid := service.GetUid(ctx)
	outgoingCtx := service.NewOutgoingCtxFromInComingCtx(ctx)

	reply, err := l.didResolver.GetOwnerByDID(outgoingCtx, &pb.GetOwnerByDIDReq{Did: recipientDID})
	if err != nil {
		log.Errorf("ExistDID | ERR:%v", err)
		return status.Error(codes.Internal, err.Error())
	}
	if reply.Uid == 0 {
		return errors.New("不存在接收者DID对应的用户")
	}

	_ = l.dao.AddChallenge(senderUid, reply.Uid, senderDID, recipientDID, recipientPubkey, ciphertext, plaintext)

	// 推送通知
	_, err = l.cometClient.PushMsg(context.Background(), &pb.PushMsgReq{
		Uid:  reply.Uid,
		Type: pb.PushMsgReq_Challenge_REQ,
		Msg:  fmt.Sprintf("%s 向你发起挑战", senderDID),
	})
	if err != nil {
		log.Errorf("PushMsg | ERR:%v", err)
		return status.Error(codes.Internal, "推送失败")
	}

	return nil
}

func (l *Logic) ChallengeReply(ctx context.Context, id uint, plaintext string) error {
	uid := service.GetUid(ctx)

	challenge, err := l.dao.GetChallengeById(id)
	if err != nil {
		return err
	}
	if challenge == nil {
		return status.Error(codes.NotFound, "not found")
	}
	if uid != challenge.RecipientUid {
		return status.Error(codes.InvalidArgument, "")
	}

	pushMsg := ""
	if challenge.Plaintext == plaintext {
		_ = l.dao.SetChallengeStatus(id, 2)
		pushMsg = fmt.Sprintf("%s 挑战成功", challenge.RecipientDid)
	} else {
		_ = l.dao.SetChallengeStatus(id, 1)
		pushMsg = fmt.Sprintf("%s 挑战失败", challenge.RecipientDid)
	}

	// 推送通知
	_, err = l.cometClient.PushMsg(context.Background(), &pb.PushMsgReq{
		Uid:  challenge.SenderUid,
		Type: pb.PushMsgReq_Challenge_REQ,
		Msg:  pushMsg,
	})
	if err != nil {
		log.Errorf("PushMsg | ERR:%v", err)
		return status.Error(codes.Internal, "推送失败")
	}

	return nil
}

func (l *Logic) ChallengeSent(ctx context.Context) ([]*pb.ChallengeSentReply_Challenge, error) {
	uid := service.GetUid(ctx)

	challenges, err := l.dao.GetChallengeSent(uid)
	if err != nil {
		return nil, err
	}

	vo := make([]*pb.ChallengeSentReply_Challenge, 0, len(challenges))
	for _, c := range challenges {
		status := "待验证"
		if c.Status == 1 {
			status = "验证失败"
		} else if c.Status == 2 {
			status = "验证成功"
		}

		vo = append(vo, &pb.ChallengeSentReply_Challenge{
			Id:              uint64(c.ID),
			SenderDid:       c.SenderDid,
			RecipientDid:    c.RecipientDid,
			RecipientPubKey: c.RecipientPubKey,
			Plaintext:       c.Plaintext,
			CreatedTime:     c.CreatedAt.Format("2006-01-02 15:04:05"),
			Status:          status,
		})
	}

	return vo, nil
}

func (l *Logic) ChallengeDoing(ctx context.Context) ([]*pb.ChallengeDoingReply_Challenge, error) {
	uid := service.GetUid(ctx)

	challenges, err := l.dao.GetChallengeDoing(uid)
	if err != nil {
		return nil, err
	}

	vo := make([]*pb.ChallengeDoingReply_Challenge, 0, len(challenges))
	for _, c := range challenges {
		vo = append(vo, &pb.ChallengeDoingReply_Challenge{
			Id:          uint64(c.ID),
			SenderDid:   c.SenderDid,
			CreatedTime: c.CreatedAt.Format("2006-01-02 15:04:05"),
			Ciphertext:  c.Ciphertext,
		})
	}

	return vo, nil
}

func (l *Logic) ChallengeDone(ctx context.Context) ([]*pb.ChallengeDoneReply_Challenge, error) {
	uid := service.GetUid(ctx)

	challenges, err := l.dao.GetChallengeDone(uid)
	if err != nil {
		return nil, err
	}

	vo := make([]*pb.ChallengeDoneReply_Challenge, 0, len(challenges))
	for _, c := range challenges {
		vo = append(vo, &pb.ChallengeDoneReply_Challenge{
			Id:          uint64(c.ID),
			SenderDid:   c.SenderDid,
			ReceiptDid:  c.RecipientDid,
			CreatedTime: c.CreatedAt.Format("2006-01-02 15:04:05"),
			Ciphertext:  c.Ciphertext,
			Status:      int32(c.Status),
		})
	}

	return vo, nil
}
