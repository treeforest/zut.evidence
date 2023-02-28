package dao

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Did did列表
type Did struct {
	Id       string `gorm:"primaryKey"`
	Owner    int64
	Document string
	Created  int64 `gorm:"autoCreateTime:milli"`
}

// Revoke did吊销列表
type Revoke struct {
	Id        string `gorm:"primaryKey"`
	Owner     int64
	Timestamp int64
	Operation string
	Signature string
}

func (d *Dao) createTables() {
	_ = d.mysql.AutoMigrate(&Did{})
	_ = d.mysql.AutoMigrate(&Revoke{})
}

// CreateDID 创建一个新的did
func (d *Dao) CreateDID(uid int64, did, document string) error {
	err := d.mysql.Create(&Did{Id: did, Owner: uid, Document: document}).Error
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

func (d *Dao) GetDIDs(uid int64) ([]Did, error) {
	dids := make([]Did, 0)
	err := d.mysql.Where("owner = ?", uid).Find(&dids).Error
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return dids, nil
}

func (d *Dao) GetDIDDocument(did string) (string, error) {
	var o Did
	err := d.mysql.Where("id = ?", did).First(&o).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", status.Error(codes.InvalidArgument, "DID 不存在")
		}
		return "", err
	}
	return o.Document, nil
}

func (d *Dao) GetOwner(did string) (int64, error) {
	var o Did
	err := d.mysql.Where("id = ?", did).Find(&o).Error
	return o.Owner, err
}

func (d *Dao) ExistDID(did string, uid int64) error {
	var o Did
	err := d.mysql.Select("owner").Where("id = ?", did).First(&o).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return status.Error(codes.InvalidArgument, "DID 不存在")
		}
		return err
	}
	if o.Owner != uid {
		return status.Error(codes.InvalidArgument, "不具备对该DID的操作权限")
	}
	return nil
}

func (d *Dao) RevokeDID(did string, uid, timestamp int64, op, sig string) error {
	return d.mysql.Transaction(func(tx *gorm.DB) error {
		tx.Where("id = ?", did).Delete(&Did{})
		return tx.Create(&Revoke{
			Id:        did,
			Owner:     uid,
			Timestamp: timestamp,
			Operation: op,
			Signature: sig,
		}).Error
	})
}
