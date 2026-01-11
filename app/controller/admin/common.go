package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-novel/app/service/common/upload_service"
	"go-novel/utils"
	"strconv"
	"strings"
)

type Common struct{}

func (common *Common) UploadClassPic(c *gin.Context) {
	url, err := upload_service.UploadFile(c, "class", "")
	if err != nil {
		utils.Fail(c, err, "上传图片失败")
		return
	}
	//替换对应的的长路经
	newPath := strings.ReplaceAll(url, utils.REPLACEFOLDER, "")
	res := gin.H{
		"domain":   utils.GetAdminUrl(),
		"url":      url,
		"show_url": newPath,
	}
	utils.Success(c, res, "ok")
	return
}

// 上传apk文件信息
func (common *Common) UploadAPkFile(c *gin.Context) {
	url, err := upload_service.UploadFile(c, "apk", "")
	if err != nil {
		utils.Fail(c, err, "上传图片失败")
		return
	}
	//替换对应的的长路经
	newPath := strings.ReplaceAll(url, utils.REPLACEAPK, "")
	res := gin.H{
		"domain":   utils.GetDownUrl(),
		"url":      url,
		"show_url": newPath,
	}
	utils.Success(c, res, "ok")
	return
}

func (common *Common) UploadBookPic(c *gin.Context) {
	url, err := upload_service.UploadFile(c, "book", "")
	if err != nil {
		utils.Fail(c, err, "上传图片失败")
		return
	}
	newPath := strings.ReplaceAll(url, utils.REPLACEFOLDER, "")
	fmt.Println(url)
	res := gin.H{
		"domain":   utils.GetAdminUrl(),
		"url":      url,
		"show_url": newPath,
	}
	utils.Success(c, res, "ok")
	return
}

func (common *Common) UploadAdverPic(c *gin.Context) {
	width := c.PostForm("width")
	height := c.PostForm("height")
	widthInt, _ := strconv.Atoi(width)
	heightInt, _ := strconv.Atoi(height)
	if widthInt > 0 && heightInt > 0 {
		err := upload_service.LimitFileWH(c, widthInt, heightInt)
		if err != nil {
			utils.Fail(c, err, "")
			return
		}
	}

	url, err := upload_service.UploadFile(c, "adver", "")
	if err != nil {
		utils.Fail(c, err, "上传图片失败")
		return
	}
	newPath := strings.ReplaceAll(url, utils.REPLACEFOLDER, "")
	res := gin.H{
		"domain":   utils.GetAdminUrl(),
		"url":      url,
		"show_url": newPath,
	}
	utils.Success(c, res, "ok")
	return
}
