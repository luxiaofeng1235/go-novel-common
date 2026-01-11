package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/utils"
)

type Pay struct{}

func (pay *Pay) AliPayNotify(c *gin.Context) {
	return
}

func (pay *Pay) AliPayCallback(c *gin.Context) {
	utils.Success(c, "", "alipay callback")
	return
}
