/*
 * @Descripttion: API 通用控制器（无需登录：本地上传等）
 * @Author: red
 * @Date: 2026-01-12 12:20:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 12:20:00
 */
package api

import (
	"go-novel/app/service/common/file_service"
	"go-novel/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Common struct{}

// Upload 本地上传（无需登录、无需加密参数）
// multipart/form-data:
// - file: 文件/图片/视频
// - dir: (可选) 相对目录，如 avatar、video/2026
func (common *Common) Upload(c *gin.Context) {
	maxSizeMB := viper.GetInt("upload.maxSizeMB")
	if maxSizeMB <= 0 {
		maxSizeMB = 50
	}
	maxBytes := int64(maxSizeMB) * 1024 * 1024
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.FailEncrypt(c, err, "读取上传文件失败")
		return
	}
	if maxBytes > 0 && fileHeader.Size > maxBytes {
		utils.FailEncrypt(c, nil, "文件过大")
		return
	}
	subDir := c.PostForm("dir")
	//上传图片
	res, err := file_service.LocalUpload(fileHeader, subDir)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, res, "ok")
}
