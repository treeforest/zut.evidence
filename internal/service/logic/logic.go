package logic

import (
	"context"
	"github.com/treeforest/zut.evidence/blockchain/contracts/evidence"
	"github.com/treeforest/zut.evidence/internal/service/logic/conf"
	"github.com/treeforest/zut.evidence/internal/service/logic/dao"
	"strings"

	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service"
)

type Logic struct {
	dao                *dao.Dao
	didResolver        pb.DIDResolverClient
	cometClient        pb.CometClient
	walletClient       pb.WalletClient
	fileClient         pb.FileClient
	getEvidenceSession func() (*evidence.EvidenceSession, error)
}

func New(c *conf.Config, didResolver pb.DIDResolverClient, walletClient pb.WalletClient,
	cometClient pb.CometClient, fileClient pb.FileClient, getEvidenceSession func() (*evidence.EvidenceSession, error)) *Logic {

	return &Logic{
		dao:                dao.New(c),
		didResolver:        didResolver,
		cometClient:        cometClient,
		walletClient:       walletClient,
		fileClient:         fileClient,
		getEvidenceSession: getEvidenceSession,
	}
}

func (l *Logic) Close() error {
	l.dao.Close()
	return nil
}

func (l *Logic) ApplyKYC(ctx context.Context, typ int, name, idCard string, cids []string) error {
	uid := service.GetUid(ctx)
	return l.dao.AddKYC(uid, typ, name, idCard, strings.Join(cids, ";"))
}
