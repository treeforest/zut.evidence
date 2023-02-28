package logic

import (
	"context"
	"fmt"
	log "github.com/treeforest/logger"
	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service"
	"github.com/treeforest/zut.evidence/pkg/did"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

// Apply 申请 proof claim
func (l *Logic) Apply(ctx context.Context, typ int, did, issuer, reason string, cids []string) error {

	outgoingCtx := service.NewOutgoingCtxFromInComingCtx(ctx)
	uid := service.GetUid(ctx)

	// 判断 did 是否属于当前 uid
	getOwnerByDIDReply, err := l.didResolver.GetOwnerByDID(outgoingCtx, &pb.GetOwnerByDIDReq{Did: did})
	if err != nil {
		log.Errorf("ExistDID | ERR:%v", err)
		return status.Error(codes.Internal, err.Error())
	}
	if getOwnerByDIDReply.Uid != uid {
		log.Debugf("required: %d actual: %d", uid, getOwnerByDIDReply.Uid)
		return status.Error(codes.InvalidArgument, "申请人的DID无效")
	}

	// 判断 issuer 是否存在
	getOwnerByDIDReply, err = l.didResolver.GetOwnerByDID(outgoingCtx, &pb.GetOwnerByDIDReq{Did: issuer})
	if err != nil {
		log.Errorf("GetOwnerByDID | ERR:%v", err)
		return status.Error(codes.Internal, err.Error())
	}
	if getOwnerByDIDReply.Uid == 0 {
		return status.Error(codes.InvalidArgument, "发行人的DID无效")
	}

	// 存入DB
	senderUid, recipientUid, senderDID, recipientDID := uid, getOwnerByDIDReply.Uid, did, issuer
	applyId, err := l.dao.CreateApply(senderUid, recipientUid, senderDID, recipientDID, typ, reason, strings.Join(cids, ";"))
	if err != nil {
		log.Errorf("CreateApply | ERR:%v", err)
		return status.Error(codes.Internal, "插入失败")
	}

	// 推送通知
	_, err = l.cometClient.PushMsg(context.Background(), &pb.PushMsgReq{
		Uid:  recipientUid,
		Type: pb.PushMsgReq_Audit_DOING,
		Msg:  fmt.Sprintf("您有一条编号为 %d 的申请待审核", applyId),
	})
	if err != nil {
		log.Errorf("PushMsg | ERR:%v", err)
		return status.Error(codes.Internal, "推送失败")
	}

	return nil
}

// ApplyDoing 申请者查询其未审核的 proof claim
func (l *Logic) ApplyDoing(ctx context.Context) ([]*pb.DoingItem, error) {
	uid := service.GetUid(ctx)

	contents, err := l.dao.GetApplyDoing(uid)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询失败")
	}

	items := make([]*pb.DoingItem, 0)

	for _, c := range contents {
		items = append(items, &pb.DoingItem{
			ApplyId: uint64(c.ID),
			Did:     c.Applicant,
			Issuer:  c.Issuer,
			Type:    pb.ProofClaimType(c.Type),
			Reason:  c.Reason,
			Cids:    strings.Split(c.Cids, ";"),
			Time:    c.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return items, nil
}

// ApplyFailed 申请者查询其审核失败的 proof claim
func (l *Logic) ApplyFailed(ctx context.Context) ([]*pb.FailedItem, error) {
	uid := service.GetUid(ctx)

	contents, err := l.dao.GetApplyFailed(uid)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询失败")
	}

	items := make([]*pb.FailedItem, 0)

	for _, c := range contents {
		items = append(items, &pb.FailedItem{
			ApplyId: uint64(c.ID),
			Did:     c.Applicant,
			Issuer:  c.Issuer,
			Type:    pb.ProofClaimType(c.Type),
			Reason:  c.Reason,
			Cids:    strings.Split(c.Cids, ";"),
			Time:    c.CreatedAt.Format("2006-01-02 15:04:05"),
			Why:     c.Why,
		})
	}

	return items, nil
}

func (l *Logic) ApplyDone(ctx context.Context) ([]*pb.DoneItem, error) {
	uid := service.GetUid(ctx)

	contents, err := l.dao.GetApplyDone(uid)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询失败")
	}

	items := make([]*pb.DoneItem, 0)

	for _, c := range contents {
		pc := did.ProofClaim{}
		if err := pc.Unmarshal([]byte(c.ProofClaim)); err != nil {
			log.Error("ProofClaim unmarshal: ", err)
			continue
		}
		items = append(items, &pb.DoneItem{
			ApplyId: uint64(c.ID),
			Did:     c.Applicant,
			Issuer:  c.Issuer,
			Type:    pb.ProofClaimType(c.Type),
			Reason:  c.Reason,
			Cids:    strings.Split(c.Cids, ";"),
			Time:    c.CreatedAt.Format("2006-01-02 15:04:05"),
			ProofClaim: &pb.ProofClaim{
				Context:        pc.Context,
				Id:             pc.Id,
				Type:           pc.Type,
				Issuer:         pc.Issuer,
				IssuanceData:   pc.IssuanceDate,
				ExpirationData: pc.ExpirationDate,
				CredentialSubject: &pb.ProofClaim_CredentialSubject{
					Id:               pc.CredentialSubject.Id,
					ShortDescription: pc.CredentialSubject.ShortDescription,
					LongDescription:  pc.CredentialSubject.LongDescription,
					Type:             pc.CredentialSubject.Type,
				},
				Proof: &pb.ProofClaim_Proof{
					Type:      pc.Proof.Type,
					Creator:   pc.Proof.Creator,
					Signature: pc.Proof.Signature,
				},
			},
			Transaction: string(c.Transaction),
			Receipt:     string(c.Receipt),
			PdfCid:      c.PdfCid,
		})
	}

	return items, nil
}
