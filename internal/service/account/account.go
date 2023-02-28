package account

import (
	"context"
	"github.com/treeforest/zut.evidence/internal/service/account/conf"
	"github.com/treeforest/zut.evidence/internal/service/account/dao"
	"github.com/treeforest/zut.evidence/pkg/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Account struct {
	dao *dao.Dao
}

func New(c *conf.Config, jwtMgr *jwt.JWTManager) *Account {
	ac := &Account{
		dao: dao.New(c, jwtMgr),
	}
	return ac
}

// Login 登录
func (ac *Account) Login(ctx context.Context, nick string, password, platform string) (token string, role int, err error) {
	// 参数合法性检查
	if err = verifyNick(nick); err != nil {
		return "", 0, err
	}
	if err = verifyPassword(password); err != nil {
		return "", 0, err
	}
	if err = verifyPlatform(platform); err != nil {
		return "", 0, err
	}

	uid, role, err := ac.dao.Login(nick, password)
	if err != nil {
		return "", 0, err
	}

	token, err = ac.dao.GenerateToken(uid, role, platform, []byte{}) // 生成token
	if err != nil {
		return "", 0, status.Error(codes.Internal, "token生成失败")
	}

	return token, role, nil
}

// Register 注册
func (ac *Account) Register(ctx context.Context, nick, phone, email, password string, role int) error {
	// 参数合法性检查
	if err := verifyNick(nick); err != nil {
		return err
	}
	if err := verifyPhone(phone); err != nil {
		return err
	}
	//if err := verifyEmail(email); err != nil {
	//	return err
	//}
	if err := verifyPassword(password); err != nil {
		return err
	}
	if err := verifyRole(role); err != nil {
		return err
	}

	return ac.dao.Register(nick, phone, email, password, role)
}

// Close 关闭
func (ac *Account) Close() error {
	ac.dao.Close()
	return nil
}
