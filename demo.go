package main

import (
	"fmt"
	"go-novel/app/service/api/book_service"
	"go-novel/app/service/common/common_service"
	"go-novel/db"
	"log"
	"time"
)

func main() {

	// 添加一个字符串值

	// 添加一个整数值
	// 添加一个嵌套的 map

	//nestedData := make(map[string]interface{})
	//nestedData["id"] = 123
	//nestedData["text"] = "CA"
	//angle := make(map[string]interface{})
	//angle["book_id"] = 123
	//angle["book_name"] = "22334"
	//angle["author"] = "测试"
	//nestedData["document"] = angle
	//fmt.Println(nestedData)
	//jsonbyte, err := json.Marshal(nestedData)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(jsonbyte))
	////data["address"] = nestedData
	////// 添加一个切片
	////data["hobbies"] = []string{"reading", "hiking", "cooking"}
	//
	//return
	//获取当前的配置信息
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	db.InitDB()
	db.InitZapLog()
	db.InitNsqProducer()
	db.InitNsqConsumer()
	db.InitGeoReadre()
	//redisKey := "testdata"
	//fmt.Println(redisKey)

	////缓存获取的书籍信息

	//var user book_service.User
	//user.Id = 1
	//user.Age = 3
	//user.Username = "测试或者哪个好啦"
	//jsonData, err := json.Marshal(user)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//callback := models.ClickCallback{}

	//list, err := book_service.GetBookCityList()
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(list)
	//cityName := "深圳市"
	//index := book_service.IsContainCity(list, cityName)
	//if index {
	//	fmt.Println("包含版权城市")
	//} else {
	//	fmt.Println("未包含版权城市")
	//}
	//return
	status, _ := book_service.GetBookCopyright("android", "com.realbest.novelread", "223.64.28.207", "android-sm")
	fmt.Println(status)
	return

	// 示例时间戳（Unix 时间戳，单位为秒）
	// 示例时间戳（Unix 时间戳，单位为秒）
	timestamp := int64(1726277964) // 例如：2023年9月14日的时间戳

	// 将时间戳转换为 Time 类型
	t := time.Unix(timestamp, 0)

	// 获取当前时间
	now := time.Now()

	// 计算从时间戳到当前时间的持续时间
	duration := now.Sub(t)

	// 将持续时间转换为小时
	hours := duration.Hours()

	// 打印结果
	fmt.Printf("从时间戳到现在的小时数: %.0f 小时\n", hours)
	return

	//callback, err := callback.GetCountByImeiAndOaidType("0561cecbe707802f", "shenma")
	callback, _ := common_service.GetShenmaClickInfo("0561cecbe707802f", "", "shenma")
	//if err != nil {
	//	fmt.Println("111", err)
	//	return
	//}
	fmt.Println(callback)
	fmt.Println(callback.Id)
	//fmt.Println(callback.Oaid)
	//key := "testarr"
	//contents := gredis.Get(key)
	//fmt.Println(contents)
	////cacheData := string(jsonData)
	////err = gredis.Set(key, cacheData, time.Hour*1)
	////if err != nil {
	////	fmt.Println(err)
	////	return
	////}
	//
	//var response book_service.User
	//err := json.Unmarshal([]byte(contents), &response)
	//if err != nil {
	//	log.Printf("Error parsing JSON results = %+v", err)
	//	return
	//}
	//fmt.Println(response, response.Id, response.Age, response.Username)
	//fmt.Println(122)
	////fmt.Println(1111)
	////fmt.Println(user)
	//res := book_service.SetBookRankList(key, user, 0)
	////获取redis中的缓存信息
	//vv := book_service.GetBookListByCacheKey(key)
	//fmt.Println(res)
	//fmt.Println(vv)
	//return

	//redisKey := "testarr"
	//chaptersVal := redis_service.Get(redisKey)
	//fmt.Println(chaptersVal)
	//return
	//err := redis_service.Set(redisKey, "1111111", 0)
	//if err != nil {
	//	err = fmt.Errorf("redis缓存失败 err=%v", err.Error())
	//	return
	//}
	//fmt.Println(redisKey)
	//return
	//aa, err := book_service.GetBookCopyright("", "", "127.0.0.1", "android-test")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(aa)
	//return

	////d41d8cd98f00b204e9800998ecf8427e
	//aaa := utils.Md5("123!@#Gg")
	//fmt.Println(aaa)
	//return
	//
	//url := "https://huichuan.uc.cn/callback/appapi?click_id=5965565233090242536&a_ty=bDo1MTt2OjExNjAyMTI5&event_type=0&sid=15097363324016059946&uctrackid=czoxNTA5NzM2MzMyNDAxNjA1OTk0NjtjOjE1MTU0NTc3ODtkOmRtcF8yNjU4OTE2ODY0NTU2NjIzMDg1O3A6aGM&act_type=&uid=210908896"
	//aa := utils.GetReplaceChaojihuiCallbak(url, 1001, "")
	//fmt.Println(aa)
	//return
	//
	//string_number := 645270
	//fmt.Println(string_number / 1000)
	//result := strconv.Itoa(string_number / 1000)
	//
	//fmt.Println(result)
	//fmt.Printf("%T\n", result)
	//return

	//kdAdIo0CNspnYHwDKXZPT+l1ATSRcr6x1o/gbbjxUChMY7U9n/oH04RcAODNgZpU 这个可以
	//ycr+8hXsyU5AzGROdIrLisNkRbpnY/9fV4xNNGQg7Xg= 这个不可用
	//res, err := base64.RawStdEncoding.DecodeString("ycr+8hXsyU5AzGROdIrLisNkRbpnY/9fV4xNNGQg7Xg=")
	// base64.RawStdEncoding.DecodeString(src)
	//fmt.Println(res)
	//return

	//url := "https://book.prod-book.iolty.xyz/book/source/v4/920/920741.html"
	//res, err := utils.ComicApiDecryptV1("U0Q4WTNINENXMkZKUFQ3RunJLYyYi0L6SMYZWRCJTI0zAFnTkpvLr0+CbQi4oNGoeo+z1cZoWt/hNfYshdGPKt3S1tQWbiW1EGL1UzvS9jNDiTenUwybmHQ/S9a2cgVxINcy78UTaeH+Z8qyJnulSspkk4SuBstThzQzxCie3G9wXpTRtNm2w9m5VLi5R4nAYTZkM4AHX92ozpIZv8oO/woPAe6vrZXdnobAmhdJFOzaTpZ3Lp36b0/S8eBvISV4D1zAkwEDqqbmYw254y74Zoy7yX27/T0+rYkjf43J0yl+ex0d/BGgLmXdIQYZqnszdoEedW7Xzn96FdnSXgT0xbiuirLInR5O4Yk2sJdvSQYaN4K2C6Q8SJLDNPIGc9CZDsnLey+rYkW2Z2LZw8/kT1y+ajodsgXkEgDj92nDlZ1s2xSUc45bh6KFV4007aaFfFmZ9PybbVyNa+p4yyR3tY7INc2efI+w9QW3RbUCNGqr72ssS6HkVlY4uyCtRnPFbP5AXBLGYAdl7ykcYYqeVji/d11Xze0TY9N7a7J4HzqiZbu76VBHnQd7URIrn7cJfZpMuX45B7bSj1A76yrVY43FqtVy97Sd5HLWSRcR9LBWL/b8xnceBz9tQmNqni8weo8BDjM9mHRh5pf+CDvWyI/c2eJ+bolhtt3K7WsLsL94UEKek4+lHWHsOiV89DMljZfOboLdmiQKoOgazfv/Dly3ZCpfUwYKgs/MQridn6xf7roc2SRcuqsx43RJDFrfV6RNWC6FDFazgTIQPx+rcVo8G5lA119Shji/H8qjCbRbUUvlEaih0lk0C6GHsiTQXFAmFW9HniZZyHnK+0SitEVRQ1MzQ1FVSlg0WjhWQlk=")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(res)
	//return

	//url := "https://chapter.chuangke.tv/starfire/97/7c/ce/1048/2.html"
	//contents := utils.GetBaiduResponse(url)
	//storeContents := gjson.Get(contents, "data.content").String()
	//// 原始二进制数据

	//jsonData := `{
	//    "code": 1,
	//    "data": [
	//        {
	//            "name": "SwQ85IjDQt/1pbQs8zjIBnEaqz/dPWRoUblSPGln1SA=",
	//            "url": "https://chapter.chuangke.tv/starfire/97/7c/ce/1048/1.html",
	//            "is_content": true,
	//            "path": "starfire/97/7c/ce/1048/1.html",
	//            "updated_at": 1679972356
	//        }
	//    ],
	//    "updated_at": 1719163624
	//}`
	//
	//var response utils.BiqugeChapterListItem
	//err := json.Unmarshal([]byte(jsonData), &response)
	//if err != nil {
	//	log.Fatalf("Error parsing JSON: %s", err)
	//}
	//fmt.Printf("Code: %d\n", response.Code)
	//for _, item := range response.Data {
	//	fmt.Printf("Name: %s\n", item.Name)
	//	fmt.Printf("URL: %s\n", item.URL)
	//	fmt.Printf("Is Content: %t\n", item.IsContent)
	//	fmt.Printf("Path: %s\n", item.Path)
	//	fmt.Printf("Updated At: %d\n", item.UpdatedAt)
	//}
	//fmt.Printf("Updated At (Response): %d\n", response.UpdatedAt)
	//return

	//str := "Cr2iOmUL0hEx/4uyNk3XhC2bdp5zhbJ9K15nKk/mo4o="

	//area := utils.GetIpDbNameByIp("223.102.17.150")
	//fmt.Println(area)
	//return
	//
	//url := utils.GetApkFileUrl("/www/wwwroot/down.mnjkfup.cn/apk/z41n7ivrz9lybkurz9.jpg")
	//fmt.Println(url)
	//return
	//推荐类型  rec_serialize 推荐完结（完成）  -  rec_new 推荐新书 (完成) - hot_serialize 热门完结 - hot_rank 热门排行 - hot_new 热门新书 - hot_search 热门搜索 - classic_search 经典热搜 - classic_hight  经典高分 - classic_rq 经典人气 - classic_serialize 经典完结 - classic_new 经典新书
	//res, err := book_service.GetApiBookRecByType("hot_serialize")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(res)
	//return

	//img := "http://picc.mnjkfup.cn/data/pic/202405/jmmdnrhfg-bl.jpg"
	//res := utils.GetFileUrl(img)
	//fmt.Println(res)
	//return
	//
	//ip := "120.197.198.79"
	//data := book_service.GetIpString(ip)
	//fmt.Println(data)
	//return

	//res := common_service.XunDelBookById("delIndex", 9994567)
	//fmt.Println(res)

	//res := common_service.XunAddBookInfo("addIndex", 9994567, "卢晓峰测试", "哈哈我得呢")
	//fmt.Println(res)
	//jsonData := common_service.XunSearchByBookName("searchList", "风水女道士寻女复仇记")
	//var resp *common_service.XunResponse
	//err := json.Unmarshal([]byte(jsonData), &resp)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//// 打印 data 里的数据
	//for _, book := range resp.Data {
	//	fmt.Printf("Book ID: %d\n", book.BookID)
	//	fmt.Printf("Book Name: %s\n", book.BookName)
	//	fmt.Printf("Author: %s\n", book.Author)
	//	fmt.Printf("Chrono: %d\n", book.Chrono)
	//	fmt.Println()
	//}
	//return

	//aaa := common_service.AddSearchBookInfo(357735)
	//fmt.Println(aaa)
	//return
	//
	////选择默认的url信息
	//flag := common_service.DelSearchDataById(3115)
	//fmt.Println(flag)
	//return
	//content := common_service.GoqueryByBookName("九曜", 1, 20, "[document.book_id]*10", "desc")
	//fmt.Println("11112", content)
	//var resp *common_service.Response
	//err := json.Unmarshal([]byte(content), &resp)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	////遍历获取对应的搜索采集信息
	//for _, val := range resp.Data.Documents {
	//	log.Printf("Document BookId: %v Document Title: %v Document Author:%v\n", val.Document.BookId, val.Document.BookName, val.Document.Author)
	//}
	//return
	//
	//ftpClient, err := setting_service.NewFTPClient("103.36.91.36:21", "kyks", "5JChEJ4Ztk82i8aR")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = ftpClient.CreateFolder("/aaa/tewt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("目录创建完成")
	//return
	//
	////sdata := setting_service.PrivacyTemplate("测试哈哈") //隐私协议内容
	//sdata := setting_service.UserAgreementTemplate("测试123")
	//fmt.Println(sdata)
	//return
	//
	//res := utils.GetAdminFileUrl("/mnt/upload/apk/r5j6j0dve6xve6av6o.jpg")
	//fmt.Println(res)
	//return
	////nowTime := time.Now()
	////var ExpireTime int64 = 24
	////expireTime := nowTime.Add(time.Duration(ExpireTime) * time.Hour * 90)
	////fmt.Println(expireTime)
	////return
	////originalString := "Hello, World! Hello, Gopher!"
	////replacedString := strings.ReplaceAll(originalString, "Hello", "Hi")
	////fmt.Println("Original string:", originalString)
	////fmt.Println("Replaced string:", replacedString)
	////ss := utils.GetCityByNetwork("125.37.9.144")
	////fmt.Println(ss)
	//
	////t := book_service.GetRandTest()
	////fmt.Println(t)
	////return
	//
	////str := "mark"
	////capitalizedStr := utils.CapitalizeFirstChar(str)
	//
	//mark := "android-xiaomi"
	//prefix := "xiaomi"
	//if strings.Contains(mark, prefix) != false {
	//	fmt.Println("ok")
	//} else {
	//	fmt.Println("error")
	//}
	//return
	//
	//aa, err := user_service.GetUserInfoByMailOrTel("513072539@qq.com", "")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(aa)
	//return
	//var getIp = "2408:8221:32b:2e81:e0d1:2ed8:517c:8454"
	//ipData := book_service.GetIpString(getIp)
	//fmt.Println(ipData)
	//
	////aa := utils.GetParseIpByV4("39.144.177.37")
	////aa := utils.GetParseIpByV6("2408:8221:32b:2e81::a25")
	////fmt.Println(aa)
	////return
	//
	////encodedURL := "/api/common/baiduStatistics?imei_md5=cfcd208495d565ef66e7dff9f98764da&os=2&ip=117.150.211.15&ua=Dalvik%2F2.1.0+%28Linux%3B+U%3B+Android+9%3B+PDBM00+Build%2FPPR1.180610.011%29+baiduboxapp%2F13.62.0.11+%28Baidu%3B+P1+9%29&ts=1722923165000&userid=57355236&pid=583260523&uid=10175191363&aid=99035857062&click_id=6a979670c6c895c7_1722923166&oaid=29B7F96B1ED148639FF99C958E41CDEC35584cb841f33cc4c4c272a414d40752&callback_url=http%3A%2F%2Focpc.baidu.com%2Focpcapi%2Fcb%2FactionCb%3Fa_type%3D%7B%7BATYPE%7D%7D%26a_value%3D%7B%7BAVALUE%7D%7D%26s%3D7680773100619208135%26ext_info%3DH-wWXb4RNg-PX-bq0HRLn1RdnWnv0HR3n1cvnjRznsDkrH0v0HcKnj6LfWmvnbnLfWKaPjTdwbD3nbm4wjNAnHKarH6vfW0KwDcvrRP7njnYf1cLfYn4wjbvnj6dPjN7wbPjnjD4rRNtHsDkP1czrHc1nHmd0HDs0HmvPj0Y0HDKn0Ds0HTvrj0LP1nknj0vnHbznj6kn1RKrjRzPWDLPWTkrH0Y0HnK&akey=NTczNTUyMzY=&ip_type=v4&interactionsType=1&sign=59b8cbfe4b8c208174ae6c5bc29be3c8"
	////使用URLcode来进行解码
	////encodedURL := "https://example.com/search?q=hello%20world"
	////decodedURL, err := url.QueryUnescape(encodedURL)
	////if err != nil {
	////	fmt.Println("Error decoding URL:", err)
	////	return
	////}
	//u, err := url.Parse(encodedURL)
	//if err != nil {
	//	fmt.Println("Error parsing URL:", err)
	//	return
	//}
	//callbackURL := u.Query().Get("callback_url")
	//if callbackURL != "" && strings.Contains(callbackURL, "baidu") != false {
	//	fmt.Println("parse callback url:", callbackURL)
	//	activateUrl := utils.GetReplaceBaiduCallbak(callbackURL, "activate", "0")
	//	fmt.Println("activate url:", activateUrl)
	//	//获取百度请求的地址信息
	//	aa := utils.GetBaiduResponse("http://www.baidu.com")
	//	fmt.Println(aa)
	//} else {
	//	fmt.Println("解析callback参数为空或者回调异常")
	//}
	//
	//return
	//fmt.Println("Decoded URL:", decodedURL)

	////处理获取的url信息
	//callback_url := "http://ocpc.baidu.com/ocpcapi/cb/actionCb?a_type={{ATYPE}}&a_value={{AVALUE}}&s=1412327566783400035&ext_info=H-wWXb4RNg-PX-bq0HRLn1RdnWnv0HR3n1cvnjRzP0DkrH0v0HcKrjfswWwAnYDLnjF7nbwDwRFKPWF7nHK7wj6LPHwjfHbKPHKAPj0vwRNDPjIDwHPjfRDzPHRzwDFjwbn3P1uAP1FtN-KKwY7fP-FK0HDLnWc4nWDknHnKnH0Kn1mYnjfKnfDs0H0KnHfknWnzP1RvPWT3n1fsnj01PfD3PHcvnHT3rHnsnHmKnsD&akey=NTczNTUyMzY=&ip_type=v4&interactionsType=1&sign=b6add2236edbf050d9d41fab0d102270"
	////获取百度的回调地址并进行替换
	//registerUrl := utils.GetReplaceBaiduCallbak(callback_url, "register", "0")
	//fmt.Println("Register:", registerUrl)
	//return
	//fmt.Println(callback_url)
	//return
	//
	//res, err := adver_service.GetPackageByCondition("测试", "")
	//if err != nil {
	//	log.Println(err)
	//}
	////str := strings.Join(res, ",")
	//log.Println(res)
	////fmt.Println(str, reflect.TypeOf(str))
	//return
	////log.Println(utils.GetGeoLite2CityByIp("116.235.238.47"))
	//return

	//url := "http://ip-api.com/json/116.235.238.47?lang=zh-CN"
	//resp, err := http.Get(url)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer resp.Body.Close()
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println(string(body))
	//return
	//
	//res := utils.GetGeoCiyByIp("114.251.193.153")
	//fmt.Println(res)
	//return
	//geodb, err := geoip2.Open("public/resource/GeoLite2-City_20221007/GeoLite2-City.mmdb")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer geodb.Close()
	//
	//ip := net.ParseIP("121.8.215.106")
	//record, err := geodb.City(ip)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("city_name：%v", record.City.Names["zh-CN"])
	//fmt.Printf("City: %vn", record.City.Names["en"])
	//fmt.Printf("Country: %vn", record.Country.Names["en"])
	//fmt.Printf("Latitude: %v, Longitude: %vn", record.Location.Latitude, record.Location.Longitude)
	//return

	//fmt.Printf("当前解析的IP=%v对应的城市为：%v\n", getIp, ipData)
	//return
	//res := utils.GetCityByNetwork("39.144.26.93")
	//fmt.Println(res)
	//return
	//
	//ss := utils.GetGeoCiyByIp("39.144.26.93")
	//fmt.Println(ss)
	//return
	//
	//ip := "121.32.198.45"
	//mark := "huawei"
	//check_status, err := book_service.GetBookCopyright("ios", "哈哈14", ip, mark)
	//if err != nil {
	//
	//	fmt.Println("11111", err)
	//}
	//fmt.Println(check_status)
	//////return
	////ip := "223.104.255.55"
	////cityName := utils.GetCityNameByIp(ip)
	////fmt.Println(cityName)
	//fmt.Println(ip)
	//return

	//jsonData := `{"id":4,"project_name":"123","app_id":"12334","package_name":"哈哈1","device_type":"ios","push_data":[{"adver_type":1,"adver_type_name":"书架广告","adver_value_string":"23451"},{"adver_type":2,"adver_type_name":"开屏广告","adver_value_string":"31234"}]}`
	//
	//var data map[string]interface{}
	//err := json.Unmarshal([]byte(jsonData), &data)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//
	//pushData := data["push_data"].([]interface{})
	//packge_id := 3
	//if len(pushData) != 0 {
	//	for _, item := range pushData {
	//		adItem := item.(map[string]interface{})
	//		adverType := int(adItem["adver_type"].(float64))
	//		adverTypeName := adItem["adver_type_name"].(string)
	//		adverValueString := adItem["adver_value_string"].(string)
	//		//projectName := adItem["project_name"].(string)
	//		//fmt.Println(projectName)
	//		fmt.Println(adverType)
	//		fmt.Printf("Adver Type: %d, Adver Type Name: %s, Adver Value String: %s\n", adverType, adverTypeName, adverValueString)
	//		fmt.Println(packge_id)
	//	}
	//}
	//return
	//
	//bookName := "心声泄漏后，真的可以为所欲为！"
	//author := "青烟往"
	//bookFile := utils.GetBookMd5(bookName, author) //获取小说和作者的加密值
	//uploadBookChapterPath := "/data/chapter/" + bookFile[0:2] + "/"
	//fmt.Println(bookName, author)
	//fmt.Println(bookFile)
	//fmt.Println(uploadBookChapterPath)
	//return
	////目录组成结构：小说章节+MD5的前两个字符作为存储的路径信息
	////var uploadBookChapterPath string
	////uploadBookChapterPath = "/data/chapter/"
	//
	////mobile := "13146899753"
	////fmt.Printf("发送的短信手机号为：%s\r", mobile)
	//////发送验证码
	////aaa := utils.GenValidateCode(6)
	////fmt.Println(aaa)
	////return
	////smsRes, err := common_service.TencentSmsSend(mobile)
	////if err != nil {
	////	log.Println(err)
	////	return
	////}
	////fmt.Println(smsRes)
	////////判断是否发送成功，给一个提示信息
	////result := utils.GetSmsSendCode(smsRes)
	////if result != "Ok" {
	////	fmt.Println("发送失败")
	////} else {
	////	fmt.Println("发送成功")
	////}
	////return
	////bookFile := utils.GetBookMd5("哈哈哈", "侧事故")
	////firstPath := bookFile[0:2]
	////fmt.Println(bookFile)
	////fmt.Println(firstPath)
	////mondayDate := utils.GetThisWeekFirstDate() //获取每周的周一的日期
	////fmt.Println(mondayDate)
	////fmt.Println(reflect.TypeOf(mondayDate))
	////savePicPath := "/data/pic/" + mondayDate + "/" //拼装存储的图片路径
	////
	//////mondayDate := utils.GetThisWeekFirstDate()
	//////fmt.Println(reflect.TypeOf(mondayDate))
	////fmt.Println(savePicPath)
	////return
	//////monday, err := GetMondayOfWeek()
	//////if err != nil {
	//////	log.Fatal(err)
	//////}
	//////res := monday.Format("20060102")
	//////fmt.Println(res)
	//////return
	//////fmt.Println(monday)
	//////return
	//////fmt.Println("本周的星期一是：", monday.Format("20060102"))
	////return
	////log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ltime)
	////host, name, user, passwd := db.GetDB()
	////db.InitDB(host, name, user, passwd)
	////addr, passwd, defaultdb := db.GetRedis()
	////db.InitRedis(addr, passwd, defaultdb)
	////db.InitZapLog()
	////db.InitNsqProducer()
	////db.InitNsqConsumer()
	////db.InitKeyLock()
	////
	////email_to := "513072539@qq.com"
	////content := "这是测试的一个邮件，发送系统测试用"
	////
	////ret := utils.SendEmail(email_to, content)
	////fmt.Println(ret)
	////return
	//////var bookId int64
	//////bookId = 318077
	////////获取单条的数据信息
	////pageSize := 5
	//
	////bookInfo, err := book_service.GetCommentUserList(pageSize)
	////if err != nil {
	////	log.Println(err)
	////}
	////fmt.Println(bookInfo)
	////return
	//
	////fmt.Println(bookInfo)
	////fmt.Printf("p1=%#v\n", bookInfo)
	//
	////keywords := "全球复苏"
	////str := fmt.Sprintf("MATCH (book_name) AGAINST ('\"%s\"' IN BOOLEAN MODE)", keywords)
	////fmt.Println(str)
	////return
	////name_list := []int{1, 2, 3, 4, 5}
	////target1 := 6
	////index := adver_service.IsContainInt(name_list, target1)
	////if index {
	////	fmt.Println("对象在数组中")
	////} else {
	////	fmt.Println("对象不在数组中")
	////}
	////return
	//
	////madver := make(map[string]interface{})
	//// 填充 madver
	//
	////var slice []Type
	////for _, value := range madver {
	////	if converted, ok := value.(Type); ok {
	////		slice = append(slice, converted)
	////	} else {
	////		fmt.Printf("无法将值：%v转为Type类型\n", value)
	////	}
	////}
	//
	////b := []byte{1, 2, 3}
	////b = append(b, 4, 5, 6)
	////fmt.Println(b)
	////id := 5
	////adverInfo, err := adver_service.GetAdverInfoById(int64(id))
	////if err != nil {
	////	fmt.Println(err)
	////	return
	////}
	////fmt.Println(adverInfo)
	//////
	////if adverInfo.Pic != "" {
	////	adverInfo.Pic = utils.GetFileUrl(adverInfo.Pic)
	////}
	////
	//////判断如果满足对应的杀进程次数的配置条件后，直接返回总的杀进程总次数
	////adverPosition := adverInfo.AdverPosition
	////killProcessTimes := 0
	////if adverPosition > 0 && adverPosition == 3 {
	////	killProcessTimes = 4
	////}
	////
	////typeInfo := map[string]interface{}{
	////	"id":             adverInfo.Id,
	////	"adver_type":     adverInfo.AdverType,
	////	"adver_name":     adverInfo.AdverName,
	////	"status":         adverInfo.Status,
	////	"pic":            adverInfo.Pic,
	////	"adver_link":     adverInfo.AdverLink,
	////	"weight":         adverInfo.Weight,
	////	"adver_position": adverInfo.AdverPosition,
	////	"error_num":      adverInfo.ErrorNum,
	////	"addtime":        adverInfo.Addtime,
	////	"uptime":         adverInfo.Uptime,
	////	"kill_times":     killProcessTimes,
	////}
	////
	////fmt.Println(utils.JSONString(typeInfo))
	//
	////////////////////////这个是获取对应的接口信息的剩余新手保护期的时间
	//userId := 111 //自定义的配置信息
	//newUserId := int64(userId)
	//beyondState := book_service.GetReadToUserBeyond(newUserId)
	//fmt.Println(beyondState)
	////if beyondState > 0 {
	////	fmt.Printf("当前用户还是新手保护，剩余时间为 %d", beyondState)
	////} else {
	//	fmt.Printf("当前用户已经过手保护期限,剩余时间为 %d", beyondState)
	//}
}
