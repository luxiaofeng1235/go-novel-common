package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"
)

func IsAuthProxy(ip, port, username, password string) (err error) {
	// 用户名密码认证(私密代理/独享代理)
	auth := proxy.Auth{
		User:     username,
		Password: password,
	}
	proxyStr := fmt.Sprintf("%v:%v", ip, port)

	// 目标网页
	page_url := "http://myip.ipip.net"

	// 设置代理
	dialer, err := proxy.SOCKS5("tcp", proxyStr, &auth, proxy.Direct)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// 请求目标网页
	client := &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}
	req, _ := http.NewRequest("GET", page_url, nil)
	req.Header.Add("Accept-Encoding", "gzip") //使用gzip压缩传输数据让访问更快
	res, err := client.Do(req)

	if err != nil {
		// 请求发生异常
		return
	}
	defer res.Body.Close() //保证最后关闭Body

	if res.StatusCode == http.StatusOK {
		return
	}

	// 有gzip压缩时,需要解压缩读取返回内容
	//if res.Header.Get("Content-Encoding") == "gzip" {
	//	reader, _ := gzip.NewReader(res.Body) // gzip解压缩
	//	defer reader.Close()
	//	io.Copy(os.Stdout, reader)
	//	os.Exit(0) // 正常退出
	//}

	// 无gzip压缩, 读取返回内容
	//body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	return
}

//func GetS5() {
//	S5RankUrl = fmt.Sprintf("%v", S5RankUrl)
//	s5html := HttpGetRequest(S5RankUrl, nil)
//	if s5html == "" {
//		time.Sleep(time.Second)
//		S5RankUrl = fmt.Sprintf("%v&token=%v", S5RankUrl, Token2)
//		GetS5()
//		return
//	}
//	s5Res := models.Socket5Res{}
//	_ = json.Unmarshal([]byte(s5html), &s5Res)
//	if s5Res.Code != 1 {
//		time.Sleep(time.Second * 5)
//		GetS5()
//		return
//	}
//	if len(s5Res.Data.List) <= 0 {
//		time.Sleep(time.Second * 5)
//		GetS5()
//		return
//	}
//	s5Proxy := s5Res.Data.List[0]
//	S5Domain = s5Proxy.Ip
//	S5Port = s5Proxy.Port
//	S5Username = s5Proxy.Username
//	S5Passwd = s5Proxy.Password
//	log.Println(S5Domain, S5Port, S5Port, S5Passwd)
//	return
//}

func IsProxy(ip, port string) (err error) {
	// 代理服务器
	proxyStr := fmt.Sprintf("%v:%v", ip, port)

	// 目标网页
	page_url := "http://myip.ipip.net"
	auth := proxy.Auth{}

	// 设置代理
	dialer, err := proxy.SOCKS5("tcp", proxyStr, &auth, proxy.Direct)
	if err != nil {
		return
	}

	// 请求目标网页
	client := &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}
	req, _ := http.NewRequest("GET", page_url, nil)
	req.Header.Add("Accept-Encoding", "gzip") //使用gzip压缩传输数据让访问更快
	res, err := client.Do(req)

	if err != nil {
		// 请求发生异常
		return
	}

	defer res.Body.Close() //保证最后关闭Body

	if res.StatusCode == http.StatusOK {
		return
	}
	return
	// 有gzip压缩时,需要解压缩读取返回内容
	//if res.Header.Get("Content-Encoding") == "gzip" {
	//	reader, _ := gzip.NewReader(res.Body) // gzip解压缩
	//	defer reader.Close()
	//	io.Copy(os.Stdout, reader)
	//	return
	//}
	//
	//// 无gzip压缩, 读取返回内容
	//body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
}

//func GetS5() {
//	if !IsS5 {
//		return
//	}
//	if GetUnix() < S5ExpireTime {
//		global.Collectlog.Errorf("时间未到 现在时间%v 过期时间%v", UnixToDatetime(GetUnix()), UnixToDatetime(S5ExpireTime))
//		return
//	}
//	if S5Type != Rank {
//		global.Errlog.Errorf("%v", "代理非随机")
//		return
//	}
//	s5html := HttpGetRequest(S5RankUrl, nil)
//	if s5html == "" {
//		time.Sleep(time.Second * 3)
//		return
//	}
//	s5Res := models.ZhimaSocket5Res{}
//	var err error
//	err = json.Unmarshal([]byte(s5html), &s5Res)
//	if err != nil {
//		global.Errlog.Errorf("%v", err.Error())
//		time.Sleep(time.Second * 3)
//		GetS5()
//		return
//	}
//	if s5Res.Code != 0 {
//		global.Errlog.Errorf("%v", s5Res.Msg)
//		time.Sleep(time.Second * 5)
//		return
//	}
//	if len(s5Res.Data) <= 0 {
//		time.Sleep(time.Second * 5)
//		return
//	}
//	s5Proxy := s5Res.Data[0]
//	S5Domain = s5Proxy.Ip
//	S5Port = fmt.Sprintf("%v", s5Proxy.Port)
//	S5ExpireTime = DateToUnix(s5Proxy.ExpireTime)
//	log.Println("芝麻获取代理", S5Domain, S5Port, s5Proxy.ExpireTime, S5ExpireTime)
//	return
//}

func GetS5() {
	if !IsS5 {
		return
	}
	if S5Type != Rank {
		global.Errlog.Errorf("%v", "代理非随机")
		return
	}
	s5html := HttpGetRequest(S5RankUrl2, nil)
	if s5html == "" {
		time.Sleep(time.Second * 3)
		return
	}
	s5Res := models.YilianSocket5Res{}
	var err error
	err = json.Unmarshal([]byte(s5html), &s5Res)
	if err != nil {
		global.Errlog.Errorf("%v", err.Error())
		time.Sleep(time.Second * 3)
		GetS5()
		return
	}
	if s5Res.Errcode != 0 {
		global.Errlog.Errorf("%v", s5Res.Errmsg)
		time.Sleep(time.Second * 5)
		return
	}
	if len(s5Res.Data) <= 0 {
		time.Sleep(time.Second * 5)
		return
	}
	s5Proxy := s5Res.Data[0]
	S5Domain = s5Proxy.Ip
	S5Port = fmt.Sprintf("%v", s5Proxy.Port)
	S5Username = s5Proxy.Username
	S5Passwd = s5Proxy.Passwd
	S5ExpireTime = DateToUnix(s5Proxy.ExpireTime)
	log.Println("一连获取代理", S5Domain, S5Port, s5Proxy.ExpireTime, S5ExpireTime)
	return
}

//func GetS5() {
//	if !IsS5 {
//		return
//	}
//	//if GetUnix() < S5ExpireTime {
//	//	global.Collectlog.Errorf("时间未到 现在时间%v 过期时间%v", UnixToDatetime(GetUnix()), UnixToDatetime(S5ExpireTime))
//	//	return
//	//}
//	if S5Type != Rank {
//		global.Errlog.Errorf("%v", "代理非随机")
//		return
//	}
//	jsonFile := "./S5Proxys.json"
//	var err error
//	if len(S5Proxys) <= 0 {
//		// 读取文件内容
//		if CheckNotExist(jsonFile) {
//			err = fmt.Errorf("%v", "内容不存在")
//			return
//		}
//		var conByte []byte
//		conByte, err = ioutil.ReadFile(jsonFile)
//		if err != nil {
//			return
//		}
//		content := string(conByte)
//		s5ips := []models.Socket5Proxy{}
//		err = json.Unmarshal([]byte(content), &s5ips)
//		if err != nil {
//			fmt.Println("Failed to parse JSON:", err)
//			return
//		}
//		S5Proxys = s5ips
//		log.Println(len(S5Proxys))
//		return
//	}
//
//	var info models.Socket5Proxy
//	if S5ProxyIndex > len(S5Proxys)-1 {
//		S5ProxyIndex = 0
//	} else {
//		S5ProxyIndex++
//	}
//	info = S5Proxys[S5ProxyIndex]
//	S5Domain = info.Ip
//	S5Port = info.Port
//	S5Username = info.Username
//	S5Passwd = info.Password
//	S5ExpireTime = DateToUnix(info.OnlineDate)
//	log.Println("获取代理", S5ProxyIndex, S5Domain, S5Port, S5Username, S5Passwd, UnixToDatetime(S5ExpireTime), S5ExpireTime)
//	return
//}

func RefreshS5Proxys() {
	jsonFile := "./S5Proxys.json"
	var err error
	var proxys []models.Socket5Proxy
	keys := GetPostalKeys()
	for _, key := range keys {
		mvals := global.Redis.HGetAll(context.Background(), key).Val()
		for _, val := range mvals {
			var proxy models.Socket5Proxy
			err = json.Unmarshal([]byte(val), &proxy)
			if err != nil {
				log.Println(val, err.Error())
			} else {
				//OnlineUnix := DateToUnix(proxy.OnlineDate)
				//if OnlineUnix < GetTodayUnix() {
				//	continue
				//}
				proxys = append(proxys, proxy)
			}
		}
	}
	S5Proxys = proxys
	var jsonData []byte
	jsonData, err = json.MarshalIndent(S5Proxys, "", "  ")
	if err != nil {
		err = fmt.Errorf("获取章节信息失败%v", err.Error())
		return
	}
	err = WriteFile(jsonFile, string(jsonData))
	if err != nil {
		global.Errlog.Errorf("%v", err.Error())
		return
	}

	// 读取文件内容
	if CheckNotExist(jsonFile) {
		err = fmt.Errorf("%v", "内容不存在")
		return
	}
	var conByte []byte
	conByte, err = ioutil.ReadFile(jsonFile)
	if err != nil {
		return
	}
	content := string(conByte)
	s5ips := []models.Socket5Proxy{}
	err = json.Unmarshal([]byte(content), &s5ips)
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
		return
	}
	S5Proxys = s5ips
	log.Println(len(S5Proxys))
	return
}

func GetPostalKeys() (keys []string) {
	//keys = []string{"CN:110000", "CN:130000", "CN:140000", "CN:310000", "CN:440000"}
	//return
	ctx := context.Background()
	var redis = global.Redis
	var cursor uint64
	iter := redis.Scan(ctx, cursor, "CN:*", 20).Iterator()
	for iter.Next(ctx) {
		val := iter.Val()
		if val == "" {
			continue
		}
		if len(val) != 9 {
			continue
		}
		if !unicode.IsUpper([]rune(val)[0]) {
			continue
		}
		if redis.Type(ctx, val).Val() != "hash" {
			continue
		}
		if global.Redis.HLen(ctx, val).Val() <= 0 {
			continue
		}
		keys = append(keys, val)
	}
	if len(keys) <= 0 {
		return
	}
	return
}

func GetRandPostalKeys() (postalKey, postal string) {
	keys := GetPostalKeys()
	//countrys := []string{"MX", "RU", "EC", "CA", "IT", "DE", "GB", "PH", "ZA", "AR", "UY", "PT", "BH"}
	if len(keys) <= 0 {
		return
	}
	index := RangeNum(0, len(keys))
	postalKey = keys[index]
	if strings.Contains(postalKey, ":") {
		parts := strings.Split(postalKey, ":")
		if len(parts) > 1 {
			postal = parts[len(parts)-1]
		}
	}
	return
}
