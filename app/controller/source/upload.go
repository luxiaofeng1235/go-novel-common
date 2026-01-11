package source

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/service/common/upload_service"
	"go-novel/utils"
)

type Upload struct{}

func (upload *Upload) UploadBookPic(c *gin.Context) {
	url, err := upload_service.UploadFile(c, "book", "")
	if err != nil {
		utils.Fail(c, err, "上传图片失败")
		return
	}
	res := gin.H{
		"domain": utils.GetApiUrl(),
		"url":    url,
	}
	utils.Success(c, res, "ok")
	return
}
