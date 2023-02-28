package dao

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

const (
	tableApplyRelation = "apply_relations"
	tableApplyContent  = "apply_contents"
)

// ApplyStatus 申请状态
type ApplyStatus = int

const (
	DOING  ApplyStatus = 1 // 待审核
	FAILED ApplyStatus = 2 // 审核失败
	DONE   ApplyStatus = 3 // 审核成功
)

// ApplyType 申请类型
type ApplyType = int

const (
	Unknown            ApplyType = 0
	Education          ApplyType = 1 // 学历
	Degree             ApplyType = 2 // 学位
	EducationAndDegree ApplyType = 3 // 学历与学位
	LeaveOffice        ApplyType = 4 // 离职
	InOffice           ApplyType = 5 // 在职
)

var mapping = map[ApplyType]string{
	Education:          "学历证书",
	Degree:             "学位证书",
	EducationAndDegree: "学历证书、学位证书",
	LeaveOffice:        "离职证明",
	InOffice:           "在职证明",
}

func GetApplyTypeText(typ ApplyType) string {
	return mapping[typ]
}

// ApplyContent 申请内容表
type ApplyContent struct {
	gorm.Model
	Applicant   string      // 申请者的 DID
	Issuer      string      // 发行人的 DID
	Type        ApplyType   // 申请类型
	Reason      string      // 申请原因
	Cids        string      // 申请材料
	Status      ApplyStatus // 状态：1: doing; 2:failed; 3:done
	Why         string      // 审核失败的原因
	ProofClaim  string      // 审核成功颁发的凭证
	Transaction []byte      // 交易信息
	Receipt     []byte      // 收据信息
	PdfCid      string
}

// ApplyRelation 申请关系关联表
type ApplyRelation struct {
	ApplyId  uint  // 申请ID
	OwnerUid int64 // 拥有者的 UID
	OtherUid int64 // 关联人的 UID
	Type     int   // 0: 发件箱；1：收件箱
}

// KYC kyc 认证表
type KYC struct {
	Uid       int64 `gorm:"primarykey"`
	CreatedAt time.Time
	Type      int
	Name      string
	IdCard    string
	Cids      string
}

func (d *Dao) createTables() {
	// 如果表不存在，则进行创建
	_ = d.mysql.AutoMigrate(&ApplyContent{})
	_ = d.mysql.AutoMigrate(&ApplyRelation{})
	_ = d.mysql.AutoMigrate(&Challenge{})
	_ = d.mysql.AutoMigrate(&Issuer{})
	_ = d.mysql.AutoMigrate(&KYC{})
}

func (d *Dao) AddKYC(uid int64, typ int, name, idCard, cids string) error {
	return d.mysql.Create(&KYC{Uid: uid, Type: typ, Name: name, IdCard: idCard, Cids: cids}).Error
}

func (d *Dao) GetKYCByUid(uid int64) (*KYC, error) {
	var kyc KYC
	err := d.mysql.Where("uid = ?", uid).Find(&kyc).Error
	return &kyc, err
}

// CreateApply 创建申请
func (d *Dao) CreateApply(senderUid, recipientUid int64, senderDID, recipientDID string, applyType int,
	reason string, cids string) (uint64, error) {

	content := &ApplyContent{
		Applicant:  senderDID,
		Issuer:     recipientDID,
		Type:       applyType,
		Reason:     reason,
		Cids:       cids,
		Status:     DOING,
		Why:        "",
		ProofClaim: "",
	}

	applyRelationSender := &ApplyRelation{
		OwnerUid: senderUid,
		OtherUid: recipientUid,
		Type:     0,
	}

	applyRelationRecipient := &ApplyRelation{
		OwnerUid: recipientUid,
		OtherUid: senderUid,
		Type:     1,
	}

	err := d.mysql.Transaction(func(tx *gorm.DB) error {
		// 存申请内容
		if err := tx.Create(content).Error; err != nil {
			return errors.WithStack(err)
		}

		// 存发件人的发件箱
		applyRelationSender.ApplyId = content.ID
		if err := tx.Create(applyRelationSender).Error; err != nil {
			return errors.WithStack(err)
		}

		// 存收件人的收件箱
		applyRelationRecipient.ApplyId = content.ID
		if err := tx.Create(applyRelationRecipient).Error; err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return uint64(content.ID), nil
}

/** 申请人接口 **/

// GetApplyDoing 获取用户待审核条目
func (d *Dao) GetApplyDoing(uid int64) ([]ApplyContent, error) {
	return d.getContents(DOING, uid, 0)
}

// GetApplyFailed 获取用户审核失败的条目
func (d *Dao) GetApplyFailed(uid int64) ([]ApplyContent, error) {
	return d.getContents(FAILED, uid, 0)
}

// GetApplyDone 获取用户审核成功的条目
func (d *Dao) GetApplyDone(uid int64) ([]ApplyContent, error) {
	return d.getContents(DONE, uid, 0)
}

func (d *Dao) GetApplyDoingCount(uid int64) (int64, error) {
	return d.getContentCount(DOING, uid, 0)
}

func (d *Dao) GetApplyFailedCount(uid int64) (int64, error) {
	return d.getContentCount(FAILED, uid, 0)
}

func (d *Dao) GetApplyDoneCount(uid int64) (int64, error) {
	return d.getContentCount(DONE, uid, 0)
}

/** 发行人接口 **/

func (d *Dao) GetDoingByID(applyId uint) (*ApplyContent, error) {
	var content ApplyContent
	err := d.mysql.Where("id = ? AND status = ?", applyId, DOING).First(&content).Error
	return &content, err
}

func (d *Dao) AuditFailed(applyId uint, uid int64, why string) error {
	var cnt int64

	// 检查用户 uid 是否是关联表中的发行人
	err := d.mysql.Model(&ApplyRelation{}).Where("apply_id = ? AND owner_uid = ? AND type = ?",
		applyId, uid, 1).Count(&cnt).Error
	if err != nil {
		return errors.WithStack(err)
	}

	if cnt == 0 {
		return errors.New("not found")
	}

	return d.mysql.Model(&ApplyContent{}).Where("id = ? AND status = ?", applyId, DOING).
		Select("status", "why").
		Updates(&ApplyContent{Status: FAILED, Why: why}).Error
}

func (d *Dao) AuditDone(applyId uint, proofClaim string, transaction []byte, receipt []byte, pdfCid string) error {
	return d.mysql.Model(&ApplyContent{}).Where("id = ? AND status = ?", applyId, DOING).
		Select("status", "proof_claim", "transaction", "receipt", "pdf_cid").
		Updates(&ApplyContent{Status: DONE, ProofClaim: proofClaim,
			Transaction: transaction, Receipt: receipt, PdfCid: pdfCid}).Error
}

// GetAuditDoing 获取发行人待审核条目
func (d *Dao) GetAuditDoing(uid int64) ([]ApplyContent, error) {
	return d.getContents(DOING, uid, 1)
}

// GetAuditFailed 获取发行人审核失败的条目
func (d *Dao) GetAuditFailed(uid int64) ([]ApplyContent, error) {
	return d.getContents(FAILED, uid, 1)
}

// GetAuditDone 获取发行人审核成功的条目
func (d *Dao) GetAuditDone(uid int64) ([]ApplyContent, error) {
	return d.getContents(DONE, uid, 1)
}

func (d *Dao) GetAuditDoingCount(uid int64) (int64, error) {
	return d.getContentCount(DOING, uid, 1)
}

func (d *Dao) GetAuditFailedCount(uid int64) (int64, error) {
	return d.getContentCount(FAILED, uid, 1)
}

func (d *Dao) GetAuditDoneCount(uid int64) (int64, error) {
	return d.getContentCount(DONE, uid, 1)
}

/** 提炼出来的mysql接口 **/

func (d *Dao) getContents(status ApplyStatus, uid int64, mailType int) ([]ApplyContent, error) {
	contents := make([]ApplyContent, 0)

	subQuery := d.mysql.Select("apply_id as id").
		Where("owner_uid = ? AND type = ?", uid, mailType).Table(tableApplyRelation)

	err := d.mysql.Where("status = ? AND id IN (?)", status, subQuery).Find(&contents).Error

	if err == gorm.ErrRecordNotFound {
		return []ApplyContent{}, nil
	}

	return contents, err
}

func (d *Dao) getContentCount(status ApplyStatus, uid int64, mailType int) (int64, error) {
	var count int64

	subQuery := d.mysql.Select("apply_id as id").
		Where("owner_uid = ? AND type = ?", uid, mailType).Table(tableApplyRelation)

	err := d.mysql.Where("status = ? AND id IN (?)", status, subQuery).Find(&ApplyContent{}).Count(&count).Error

	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}

	return count, err
}
