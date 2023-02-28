package dao

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// User 用户表ORM
type User struct {
	Uid      int64  `gorm:"primaryKey;uniqueIndex;"`
	Nick     string `gorm:"unique"`
	Role     int
	Phone    string `gorm:"unique"`
	Email    string
	Password string
	Created  int64 `gorm:"autoCreateTime:milli"`
	Updated  int64 `gorm:"autoUpdateTime:milli"`
}

// Unique 唯一ID ORM
type Unique struct {
	Type string `gorm:"primaryKey"`
	Id   int64
}

func hash(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func (d *Dao) createTables() {
	// 如果表不存在，则进行创建
	_ = d.mysql.AutoMigrate(&User{})
	_ = d.mysql.AutoMigrate(&Unique{})
}

// Login 登录 用户名+密码+角色
func (d *Dao) Login(nick, password string) (int64, int, error) {
	var user User
	err := d.mysql.Select("uid", "password", "role").
		Where("nick = ?", nick).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return 0, 0, status.Error(codes.InvalidArgument, "用户名不存在")
	}
	if user.Password != hash([]byte(password)) {
		return 0, 0, status.Error(codes.InvalidArgument, "密码错误")
	}
	return user.Uid, user.Role, nil
}

// Register 注册
func (d *Dao) Register(nick, phone, email, password string, role int) error {
	return d.mysql.Transaction(func(tx *gorm.DB) error {
		var tmp User

		if err := tx.Select("uid").Where("nick = ?", nick).First(&tmp).Error; err != gorm.ErrRecordNotFound {
			if err == nil {
				return status.Error(codes.InvalidArgument, "用户名已注册")
			}
			return status.Error(codes.Internal, err.Error())
		}

		if err := tx.Select("uid").Where("phone = ?", phone).First(&tmp).Error; err != gorm.ErrRecordNotFound {
			if err == nil {
				return status.Error(codes.InvalidArgument, "手机号已注册")
			}
			return status.Error(codes.Internal, err.Error())
		}

		//if err := tx.Select("uid").Where("email = ?", email).First(&tmp).Error; err != gorm.ErrRecordNotFound {
		//	if err == nil {
		//		return status.Error(codes.InvalidArgument, "邮箱已注册")
		//	}
		//	return status.Error(codes.Internal, err.Error())
		//}

		uniqueUid := &Unique{Type: "uid"}
		if err := tx.First(uniqueUid).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return status.Error(codes.Internal, err.Error())
			}
			// 创建
			uniqueUid.Id = 10000
			if err = tx.Create(uniqueUid).Error; err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		}

		uniqueUid.Id++

		if err := tx.Model(&uniqueUid).Update("id", uniqueUid.Id).Error; err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		user := User{
			Uid:      uniqueUid.Id,
			Nick:     nick,
			Phone:    phone,
			Email:    email,
			Password: hash([]byte(password)),
			Role:     role,
		}

		if err := d.mysql.Create(&user).Error; err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
}
