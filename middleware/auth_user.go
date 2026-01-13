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

// 兼容原项目约定：在 Context 中注入 user_id，供 controller 层直接 c.Get("user_id") 获取
const ctxKeyUserID = "user_id"
const ctxKeyUsername = "username"

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

		c.Set(ctxKeyUsername, strings.TrimSpace(claims.Username))
		c.Set(ctxKeyUserID, claims.UserID)
		c.Next()
	}
}

// ApiJwt 兼容原项目命名：API 侧 JWT 鉴权（注入 user_id）
func ApiJwt() gin.HandlerFunc {
	return AuthUser()
}

func GetAuthUsername(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if val, ok := c.Get(ctxKeyUsername); ok {
		if s, ok := val.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
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
