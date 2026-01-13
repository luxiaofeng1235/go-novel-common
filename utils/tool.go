/*
 * @Descripttion: 工具方法集合
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 10:25:00
 */
package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/mail"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/ant-libs-go/ip_parser"
	"github.com/axgle/mahonia"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ipipdotnet/ipdb-go"
	"github.com/longbridgeapp/opencc"
	"github.com/mssola/user_agent"
	"github.com/olahol/melody"
	"github.com/oschwald/geoip2-golang"
	uuid "github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
)

func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result[i] = bytes[r.Intn(len(bytes))]
	}
	return string(result)
}

// 获取当前时间戳
func GetNowUnix() int64 {
	return time.Now().Unix()
}

func StrToUnix(str string) int64 {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// 字符串转int
func StrToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// 字符串转int64
func StrToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// 手机号中间4位替换为*号
func FormatMobileStar(mobile string) string {
	if len(mobile) <= 10 {
		return mobile
	}
	return mobile[:3] + "****" + mobile[7:]
}

// 获取请求的IP地址  返回远程客户端的 IP 可以自动获取IPV4的地址
func GetRequestIP(c *gin.Context) string {
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

// RemoteIp 返回远程客户端的 IP，如 192.168.1.1
func RemoteIp(req *gin.Context) string {
	remoteAddr := req.Request.RemoteAddr
	if ip := req.Request.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Request.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	log.Printf("get remote_addr :%v", remoteAddr)
	return remoteAddr
}

func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		// 检查ip地址判断是否回环地址
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}

		if ipAddr.IP.To4() != nil {
			return ipAddr.IP.To4().String(), nil
		}
		return ipAddr.IP.String(), nil
	}
	return
}

func WsRemoteIp(req *melody.Session) string {
	remoteAddr := req.Request.RemoteAddr
	if ip := req.Request.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Request.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(conn.RemoteAddr())
	return localAddr.IP.String()
}

// 首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// CheckMobile 检验手机号
func CheckMobile(phone string) bool {
	// 匹配规则
	// ^1第一位为一
	// [345789]{1} 后接一位345789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1([356789][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\\d{8}$"
	// 正则调用规则
	reg := regexp.MustCompile(regRuler)
	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)
}

func Checkemail(email string) bool {
	regRuler := "^([A-Za-z0-9_.])+@([A-Za-z0-9_.])+.([A-Za-z]{2,4})$"
	// 正则调用规则
	reg := regexp.MustCompile(regRuler)
	// 返回 MatchString 是否匹配
	return reg.MatchString(email)
}

// 验证邮箱
func CheckEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// 验证eth地址
func CheckEthAddress(address string) bool {
	address = strings.ToLower(address)
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

// JSON序列化方式
func StructToMap(stuObj interface{}) (map[string]interface{}, error) {
	// 结构体转json
	strRet, err := json.Marshal(stuObj)
	if err != nil {
		return nil, err
	}
	// json转map
	var mRet map[string]interface{}
	err1 := json.Unmarshal(strRet, &mRet)
	if err1 != nil {
		return nil, err1
	}
	return mRet, nil
}

// 截取字符串，支持多字节字符
// start：起始下标，负数从从尾部开始，最后一个为-1
// length：截取长度，负数表示截取到末尾
func SubStr(str string, start int, length int) (result string) {
	s := []rune(str)
	total := len(s)
	if total == 0 {
		return
	}
	// 允许从尾部开始计算
	if start < 0 {
		start = total + start
		if start < 0 {
			return
		}
	}
	if start > total {
		return
	}
	// 到末尾
	if length < 0 {
		length = total
	}

	end := start + length
	if end > total {
		result = string(s[start:])
	} else {
		result = string(s[start:end])
	}

	return
}

// 获取ip所属城市
func GetCityByIp(ip string) string {
	if ip == "" {
		return ""
	}
	if ip == "[::1]" || ip == "127.0.0.1" {
		return "内网IP"
	}
	url := "http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip
	res, err := http.Get(url)
	if err != nil {
		return ""
	}
	bytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return ""
	}
	src := string(bytes)
	tmp := ConvertToString(src, "gbk", "utf-8")
	//log.Println(tmp)

	if gjson.Get(tmp, "code").Int() == 0 {
		city := gjson.Get(tmp, "city").String()
		return city
	} else {
		return ""
	}
}

// gbk转utf8
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// 获取浏览器userAgent
func GetBrowser(userAgent string) (string, string) {
	ua := user_agent.New(userAgent)
	Browser, _ := ua.Browser()
	Os := ua.OS()
	return Browser, Os
}

// 求并集
func Union(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// 求交集
func Intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// 求差集
func Difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

// 禁用json中的转义字符
func EscapeHtml(data interface{}) string {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(data); err != nil {
		return ""
	}
	return bf.String()
}

// 第一个参数name 为 cookie 名；
// 第二个参数value 为 cookie 值；
// 第三个参数maxAge 为 cookie 有效时长，当 cookie 存在的时间超过设定时间时，cookie 就会失效，它就不再是我们有效的 cookie，他的时间单位是秒second；
// 第四个参数path 为 cookie 所在的目录；
// 第五个domain 为所在域，表示我们的 cookie 作用范围，里面可以是localhost也可以是你的域名，看自己情况；
// 第六个secure 表示是否只能通过 https 访问，为true只能是https；
// 第七个httpOnly 表示 cookie 是否可以通过 js代码进行操作，为true时不能被js获取
func SetCookie(c *gin.Context, key string, val string, expireSecond int) {
	c.SetCookie(key, val, expireSecond, "/", "127.0.0.1",
		false, false)
}

func GetCookie(c *gin.Context, key string) (string, error) {
	cookie, err := c.Cookie(key)
	if err != nil {
		errors.New("Cookie does not exist")
		return "", err
	}
	return cookie, nil
}

// 设置session  utils.SetSession(c,"name","张三")
func SetSession(c *gin.Context, key interface{}, value interface{}) error {
	session := sessions.Default(c)
	if session == nil {
		return nil
	}
	session.Set(key, value)
	return session.Save()
}

// 获取session utils.GetSession(c,"name")
func GetSession(c *gin.Context, key interface{}) interface{} {
	session := sessions.Default(c)
	return session.Get(key)
}

// 删除某个session
func DeleteSession(c *gin.Context, key string) {
	session := sessions.Default(c)
	session.Delete(key)
}

// 清空session
func ClearSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
}

// 获取uuid
func GetUUID() uuid.UUID {
	return uuid.NewV4()
}

func GetRandomUsername() string {
	// 生成8位uuid
	uuidStr := GetUUID().String()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	shortUsername := uuidStr[:9]
	// 随机生成数字，避免以0开头
	digit := rand.Intn(9) + 1
	// 拼接生成的ID
	uniqueID := fmt.Sprintf("%v%v", digit, shortUsername)
	return uniqueID
}

func FormatDecimal(amount string) string {
	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return amount // 转换失败，返回原始字符串
	}

	// 判断是否有小数部分
	hasDecimal := value != float64(int(value))

	if hasDecimal {
		// 保留3位小数
		formatted := fmt.Sprintf("%.3f", value)
		return formatted
	}

	// 没有小数部分，返回整数形式
	return strconv.Itoa(int(value))
}

func FormatInt64(str string) (num int64) {
	var err error
	num, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		return
	}
	return
}

func GetWanFormatted(hits int64) string {
	hitsFloat := float64(hits)
	absHits := math.Abs(hitsFloat)

	switch {
	case absHits >= 1e8:
		tempStr := fmt.Sprintf("%.2f", hitsFloat/1e8)
		return fmt.Sprintf("%v 亿", trimZero(tempStr))
	case absHits >= 1e4:
		tempStr := fmt.Sprintf("%.2f", hitsFloat/1e4)
		return fmt.Sprintf("%v 万", trimZero(tempStr))
	default:
		return fmt.Sprintf("%d", hits)
	}
}

func trimZero(tempStr string) string {
	tempStr = strings.TrimRight(tempStr, "0")
	tempStr = strings.TrimRight(tempStr, ".")
	return tempStr
}

func GetWords(html string) (words int) {
	words = utf8.RuneCountInString(html)
	return
}

func GetAdminPic(imgs string) (imgArr []string) {
	pics := strings.Split(imgs, ",")
	for _, pic := range pics {
		pic = GetAdminFileUrl(pic)
		imgArr = append(imgArr, pic)
	}
	return
}

func GetApiPic(imgs string) (imgArr []string) {
	pics := strings.Split(imgs, ",")
	for _, pic := range pics {
		pic = GetFileUrl(pic)
		imgArr = append(imgArr, pic)
	}
	return
}

func PrintFormat(data interface{}) string {
	var str bytes.Buffer
	_ = json.Indent(&str, []byte(JSONString(data)), "", "    ")
	return str.String()
}

func TrimDotTag(tagName string) (tag string) {
	tagName = strings.TrimSpace(tagName)
	if strings.Contains(tagName, "新书") {
		tagName = strings.Replace(tagName, "新书", "", -1)
		tagName = strings.Replace(tagName, RoundDot, "", -1)
		tagName = strings.TrimSpace(tagName)
	}
	tag = tagName
	return
}

// 根据V4的接口解析当前的接口状态
func GetParseIpByV4(client_ip string) (city_name string) {
	if client_ip == "" {
		return
	}
	key := "dgXAYv0e14cgJ42ZQZPSAwUCGKQdUT1cVZDjIQSHmbBFBcpsoEwIR2JkYrmI87eu"                                      //秘钥
	apiurl := fmt.Sprintf("https://api.ipplus360.com/ip/geo/v1/city/?key=%v&ip=%v&coordsys=WGS84", key, client_ip) //请求接口地址
	response, err := http.Get(apiurl)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer response.Body.Close()
	// 读取响应数据
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}
	// 将响应数据转换为字符串并打印输出
	data := string(body)
	log.Printf("V4接口返回内容：%v", data)
	code := gjson.Get(data, "code").String()
	if code == "Success" {
		//获取所在地区编码
		location := gjson.Get(data, "data.city").String()
		return location
	} else {
		return ""
	}
}

// 获取V6接口
func GetParseIpByV6(client_ip string) (city_name string) {
	if client_ip == "" {
		return
	}
	key := "mlkZnecLQyTR2bDtrUaW4IP4hYfl5IP2Aw2KIn3oNTFs4veFj9nRLTL7XmZdpWTz"
	apiurl := fmt.Sprintf("https://api.ipplus360.com/ip/geo/v1/ipv6/?key=%v&ip=%v&coordsys=WGS84", key, client_ip)
	response, err := http.Get(apiurl)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer response.Body.Close()
	// 读取响应数据
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}
	// 将响应数据转换为字符串并打印输出
	data := string(body)
	log.Printf("V6接口返回内容：%v", data)
	code := gjson.Get(data, "code").String()
	if code == "Success" {
		//获取所在地区编码
		location := gjson.Get(data, "data.city").String()
		return location
	} else {
		return ""
	}
}

// 通过接口来获取所在地
func GetCityByNetwork(client_ip string) (city_name string) {

	//url := "https://searchplugin.csdn.net/api/v1/ip/get?ip=" + client_ip
	url := "http://ip-api.com/json/" + client_ip + "?lang=zh-CN"
	log.Printf("请求的接口地址 url = %v\n", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer response.Body.Close()
	// 读取响应数据
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}
	// 将响应数据转换为字符串并打印输出
	data := string(body)
	log.Printf("返回内容：%v", data)

	//获取状态码
	code := gjson.Get(data, "status").String()
	if code == "success" {
		city_name = gjson.Get(data, "city").String()
		prefix := "市"
		//处理未返回的市的字段信息
		if strings.Contains(city_name, string(prefix)) == false {
			city_name = fmt.Sprintf("%v%v", city_name, prefix)
		}
		return city_name
	} else {
		return "北京市"
	}
	//if code == 200 {
	//	location := gjson.Get(data, "data.address").String()
	//	if location != "" {
	//		words := strings.Split(location, " ")
	//		for key, val := range words {
	//			if key == 2 && val != "" {
	//				city_name = fmt.Sprintf("%s%s", val, "市")
	//				break
	//			}
	//		}
	//		//处理城市问题
	//		if city_name == "" {
	//			city_name = "北京市"
	//			log.Printf("ip = 【%s】未获取到城市信息，给一个默认的城市为：【%s】", client_ip, city_name)
	//		}
	//	}
	//} else {
	//	city_name = "北京市"
	//	log.Printf("未获取到城市信息，一个默认的城市为11：【%s】", city_name)
	//}
	//return city_name
}

// 根据IP获取对应的城市信息
func GetGeoCiyByIp(client_ip string) (city_name string) {
	geodb, err := geoip2.Open("public/resource/GeoLite2-City_20221007/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer geodb.Close()

	ip := net.ParseIP(client_ip)
	record, err := geodb.City(ip) //解析对应的城市信息
	if err != nil {
		log.Printf("GeoLite2 City Error:%s", err.Error())
		return
	}
	fmt.Println(record)
	city_name = record.City.Names["zh-CN"]
	if city_name != "" {
		//特殊判断下防止出问题
		prefix := "市"
		if strings.Contains(city_name, string(prefix)) == false {
			city_name = fmt.Sprintf("%v%v", city_name, prefix)
		}
		//en是英文缩写
	} else {
		log.Printf("本地解析IP未获取导数据，给一个默认的 *********** 北京市")
		city_name = "北京市" //给一个默认的北京市区
	}
	return
}

// 获取渠道号
func GetDeviceQdhInfo(c *gin.Context) (mark string) {
	headers := c.Request.Header
	for key, values := range headers {
		for _, value := range values {
			//获取渠道的统计标识
			if key == "mark" || key == "Mark" {
				log.Printf("匹配到mark的对应参数 key =%s value =%s \n", key, value)
				mark = value
				break
			}
		}
	}
	return mark
}

// 把第一个字母转换大写
func CapitalizeFirstChar(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// 根据特定的name获取对应的header信息
func GetRequestHeaderByName(c *gin.Context, name string) (mark string) {
	headers := c.Request.Header
	for key, values := range headers {
		if key == name || key == CapitalizeFirstChar(name) {
			tval := values[0]
			//log.Printf("匹配对应的header的对应参数 key =%s value =%s \n", key, tval)
			mark = tval
			break
		}
	}
	return mark
}

func GetGeoLite2CityByIp(ip string) (cityEnName, cityZhName, postal string) {
	cityEnName = "unknown"
	cityZhName = "unknown"
	// If you are using strings that may be invalid, check that ip is not nil
	netIp := net.ParseIP(ip)
	if netIp == nil || global.GeoCityReader == nil {
		return
	}
	record, err := global.GeoCityReader.City(netIp)
	if err != nil {
		log.Printf("GeoLite2 City Error:%s", err.Error())
		return
	}
	cityEnName = record.City.Names["en"]
	cityZhName = record.City.Names["zh-CN"]
	postal = record.Postal.Code
	return
}

// 使用商业版的IP解析
func GetIpDbNameByIp(ip string) (city_name string) {
	if ip == "" {
		return
	}
	db, err := ipdb.NewCity("public/resource/ipv4_china.ipdb")
	if err != nil {
		log.Fatal(err)
	}
	info, err := db.FindMap(ip, "CN")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("获取转换的ip = 【%v】 对应的结果为:%v\n", ip, info)
	city_name = info["city_name"]
	if city_name != "" {
		prefix := "市"
		//处理未返回的市的字段信息
		if strings.Contains(city_name, string(prefix)) == false {
			city_name = fmt.Sprintf("%v%v", city_name, prefix)
		}
	} else {
		log.Printf("未获取到IP = 【%v】 的结果，给一个默认的城市 ：北京\n", ip)
		city_name = "北京市" //默认给一个北京的
	}
	return
}

func GetCityNameByIp(ip string) (cityName string) {
	if !isValidIP(ip) {
		return
	}
	dat, _ := ioutil.ReadFile("public/resource/qqwry.dat")
	info := ip_parser.NewIpParser(ip, dat).Parse()
	cityName = strings.TrimSpace(info.City)
	return
}

func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

func EncodeImage(imgsrc, newsrc string, xorNum byte) error {
	imgData, err := ioutil.ReadFile(imgsrc)
	if err != nil {
		return err
	}

	encodedData := make([]byte, len(imgData))
	for i, value := range imgData {
		encodedData[i] = value ^ xorNum
	}

	encodedHex := hex.EncodeToString(encodedData)
	decodedData, err := hex.DecodeString(encodedHex)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(newsrc, decodedData, 0644)
}

func DecodeImage(imgsrc, newsrc string, xorNum byte) error {
	imgData, err := ioutil.ReadFile(imgsrc)
	if err != nil {
		return err
	}

	decodedData := make([]byte, len(imgData))
	for i, value := range imgData {
		decodedData[i] = value ^ xorNum
	}

	return ioutil.WriteFile(newsrc, decodedData, 0644)
}

func GetEncodeEncImage(pic string) string {
	suffix := strings.Split(pic, ".")
	dstSuffix := "_xfile." + suffix[1]
	dstPic := strings.Replace(pic, "."+suffix[1], dstSuffix, 1)
	return dstPic
}

//func ScanDirPic(dirPath string) (dirPics []models.DirPics, err error) {
//	imageExtensions := []string{".jpg", ".jpeg", ".png"} // 可根据需要添加其他图片文件扩展名
//
//	var result []models.DirPics
//
//	files, err := ioutil.ReadDir(dirPath)
//	if err != nil {
//		return
//	}
//
//	for _, file := range files {
//		if file.IsDir() {
//			subDir := filepath.Join(dirPath, file.Name())
//			subDirImages, err := ScanDirPic(subDir)
//			if err != nil {
//				fmt.Printf("遍历子目录时出错: %v\n", err)
//				continue
//			}
//
//			if len(subDirImages) > 0 {
//				result = append(result, subDirImages...)
//			}
//		} else if IsImageFile(file.Name(), imageExtensions) {
//			imagePath := filepath.Join(dirPath, file.Name())
//
//			// 查找是否已存在该子目录的记录
//			var existingDir *models.DirPics
//			for i := range result {
//				if result[i].DirPath == dirPath {
//					existingDir = &result[i]
//					break
//				}
//			}
//
//			// 如果已存在记录，则将图片路径添加到已存在的记录中
//			if existingDir != nil {
//				existingDir.DirPics = append(existingDir.DirPics, imagePath)
//			} else {
//				// 否则创建新的记录
//				result = append(result, models.DirPics{
//					DirName: filepath.Base(dirPath),
//					DirPath: dirPath,
//					DirPics: []string{imagePath},
//				})
//			}
//		}
//	}
//
//	return result, nil
//}

//func ScanDirPic(dirPath string) (dirPics []models.DirPics, err error) {
//	imageExtensions := []string{".jpg", ".jpeg", ".png"} // 可根据需要添加其他图片文件扩展名
//
//	files, err := ioutil.ReadDir(dirPath)
//	if err != nil {
//		return
//	}
//
//	var subDirs []string
//	for _, file := range files {
//		if file.IsDir() {
//			subDirs = append(subDirs, file.Name())
//		}
//	}
//
//	sort.Slice(subDirs, func(i, j int) bool {
//		return getChapterNumber(subDirs[i]) < getChapterNumber(subDirs[j])
//	}) // 按子目录中的章节号排序
//
//	for _, subDir := range subDirs {
//		subDirPath := filepath.Join(dirPath, subDir)
//
//		subDirImages, err := ScanDirPic(subDirPath)
//		if err != nil {
//			fmt.Printf("遍历子目录时出错: %v\n", err)
//			continue
//		}
//
//		if len(subDirImages) > 0 {
//			dirPics = append(dirPics, subDirImages...)
//		}
//	}
//
//	for _, file := range files {
//		if !file.IsDir() && IsImageFile(file.Name(), imageExtensions) {
//			imagePath := filepath.Join(dirPath, file.Name())
//
//			// 查找是否已存在该子目录的记录
//			var existingDir *models.DirPics
//			for i := range dirPics {
//				if dirPics[i].DirPath == dirPath {
//					existingDir = &dirPics[i]
//					break
//				}
//			}
//
//			// 如果已存在记录，则将图片路径添加到已存在的记录中
//			if existingDir != nil {
//				existingDir.DirPics = append(existingDir.DirPics, imagePath)
//			} else {
//				// 否则创建新的记录
//				dirPics = append(dirPics, models.DirPics{
//					DirName: filepath.Base(dirPath),
//					DirPath: dirPath,
//					DirPics: []string{imagePath},
//				})
//			}
//		}
//	}
//
//	return dirPics, nil
//}
//
//func getChapterNumber(dirName string) int {
//	// 从目录名称中提取章节号
//	chapterStr := strings.TrimPrefix(dirName, "第")
//	chapterStr = strings.TrimSuffix(chapterStr, "话")
//
//	chapterNum, err := strconv.Atoi(chapterStr)
//	if err != nil {
//		return 0
//	}
//
//	return chapterNum
//}

func ScanDirPic(dirPath string) (dirPics []models.DirPics, err error) {
	imageExtensions := []string{".jpg", ".jpeg", ".png"} // 可根据需要添加其他图片文件扩展名

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}

	var subDirs []string
	for _, file := range files {
		if file.IsDir() {
			subDirs = append(subDirs, file.Name())
		}
	}
	sort.Slice(subDirs, func(i, j int) bool {
		return getChapterNumber(subDirs[i]) < getChapterNumber(subDirs[j])
	}) // 按子目录中的章节号排序

	for _, subDir := range subDirs {
		subDirPath := filepath.Join(dirPath, subDir)
		subDirImages, err := ScanDirPic(subDirPath)
		if err != nil {
			fmt.Printf("遍历子目录时出错: %v\n", err)
			continue
		}

		if len(subDirImages) > 0 {
			dirPics = append(dirPics, subDirImages...)
		}
	}
	for _, file := range files {
		if !file.IsDir() && IsImageFile(file.Name(), imageExtensions) {
			imagePath := filepath.Join(dirPath, file.Name())
			//log.Println(file.Name())

			// 查找是否已存在该子目录的记录
			var existingDir *models.DirPics
			for i := range dirPics {
				if dirPics[i].DirPath == dirPath {
					existingDir = &dirPics[i]
					break
				}
			}

			// 如果已存在记录，则将图片路径添加到已存在的记录中
			if existingDir != nil {
				existingDir.DirPics = append(existingDir.DirPics, imagePath)
			} else {
				// 否则创建新的记录
				dirPics = append(dirPics, models.DirPics{
					DirName: filepath.Base(dirPath),
					DirPath: dirPath,
					DirPics: []string{imagePath},
				})
			}
		}
	}
	return dirPics, nil
}

func getChapterNumber(dirName string) int {
	// 从目录名称中提取章节号
	chapterNum, _ := getMiddleNumber(dirName)
	return chapterNum
}

func getMiddleNumber(dirName string) (int, error) {
	re := regexp.MustCompile(`第(\d+)`)
	matches := re.FindStringSubmatch(dirName)
	if len(matches) < 2 {
		return 0, fmt.Errorf("未找到匹配的数字")
	}

	numberStr := matches[1]
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return 0, fmt.Errorf("无法解析数字: %v", err)
	}

	return number, nil
}
func getChapterNumber1(dirName string) int {
	// 从目录名称中提取章节号
	chapterStr := strings.TrimPrefix(dirName, "第")
	chapterStr = strings.TrimSuffix(chapterStr, "章")

	chapterNum, err := strconv.Atoi(chapterStr)
	if err != nil {
		return 0
	}

	return chapterNum
}

func ScanBaseDir(dirPath string) (dirNames []string, err error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			dirName := file.Name()
			if dirName == "debian" || dirName == "test" {
				continue
			}
			dirNames = append(dirNames, dirName)
		}
	}
	return
}

// 判断文件是否为图片文件
func IsImageFile(filename string, imageExtensions []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}

func SerializeT(value interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	gob.Register(value)

	err := enc.Encode(&value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DeserializeT(valueBytes []byte) (interface{}, error) {
	var value interface{}
	buf := bytes.NewBuffer(valueBytes)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(&value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// @Summary 切片分页
// @Param page 当前页
// @Param pageSize 每页显示数量
// @Param nums 数据总数
// @return sliceStart 切片开始
// @return sliceEnd 切片结尾
func SlicePage(page, pageSize int, nums int) (sliceStart, sliceEnd int) {
	if page < 0 {
		page = 1
	}

	if pageSize < 0 {
		pageSize = 20
	}

	if pageSize > nums {
		if page > 1 {
			return
		}
		return sliceStart, nums
	}

	// 总页数
	pageCount := int(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize
	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}

func GetBookMd5(bookName, author string) (md5Name string) {
	md5Name = Md5(fmt.Sprintf("%v%v", strings.TrimSpace(bookName), strings.TrimSpace(author)))
	return
}

func GetChapterMd5(chapterName string) (md5Name string) {
	md5Name = Md5(strings.TrimSpace(chapterName))
	return
}

func GetSimpleHtml(strIn string) (strOut string) {
	s2t, err := opencc.New("t2s")
	if err != nil {
		log.Fatal(err)
	}
	strOut, err = s2t.Convert(strIn)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func IsHan(str string) (isHan bool) {
	// 使用正则表达式检测汉字
	re := regexp.MustCompile(`[\p{Han}]+`)
	match := re.FindString(str)
	// 如果找到汉字，则打印包含汉字
	if match != "" {
		isHan = true
	}
	return
}

func IsChapterName(chapterName string) (isChapterName bool) {
	if !IsHan(chapterName) {
		return
	}
	reg1 := []string{"第", "章", "节", "后记", "番外"}
	reg2 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	reg3 := []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "十", "百", "千", "万", "亿"}
	for _, val := range reg1 {
		if strings.Contains(chapterName, val) {
			isChapterName = true
			return
		}
	}
	for _, val := range reg2 {
		if strings.Contains(chapterName, val) {
			isChapterName = true
			return
		}
	}
	for _, val := range reg3 {
		if strings.Contains(chapterName, val) {
			isChapterName = true
			return
		}
	}
	return
}
