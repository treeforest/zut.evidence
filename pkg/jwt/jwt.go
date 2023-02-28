package jwt

import (
	"fmt"
	"github.com/pkg/errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secret                 = "197346825zoqlapskxjdmcnfhbeiurys"
	DefaultTokenExpiration = 60 * 60 * 24 // 单位：秒, 24h
)

// JWTManager is a JSON web token manager
type JWTManager struct {
	// tokenDuration token过期时间
	tokenDuration time.Duration
}

// UserClaims jwt令牌中所包含的自定义信息
type UserClaims struct {
	jwt.StandardClaims
	Uid      int64  `json:"uid"`
	Role     int    `json:"role"`
	Platform string `json:"platform"`
	Extra    []byte `json:"extra"`
}

// New returns a new JWT manager
func New(tokenDuration time.Duration) *JWTManager {
	return &JWTManager{tokenDuration}
}

// Generate generates and signs a new token for a account
func (m *JWTManager) Generate(uid int64, role int, platform string, extra []byte) (string, error) {
	expiresAt := time.Now().Add(m.tokenDuration).Unix() // expire time
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		Uid:      uid,
		Role:     role,
		Platform: platform,
		Extra:    extra,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// Verify verifies the signed token string and return a account claim detail if the token is valid
func (m *JWTManager) Verify(signedToken string) (
	uid int64, role int, platform string, extra []byte, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return 0, 0, "", nil, errors.WithStack(err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return 0, 0, "", nil, errors.New("unexpected Claims")
	}
	return claims.Uid, claims.Role, claims.Platform, claims.Extra, nil
}
