/*
 * @Descripttion: API Token 鉴权（从请求头解析 JWT，并注入到 gin.Context）
 * @Author: red
 * @Date: 2026-01-13 09:40:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-13 09:40:00
 */
package middleware

import (
	"strings"

	"go-novel/utils"

	"github.com/gin-gonic/gin"
)

const ctxKeyAuthUsername = "authUsername"
const ctxKeyAuthUserID = "authUserID"

// AuthUser 必须携带有效 token（支持 Authorization: Bearer xxx 或 Token: xxx）
func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			utils.FailEncrypt(c, nil, "缺少token")
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil || claims == nil || strings.TrimSpace(claims.Username) == "" {
			utils.FailEncrypt(c, nil, "token无效")
			c.Abort()
			return
		}
		if claims.UserID <= 0 {
			utils.FailEncrypt(c, nil, "token无效")
			c.Abort()
			return
		}

		c.Set(ctxKeyAuthUsername, strings.TrimSpace(claims.Username))
		c.Set(ctxKeyAuthUserID, claims.UserID)
		c.Next()
	}
}

func GetAuthUsername(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if val, ok := c.Get(ctxKeyAuthUsername); ok {
		if s, ok := val.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func GetAuthUserID(c *gin.Context) int64 {
	if c == nil {
		return 0
	}
	if val, ok := c.Get(ctxKeyAuthUserID); ok {
		switch v := val.(type) {
		case int64:
			return v
		case int:
			return int64(v)
		case float64:
			return int64(v)
		}
	}
	return 0
}

func extractToken(c *gin.Context) string {
	if c == nil {
		return ""
	}

	// Authorization: Bearer <token>
	auth := strings.TrimSpace(c.GetHeader("Authorization"))
	if auth != "" {
		lower := strings.ToLower(auth)
		if strings.HasPrefix(lower, "bearer ") {
			return strings.TrimSpace(auth[len("bearer "):])
		}
		// 兼容直接传 token
		return auth
	}

	// 兼容历史 Token 头
	if t := strings.TrimSpace(c.GetHeader("Token")); t != "" {
		return t
	}
	return ""
}
