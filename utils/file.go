/*
 * @Descripttion: 文件处理工具（包含本地上传）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 12:15:00
 */
package utils

import (
	"crypto/rand"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mozillazg/go-pinyin"
	"go-novel/app/models"
	"go-novel/global"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"html"
	"image"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GetImageName 获取图片名称
func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = Md5(fileName)

	return fileName + ext
}

func GetFileName(filePath string) string {
	return filepath.Base(filePath)
}

func GetFileBase(fileName string) (fileBase string) {
	fileName = filepath.Base(fileName)
	fileBase = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	return
}

func GetFilePath(filePath string) (path string) {
	dirPath := filepath.Dir(filePath)
	path = filepath.ToSlash(dirPath)
	return
}

// GetSize 获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)
	return len(content), err
}

// GetExt 获取文件后缀
func GetExt(fileName string) string {
	return strings.TrimLeft(path.Ext(fileName), ".")
}

// SaveMultipartFile 保存上传文件到本地目录（dstDir 必须是本地目录路径，不是 URL）
func SaveMultipartFile(fileHeader *multipart.FileHeader, dstDir, dstFilename string) (savedFullPath string, err error) {
	if fileHeader == nil {
		return "", fmt.Errorf("上传文件不能为空")
	}
	if dstDir == "" {
		return "", fmt.Errorf("dstDir 不能为空")
	}
	if dstFilename == "" {
		return "", fmt.Errorf("dstFilename 不能为空")
	}

	if err = IsNotExistMkDir(dstDir); err != nil {
		return "", err
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	savedFullPath = filepath.Join(dstDir, dstFilename)
	dst, err := os.Create(savedFullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return savedFullPath, nil
}

// NewRandomFilename 基于原始文件名生成随机文件名（保留扩展名）
func NewRandomFilename(originalName string) string {
	ext := ""
	if originalName != "" {
		ext = filepath.Ext(originalName)
	}

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// 兜底：时间戳 + pid
		return strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	}
	return fmt.Sprintf("%x%s", b, ext)
}

// CheckNotExist 检查目录是否存在 返回true不存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// 检查文件是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// IsNotExistMkDir 如果不存在则新建文件夹
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

// MkDir 新建文件夹
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// 创建路径
func CreateFilePath(filePath string) error {
	// 路径不存在创建路径
	path, _ := filepath.Split(filePath) // 获取路径
	_, err := os.Stat(path)             // 检查路径状态，不存在创建
	if err != nil || os.IsExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
	}
	return err
}

// 获取上传文件字节流
func GetFileByte(f *multipart.FileHeader) (fileByte []byte, err error) {
	var fileHandle multipart.File
	fileHandle, err = f.Open() //打开上传文件
	if err != nil {
		return
	}
	defer fileHandle.Close()

	fileByte, err = ioutil.ReadAll(fileHandle) //获取上传文件字节流
	if err != nil {
		return
	}
	return
}

// 删除图片
func RemoveFile(path string) (err error) {
	if strings.Contains(path, "http") == false {
		err = os.Remove(path)
	}
	return err
}

func RemoveDir(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

// 删除图片
func RemoveFiles(paths []string) error {
	if len(paths) > 0 {
		for i, _ := range paths {
			_ = RemoveFile(paths[i])
		}
	}
	return nil
}

// 删除富文本图片
func RemoveEditorFile(content string) (err error) {
	oldArr := GetEditorImage(content)
	err = RemoveFiles(oldArr)
	return err
}

// 获取富文本图片地址
func GetEditorImage(content string) []string {
	imgUrl := `<img[\s\S]*?src\s*=\s*[\"|\'](.*?)[\"|\'][\s\S]*?>`
	//reSuperUrl := `<img.*?src\=\"(.*?)\"[^>]*>`
	//videoUrl :=<video.*?src\=\"(.*?)\"[^>]*>
	reg := regexp.MustCompile(imgUrl)
	list := reg.FindAllStringSubmatch(content, -1)
	//log.Println("总共: ", len(list))
	adminUrl := GetAdminUrl()
	var imgs []string
	if len(list) > 0 {
		for _, v := range list {
			if len(v) > 1 {
				imgPath := v[1]
				if strings.Contains(imgPath, "localhost") || strings.Contains(imgPath, "127.0.0.1") || strings.Contains(imgPath, adminUrl) {
					index := strings.Index(imgPath, "public")
					imgs = append(imgs, imgPath[index:])
				} else {
					imgs = append(imgs, imgPath)
				}
			}
		}
	}
	//log.Println(imgs)
	return imgs
}

// 删除富文本图片
func DelEditorImage(oldcontent, newcontent string) (err error) {
	oldArr := GetEditorImage(oldcontent)
	newArr := GetEditorImage(newcontent)
	//log.Println("oldArr:",oldArr)
	//log.Println("newArr:",newArr)
	tempArr := Intersect(oldArr, newArr)
	if len(tempArr) > 0 {
		delArr := Difference(oldArr, tempArr)
		//log.Println("delArr:",delArr)
		err = RemoveFiles(delArr)
	} else {
		//log.Println("oldArr:",oldArr)
		err = RemoveFiles(oldArr)
	}
	return err
}

func WriteFile(fileUrl string, data string) (err error) {
	err = IsNotExistMkDir(filepath.Dir(fileUrl))
	if err != nil {
		return
	}
	f, err := os.OpenFile(fileUrl, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	n, err := f.Write([]byte(data))
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func GetFile(filename string) (content string, err error) {
	if !FileExist(filepath.Dir(filename)) {
		err = fmt.Errorf("%s%s", filename, "文件不存在")
		return
	}
	var buffer []byte
	buffer, err = os.ReadFile(filename)
	if err != nil {
		return
	}
	content = string(buffer)
	return
}

// 下载图片
func DownImg(bookName, author, picUrl, pathDir string) (filePath string, err error) {
	// 检查URL是否合法
	if !strings.HasPrefix(picUrl, "http://") && !strings.HasPrefix(picUrl, "https://") {
		err = fmt.Errorf("图片链接错误 picUrl=%v", picUrl)
		return
	}

	//fileName := GetFileName(picUrl)
	//bookName = strings.Join(pinyin.LazyPinyin(bookName, pinyin.NewArgs()), "")
	bookName = GetFirstLetter(bookName)
	author = GetFirstLetter(author)
	// 获取文件扩展名
	fileExt := strings.ToLower(GetExt(picUrl))
	fileName := fmt.Sprintf("%v-%v.%v", bookName, author, fileExt)

	// 验证扩展名是否合法
	validExts := map[string]bool{"jpg": true, "jpeg": true, "png": true, "gif": true, "bmp": true, "webp": true}
	if !validExts[fileExt] {
		err = fmt.Errorf("图片后缀名不合法 fileExt=%v", fileExt)
		return
	}

	//4、创建图片保存目录  static/upload/20200623
	// 获取当前时间
	//now := time.Now()
	// 构建目标目录路径
	//dayPath := fmt.Sprintf("%d/%02d/%02d/", now.Year(), now.Month(), now.Day())
	//uploadPath := fmt.Sprintf("%s%s", pathDir, dayPath)
	uploadPath := fmt.Sprintf("%s", pathDir)
	err = IsNotExistMkDir(uploadPath)
	if err != nil {
		err = fmt.Errorf("%v", "创建目录失败")
		return
	}

	// 文件完整路径
	filePath = fmt.Sprintf("%s%s", uploadPath, fileName)
	if FileExist(filePath) {
		//global.Collectlog.Errorf("%v 文件已存在", filePath)
		return
	}

	httpClient := &http.Client{}

	//if IsS5 {
	//	httpTransport := getHttpTransport()
	//	httpClient = &http.Client{Transport: httpTransport}
	//}

	// 发送HTTP请求获取文件内容
	resp, err := httpClient.Get(picUrl)
	if err != nil {
		err = fmt.Errorf("%v err=%v获取图片链接失败", picUrl, err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.Status == "404 Not Found" {
		global.Lydlog.Errorf("%v 404 Not Found", picUrl)
		err = fmt.Errorf("%v", "404 Not Found")
		return
	}
	// 读取文件内容
	fileContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("%v 读取图片失败", picUrl)
		return
	}

	// 写入文件
	err = ioutil.WriteFile(filePath, fileContent, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("%v 保存图片失败", picUrl)
		return
	}
	return
}

func DownloadImage(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func ReplaceWords(text string, rules []*models.TextReplace) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	text = strings.Join(lines, "\n")
	for _, val := range rules {
		text = strings.ReplaceAll(text, val.Find, val.Replace)
	}
	text = strings.TrimSpace(text)
	return text
}

func ReplaceText(text string) (newText string) {
	newText = html.UnescapeString(text)
	newText = strings.ReplaceAll(newText, "&nbsp;", " ")
	newText = strings.ReplaceAll(newText, "NBSP", " ")
	newText = strings.ReplaceAll(newText, "<br>", "\n")
	newText = strings.ReplaceAll(newText, "<br/>", "\n")
	newText = strings.ReplaceAll(newText, "<br />", "\n")
	newText = regexp.MustCompile(`<br\s*/?>`).ReplaceAllString(newText, "\n")
	return
}

func ReadTextReplace(content string) (text string) {
	// 使用正则表达式替换多个换行符为一个换行符
	re := regexp.MustCompile(`\r\n+`)
	text = re.ReplaceAllString(content, "\n")
	re = regexp.MustCompile(`\n{2,}`)
	text = re.ReplaceAllString(text, "\n")
	re = regexp.MustCompile(`^\s+|\s+$`)
	text = re.ReplaceAllString(text, "")
	return
}

// 净化描述、去除html标签、去除头尾空格
func desHandle(s string) string {
	return regexp.MustCompile(`<[^>]+>|(^\s*)   |(\s*$)|&nbsp;`).ReplaceAllString(s, "")
}

func DecodeGBKtoUTF8(contentBody io.Reader) (string, error) {
	reader := transform.NewReader(contentBody, simplifiedchinese.GBK.NewDecoder())
	decodedContent, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decodedContent), nil
}

func GetFirstLetter(text string) (newText string) {
	if text == "" {
		return
	}
	pinyinSlice := pinyin.LazyPinyin(text, pinyin.NewArgs())
	firstLetterSlice := make([]string, len(pinyinSlice))
	for i, pinyinStr := range pinyinSlice {
		firstLetterSlice[i] = string([]rune(pinyinStr)[0])
	}
	newText = strings.Join(firstLetterSlice, "")
	return
}

func BucketUpload(isEnc bool, filePath, uploadPath string) (encFilePath, encFileUrl string, err error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("打开文件失败:%v", err)
		return
	}
	defer file.Close()
	if isEnc {
		randName := GetImageName(filePath)
		encName := GetEncodeEncImage(randName)

		encPath := fmt.Sprintf("attachment/%v", uploadPath)
		err = IsNotExistMkDir(encPath)
		if err != nil {
			err = fmt.Errorf("%v", "创建加密图片目录失败")
			return
		}
		encFilePath = fmt.Sprintf("%v/%v", encPath, encName)
		encFileUrl = fmt.Sprintf("%v/%v", BucketDomain, encFilePath)
		err = EncodeImage(filePath, encFilePath, ImgEncry)
		if err != nil {
			err = fmt.Errorf("加密图片失败：%v", err.Error())
			return
		}
		// 打开文件
		file, err = os.Open(encFilePath)
		if err != nil {
			err = fmt.Errorf("打开加密文件失败:%v", err)
			return
		}
		defer file.Close()
	}
	// 创建AWS会话
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("auto"), // 替换为你的AWS区域
		Endpoint:    aws.String(fmt.Sprintf("https://%v.r2.cloudflarestorage.com", ACCOUNT_ID)),
		Credentials: credentials.NewStaticCredentials(ACCESS_KEY_ID, ACCESS_KEY_SECRET, ""),
	})
	if err != nil {
		err = fmt.Errorf("创建AWS回话失败:%v", err.Error())
		return
	}

	svc := s3.New(sess)
	// 执行文件上传
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(encFilePath),
		Body:   file,
	})

	if err != nil {
		err = fmt.Errorf("上传文件失败:%v", err.Error())
		return
	}
	return
}

func DeleteBucketPic(filePath string) (err error) {
	// 创建AWS会话
	var sess *session.Session
	sess, err = session.NewSession(&aws.Config{
		Region:      aws.String("auto"), // 替换为你的AWS区域
		Endpoint:    aws.String(fmt.Sprintf("https://%v.r2.cloudflarestorage.com", ACCOUNT_ID)),
		Credentials: credentials.NewStaticCredentials(ACCESS_KEY_ID, ACCESS_KEY_SECRET, ""),
	})
	if err != nil {
		err = fmt.Errorf("创建AWS回话失败:%v", err.Error())
		return
	}

	svc := s3.New(sess)
	// 执行删除操作
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(filePath),
	})
	if err != nil {
		err = fmt.Errorf("删除仓库图片失败:%v", err.Error())
		return
	}
	return
}

func IsPic(picPath string) (isPic bool) {
	// 打开图片文件
	file, err := os.Open(picPath)
	if err != nil {
		fmt.Println("无法打开图片文件:", err)
		return
	}
	defer file.Close()

	// 解码图片
	_, _, err = image.Decode(file)
	if err != nil {
		return
	}

	// 获取图片的大小
	//width := img.Bounds().Dx()
	//height := img.Bounds().Dy()
	//fmt.Printf("图片宽度: %d\n", width)
	//fmt.Printf("图片高度: %d\n", height)
	isPic = true
	return
}
