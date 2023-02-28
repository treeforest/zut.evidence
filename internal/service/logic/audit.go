package logic

import (
	"context"
	"encoding/json"
	log "github.com/treeforest/logger"
	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service"
	"github.com/treeforest/zut.evidence/internal/service/logic/dao"
	"github.com/treeforest/zut.evidence/pkg/did"
	"github.com/treeforest/zut.evidence/pkg/pdf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
	"strings"
)

func (l *Logic) Audit(ctx context.Context, req *pb.AuditReq) error {
	uid := service.GetUid(ctx)

	if req.Result == 1 {
		// 审核失败
		return l.dao.AuditFailed(uint(req.ApplyId), uid, req.Why)
	}

	// log.Debug(*req)

	// 审核成功
	// 1. 获取对应申请
	var content *dao.ApplyContent
	content, err := l.dao.GetDoingByID(uint(req.ApplyId))
	if err != nil {
		log.Warn(err)
		return status.Error(codes.InvalidArgument, "审核失败")
	}

	// 2. 获取到公钥id
	outgoingCtx := service.NewOutgoingCtxFromInComingCtx(ctx)
	getDidDocReply, err := l.didResolver.GetDIDDocument(outgoingCtx, &pb.GetDIDDocumentReq{Did: content.Issuer})
	if err != nil {
		log.Warn(err)
		return err
	}

	pubKeyId := getDidDocReply.DidDocument.PublicKey[0].Id

	// 3. 生成 base proof claim
	s := strings.Replace(strings.Replace(req.Expiration, `"`, "", -1), "T", " ", -1)
	expiration := strings.Split(s, ".")[0]
	proofClaim := did.BaseProofClaim(expiration, dao.GetApplyTypeText(content.Type), content.Issuer, pubKeyId,
		content.Applicant, req.ShortDesc, req.LongDesc)

	// 4. 使用私钥签名
	signReq, err := l.walletClient.Sign(outgoingCtx, &pb.SignReq{Data: string(proofClaim.Marshal())})
	if err != nil {
		log.Warn(err)
		return err
	}
	proofClaim.Proof.Signature = signReq.Signature

	// 5. 生成pdf证书
	pdfData := pdf.GenEvidencePdf(&pdf.PdfInfo{ProofClaim: proofClaim})
	uploadReply, err := l.fileClient.Upload(context.Background(), &pb.UploadReq{Uid: 101, Filename: "certificate.pdf", Data: pdfData}) // TODO: 这里的UID仅用于测试
	if err != nil {
		log.Warn(err)
		return err
	}

	// 6. TODO: 生成证书
	session, err := l.getEvidenceSession()
	if err != nil {
		log.Error(err)
		return err
	}
	transaction, receipt, err := session.Add(big.NewInt(uid), proofClaim.Id)
	if err != nil {
		log.Error(err)
		return status.Error(codes.Internal, "上链失败")
	}
	transactionData, _ := transaction.MarshalJSON()
	receiptData, _ := json.Marshal(receipt)

	// 6. 存入数据库
	err = l.dao.AuditDone(uint(req.ApplyId), string(proofClaim.Marshal()), transactionData, receiptData, uploadReply.Cid)

	return err
}

func (l *Logic) ApplyCount(ctx context.Context, typ int) (int64, error) {
	uid := service.GetUid(ctx)

	switch typ {
	case 1:
		return l.dao.GetApplyDoingCount(uid)
	case 2:
		return l.dao.GetApplyFailedCount(uid)
	case 3:
		return l.dao.GetApplyDoneCount(uid)
	}

	return 0, status.Error(codes.InvalidArgument, "参数错误")
}

func (l *Logic) AuditCount(ctx context.Context, typ int) (int64, error) {
	uid := service.GetUid(ctx)

	switch typ {
	case 1:
		return l.dao.GetAuditDoingCount(uid)
	case 2:
		return l.dao.GetAuditFailedCount(uid)
	case 3:
		return l.dao.GetAuditDoneCount(uid)
	}

	return 0, status.Error(codes.InvalidArgument, "参数错误")
}

func (l *Logic) AuditDoing(ctx context.Context) ([]*pb.DoingItem, error) {
	uid := service.GetUid(ctx)

	log.Debug("uid: ", uid)

	contents, err := l.dao.GetAuditDoing(uid)
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

func (l *Logic) AuditDone(ctx context.Context) ([]*pb.DoneItem, error) {
	uid := service.GetUid(ctx)

	contents, err := l.dao.GetAuditDone(uid)
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

func (l *Logic) AuditFailed(ctx context.Context) ([]*pb.FailedItem, error) {
	uid := service.GetUid(ctx)

	contents, err := l.dao.GetAuditFailed(uid)
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
