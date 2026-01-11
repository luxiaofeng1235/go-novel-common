package main

import (
	"bytes"
	"crypto/tls"
	"go-novel/app/models"
	"go-novel/db"
	"go-novel/utils"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func main() {
	host, name, user, passwd := db.GetDB()
	db.InitDB(host, name, user, passwd)
	db.InitZapLog()
	html, _ := utils.GetHtml("https://www.paoshuba.cc/Partlist/452/26837.html", "utf-8", 1, 0)
	if html == "" {
		//err = fmt.Errorf("获取小说详情页面失败 bookUrl=%v", "bookUrl")
		return
	}
	log.Println(html)
	//var collect *models.McCollect
	//global.DB.Model(models.McCollect{}).Where("id = 1").First(&collect)
	//_ = fieldContent(collect, "https://www.paoshuba.cc/Partlist/452/26837.html")
}

func fieldContent(collect *models.McCollect, bookUrl string) (err error) {
	// GET请求
	//html, err := curlGet(bookUrl)
	//if err != nil {
	//	fmt.Println("GET请求失败:", err)
	//} else {
	//	fmt.Println("GET请求结果:", html)
	//}
	//Doc, err := goquery.NewDocumentFromReader(Resp.Body)
	if err != nil {
		return
	}
	// 获取所属采集类号
	//class := gjson.Get(html, "og:novel:status")
	//log.Println(class)
	return
}

//func fieldContent(collect *models.McCollect, bookUrl string) (err error) {
//	var count int64
//	global.DB.Model(models.McBook{}).Where("source_url", bookUrl).Count(&count)
//	if count > 0 {
//		err = fmt.Errorf("%v", "该小说已存在")
//		return
//	}
//
//	// GET请求
//	html, err := curlGet(bookUrl)
//	if err != nil {
//		fmt.Println("GET请求失败:", err)
//	} else {
//		fmt.Println("GET请求结果:", html)
//	}
//	//log.Println(html)
//	//[{"url":"http://www.biquw.com/xs/quanbu-default-0-0-0-0-0-0-[内容].html","type":"1","param":["1","6678","1",0]}]
//	dd := make(map[string]interface{})
//	//网站标题
//	//网址
//	//编码
//	//网址补全
//	//倒叙采集
//	//图片本地化
//	//采集小说列表页
//	//	采集范围正则
//	//	采集列表页书名和链接正则
//	//	采集列表分页规则 (json 网址类型（序列网址 多网址 单网址）采集列表地址 从第几页到第几页 每次增加几页 正序倒叙)
//	//
//	//采集小说详情页
//	//	获取详情页分类正则
//	//	分类转换规则 对应关系
//	//	获取小说名称正则
//	//	获取作者正则
//	//	获取连载状态
//	//	获取小说图片
//	//	小说目录正则
//
//
//	//[
//	//{"target":"玄幻小说","local":"18"},
//	//{"target":"仙侠小说","local":"19"},
//	//{"target":"都市小说","local":"21"},
//	//{"target":"军史小说","local":"20"},
//	//{"target":"网游小说","local":"34"},
//	//{"target":"科幻小说","local":"22"},
//	//{"target":"灵异小说","local":"23"},
//	//{"target":"言情小说","local":"26"},
//	//{"target":"其他小说","local":"35"}
//	//]
//	列表页盒子(获取区间范围)：<div class="sitebox">[内容]<div id="pages">
//	列表页书名和链接正则 <dt><a href="[内容1]" target="_blank">
//
//	小说目录：{"title":"章节页","page":"default","chapter":"1","section":"<div class=\"book_list\">[内容]<div class=\"cr\"></div>","url_rule":"<li><a href=\"[内容1]\">[章节标题]</a></li>","url_merge":""}]
//	小说详情页分类名称：
//	"category":{"field":"category","source":"default","rule":"<meta property=\"og:novel:category\" content=\"[内容1]\" \/>","merge":"","strip":""},
//	小说详情页标题：
//	"title":{"field":"title","source":"default","rule":"<meta property=\"og:novel:book_name\" content=\"[内容1]\" \/>","merge":"","strip":"","replace":""}
//	小说详情页作者：
//	"author":{"field":"author","source":"default","rule":"<meta property=\"og:novel:author\" content=\"[内容1]\" \/>","merge":"","strip":"","replace":""}
//	小说连载状态：
//	"serialize":{"field":"serialize","source":"default","rule":"<meta property=\"og:novel:status\" content=\"[内容1]\" \/>","merge":"","serial":"连载中","over":"完本","strip":"","replace":""}
//	小说图片：
//	"pic":{"field":"pic","source":"default","rule":"<meta property=\"og:image\" content=\"[内容1]\"\/>","merge":"","strip":"","replace":""}
//	章节标题：
//	"chapter_title":{"field":"chapter_title","source":"0","rule":"<h1>[内容1]<\/h1>","merge":"","strip":"","replace":""}
//	标签
//	"tag":{"field":"tag","source":"default","rule":"<meta property=\"og:title\" content=\"[内容1]\" \/>","merge":"","strip":"","replace":""},
//	章节内容
//	"chapter_content":{"field":"chapter_content","source":"0","rule":"<div id=\"htmlContent\" class=\"contentbox clear\">[内容1]<\/div>","merge":"","strip":"a,iframe,form"}
//	替换
//"replace":"[
//{"find":"恋上你看书网","replaces":"狂雨小说网"},
//{"find":"http://www.qidian.com/Book/1644509.aspx","replaces":""},
//{"find":"喜欢的朋友可以加一下QQ群：VIP群：560137817普通群：298412581","replaces":""},
//{"find":"微信公众号：水冷酒家","replaces":""},
//{"find":"微信：shuilengjiujia","replaces":""},
//{"find":"www.biquw.com","replaces":""},
//{"find":"关注九灯微信：jiudengheshan","replaces":""},
//{"find":"【读者群：487821318（龙武盟）】","replaces":""},
//{"find":"聊天加群：750610765","replaces":""},
//{"find":"【粉丝群已开放，群号：27414o891】","replaces":""},
//{"find":"读者交流++VIP书友群：450416188（需全订），普通书友群：392767347（非全订）","replaces":""},
//{"find":"寡头书友群：67227487（满）    42305786    35741501    40247521    87686729（新）","replaces":""}]"
//{"find":"read3;<更新更快就在笔趣网www.biquw.com>","replaces":""},
//{"find":"作者大眼猫神说：订阅加v群：560137817，鲜花榜奖励每月4号准时群内全部发放！凭订阅截图入群！老猫欢迎亲们进入！普通群：298412581","replaces":""},
//{"find":"喜欢的朋友可以加一下QQ群：VIP群：560137817普通群：298412581","replaces":""},
//{"find":"http://www.qidian.com/Book/1644509.aspx","replaces":""},
//{"find":"<更新更快就在笔趣网www.biquw.com>","replaces":""},
//{"find":"http://www.biquw.com","replaces":""},
//{"find":"Ps:书友们，我是孤单地飞，推荐一款免费小说App，支持小说下载、听书、零广告、多种阅读模式。请您关注微信公众号：dazhuzaiyuedu（长按三秒复制）书友们快关注起来吧！","replaces":""},
//{"find":"www.biquw.com","replaces":""},
//{"find":"笔趣网","replaces":""},
//{"find":"请关注微信公众号在线看:meinvxuan1(长按三秒复制)!!","replaces":""},
//{"find":"Ps:书友们，我是青莲剑仙，推荐一款免费小说App，支持小说下载、听书、零广告、多种阅读模式。请您关注微信公众号：dazhuzaiyuedu（长按三秒复制）书友们快关注起来吧！","replaces":""},
//{"find":"dazhuzaiyuedu","replaces":""}

//	//dd["field"] = "category"
//	//dd["source"] = "default"
//	//dd["rule"] = "<p><span>分类：[内容1]<span><span>大小："
//	//dd["merge"] = ""
//	//dd["strip"] = ""
//
//	//dd["field"] = "serialize"
//	//dd["source"] = "default"
//	//dd["rule"] = "<meta property=\"og:novel:status\" content=\"[内容1]\" />"
//	//dd["merge"] = ""
//	//dd["serial"] = "连载中"
//	//dd["over"] = "完本"
//	//dd["strip"] = ""
//	//dd["replace"] = ""
//
//	dd["field"] = "pic"
//	dd["source"] = "default"
//	dd["rule"] = "<meta property=\"og:image\" content=\"[内容1]\"/>"
//	dd["merge"] = ""
//	dd["strip"] = ""
//	dd["replace"] = ""
//1	www.biquw.com 笔趣网	auto	novel	1	0	1	[{"url":"http://www.biquw.com/xs/quanbu-default-0-0-0-0-0-0-[内容].html","type":"1","param":["1","6678","1",0]}]	<div class="sitebox">[内容]<div id="pages">	<dt><a href="[内容1]" target="_blank">				[{"title":"章节页","page":"default","chapter":"1","section":"<div class=\"book_list\">[内容]<div class=\"cr\"></div>","url_rule":"<li><a href=\"[内容1]\">[章节标题]</a></li>","url_merge":""}]	{"category":{"field":"category","source":"default","rule":"<meta property=\"og:novel:category\" content=\"[内容1]\" \/>","merge":"","strip":""},"title":{"field":"title","source":"default","rule":"<meta property=\"og:novel:book_name\" content=\"[内容1]\" \/>","merge":"","strip":"","replace":""},"author":{"field":"author","source":"default","rule":"<meta property=\"og:novel:author\" content=\"[内容1]\" \/>","merge":"","strip":"","replace":""},"serialize":{"field":"serialize","source":"default","rule":"<meta property=\"og:novel:status\" content=\"[内容1]\" \/>","merge":"","serial":"连载中","over":"完本","strip":"","replace":""},"pic":{"field":"pic","source":"default","rule":"<meta property=\"og:image\" content=\"[内容1]\"\/>","merge":"","strip":"","replace":""},"content":{"field":"content","source":"default","rule":"<\/strong><br \/>[内容1]<br\/>","merge":"","strip":"","replace":"[{\"find\":\"恋上你看书网\",\"replaces\":\"狂雨小说网\"},{\"find\":\"http:\/\/www.qidian.com\/Book\/1644509.aspx\",\"replaces\":\"\"},{\"find\":\"喜欢的朋友可以加一下QQ群：VIP群：560137817普通群：298412581\",\"replaces\":\"\"},{\"find\":\"微信公众号：水冷酒家\",\"replaces\":\"\"},{\"find\":\"微信：shuilengjiujia\",\"replaces\":\"\"},{\"find\":\"www.biquw.com\",\"replaces\":\"\"},{\"find\":\"关注九灯微信：jiudengheshan\",\"replaces\":\"\"},{\"find\":\"【读者群：487821318（龙武盟）】\",\"replaces\":\"\"},{\"find\":\"聊天加群：750610765\",\"replaces\":\"\"},{\"find\":\"【粉丝群已开放，群号：27414o891】\",\"replaces\":\"\"},{\"find\":\"读者交流++VIP书友群：450416188（需全订），普通书友群：392767347（非全订）\",\"replaces\":\"\"},{\"find\":\"寡头书友群：67227487（满）    42305786    35741501    40247521    87686729（新）\",\"replaces\":\"\"}]"},"tag":{"field":"tag","source":"default","rule":"<meta property=\"og:title\" content=\"[内容1]\" \/>","merge":"","strip":"","replace":""},"chapter_title":{"field":"chapter_title","source":"0","rule":"<h1>[内容1]<\/h1>","merge":"","strip":"","replace":""},"chapter_content":{"field":"chapter_content","source":"0","rule":"<div id=\"htmlContent\" class=\"contentbox clear\">[内容1]<\/div>","merge":"","strip":"a,iframe,form","replace":"[{\"find\":\"read3;<更新更快就在笔趣网www.biquw.com>\",\"replaces\":\"\"},{\"find\":\"作者大眼猫神说：订阅加v群：560137817，鲜花榜奖励每月4号准时群内全部发放！凭订阅截图入群！老猫欢迎亲们进入！普通群：298412581\",\"replaces\":\"\"},{\"find\":\"喜欢的朋友可以加一下QQ群：VIP群：560137817普通群：298412581\",\"replaces\":\"\"},{\"find\":\"http:\/\/www.qidian.com\/Book\/1644509.aspx\",\"replaces\":\"\"},{\"find\":\"<更新更快就在笔趣网www.biquw.com>\",\"replaces\":\"\"},{\"find\":\"http:\/\/www.biquw.com\",\"replaces\":\"\"},{\"find\":\"Ps:书友们，我是孤单地飞，推荐一款免费小说App，支持小说下载、听书、零广告、多种阅读模式。请您关注微信公众号：dazhuzaiyuedu（长按三秒复制）书友们快关注起来吧！\",\"replaces\":\"\"},{\"find\":\"www.biquw.com\",\"replaces\":\"\"},{\"find\":\"笔趣网\",\"replaces\":\"\"},{\"find\":\"请关注微信公众号在线看:meinvxuan1(长按三秒复制)!!\",\"replaces\":\"\"},{\"find\":\"Ps:书友们，我是青莲剑仙，推荐一款免费小说App，支持小说下载、听书、零广告、多种阅读模式。请您关注微信公众号：dazhuzaiyuedu（长按三秒复制）书友们快关注起来吧！\",\"replaces\":\"\"},{\"find\":\"dazhuzaiyuedu\",\"replaces\":\"\"}]"}}	0	0	[{"target":"玄幻小说","local":"18"},{"target":"仙侠小说","local":"19"},{"target":"都市小说","local":"21"},{"target":"军史小说","local":"20"},{"target":"网游小说","local":"34"},{"target":"科幻小说","local":"22"},{"target":"灵异小说","local":"23"},{"target":"言情小说","local":"26"},{"target":"其他小说","local":"35"}]
//	categoryMB := FieldRule(dd, html, false)
//	log.Println(categoryMB)
//	return
//}

func FieldRule(fieldParams map[string]interface{}, html string, isLoop bool) interface{} {
	fieldParams["rule"] = ConvertSignMatch(fieldParams["rule"].(string))
	fieldParams["merge"] = SetMergeDefault(fieldParams["rule"].(string), fieldParams["merge"].(string))

	if fieldParams["chapter"] != nil && fieldParams["chapter"].(bool) {
		fieldParams["rule"] = strings.Replace(fieldParams["rule"].(string), "[章节标题]", "(?P<title>[\\s\\S]*?)", -1)
	}

	signMatch := fieldParams["rule"].(string)
	signMatchRegex := regexp.MustCompile(signMatch)
	matchSigns := signMatchRegex.FindAllStringSubmatch(fieldParams["merge"].(string), -1)

	if len(matchSigns) > 0 {
		if isLoop {
			matchConts := regexp.MustCompile(fieldParams["rule"].(string)).FindAllStringSubmatch(html, -1)
			curI := 0
			val := make([]interface{}, 0)

			for _, matchCont := range matchConts {
				curI++
				reMatch := make(map[string]string)

				for i, submatch := range signMatchRegex.FindAllStringSubmatch(fieldParams["merge"].(string), -1) {
					reMatch[submatch[1]] = matchCont[i+1]
				}

				contVal := fieldParams["merge"].(string)
				for k, v := range reMatch {
					contVal = strings.Replace(contVal, "{"+k+"}", v, -1)
				}

				if fieldParams["strip"] != nil {
					if strings.Contains(fieldParams["strip"].(string), "all") {
						contVal = StripTagsContent(contVal, "style,script,object")
						contVal = stripTags(contVal)
					} else {
						contVal = StripTagsContent(contVal, fieldParams["strip"].(string))
					}
				}

				if fieldParams["replace"] != nil {
					replaces := fieldParams["replace"].([]interface{})
					for _, replace := range replaces {
						find := replace.(map[string]interface{})["find"].(string)
						replacements := replace.(map[string]interface{})["replaces"].(string)
						contVal = strings.Replace(contVal, find, replacements, -1)
					}
				}

				//if isLoop {
				//	if fieldParams["chapter"] != nil && fieldParams["chapter"].(bool) {
				//		val = append(val, map[string]interface{}{"title": matchCont[1], "url": contVal})
				//	} else {
				//		val = append(val, strings.TrimSpace(contVal))
				//	}
				//} else {
				//	val = strings.TrimSpace(contVal)
				//}
			}

			return val
		} else {
			matchCont := regexp.MustCompile(fieldParams["rule"].(string)).FindStringSubmatch(html)
			if len(matchCont) > 0 {
				reMatch := make(map[string]string)

				for i, submatch := range signMatchRegex.FindAllStringSubmatch(fieldParams["merge"].(string), -1) {
					reMatch[submatch[1]] = matchCont[i+1]
				}

				contVal := fieldParams["merge"].(string)
				for k, v := range reMatch {
					contVal = strings.Replace(contVal, "{"+k+"}", v, -1)
				}

				if fieldParams["strip"] != nil {
					if strings.Contains(fieldParams["strip"].(string), "all") {
						contVal = StripTagsContent(contVal, "style,script,object")
						contVal = stripTags(contVal)
					} else {
						contVal = StripTagsContent(contVal, fieldParams["strip"].(string))
					}
				}

				if fieldParams["replace"] != nil {
					replaces := fieldParams["replace"].([]interface{})
					for _, replace := range replaces {
						find := replace.(map[string]interface{})["find"].(string)
						replacements := replace.(map[string]interface{})["replaces"].(string)
						contVal = strings.Replace(contVal, find, replacements, -1)
					}
				}

				val := strings.TrimSpace(contVal)
				return val
			}
		}
	}

	return nil
}

func ConvertSignMatch(rule string) string {
	// 实现自定义的转换逻辑
	// ...
	return rule
}

func SetMergeDefault(rule string, merge string) string {
	// 实现自定义的设置逻辑
	// ...
	return merge
}

func StripTagsContent(contVal string, strip string) string {
	// 实现自定义的去除标签逻辑
	// ...
	return contVal
}

func stripTags(contVal string) string {
	// 实现自定义的去除标签逻辑
	// ...
	return contVal
}

func GetHTML(url string, encode string, urlComplete bool) string {

	client := &http.Client{
		Timeout: time.Second * 60,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	response, err := client.Do(request)
	if err != nil {
		return ""
	}
	defer response.Body.Close()

	html, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}

	if urlComplete {
	}

	return autoConvertToUtf8(string(html), encode)
}

func urlComplete(html string, baseURL string) string {
	html = regexp.MustCompile(`(?<=\bhref\=[\'\"])([^\'\"]*)(?=[\'\"])`).ReplaceAllStringFunc(html, func(match string) string {
		return createURL(match, baseURL)
	})
	html = regexp.MustCompile(`(?<=\bsrc\=[\'\"])([^\'\"]*)(?=[\'\"])`).ReplaceAllStringFunc(html, func(match string) string {
		return createURL(match, baseURL)
	})
	return html
}

func createURL(match string, baseURL string) string {
	if strings.HasPrefix(match, "http://") || strings.HasPrefix(match, "https://") {
		return match
	}
	if strings.HasPrefix(match, "/") {
		return baseURL + match
	}
	lastSlashIndex := strings.LastIndex(baseURL, "/")
	if lastSlashIndex == -1 {
		return baseURL + "/" + match
	}
	baseURL = baseURL[:lastSlashIndex]
	return baseURL + "/" + match
}

func autoConvertToUtf8(str string, encode string) string {
	if encode == "" || encode == "auto" {
		encodings := []string{"ASCII", "UTF-8", "GB2312", "GBK", "BIG5"}
		for _, encoding := range encodings {
			if strings.EqualFold(encoding, "UTF-8") {
				continue
			}
			//if strings.ToLower(encodings) != 0 {
			//	str = iconv.ConvertString(str, encode, "utf-8//IGNORE")
			//}
			return str
		}
	}
	return str
}

func curlGet(reqUrl string) (string, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func curlPost(reqUrl string, postData map[string]string, timeout time.Duration, header string, proxy string) (string, error) {
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, value := range postData {
		_ = writer.WriteField(key, value)
	}
	err := writer.Close()
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", reqUrl, body)
	if err != nil {
		return "", err
	}

	if header != "" {
		request.Header.Set("Content-Type", header)
	}

	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return "", err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}

//func getList() []gjson.Result {
//	// fmt.Println("采集资源站“", here.name, "”，第", here.pg, "页")
//	c := newHttpHandle()
//	res, err := c.Get(here.url + "?ac=list&pg=" + strconv.Itoa(pgCount-here.pg))
//	if err != nil {
//		// panic("采集资源站“" + here.name + "“获取采集页数失败")
//		util.Logger.Panicf("getter get the resource station called %s, getting page failed, err is %s", here.name, err)
//	}
//	defer res.Body.Close()
//	body, _ := io.ReadAll(res.Body)
//	list := gjson.Get(string(body), "list.#.vod_id").Array()
//	return list
//}

// 采集content
//func getContent(id int) {
//	res, err := c.Get(here.url + "?ac=detail&ids=" + strconv.Itoa(id))
//	if err != nil {
//		util.Logger.Panicf("getter get content failed, err is %s", err)
//		// panic后通过外部的recover来重新获取json
//	}
//	defer res.Body.Close()
//	// 获取body
//	body, _ := io.ReadAll(res.Body)
//
//	// 获取所属采集类号
//	class := int(gjson.Get(string(body), "list.0.type_id").Value().(float64))
//
//	if !db.JudgeClass(here.id, uint(class)) {
//		return
//	}
//
//	// 获取影片名
//	name := gjson.Get(string(body), "list.0.vod_name").Value().(string)
//
//	// 获取图片封面地址
//	pic := gjson.Get(string(body), "list.0.vod_pic").Value().(string)
//	pic = urlHandle(pic)
//
//	// 获取主演列表
//	actor := ""
//	actor_val := gjson.Get(string(body), "list.0.vod_actor").Value()
//	if actor_val != nil {
//		actor = actor_val.(string)
//	}
//
//	// 获取导演
//	director := ""
//	director_val := gjson.Get(string(body), "list.0.vod_director").Value()
//	if director_val != nil {
//		director = director_val.(string)
//	}
//
//	// 获取时长
//	duration := ""
//	duration_val := gjson.Get(string(body), "list.0.vod_duration").Value()
//	if duration_val != nil {
//		duration = duration_val.(string)
//	}
//
//	// 获取简介
//	description := ""
//	description_val := gjson.Get(string(body), "list.0.vod_content").Value()
//	if description_val != nil {
//		description = description_val.(string)
//	}
//	description = desHandle(description) // 净化功能
//
//	// 获取播放链接
//	url := gjson.Get(string(body), "list.0.vod_play_url").Value().(string)
//	url = urlHandle(url)
//
//	// 获取属于的source
//	belong := here.id
//	util.Logger.Infof("collect resource station called %s, get a film called %s", here.name, name)
//	err = db.AddContent(id, name, pic, actor, director, duration, description, url, class, belong)
//	if err != nil {
//		util.Logger.Errorf("getter get content, store the data to database failed, err is %s", err)
//	}
//	// 每当获取完一条信息后就尝试休眠一秒
//	time.Sleep(1 * time.Second)
//}
