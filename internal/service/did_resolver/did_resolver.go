package did_resolver

import (
	"context"
	"encoding/json"
	"github.com/treeforest/zut.evidence/api/pb"
	"github.com/treeforest/zut.evidence/internal/service"
	"github.com/treeforest/zut.evidence/internal/service/did_resolver/conf"
	"github.com/treeforest/zut.evidence/internal/service/did_resolver/dao"
	didTool "github.com/treeforest/zut.evidence/pkg/did"
	rsaHelper "github.com/treeforest/zut.evidence/pkg/rsa"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type DIDResolver struct {
	walletClient pb.WalletClient
	dao          *dao.Dao
}

// operationItem 对did的操作记录元素
type operationItem struct {
	DID       string `json:"did"`
	Operation string `json:"operation"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

func New(c *conf.Config, walletClient pb.WalletClient) *DIDResolver {
	r := &DIDResolver{
		dao:          dao.New(c),
		walletClient: walletClient,
	}
	return r
}

func (r *DIDResolver) Close() error {
	r.dao.Close()
	return nil
}

func (r *DIDResolver) CreateDID(ctx context.Context, uid int64) (string, string, error) {
	getPrivateKeyReply, err := r.walletClient.GetPrivateKey(service.NewOutgoingCtxFromInComingCtx(ctx), &emptypb.Empty{})
	if err != nil {
		return "", "", err
	}

	key, _ := rsaHelper.UnmarshalKey(getPrivateKeyReply.PrivateKey)
	doc := didTool.CreateDocumentWithPrivateKey(key)
	docData := doc.Marshal()

	if err = r.dao.CreateDID(uid, doc.Id, string(docData)); err != nil {
		return "", "", err
	}

	// TODO: 将 did document 的哈希上链
	return doc.Id, doc.Created, nil
}

func (r *DIDResolver) RevokeDID(ctx context.Context, did string) error {
	uid := service.GetUid(ctx)

	err := r.dao.ExistDID(did, uid)
	if err != nil {
		return err
	}

	// create signature base data
	op := operationItem{DID: did, Operation: "delete", Timestamp: time.Now().UnixNano()}
	data, _ := json.Marshal(&op)

	// call wallet service and sign base data
	signReply, err := r.walletClient.Sign(service.NewOutgoingCtxFromInComingCtx(ctx), &pb.SignReq{Data: string(data)})
	if err != nil {
		return err
	}

	return r.dao.RevokeDID(did, uid, op.Timestamp, op.Operation, signReply.Signature)
}

func (r *DIDResolver) GetDIDs(ctx context.Context, uid int64) ([]*pb.GetDIDItem, error) {
	dids, err := r.dao.GetDIDs(uid)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.GetDIDItem, 0)
	for _, d := range dids {
		var doc didTool.Document
		_ = json.Unmarshal([]byte(d.Document), &doc)
		items = append(items, &pb.GetDIDItem{Did: d.Id, Created: doc.Created})
	}

	return items, nil
}

func (r *DIDResolver) GetDIDDocument(ctx context.Context, did string) (*didTool.Document, error) {
	docStr, err := r.dao.GetDIDDocument(did)
	if err != nil {
		return nil, err
	}
	var doc didTool.Document
	_ = json.Unmarshal([]byte(docStr), &doc)
	return &doc, nil
}

func (r *DIDResolver) GetPublicKeyByDID(ctx context.Context, did string) (string, error) {
	doc, err := r.GetDIDDocument(ctx, did)
	if err != nil {
		return "", err
	}
	return doc.PublicKey[0].PublicKeyHex, nil
}

func (r *DIDResolver) ExistDiD(ctx context.Context, uid int64, did string) bool {
	return r.dao.ExistDID(did, uid) == nil
}

func (r *DIDResolver) GetOwnerByDID(ctx context.Context, did string) (int64, error) {
	return r.dao.GetOwner(did)
}
