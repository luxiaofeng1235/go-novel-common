//Package middleware ...
/*
 * @Descripttion:
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: congz
 * @LastEditTime: 2020-08-08 16:33:29
 */
package middleware

import (
	"github.com/gin-gonic/gin"
	"go-novel/utils"
	"go-novel/utils/e"
	"net/http"
	"strings"
	"time"
)

// jwt token验证中间件1
func AdminJwt() gin.HandlerFunc {
	return func(c *gin.Context) {

		var code int
		code = e.Success

		//效验token
		//token := c.GetHeader("Token")
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusOK, gin.H{
				"code": e.NotPermission,
				"msg":  "请先进行登录",
			})
			c.Abort()
			return
		}

		token := tokenString[7:]
		//parts := strings.Split(tokenString, " ")
		//token := parts[1]

		if token == "" {
			code = e.NotPermission
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": e.NotPermission,
					"msg":  "登陆已过期，请重新登陆",
					"err":  err.Error(),
				})

				c.Abort()
				return
			}

			if time.Now().Unix() > claims.ExpiresAt {
				//token过期
				code = e.NotPermission
				c.JSON(http.StatusOK, gin.H{
					"code": code,
					"msg":  "登陆已过期，请重新登陆！",
					"err":  err.Error(),
				})

				c.Abort()
				return
			}

			username := claims.Username
			//log.Println(username)
			c.Set("username", username)
		}

		if code != e.Success {
			c.JSON(http.StatusOK, gin.H{
				"code": e.NotPermission,
				"msg":  "请先进行登录",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
