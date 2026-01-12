/*
 * @Descripttion: JWT 生成与校验
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 10:45:00
 */
package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

var ExpireTime int64 = 24

func getJwtSecret() []byte {
	// 不能在 package init 时读取：config/viper 可能尚未初始化，导致 secret 为空，
	// 从而出现“看似没多久就过期”（实际上是签名校验失败）的现象。
	secret := viper.GetString("jwt.secret")
	return []byte(secret)
}

// Claims ...
type Claims struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Authority int    `json:"authority"`
	jwt.StandardClaims
}

// GenerateToken 签发用户Token
func GenerateToken(username, password string, authority int) (string, int64, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(ExpireTime) * time.Hour * 90)

	claims := Claims{
		username,
		password,
		authority,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "cmall",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(getJwtSecret())

	return token, expireTime.Unix(), err
}

// ParseToken 验证用户token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 防止 alg 混淆（只允许 HMAC 系列）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getJwtSecret(), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// 刷新产生新的token
func RefreshToken(tokenString string) (string, int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getJwtSecret(), nil
	})
	if err != nil {
		return "", 0, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return GenerateToken(claims.Username, claims.Password, 1)
	}
	return "", 0, errors.New("token刷新失败")
}

// EmailClaims ...
type EmailClaims struct {
	UserID        uint   `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OperationType uint   `json:"operation_type"`
	jwt.StandardClaims
}

// GenerateEmailToken 签发邮箱验证Token
func GenerateEmailToken(userID, Operation uint, email, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(15 * time.Minute)

	claims := EmailClaims{
		userID,
		email,
		password,
		Operation,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "cmall",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(getJwtSecret())

	return token, err
}

// ParseEmailToken 验证邮箱验证token
func ParseEmailToken(token string) (*EmailClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getJwtSecret(), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*EmailClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
