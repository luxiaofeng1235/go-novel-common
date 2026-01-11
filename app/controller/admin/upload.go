package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-novel/app/service/common/upload_service"
	"go-novel/utils"
	"net/http"
	"os"
	"path"
	"strings"
)

type Upload struct{}

func (upload *Upload) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, err, "上传文件不存在")
		return
	}

	//2、获取后缀名 判断类型是否正确  .jpg .png .gif .jpeg
	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
		".apk":  true,
	}

	if _, ok := allowExtMap[extName]; !ok {
		utils.Fail(c, errors.New("图片后缀名不合法1111"), "上传出错了")
		return
	}

	//4、创建图片保存目录  static/upload/20200623
	day := utils.GetDay()
	dir := "public/resource/upload/" + day

	if err := os.MkdirAll(dir, 0666); err != nil {
		utils.Fail(c, errors.New(err.Error()), "上传目录不存在")
		return
	}

	//5、生成文件名称   144325235235.png
	fileName := fmt.Sprintf("%d%s", utils.GetUnixNano(), file.Filename)
	path := path.Join(dir, fileName)

	if utils.OssUpload {
		var url string
		url, err = utils.UploadOss(file, path)
		if err != nil {
			utils.Fail(c, err, "上传失败")
			return
		}
		ret := make(map[string]string)
		ret["url"] = url
		utils.Success(c, ret, "上传成功")
	} else {
		if err := c.SaveUploadedFile(file, path); err != nil {
			utils.Fail(c, err, "上传失败")
			return
		}
		ret := make(map[string]string)
		ret["url"] = utils.GetSite(c) + "/" + path
		utils.Success(c, ret, "上传成功")
	}
}

// 上传富文本图片
func (upload *Upload) CkEditorUp(c *gin.Context) {
	ftype := strings.ToLower(c.Query("type"))

	var fileName string
	var url string
	var err error
	if ftype == "images" {
		fileName, url, err = upload_service.UploadEditorImg(c)
	} else if ftype == "files" {
		fileName, url, err = upload_service.UploadEditorFile(c)
	}

	if err != nil {
		utils.Fail(c, err, "上传文件失败")
		return
	}
	if err != nil {
		response := gin.H{
			"error": gin.H{"message": "上传失败，" + err.Error(), "number": 105},
		}
		c.JSON(http.StatusOK, response)
	} else {
		if strings.Contains(url, "http") == false {
			url = utils.GetSite(c) + "/" + url
		}
		response := gin.H{
			"fileName": fileName,
			"uploaded": 1,
			"url":      url,
		}
		c.JSON(http.StatusOK, response)
	}
}
