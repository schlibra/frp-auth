package utils

import (
	"errors"
	"fmt"
	"frp-auth/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func JwtUserGenerate(userId int, username string, admin bool, enable bool, tokenVersion string) (string, error) {
	now := time.Now()
	cfg := LoadConfig()
	sk := []byte(cfg.Jwt.Key)
	claims := model.JwtToken{
		UserID:       userId,
		Username:     username,
		Admin:        admin,
		Enable:       enable,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "frp-auth",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(sk)
	if err != nil {
		return "", err
	}
	return signedString, nil
}
func JwtUserCheck(tokenStr string) (*model.JwtToken, error) {
	cfg := LoadConfig()
	sk := []byte(cfg.Jwt.Key)
	token, err := jwt.ParseWithClaims(tokenStr, &model.JwtToken{}, func(token *jwt.Token) (any, error) {
		// 必须检查签名算法，防御 "alg: none" 或算法混淆攻击
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return sk, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token 已过期")
		}
		return nil, fmt.Errorf("token 无效: %w", err)
	}

	// 提取并校验有效载荷
	if claims, ok := token.Claims.(*model.JwtToken); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无法解析 token claims")
}
