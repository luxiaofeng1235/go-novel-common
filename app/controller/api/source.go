package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/source_service"
	"go-novel/utils"
)

type Source struct{}

func (source *Source) SourceList(c *gin.Context) {
	var req models.BookSourceReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	sourceList, err := source_service.GetBookSourceList(req.Bid)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, sourceList, "ok")
}
