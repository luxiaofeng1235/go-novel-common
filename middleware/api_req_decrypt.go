package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/utils"
	"io/ioutil"
)

func ApiReqDecrypt() gin.HandlerFunc {
	return func(c *gin.Context) {
		encrypt := utils.GetApiEncrypt()
		if encrypt == false {
			c.Next()
			return
		}
		// 读取原始请求 Body 数据
		var err error
		var req models.ApiAesReq
		// 参数绑定
		if err = c.ShouldBind(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			c.Abort()
			return
		}
		data := req.Data
		if data != "" {
			var apiJson string
			apiJson, err = utils.AesDecryptByCFB(utils.ApiAesKey, data)
			if err != nil {
				utils.Fail(c, err, "参数解密失败")
				c.Abort()
				return
			}
			c.Request.Body = ioutil.NopCloser(bytes.NewReader([]byte(apiJson)))
			c.Request.ContentLength = int64(len(apiJson))
		}
		// 处理请求
		c.Next()
	}
}
