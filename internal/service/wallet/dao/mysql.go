package dao

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Wallet 钱包对象
type Wallet struct {
	Uid int64 `gorm:"primaryKey"`
	Key string
}

func (d *Dao) createTables() {
	// 如果表不存在，则进行创建
	_ = d.mysql.AutoMigrate(&Wallet{})
}

// ExistKey 判断是否存在key
func (d *Dao) ExistKey(uid int64) (bool, error) {
	var w Wallet
	err := d.mysql.Select("uid").Where("uid = ?", uid).First(&w).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return false, status.Errorf(codes.Internal, "查询失败: %v", err)
		}
		return false, nil
	}
	return true, nil
}

// AddKey 添加私钥
func (d *Dao) AddKey(uid int64, key string) error {
	return d.mysql.Transaction(func(tx *gorm.DB) error {
		var w Wallet
		err := tx.Where("uid = ?", uid).First(&w).Error
		if err == nil {
			return status.Error(codes.InvalidArgument, "密钥已存在")
		}
		if err != gorm.ErrRecordNotFound {
			return status.Errorf(codes.Internal, "查询失败: %v", err)
		}
		w.Uid = uid
		w.Key = key
		if err = tx.Create(&w).Error; err != nil {
			return status.Errorf(codes.Internal, "转储失败: %v", err)
		}
		return nil
	})
}

// GetKey 获取私钥
func (d *Dao) GetKey(uid int64) (key string, err error) {
	var w Wallet
	err = d.mysql.Where("uid = ?", uid).First(&w).Error
	if err == gorm.ErrRecordNotFound {
		return "", status.Error(codes.InvalidArgument, "还未创建密钥")
	}
	return w.Key, err
}
