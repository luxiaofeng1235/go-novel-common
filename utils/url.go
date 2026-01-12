/*
 * @Descripttion: URL/路径处理工具
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 11:45:00
 */
package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-novel/config"
	"go-novel/global"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// 获取Scheme  http 或https
func GetScheme(c *gin.Context) string {
	if scheme := c.Request.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if c.Request.URL.Scheme != "" {
		return c.Request.URL.Scheme
	}
	if c.Request.TLS == nil {
		return "http"
	}
	return "https"
}

// 获取site http://www.aaa.com
func GetSite(c *gin.Context) string {
	return GetScheme(c) + "://" + c.Request.Host
}

// 获取后台图片资源的地址
func GetAdminUrl() string {
	if val := viper.GetString("admin.adminUrl"); strings.TrimSpace(val) != "" {
		return val
	}
	// 兼容旧字段
	return viper.GetString("server.adminUrl")
}

// 获取提供给对外的下载域名
func GetDownUrl() string {
	return viper.GetString("server.downUrl")
}

func GetApiUrl() string {
	return viper.GetString("api.apiUrl")
}

func GetApiEncrypt() bool {
	return viper.GetBool("api.encrypt")
}

// 获取静态资源对外访问的基础域名（优先使用 source.publicBaseUrl）
func GetSourcePublicBaseUrl() string {
	return viper.GetString("source.publicBaseUrl")
}

func GetFileUrl(path string) (fileUrl string) {
	if strings.Contains(path, "http") {
		fileUrl = path
	} else if path != "" {
		fileUrl = fmt.Sprintf("%v%v", GetApiUrl(), path)
		if !FileExist(path) {
			fileUrl = fmt.Sprintf("%v%v", GetApiUrl(), DefaultPic)
		} else if strings.Contains(path, ".apk") {
			// 如果是apk文件，直接返回拼接的url
			fileUrl = fmt.Sprintf("%v%v", GetApiUrl(), path)
		} else if !IsPic(path) {
			// 不是图片文件，返回默认图片
			fileUrl = fmt.Sprintf("%v%v", GetApiUrl(), DefaultPic)
		}
	}

	// 去掉路径中的某个文件夹
	newPath := strings.ReplaceAll(fileUrl, REPLACEFOLDER, "")
	return newPath
}

//func GetFileUrl(path string) (fileUrl string) {
//	if path != "" {
//		fileUrl = fmt.Sprintf("%v%v", GetApiUrl(), path)
//	} else {
//		fileUrl = path
//	}
//	return
//}

// 解析本地路径信息
func ParseLocalUrl(path string) (newfile string) {
	if path == "" {
		return
	}
	env := config.GetString("server.env")
	//优先判断路径是否为本地，如果非本地直接返回，本地才进行逻辑判断
	if env != Local {
		return path
	}
	// 脚手架最小实现：不依赖数据库配置表，直接返回传入路径
	return path
}

// 获取对应的apk的下载链接
func GetApkFileUrl(path string) (fileUrl string) {
	if path == "" {
		return
	}
	//path = strings.ReplaceAll(path, REPLACEAPK, "") //替换对应的路径信息进行显示
	fileUrl = fmt.Sprintf("%v%v", GetDownUrl(), path)
	return
}

// 获取路径信息
func GetAdminFileUrl(path string) (fileUrl string) {
	isHttp := strings.Contains(path, "http")
	spath := path              //先保留当前的路径信息方便做存储
	path = ParseLocalUrl(path) //解析是否为本地路径，如果是替换对应的路径信息，如果不是就返回默认的
	//线上的配置信息
	if isHttp == false {
		if path != "" {
			fileUrl = fmt.Sprintf("%v%v", GetAdminUrl(), spath)
			if !FileExist(path) { //文件不存在的时候
				fileUrl = fmt.Sprintf("%v%v", GetAdminUrl(), DefaultPic)
				log.Println(111, path)
				return
			} else { //判断文件存在的情况
				//判断如果是apk文件直接返回不需要判断
				if (strings.Contains(path, ".apk")) != false {
					fileUrl = fmt.Sprintf("%v%v", GetAdminUrl(), spath)
					return
				} else {
					isPic := IsPic(path)
					if !isPic {
						log.Println(222, path)
						fileUrl = fmt.Sprintf("%v%v", GetAdminUrl(), DefaultPic)
						return
					}
				}
			}
		}
	} else {
		//走线下的配置
		fileUrl = spath
	}

	//替换的路径信息
	newPath := strings.ReplaceAll(fileUrl, REPLACEFOLDER, "")
	return newPath
}

func GetLastNumber(str string) (number int) {
	segments := strings.Split(str, "/")
	lastSegment := segments[len(segments)-2]
	number, err := strconv.Atoi(lastSegment)
	if err != nil {
		// 处理转换错误的情况
		return
	}
	return
}

func GetUrlDomain(linkUrl string) (domain string) {
	if linkUrl == "" {
		return
	}
	parsedURL, err := url.Parse(linkUrl)
	if err != nil {
		global.Errlog.Errorf("GetUrlDomain err=%v", err.Error())
		return
	}
	domain = fmt.Sprintf("%v://%v", parsedURL.Scheme, parsedURL.Hostname())
	return
}

func GetUrlBookNum(bookUrl string) (num string) {
	// 使用正则表达式提取数字部分
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(bookUrl, -1)

	if len(matches) > 0 {
		num = matches[len(matches)-1]
	}
	return
}

func GetUrlSuffix(pathUrl string) (str string) {
	str = pathUrl[strings.LastIndex(pathUrl, "/")+1:]
	return
}

// 获取神马搜索的一些拼装字段信息
func GetReplaceChaojihuiCallbak(callback_url string, atype int, source string) (str string) {
	if callback_url == "" {
		return
	}
	var err error
	//解析对应的地址信息转换为decode编码信息
	callback_url, err = url.QueryUnescape(callback_url)
	if err != nil {
		fmt.Println("Error decoding URL:", err)
	}
	var urlStr string
	if atype == 1 { //激活
		urlStr = "imei_sum"
	} else { //留存
		urlStr = "idfa"
	}

	returnUrl := fmt.Sprintf("%v&type=%v&%v=&event_time%v=&source=", callback_url, atype, urlStr, GetUnix())
	return returnUrl
}

// 获取百度的回调地址并进行替换
func GetReplaceBaiduCallbak(callback_url, atype string, avalue int) (str string) {
	if callback_url == "" {
		return
	}
	var err error
	//解析对应的地址信息转换为decode编码信息
	callback_url, err = url.QueryUnescape(callback_url)
	if err != nil {
		fmt.Println("Error decoding URL:", err)
	}
	registerUrl := strings.ReplaceAll(callback_url, "{{ATYPE}}", atype)
	registerUrl = strings.ReplaceAll(registerUrl, "{{AVALUE}}", strconv.Itoa(avalue))
	return registerUrl
}

// 获取百度请求的接口信息
func GetBaiduResponse(url string) (str string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	str = string(body)
	// 打印响应体
	log.Printf("Get url= 【%s】 Response body: %s", url, str)
	return
}
