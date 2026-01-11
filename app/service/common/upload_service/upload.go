package upload_service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
	"go-novel/app/service/admin/setting_service"
	"go-novel/config"
	"go-novel/utils"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func LimitFileWH(c *gin.Context, width, height int) (err error) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		err = errors.New("上传文件不存在")
		return
	}
	// 校验尺寸大小
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return
	}
	if img.Width != width || img.Height != height {
		err = fmt.Errorf("上传文件的图片大小不合符标准,宽需要为%dpx 高需要为%dpx 。当前上传图片的宽高分别为：%dpx和%dpx", width, height, img.Width, img.Height)
		return
	}
	return
}

// 上传富文本图片
func UploadEditorImg(c *gin.Context) (fileName string, url string, err error) {
	file, err := c.FormFile("upload")
	if err != nil {
		err = errors.New("上传文件不存在")
		return
	}
	fileName = file.Filename
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	url, err = uploadFile(c, file, "editor", "", allowExtMap)
	return
}

// 上传富文本文件
func UploadEditorFile(c *gin.Context) (fileName string, url string, err error) {
	file, err := c.FormFile("upload")
	if err != nil {
		err = errors.New("上传文件不存在")
		return
	}
	fileName = file.Filename
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
		".apk":  true,
		".txt":  true,
		".xls":  true,
		".xlxs": true,
		".pdf":  true,
	}
	url, err = uploadFile(c, file, "editor", "", allowExtMap)
	return
}

// 上传文件
func UploadFile(c *gin.Context, uploadPath, fileName string) (url string, err error) {
	file, err := c.FormFile("file")
	if err != nil {
		err = errors.New("上传文件不存在")
		return
	}
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
		".mp4":  true,
		".apk":  true, //兼容安卓的apk包
	}
	if uploadPath == "book" {
		fileName = strings.Join(pinyin.LazyPinyin(file.Filename, pinyin.NewArgs()), "")
	}
	url, err = uploadFile(c, file, uploadPath, fileName, allowExtMap)
	return
}

func uploadFile(c *gin.Context, f *multipart.FileHeader, uploadPath, fileName string, allowExtMap map[string]bool) (url string, err error) {
	//2、获取后缀名 判断类型是否正确  .jpg .png .gif .jpeg
	extName := path.Ext(f.Filename)
	if _, ok := allowExtMap[extName]; !ok {
		err = errors.New("图片后缀名不合法")
		return
	}

	currentRouter := c.FullPath() //获取当前的路由配置信息
	log.Printf("当前请求的路由配置信息:【%s】", currentRouter)
	//4、创建图片保存目录  static/upload/20200623
	//day := utils.GetDay()
	uploadPath = fmt.Sprintf("%s%s", uploadPath, "/")

	var uploadCommonPath string
	uploadCommonPath, _ = setting_service.GetValueByName("uploadCommonPath")
	//dir := "public/resource/upload/" + uploadPath + day
	var dir string
	if currentRouter == "/system/common/uploadAPkFile" { //只有这个路由用默认的上传
		dir = utils.UPLOADAPK //获取上传的路径信息
	} else {
		//走系统的配置路由信息
		if uploadCommonPath != "" {
			dir = fmt.Sprintf("%v%v", uploadCommonPath, uploadPath)
		} else {
			dir = fmt.Sprintf("/data/upload/%v", uploadCommonPath)
		}
	}
	log.Printf("上传对应的目录文件为：%s\n", dir)
	if err = os.MkdirAll(dir, 0666); err != nil {
		err = errors.New("创建目录失败")
		return
	}
	//5、生成文件名称   144325235235.png
	//fileName := fmt.Sprintf("%d%s", utils.GetUnixNano(), file.Filename)
	if fileName == "" {
		fileName = fmt.Sprintf("%s%s", utils.RandomString("rand", 18), extName)
	} else {
		fileName = fmt.Sprintf("%s%s", fileName, ".png")
	}
	path := path.Join(dir, fileName)

	//上传到阿里云
	if utils.OssUpload {
		url, err = utils.UploadOss(f, path)
		return
	}

	var uploadType string = "0"
	setting, err := setting_service.GetSetByName("uploadType")
	if err == nil && setting != nil {
		uploadType = setting.Value
	}
	//上传到七牛云
	if utils.QiNiuUpload || uploadType == "2" {
		url, err = utils.UploadQiNiu(f, path)
		return
	}

	if extName == ".apk" {
		//判断是否为apk文件单独进行上传请求
		fmt.Println("big 111111111111111111111111111111111")
		if err = uploadBigHandler(f, path); err != nil {
			err = errors.New(fmt.Sprintf("%s%s", "上传失败", err.Error()))
			return
		}
	} else {
		if err = c.SaveUploadedFile(f, path); err != nil {
			err = errors.New(fmt.Sprintf("%s%s", "上传失败", err.Error()))
			return
		}
	}

	//url = utils.GetSite(c) + "/" + path
	env := config.GetString("server.env")
	if env == utils.Local {
		//如果是本地路径就强制替换
		//替换所有的本地路径信息方便做测试
		url = strings.ReplaceAll(path, "E:\\", "/") //替换具体的某个盘
		url = strings.ReplaceAll(url, "\\", "/")
	} else {
		url = path
	}
	return
}

// 上传超文件信息
func uploadBigHandler(file *multipart.FileHeader, dst string) error {
	// 限制上传文件大小为 10MB
	const maxUploadSize = 200 * 1024 * 1024 // 200 MB
	src, err := file.Open()
	if err != nil {
		return err
	}
	fmt.Printf("file_path = %v\n", dst)
	defer src.Close()

	// 使用 LimitReader 限制读取大小
	limitedReader := io.LimitReader(src, maxUploadSize)

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}
	// 创建目标文件
	out, err := os.Create(dst)
	if err != nil {
		fmt.Println("Unable to create the file")
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, limitedReader)
	if err != nil {
		fmt.Println("Failed to save the file")
		return err
	}
	log.Printf("file_path =%v upload success!!!", dst)
	return nil
}
