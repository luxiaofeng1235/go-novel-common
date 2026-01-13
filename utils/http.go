package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-novel/global"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alexeyco/goozzle"
	"github.com/gocolly/colly/v2"
	"golang.org/x/net/html/charset"
	"golang.org/x/net/proxy"
)

// Http Get请求基础函数, 通过封装Go语言Http请求, 支持火币网REST API的HTTP Get请求
// strUrl: 请求的URL
// strParams: string类型的请求参数, user=lxz&pwd=lxz
// return: 请求结果
func HttpGetRequest(strUrl string, mapParams map[string]string) string {
	httpClient := &http.Client{}

	var strRequestUrl string
	if nil == mapParams {
		strRequestUrl = strUrl
	} else {
		strParams := Map2UrlQuery(mapParams)
		strRequestUrl = strUrl + "?" + strParams
	}
	strRequestUrl = strings.TrimSpace(strRequestUrl)
	// 构建Request, 并且按官方要求添加Http Header
	request, err := http.NewRequest("GET", strRequestUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")

	// 发出请求
	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()
	// 解析响应内容
	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}
	return string(body)
}

// GetContextResponse 兼容旧代码的“GET 请求并忽略结果”工具函数（用于上报等场景）。
// 注意：脚手架默认不依赖该返回值，失败时返回空字符串。
func GetContextResponse(strUrl string) string {
	strUrl = strings.TrimSpace(strUrl)
	if strUrl == "" {
		return ""
	}
	httpClient := &http.Client{Timeout: 3 * time.Second}
	request, err := http.NewRequest("GET", strUrl, nil)
	if err != nil {
		return ""
	}
	request.Header.Add("User-Agent", "Mozilla/5.0")
	response, err := httpClient.Do(request)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

// Http POST请求基础函数, 通过封装Go语言Http请求, 支持火币网REST API的HTTP POST请求
// strUrl: 请求的URL
// mapParams: map类型的请求参数
// return: 请求结果
func HttpPostRequest(strUrl string, mapParams map[string]interface{}, headers map[string]string) string {
	httpClient := &http.Client{}

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Language", "zh-cn")
	if len(headers) > 0 {
		for key, val := range headers {
			request.Header.Add(key, val)
		}
	}

	response, err := httpClient.Do(request)
	defer response.Body.Close()
	if nil != err {
		return err.Error()
	}

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}

// 发起post请求
func GetPostData(strUrl string, mapParams map[string]interface{}) string {

	// 编码 JSON 数据
	jsonBytes, err := json.Marshal(mapParams)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	// 创建 HTTP POST 请求
	url := strUrl
	fmt.Println(mapParams)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	fmt.Println(url)
	// 发送请求并处理响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()
	// 读取响应数据
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	responseBody := string(respBody)
	log.Printf("url = %v Status: %v Response: %v", url, resp.Status, responseBody)
	return responseBody
}

// 将map格式的请求参数转换为字符串格式的
// mapParams: map格式的参数键值对
// return: 查询字符串
func Map2UrlQuery(mapParams map[string]string) string {
	var strParams string
	for key, value := range mapParams {
		strParams += (key + "=" + value + "&")
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}

func PayPostRequest(reqUrl string, mapParams map[string]interface{}) (result string, err error) {
	httpClient := &http.Client{}

	jsonParams := ""
	if mapParams != nil {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}
	var request *http.Request
	request, err = http.NewRequest("POST", reqUrl, strings.NewReader(jsonParams))
	if err != nil {
		return
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Language", "zh-cn")

	response, err := httpClient.Do(request)
	defer response.Body.Close()
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	result = string(body)
	log.Println(reqUrl, mapParams, result)
	return
}

func GetCharsetFromContentType(contentType string) string {
	contentType = strings.ToLower(contentType)
	charset := ""

	// 查找字符集信息
	if strings.Contains(contentType, "charset") {
		charsetIndex := strings.Index(contentType, "charset")
		if charsetIndex != -1 {
			charsetStart := charsetIndex + len("charset")
			charset = strings.TrimSpace(contentType[charsetStart:])
			if strings.HasPrefix(charset, "=") {
				charset = strings.TrimSpace(charset[1:])
			}
		}
	}

	return charset
}

func DoGet(strUrl string) (html string, err error) {
	httpClient := &http.Client{}
	if IsS5 {
		httpTransport := GetHttpTransport()
		httpClient = &http.Client{Transport: httpTransport}
	}

	// 构建Request, 并且按官方要求添加Http Header
	request, err := http.NewRequest("GET", strUrl, nil)
	if err != nil {
		err = fmt.Errorf("request err=%v", err.Error())
		return
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Set("Accept-Charset", "utf-8")
	request.Host = "123.14.254.86"
	// 发出请求
	response, err := httpClient.Do(request)
	if err != nil {
		err = fmt.Errorf("response err=%v", err.Error())
		return
	}
	defer response.Body.Close()
	// 解析响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("body err=%v", err.Error())
		return
	}
	log.Println(string(body))
	var htmlByte []byte
	htmlByte, err = ConvertToUTF8(body)
	html = string(htmlByte)
	if err != nil {
		global.Collectlog.Errorf("采集html为空%v", err.Error())
	}
	if html == "" {
		global.Collectlog.Errorf("采集html为空%v")
		return
	}
	return
}

func GetHttpTransport() (httpTransport *http.Transport) {
	if S5Domain == "" || S5Port == "" {
		if S5Type == Rank {
			GetS5()
		}
	}

	proxyAddress := fmt.Sprintf("%v:%v", S5Domain, S5Port)

	var proxyAuth = proxy.Auth{}
	if S5Username != "" && S5Passwd != "" {
		proxyAuth = proxy.Auth{
			User:     S5Username,
			Password: S5Passwd,
		}
	}

	// 创建 SOCKS5 代理客户端
	var err error
	var dialer proxy.Dialer
	dialer, err = proxy.SOCKS5("tcp", proxyAddress, &proxyAuth, proxy.Direct)
	if err != nil {
		global.Collectlog.Errorf("无法连接到代理服务器:%v", err.Error())
		return
	}
	// 创建 HTTP 客户端
	httpTransport = &http.Transport{Dial: dialer.Dial}
	httpTransport.Dial = dialer.Dial
	return
}

func ConvertToUTF8(gbkBytes []byte) ([]byte, error) {
	// 将 GBK 编码的字节转换为字符串
	gbkStr := string(gbkBytes)

	// 检测并转换字符集
	utf8Reader, err := charset.NewReader(strings.NewReader(gbkStr), "")
	if err != nil {
		return nil, err
	}

	// 读取转换后的 HTML 内容
	utf8Bytes, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		return nil, err
	}
	return utf8Bytes, nil
}

func GetHtml(url, encode string, urlComplete int, sleepSecond int64) (html string, err error) {
	html, err = DoGet(url)
	if html == "" {
		return
	}
	if sleepSecond > 0 {
		time.Sleep(time.Second * time.Duration(sleepSecond))
	}
	if urlComplete > 0 {
		html = UrlComplete(html, url)
	}
	return
}

func GetHtmlByGoozzle(linkUrl string) (html string, err error) {
	u, _ := url.Parse(linkUrl)
	res, err := goozzle.Get(u).Do()
	if err != nil {
		return
	}
	decodedContent, err := ConvertToUTF8(res.Body())
	if err != nil {
		return
	}
	html = string(decodedContent)
	return
}

// 请求指定 URL 并返回 HTML 内容
func GetHtmlcolly(linkUrl string) (html string, err error) {
	// 创建一个新的 Collector
	c := colly.NewCollector()
	// 在访问页面之后执行的回调函数
	//c.OnResponse(func(res *colly.Response) {
	//	//html = string(res.Body)
	//	log.Println("body:", string(res.Body))
	//})
	c.SetRequestTimeout(time.Second * 50)
	c.OnHTML("html", func(e *colly.HTMLElement) {
		html, _ = e.DOM.Html()
	})

	// 在访问页面出错时执行的回调函数
	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == http.StatusNotFound {
			return
		}
		if err != nil {
			global.Errlog.Errorf("访问页面出错 %v", err.Error())
			//time.Sleep(time.Second * 1)
			err = TrydoRequest(c, linkUrl)
			if err != nil {
				global.Errlog.Errorf("自动重试后 第1次 %v", err.Error())
				//time.Sleep(time.Second * 2)
				err = TrydoRequest(c, linkUrl)
				if err != nil {
					global.Errlog.Errorf("自动重试后 第2次 %v", err.Error())
					//time.Sleep(time.Second * 3)
					err = TrydoRequest(c, linkUrl)
					if err != nil {
						global.Errlog.Errorf("自动重试后 第3次 %v", err.Error())
					}
				}
			}
			return
		}
	})

	// 访问指定的页面
	err = c.Visit(linkUrl)
	if err != nil {
		return
	}
	return
}

func TrydoRequest(c *colly.Collector, linkUrl string) (err error) {
	clonedCollector := c.Clone()
	clonedCollector.AllowURLRevisit = true
	err = clonedCollector.Visit(linkUrl)
	return
}
