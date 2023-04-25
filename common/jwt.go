package common

import (
	"BookcaseServer/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("bookcase")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

// 单Token：72小时 双Token：2小时 + 7 × 24小时(14 × 24小时)

// TokenExpirationTime token过期时间
const TokenExpirationTime time.Duration = 120 * time.Minute

// RefreshTokenExpirationTime refresh_token过期时间
const RefreshTokenExpirationTime time.Duration = 7 * 24 * time.Hour

// ReleaseToken 生成并发放token
func ReleaseToken(user model.User, expirationTime time.Duration) (string, error) {
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expirationTime).Unix(), //设置token的有效时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "bookcase",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}
