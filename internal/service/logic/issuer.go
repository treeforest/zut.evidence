package logic

import (
	"context"
	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service"
)

func (l *Logic) Issue(ctx context.Context, did, shortDesc, longDesc, endpoint, website string, typ int) error {
	uid := service.GetUid(ctx)
	return l.dao.AddIssuer(uid, did, website, endpoint, shortDesc, longDesc, typ)
}

func (l *Logic) GetIssued(ctx context.Context) ([]*pb.GetIssuedReply_Issuer, error) {
	uid := service.GetUid(ctx)
	issuers, err := l.dao.GetIssued(uid)
	if err != nil {
		return []*pb.GetIssuedReply_Issuer{}, err
	}

	vo := make([]*pb.GetIssuedReply_Issuer, 0, len(issuers))
	for _, issuer := range issuers {
		vo = append(vo, &pb.GetIssuedReply_Issuer{
			Id:               uint64(issuer.ID),
			Did:              issuer.Did,
			Website:          issuer.Website,
			Endpoint:         issuer.Endpoint,
			ShortDescription: issuer.ShortDesc,
			LongDescription:  issuer.LongDesc,
			Type:             pb.ProofClaimType(issuer.Type),
			CreateTime:       issuer.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return vo, nil
}

func (l *Logic) RevokeIssued(ctx context.Context, id uint64) error {
	return l.dao.DelIssued(id)
}

func (l *Logic) GetIssuerList(ctx context.Context) ([]*pb.GetIssuerListReply_Issuer, error) {
	issuers, err := l.dao.GetIssuers()
	if err != nil {
		return []*pb.GetIssuerListReply_Issuer{}, nil
	}

	vo := make([]*pb.GetIssuerListReply_Issuer, 0, len(issuers))
	for _, issuer := range issuers {
		vo = append(vo, &pb.GetIssuerListReply_Issuer{
			Did:              issuer.Did,
			Website:          issuer.Website,
			Endpoint:         issuer.Website,
			ShortDescription: issuer.ShortDesc,
			LongDescription:  issuer.LongDesc,
			Type:             pb.ProofClaimType(issuer.Type),
			CreateTime:       issuer.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return vo, nil
}
