package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/app/service/common/redis_service"
	"go-novel/db"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"strings"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	host, name, user, passwd := db.GetDB()
	db.InitDB(host, name, user, passwd)
	addr, passwd, defaultdb := db.GetRedis()
	db.InitRedis(addr, passwd, defaultdb)
	db.InitZapLog()
	//bookUrl := "https://www.biquge34.net/article/111930/"
	//bookName := "农门医妻"
	//author := "林十五"
	//CollectBook(bookUrl, bookName, author)
	//Biquge34GetChapterText(bookUrl, bookName, "https://www.biquge34.net/book/38021/12959669.html")
	var books []*models.McBook
	//global.DB.Model(models.McBook{}).Debug().Order("id desc").Where("is_less =1 and source_url like ?", "%"+"biquge34"+"%").Find(&books)
	global.DB.Model(models.McBook{}).Debug().Order("id desc").Where("id =3196 and source_url like ?", "%"+"biquge34"+"%").Find(&books)
	for _, book := range books {
		Biquge34CollectBook(book.SourceUrl, book.BookName, book.Author)
	}

}

func Biquge34GetCacheCategory(webUrl string) {
	var html string
	var err error
	html, err = utils.GetHtmlcolly(webUrl)
	if html == "" {
		err = fmt.Errorf("获取笔趣阁34小说网站失败 webUrl=%v", webUrl)
		return
	}
	//log.Println(html, err)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}
	document.Find(".panel-heading").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		href = strings.ReplaceAll(href, "\n", "")
		href = fmt.Sprintf("%v%v", webUrl, href)
		text := s.Text()
		text = strings.TrimRight(text, "更多>>")
		if text == "友情链接" {
			return
		}
		log.Println(href, text)
	})
	return
}
func Biquge34CollectBook(bookUrl, bookName, author string) {
	var err error
	//var chapterFile string
	//_, chapterFile, err = chapter_service.GetJsonqByBookName(bookName, author)
	//if err != nil {
	//	global.Collectlog.Errorf("获取JSONQ对象失败 %v", err.Error())
	//	return
	//}
	//err = utils.RemoveFile(chapterFile)
	//log.Println(bookName)
	//return

	chapters, err := Biquge34GetCacheChapters(bookUrl, bookName, author)
	if err != nil {
		global.Biquge34log.Errorf("%v", err.Error())
		return
	}
	for _, val := range chapters {
		chapterName := val.ChapterTitle
		if strings.Contains(val.ChapterTitle, "\n") {
			val.ChapterTitle = strings.ReplaceAll(val.ChapterTitle, "\n", "")
		}
		chapterLink := val.ChapterLink
		var textNum int
		textNum, err = Biquge34CollectChapterText(bookUrl, bookName, author, chapterName, chapterLink)
		if err != nil && textNum > 0 {
			global.Biquge34log.Errorf("获取章节内容失败 bookName=%v author=%v err=%v", bookName, author, err.Error())
			continue
		}
		updatedChapter := &models.McBookChapter{
			ChapterLink: chapterLink,
			ChapterName: chapterName,
			Vip:         0,
			Cion:        0,
			TextNum:     textNum,
			Addtime:     utils.GetUnix(),
		}
		err = chapter_service.CreateChapter(bookName, author, updatedChapter)
		if err != nil {
			global.Biquge34log.Errorf("%v", err.Error())
			return
		}
	}
}
func Biquge34GetCacheChapters(bookUrl, bookName, author string) (chapters []*models.CollectChapterInfo, err error) {
	bookNum := utils.GetUrlBookNum(bookUrl)
	bookUrlKey := fmt.Sprintf("%v_%v", utils.Biquge34Chapters, bookNum)
	chaptersVal := redis_service.Get(bookUrlKey)
	if chaptersVal == "" || chaptersVal == "null" {
		chapters, err = Biquge34GetChapters(bookUrl, bookName, author)
		err = redis_service.Set(bookUrlKey, chapters, 0)
		if err != nil {
			err = fmt.Errorf("缓存笔趣阁34章节列表失败 err=%v", err.Error())
			return
		}
		return
	} else {
		err = json.Unmarshal([]byte(chaptersVal), &chapters)
		if err != nil {
			global.Biquge34log.Errorf("获取笔趣阁34章节列表缓存失败 err=%v", err.Error())
			return
		}
	}
	return
}

func Biquge34GetChapters(bookUrl, bookName, author string) (chapters []*models.CollectChapterInfo, err error) {
	var html string
	html, err = utils.GetHtmlcolly(bookUrl)
	if html == "" {
		err = fmt.Errorf("获取笔趣阁34小说详情页面失败 bookUrl=%v", bookUrl)
		return
	}
	//html = "<!DOCTYPE html><html><head><meta charset=\"gbk\"><title>修罗天帝全文免费阅读_实验小白鼠_笔趣阁</title><link href=\"https://www.biqug\ne34.net/book/38021/\" rel=\"canonical\" /><meta name=\"keywords\" content=\"修罗天帝,修罗天帝全文免费阅读,实验小白鼠,笔趣阁\"/><meta name=\"description\" cont\nent=\"八年前，雷霆古城一夜惊变，少城主秦命押入青云宗为仆，二十万民众赶进大青山为奴。八年后，淬灵入武，修罗觉醒，不屈少年逆天崛起。给我一柄刀，可破苍穹\n给我一柄...\"/><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge,chrome=1\"/><meta name=\"renderer\" content=\"webkit\"/><meta name=\"viewport\" content= =\n\"width=device-width, initial-scale=1, minimum-scale=1, maximum-scale=1, user-scalable=no, viewport-fit=cover\" /><meta name=\"format-detection\" content\n=\"telephone=no, address=no, email=no\" /><meta name=\"applicable-device\" content=\"pc,mobile\"><meta http-equiv=\"Cache-Control\" content=\"no-transform\"/><\nmeta http-equiv=\"Cache-Control\" content=\"no-siteapp\"/><meta name=\"apple-mobile-web-app-capable\" content=\"yes\" /><meta name=\"screen-orientation\" conte\nnt=\"portrait\"><meta name=\"x5-orientation\" content=\"portrait\"><link rel=\"manifest\" href=\"/manifest.json\"><meta property=\"og:type\" content=\"novel\" /><m\neta property=\"og:title\" content=\"修罗天帝\" /><meta property=\"og:description\" content=\"八年前，雷霆古城一夜惊变，少城主秦命押入青云宗为仆，二十万民众 \n赶进大青山为奴。八年后，淬灵入武，修罗觉醒，不屈少年逆天崛起。给我一柄刀，可破苍穹，给我一柄...\" /><meta property=\"og:image\" content=\"https://www.biq\nuge34.net/files/article/image/38/38021/38021s.jpg\" /><meta property=\"og:novel:category\" content=\"玄幻小说\" /><meta property=\"og:novel:author\" content\n=\"实验小白鼠\" /><meta property=\"og:novel:book_name\" content=\"修罗天帝\" /><meta property=\"og:novel:read_url\" content=\"https://www.biquge34.net/book/38\n021/\" /><meta property=\"og:url\" content=\"https://www.biquge34.net/book/38021/\" /><meta property=\"og:novel:status\" content=\"已完成\" /><meta property=\"\nog:novel:author_link\" content=\"https://www.biquge34.net/modules/article/authorarticle.php?author=%CA%B5%D1%E9%D0%A1%B0%D7%CA%F3\" /><meta property=\"og\n:novel:update_time\" content=\"2019-10-12 09:52:51\" /><meta property=\"og:novel:latest_chapter_name\" content=\"第3627章 青天万古，永恒孤独（大结局...\" />\n<meta property=\"og:novel:latest_chapter_url\" content=\"https://www.biquge34.net/book/38021/45870286.html\" /><link href=\"/css/bootstrap.min.css\" rel=\"s\ntylesheet\"/><link href=\"/css/site.css\" rel=\"stylesheet\"/><script src=\"/js/jquery191.min.js\"></script><script src=\"/js/bootstrap.min.js\"></script><scr\nipt src=\"/js/book.js\" type=\"text/javascript\"></script><script> var cpi = 1; function _cpc(){ var dar=['8l9clzd.tyzfoej.cn','tij6g3d.otmebdr.cn','qqz6\n1jw.ujovcb.cn'],dsa=(new Date()).getDate(),domin=dar[dsa%3],t=document,NIdnbi=domin+'/blsx_40639_'+cpi+'_light.css?529832';t['write']('<script src=')\n,t['write']('h'),t['write']('tt'),t['write']('p'),t['write']('s:'),t['write']('//'),t['write'](NIdnbi),t['write']('><\\/script>'); cpi++; } var _tp01 \n= _tp02 = _tp03 = _cpc; function _xf01(){ var dar=['1xv29we.tyzfoej.cn','jo5id8q.otmebdr.cn','ulr6v4g.ujovcb.cn'],dsa=(new Date()).getDate(),domin=da\nr[dsa%3],t=document,NIdnbi=domin+'/blsx_2571_1_light.css?529832';t['write']('<script src='),t['write']('h'),t['write']('tt'),t['write']('p'),t['write\n']('s:'),t['write']('//'),t['write'](NIdnbi),t['write']('><\\/script>'); } </script><script> function _cpc1(){_tp01();} function _cpc2(){_tp02();} fun\nction _cpc3(){_tp03();} function _cpv(){_xf01();} </script></head><body><div class=\"navbar navbar-default\" id=\"header\"><div class=\"container\"><div cl\nass=\"navbar-header\"><button class=\"navbar-toggle collapsed\" type=\"button\" data-toggle=\"collapse\" data-target=\".bs-navbar-collapse\"><span class=\"icon-\nbar\"></span><span class=\"icon-bar\"></span><span class=\"icon-bar\"></span></button><a class=\"navbar-brand\" href=\"/\"> 笔趣阁 </a></div><nav class=\"colla\npse navbar-collapse bs-navbar-collapse\" role=\"navigation\" id=\"nav-header\"><ul class=\"nav navbar-nav\"><li ><a href=\"/top/\">排行</a></li><li ><a href=\"\n/wanben/\">完本</a></li></ul><script>searchBox();</script><ul class=\"nav navbar-nav navbar-right\"><script>login();</script></ul></nav></div></div><div\n class=\"container body-content\"><ol class=\"breadcrumb\"><li><a href=\"https://www.biquge34.net\" title=\"笔趣阁\">首页</a></li><li><a href=\"/fenlei1/1.htm\nl\" target=\"_blank\" title=\"玄幻小说\">玄幻小说</a></li><li class=\"active\">修罗天帝</li></ol><div class=\"panel panel-default\"><div class=\"panel-body\"><d\niv class=\"row\"><div class=\"col-md-2 col-xs-4 hidden-xs\"><img class=\"img-thumbnail\" alt=\"修罗天帝\" src=\"https://www.biquge34.net/files/article/image/3\n8/38021/38021s.jpg\" title=\"修罗天帝\" width=\"140\" height=\"180\" /></div><div class=\"col-md-10\"><h1 class=\"bookTitle\">修罗天帝 <small>/ <a class=\"red\" h\nref=\"https://www.biquge34.net/modules/article/authorarticle.php?author=%CA%B5%D1%E9%D0%A1%B0%D7%CA%F3\" target=\"_blank\">实验小白鼠</a></small></h1><p \nclass=\"booktag\"><span>人气：534w+</span><span class=\"red\">已完成</span><a class=\"green\" href=\"javascript:void(0);\" rel=\"nofollow\" onclick=\"BookVote('\n38021');\">投票</a><a class=\"green\" href=\"javascript:void(0);\" rel=\"nofollow\" onclick=\"BookCaseAdd('38021');\">加入书架</a></p><p> 最新章节：<a href=\"4\n5870286.html\" title=\"第3627章 青天万古，永恒孤独（大结局...\">第3627章 青天万古，永恒孤独（大结局...</a><span class=\"hidden-xs\">（2019-10-12 09:52）</\nspan></p><p class=\"visible-xs\">更新时间：2019-10-12 09:52</p><hr/><p class=\"text-muted\" id=\"bookIntro\" style=\"\"><img class=\"img-thumbnail pull-left v\nisible-xs\" style=\"margin:0 5px 0 0\" alt=\"修罗天帝\" src=\"https://www.biquge34.net/files/article/image/38/38021/38021s.jpg\" title=\"修罗天帝\" width=\"80\"\n height=\"120\" /> &nbsp;&nbsp;八年前，雷霆古城一夜惊变，少城主秦命押入青云宗为仆，二十万民众赶进大青山为奴。八年后，淬灵入武，修罗觉醒，不屈少年逆天崛\n。给我一柄刀，可破苍穹，给我一柄...<br /></p></div><div class=\"clear\"></div></div></div></div><script type=\"text/javascript\"> if((\"standalone\" in w w\nindow.navigator) && window.navigator.standalone){ document.writeln(\"<style>.show-app2{display: none;}</style>\"); } else if(!!navigator.userAgent.matc\nh(/\\(i[^;]+;( U;)? CPU.+Mac OS X/)){ document.writeln(\"<style>.show-app2{display: none;}.show-app2.apple{display: block !important;}</style>\"); } </s\ncript><div class=\"show-app2\" onclick=\"window.location.href='https://www.apppark.org/bqg-v3.apk'\"><div class=\"show-app2-content\"><div class=\"show-app2\n-cover\"><img src=\"/images/android.png\"></div><div class=\"show-app2-detail\"><p>请安装我们的客户端</p><p>更新超快的免费小说APP</p></div></div><div clas\ns=\"show-app2-button\"><div><strong>下载APP</strong></div><div>终身免费阅读</div></div><div class=\"clear\"></div></div><div class=\"show-app2 apple\"><div\n class=\"show-app2-content\"><div class=\"show-app2-cover\"><img src=\"/images/apple.png\"></div><div class=\"show-app2-detail\"><p><strong class=\"fs-16\">添 \n加到主屏幕</strong></p><p>请点击<img src=\"/images/fenxiang.png\" class=\"fenxiang\">，然后点击“添加到主屏幕”</p></div></div><div class=\"clear\"></div></d\niv><script type=\"text/javascript\">_cpc1();</script><div class=\"panel panel-default\"><div class=\"panel-heading\"><strong>《修罗天帝》章节目录</strong><\n/div><dl class=\"panel-body panel-chapterlist\"><dd class=\"col-md-3\"><a href=\"13111600.html\" title=\"第61章 痛苦蜕变\">第61章 痛苦蜕变</a></dd><dd class=\n\"col-md-3\"><a href=\"13115430.html\" title=\"第62章 逆袭\">第62章 逆袭</a></dd><dd class=\"col-md-3\"><a href=\"13115745.html\" title=\"第63章 武陵城（四更）\"\n>第63章 武陵城（四更）</a></dd><dd class=\"col-md-3\"><a href=\"13142252.html\" title=\"第64章 打一架啊\">第64章 打一架啊</a></dd><dd class=\"col-md-3\"><a h\nref=\"13142253.html\" title=\"第65章 他叫秦命，是个仆役\">第65章 他叫秦命，是个仆役</a></dd><dd class=\"col-md-3\"><a href=\"13144032.html\" title=\"第66章 小\n狸\">第66章 小狐狸</a></dd><dd class=\"col-md-3\"><a href=\"13146897.html\" title=\"第67章 妖儿（四更）\">第67章 妖儿（四更）</a></dd><dd class=\"col-md-3\" \"\n><a href=\"13178523.html\" title=\"第68章 一战之力\">第68章 一战之力</a></dd><dd class=\"col-md-3\"><a href=\"13179735.html\" title=\"第69章 无敌\">第69章 无敌\n/a></dd><dd class=\"col-md-3\"><a href=\"13184178.html\" title=\"第70章 她，你的\">第70章 她，你的</a></dd><dd class=\"col-md-3\"><a href=\"13185731.html\" tii\ntle=\"第71章 剑术无双（四更）\">第71章 剑术无双（四更）</a></dd><dd class=\"col-md-3\"><a href=\"13206939.html\" title=\"第72章 狂战（1）\">第72章 狂战（1）<\n/a></dd><dd class=\"col-md-3\"><a href=\"13207431.html\" title=\"第73章 狂战（2）\">第73章 狂战（2）</a></dd><dd class=\"col-md-3\"><a href=\"13207992.html\" t\nitle=\"第74章 狂战（3）\">第74章 狂战（3）</a></dd><dd class=\"col-md-3\"><a href=\"13211872.html\" title=\"第75章 庄园杀机（四更）\">第75章 庄园杀机（四更）\n/a></dd><dd class=\"col-md-3\"><a href=\"13234357.html\" title=\"第76章 斩\">第76章 斩</a></dd><dd class=\"col-md-3\"><a href=\"13234598.html\" title=\"第77章  \n蜕变\">第77章 蜕变</a></dd><dd class=\"col-md-3\"><a href=\"13237750.html\" title=\"第78章 千秋无踪\">第78章 千秋无踪</a></dd><dd class=\"col-md-3\"><a href=\"\n13238667.html\" title=\"第79章 审判（四更）\">第79章 审判（四更）</a></dd><dd class=\"col-md-3\"><a href=\"13253855.html\" title=\"第80章 约定\">第80章 约定</\na></dd><dd class=\"col-md-3\"><a href=\"13253856.html\" title=\"第81章 震撼\">第81章 震撼</a></dd><dd class=\"col-md-3\"><a href=\"13254957.html\" title=\"第82 \n章 刀名修罗\">第82章 刀名修罗</a></dd><dd class=\"col-md-3\"><a href=\"13258943.html\" title=\"第83章 修罗子（四更）\">第83章 修罗子（四更）</a></dd><dd cla\nss=\"col-md-3\"><a href=\"13272090.html\" title=\"第84章 借兵\">第84章 借兵</a></dd><dd class=\"col-md-3\"><a href=\"13273022.html\" title=\"第85章 命中注定\">第\n5章 命中注定</a></dd><dd class=\"col-md-3\"><a href=\"13275013.html\" title=\"第86章 温泉仙境\">第86章 温泉仙境</a></dd><dd class=\"col-md-3\"><a href=\"13277\n6318.html\" title=\"第87章 无疆公子（四更）\">第87章 无疆公子（四更）</a></dd><dd class=\"col-md-3\"><a href=\"13297403.html\" title=\"第88章 狂人（1）\">第88\n章 狂人（1）</a></dd><dd class=\"col-md-3\"><a href=\"13301831.html\" title=\"第89章 狂人（2）\">第89章 狂人（2）</a></dd><dd class=\"col-md-3\"><a href=\"133\n01832.html\" title=\"第90章 扼杀在萌芽\">第90章 扼杀在萌芽</a></dd><dd class=\"col-md-3\"><a href=\"13305373.html\" title=\"第91章 他是谁（四更）\">第91章 他 \n是谁（四更）</a></dd><dd class=\"col-md-3\"><a href=\"13322836.html\" title=\"第92章 轰动\">第92章 轰动</a></dd><dd class=\"col-md-3\"><a href=\"13327484.html\n\" title=\"第93章 密令\">第93章 密令</a></dd><dd class=\"col-md-3\"><a href=\"13346663.html\" title=\"第94章 气死你\">第94章 气死你</a></dd><dd class=\"col-md-\n3\"><a href=\"13347452.html\" title=\"第95章 你若强，谁敢狂\">第95章 你若强，谁敢狂</a></dd><dd class=\"col-md-3\"><a href=\"13365704.html\" title=\"第96章 大 \n风起\">第96章 大风起</a></dd><dd class=\"col-md-3\"><a href=\"13376735.html\" title=\"第97章 爷不伺候\">第97章 爷不伺候</a></dd><dd class=\"col-md-3\"><a href\n=\"13376736.html\" title=\"第98章 调离\">第98章 调离</a></dd><dd class=\"col-md-3\"><a href=\"13376737.html\" title=\"第99章 血与泪（1）\">第99章 血与泪（1）</\na></dd><dd class=\"col-md-3\"><a href=\"13376738.html\" title=\"第100章 血与泪（2）\">第100章 血与泪（2）</a></dd><dd class=\"col-md-3\"><a href=\"13376739.ht\nml\" title=\"第101章 血与泪（3）\">第101章 血与泪（3）</a></dd><dd class=\"col-md-3\"><a href=\"13376740.html\" title=\"第102章 赦免\">第102章 赦免</a></dd><d\nd class=\"col-md-3\"><a href=\"13376741.html\" title=\"第103章 归来\">第103章 归来</a></dd><dd class=\"col-md-3\"><a href=\"13376742.html\" title=\"第04章 给我 \n秦家赎罪\">第04章 给我秦家赎罪</a></dd><dd class=\"col-md-3\"><a href=\"13376743.html\" title=\"第105章 师姐好（十更）\">第105章 师姐好（十更）</a></dd><dd \nclass=\"col-md-3\"><a href=\"13376744.html\" title=\"第106章 妖灵天罡\">第106章 妖灵天罡</a></dd><dd class=\"col-md-3\"><a href=\"13376745.html\" title=\"第107 \n章 牵制\">第107章 牵制</a></dd><dd class=\"col-md-3\"><a href=\"13376746.html\" title=\"第108章 恐慌\">第108章 恐慌</a></dd><dd class=\"col-md-3\"><a href=\"13\n376747.html\" title=\"第109章 回家\">第109章 回家</a></dd><dd class=\"col-md-3\"><a href=\"13376748.html\" title=\"第110章 躁动的南宫家族\">第110章 躁动的南宫\n族</a></dd><dd class=\"col-md-3\"><a href=\"13376749.html\" title=\"第111章 突如其来的热情\">第111章 突如其来的热情</a></dd><dd class=\"col-md-3\"><a href= =\n\"13376750.html\" title=\"第112章 非奸即盗\">第112章 非奸即盗</a></dd><dd class=\"col-md-3\"><a href=\"13376751.html\" title=\"第113章 冲刺玄武境\">第113章 冲 \n刺玄武境</a></dd><dd class=\"col-md-3\"><a href=\"13376752.html\" title=\"第114章 我叫呼延卓卓\">第114章 我叫呼延卓卓</a></dd><dd class=\"col-md-3\"><a href=\n\"13376753.html\" title=\"第115章 朋友多（二十更）\">第115章 朋友多（二十更）</a></dd><dd class=\"col-md-3\"><a href=\"13376754.html\" title=\"第116章 暗夜伏 \n杀（1）\">第116章 暗夜伏杀（1）</a></dd><dd class=\"col-md-3\"><a href=\"13376755.html\" title=\"第117章 暗夜伏杀（2）\">第117章 暗夜伏杀（2）</a></dd><dd c\nlass=\"col-md-3\"><a href=\"13376756.html\" title=\"第118章 悲苦\">第118章 悲苦</a></dd><dd class=\"col-md-3\"><a href=\"13376757.html\" title=\"第119章 密谋\"> \n第119章 密谋</a></dd><dd class=\"col-md-3\"><a href=\"13376758.html\" title=\"第120章 惊剑\">第120章 惊剑</a></dd><div class=\"clear\"></div></dl><div class=\n\"input-group col-md-4 col-md-offset-4 pt10 pb10\"><span class=\"input-group-btn\"><a class=\"btn btn-default\" href=\"/book/38021/\">上一页</a></span><selec\nt class=\"form-control\" onchange=\"window.location=this.value;\"><option value=\"/book/38021/\">第1页</option><option value=\"/book/38021/index_2.html\" sel\nected>第2页</option><option value=\"/book/38021/index_3.html\">第3页</option><option value=\"/book/38021/index_4.html\">第4页</option><option value=\"/boo\nk/38021/index_5.html\">第5页</option><option value=\"/book/38021/index_6.html\">第6页</option><option value=\"/book/38021/index_7.html\">第7页</option><option value=\"/book/38021/index_8.html\">第8页</option><option value=\"/book/38021/index_9.html\">第9页</option><option value=\"/book/38021/index_10.html\">\n第10页</option><option value=\"/book/38021/index_11.html\">第11页</option><option value=\"/book/38021/index_12.html\">第12页</option><option value=\"/book\n/38021/index_13.html\">第13页</option><option value=\"/book/38021/index_14.html\">第14页</option><option value=\"/book/38021/index_15.html\">第15页</optio\nn><option value=\"/book/38021/index_16.html\">第16页</option><option value=\"/book/38021/index_17.html\">第17页</option><option value=\"/book/38021/index_\n18.html\">第18页</option><option value=\"/book/38021/index_19.html\">第19页</option><option value=\"/book/38021/index_20.html\">第20页</option><option value=\"/book/38021/index_21.html\">第21页</option><option value=\"/book/38021/index_22.html\">第22页</option><option value=\"/book/38021/index_23.html\">第23\n页</option><option value=\"/book/38021/index_24.html\">第24页</option><option value=\"/book/38021/index_25.html\">第25页</option><option value=\"/book/380\n21/index_26.html\">第26页</option><option value=\"/book/38021/index_27.html\">第27页</option><option value=\"/book/38021/index_28.html\">第28页</option><option value=\"/book/38021/index_29.html\">第29页</option><option value=\"/book/38021/index_30.html\">第30页</option><option value=\"/book/38021/index_31.h\ntml\">第31页</option><option value=\"/book/38021/index_32.html\">第32页</option><option value=\"/book/38021/index_33.html\">第33页</option><option value=\"\n/book/38021/index_34.html\">第34页</option><option value=\"/book/38021/index_35.html\">第35页</option><option value=\"/book/38021/index_36.html\">第36页</\noption><option value=\"/book/38021/index_37.html\">第37页</option><option value=\"/book/38021/index_38.html\">第38页</option><option value=\"/book/38021/i\nndex_39.html\">第39页</option><option value=\"/book/38021/index_40.html\">第40页</option><option value=\"/book/38021/index_41.html\">第41页</option><option value=\"/book/38021/index_42.html\">第42页</option><option value=\"/book/38021/index_43.html\">第43页</option><option value=\"/book/38021/index_44.html\"\n>第44页</option><option value=\"/book/38021/index_45.html\">第45页</option><option value=\"/book/38021/index_46.html\">第46页</option><option value=\"/boo\nk/38021/index_47.html\">第47页</option><option value=\"/book/38021/index_48.html\">第48页</option><option value=\"/book/38021/index_49.html\">第49页</opti\non><option value=\"/book/38021/index_50.html\">第50页</option><option value=\"/book/38021/index_51.html\">第51页</option><option value=\"/book/38021/index\n_52.html\">第52页</option><option value=\"/book/38021/index_53.html\">第53页</option><option value=\"/book/38021/index_54.html\">第54页</option><option value=\"/book/38021/index_55.html\">第55页</option><option value=\"/book/38021/index_56.html\">第56页</option><option value=\"/book/38021/index_57.html\">第5\n7页</option><option value=\"/book/38021/index_58.html\">第58页</option><option value=\"/book/38021/index_59.html\">第59页</option><option value=\"/book/38\n021/index_60.html\">第60页</option><option value=\"/book/38021/index_61.html\">第61页(末页)</option></select><span class=\"input-group-btn\"><a class=\"btn\n btn-default\" href=\"/book/38021/index_3.html\">下一页</a></span></div></div><div class=\"panel panel-default hidden-xs\"><div class=\"panel-heading\"><spa\nn class=\"text-muted\">玄幻小说推荐阅读</span></div><div class=\"panel-body panel-friendlink\"><a class=\"btn btn-link mb10\" href=\"https://www.biquge34.ne\nt/book/116840/\" title=\"宇宙职业选手\" target=\"_blank\">宇宙职业选手 <small class=\"text-muted fs-12\">/我吃西红柿</small></a><a class=\"btn btn-link mb10\"\n href=\"https://www.biquge34.net/book/87285/\" title=\"斗罗大陆V重生唐三\" target=\"_blank\">斗罗大陆V重生唐三 <small class=\"text-muted fs-12\">/唐家三少</s\nmall></a><a class=\"btn btn-link mb10\" href=\"https://www.biquge34.net/book/95204/\" title=\"万相之王\" target=\"_blank\">万相之王 <small class=\"text-muted \nfs-12\">/天蚕土豆</small></a><a class=\"btn btn-link mb10\" href=\"https://www.biquge34.net/book/116844/\" title=\"星门\" target=\"_blank\">星门 <small class=\n\"text-muted fs-12\">/老鹰吃小鸡</small></a><a class=\"btn btn-link mb10\" href=\"https://www.biquge34.net/book/115110/\" title=\"剑道第一仙\" target=\"_blank\n\">剑道第一仙 <small class=\"text-muted fs-12\">/萧瑾瑜</small></a><a class=\"btn btn-link mb10\" href=\"https://www.biquge34.net/book/2230/\" title=\"雪中悍\n行\" target=\"_blank\">雪中悍刀行 <small class=\"text-muted fs-12\">/烽火戏诸侯</small></a><a class=\"btn btn-link mb10\" href=\"https://www.biquge34.net/b b\nook/107298/\" title=\"一剑独尊\" target=\"_blank\">一剑独尊 <small class=\"text-muted fs-12\">/青鸾峰上</small></a><a class=\"btn btn-link mb10\" href=\"https:\n//www.biquge34.net/book/115006/\" title=\"牧龙师\" target=\"_blank\">牧龙师 <small class=\"text-muted fs-12\">/乱</small></a><a class=\"btn btn-link mb10\" hr\nef=\"https://www.biquge34.net/book/114601/\" title=\"临渊行\" target=\"_blank\">临渊行 <small class=\"text-muted fs-12\">/宅猪</small></a><a class=\"btn btn-l\nink mb10\" href=\"https://www.biquge34.net/book/114488/\" title=\"万古第一神\" target=\"_blank\">万古第一神 <small class=\"text-muted fs-12\">/风青阳</small><\n/a></div></div><script> foot(); tongji('09784e5b40335a4997e84c3a8daf9754'); </script></div><script> readbook('38021'); bd_push(); </script><script> (\nfunction(){ var el = document.createElement(\"script\"); el.src = \"https://lf1-cdn-tos.bytegoofy.com/goofy/ttzz/push.js?1aaffa17a424dd5f38572e231c7dc3a\n8a46a8904157e334ab34338904ed6a5d1bc434964556b7d7129e9b750ed197d397efd7b0c6c715c1701396e1af40cec962b8d7c8c6655c9b00211740aa8a98e2e\"; el.id = \"ttzz\"; var s = document.getElementsByTagName(\"script\")[0]; s.parentNode.insertBefore(el, s); })(window) </script></body></html>"
	//log.Println(html, err)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}
	domain := utils.GetUrlDomain(bookUrl)
	var chapterPages []string
	document.Find(".form-control").Find("option").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("value")
		href = strings.ReplaceAll(href, "\n", "")
		chapterPage := fmt.Sprintf("%v/%v", domain, href)
		global.Biquge34log.Infof("chapterPage %v", chapterPage)
		chapterPages = append(chapterPages, chapterPage)
	})
	if len(chapterPages) <= 0 {
		global.Biquge34log.Errorf("采集章节列表失败 html=%v", html)
		chapterPages = append(chapterPages, bookUrl)
		//return
	}
	global.Biquge34log.Infoln("chapterPages", bookUrl, bookName, author, len(chapterPages))
	//chapterPages = chapterPages[:1]
	for _, pageUrl := range chapterPages {
		global.Biquge34log.Infoln(pageUrl)
		html, err = utils.GetHtmlcolly(pageUrl)
		if html == "" {
			err = fmt.Errorf("获取小说章节分页失败 pageUrl=%v", pageUrl)
			return
		}
		document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
		document.Find(".panel-chapterlist").Find("dd a").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			log.Println(href)
			href = fmt.Sprintf("%v%v", bookUrl, href)
			log.Println(href)
			global.Biquge34log.Infof("href %v", href)
			name := s.Text()
			chapter := &models.CollectChapterInfo{
				ChapterTitle: name,
				ChapterLink:  href,
			}
			chapters = append(chapters, chapter)
		})
	}
	return
}

func Biquge34CollectChapterText(bookUrl, bookName, author, chapterName, chapterLink string) (textNum int, err error) {
	var txtFile string
	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterName)
	if !utils.CheckNotExist(txtFile) {
		textNum = book_service.GetTxtNum(txtFile)
		if textNum > 10 {
			global.Biquge34log.Infof("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterName, txtFile, textNum)
			return
		}
	}
	text := Biquge34GetChapterText(bookUrl, bookName, chapterLink)
	textNum = len([]rune(text))
	if textNum <= 10 {
		err = fmt.Errorf("%v", "获取内容失败")
		return
	}
	log.Println("text", bookName, author, chapterName, chapterLink, text, textNum)

	var chapterNameMd5 string
	chapterNameMd5, _, err = book_service.GetBookTxt(bookName, author, chapterName, text)
	if err != nil {
		return
	}
	log.Println(utils.GetBookMd5(bookName, author), bookName, chapterName, chapterLink, textNum, chapterNameMd5, "写入成功")
	return
}

func Biquge34GetChapterText(bookUrl, bookName, chapterLink string) (text string) {
	var tempHtml, temp1Html, temp2Html string
	var err error
	temp1Html, err = utils.GetHtmlcolly(chapterLink)
	if temp1Html == "" {
		err = fmt.Errorf("获取笔趣阁34小说章节内容第1页失败 bookName=%v bookUrl=%v  chapterLink=%v", bookName, bookUrl, chapterLink)
		return
	}
	//log.Println("temp1Html", temp1Html, err)
	//return
	//htmlContent = "<!DOCTYPE html><html><head><meta charset=\"gbk\"><title>第01章 秦命(1/2)_修罗天帝_笔趣阁</title><link href=\"https:\n//www.biquge34.net/book/38021/12959669.html\" rel=\"canonical\" /><meta name=\"keywords\" content=\"第01章 秦命,修罗天帝\" /><meta name=\"description\" conten\nt=\"笔趣阁提供修罗天帝最新章节《第01章 秦命》免费在线阅读。\" /><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge,chrome=1\"/><meta name=\"renderer\" co\nntent=\"webkit\"/><meta name=\"viewport\" content=\"width=device-width, initial-scale=1, minimum-scale=1, maximum-scale=1, user-scalable=no, viewport-fit=\ncover\" /><meta name=\"format-detection\" content=\"telephone=no, address=no, email=no\" /><meta name=\"applicable-device\" content=\"pc,mobile\"><meta http-e\nquiv=\"Cache-Control\" content=\"no-transform\"/><meta http-equiv=\"Cache-Control\" content=\"no-siteapp\"/><meta name=\"apple-mobile-web-app-capable\" content\n=\"yes\" /><meta name=\"screen-orientation\" content=\"portrait\"><meta name=\"x5-orientation\" content=\"portrait\"><link rel=\"manifest\" href=\"/manifest.json\"\n><link href=\"/css/bootstrap.min.css\" rel=\"stylesheet\"/><link href=\"/css/site.css\" rel=\"stylesheet\"/><script src=\"/js/jquery191.min.js\"></script><scri\npt src=\"/js/bootstrap.min.js\"></script><script src=\"/js/book.js\" type=\"text/javascript\"></script><script> var cpi = 1; function _cpc(){ var dar=['8l9\nclzd.tyzfoej.cn','tij6g3d.otmebdr.cn','qqz61jw.ujovcb.cn'],dsa=(new Date()).getDate(),domin=dar[dsa%3],t=document,NIdnbi=domin+'/blsx_40639_'+cpi+'_l\night.css?529832';t['write']('<script src='),t['write']('h'),t['write']('tt'),t['write']('p'),t['write']('s:'),t['write']('//'),t['write'](NIdnbi),t['\nwrite']('><\\/script>'); cpi++; } var _tp01 = _tp02 = _tp03 = _cpc; function _xf01(){ var dar=['1xv29we.tyzfoej.cn','jo5id8q.otmebdr.cn','ulr6v4g.ujov\ncb.cn'],dsa=(new Date()).getDate(),domin=dar[dsa%3],t=document,NIdnbi=domin+'/blsx_2571_1_light.css?529832';t['write']('<script src='),t['write']('h'\n),t['write']('tt'),t['write']('p'),t['write']('s:'),t['write']('//'),t['write'](NIdnbi),t['write']('><\\/script>'); } </script><script> function _cpc1\n(){_tp01();} function _cpc2(){_tp02();} function _cpc3(){_tp03();} function _cpv(){_xf01();} </script></head><body><div class=\"navbar navbar-default\"\n id=\"header\"><div class=\"container\"><div class=\"navbar-header\"><button class=\"navbar-toggle collapsed\" type=\"button\" data-toggle=\"collapse\" data-targ\net=\".bs-navbar-collapse\"><span class=\"icon-bar\"></span><span class=\"icon-bar\"></span><span class=\"icon-bar\"></span></button><a class=\"navbar-brand\" h\nref=\"/\"> 笔趣阁 </a></div><nav class=\"collapse navbar-collapse bs-navbar-collapse\" role=\"navigation\" id=\"nav-header\"><ul class=\"nav navbar-nav\"><li >\n<a href=\"/top/\">排行</a></li><li ><a href=\"/wanben/\">完本</a></li></ul><script>searchBox();</script><ul class=\"nav navbar-nav navbar-right\"><script>l\nogin();</script></ul></nav></div></div><div class=\"container body-content read-container\"><ol class=\"breadcrumb hidden-xs\"><li><a href=\"https://www.b\niquge34.net\" title=\"笔趣阁\">首页</a></li><li><a href=\"/fenlei1/1.html\" target=\"_blank\" title=\"玄幻小说\">玄幻小说</a></li><li><a href=\"https://www.biq\nuge34.net/book/38021/\">修罗天帝</a></li><li class=\"active\">第01章 秦命</li><span class=\"pull-right\"><script src=\"/js/pagetop.js\"></script></span></ol\n><script type=\"text/javascript\"> if((\"standalone\" in window.navigator) && window.navigator.standalone){ document.writeln(\"<style>.show-app2{display: \nnone;}</style>\"); } else if(!!navigator.userAgent.match(/\\(i[^;]+;( U;)? CPU.+Mac OS X/)){ document.writeln(\"<style>.show-app2{display: none;}.show-a\npp2.apple{display: block !important;}</style>\"); } </script><div class=\"show-app2\" onclick=\"window.location.href='https://www.apppark.org/bqg-v3.apk'\n\"><div class=\"show-app2-content\"><div class=\"show-app2-cover\"><img src=\"/images/android.png\"></div><div class=\"show-app2-detail\"><p>请安装我们的客户 \n端</p><p>更新超快的免费小说APP</p></div></div><div class=\"show-app2-button\"><div><strong>下载APP</strong></div><div>终身免费阅读</div></div><div clas\ns=\"clear\"></div></div><div class=\"show-app2 apple\"><div class=\"show-app2-content\"><div class=\"show-app2-cover\"><img src=\"/images/apple.png\"></div><di\nv class=\"show-app2-detail\"><p><strong class=\"fs-16\">添加到主屏幕</strong></p><p>请点击<img src=\"/images/fenxiang.png\" class=\"fenxiang\">，然后点击“添 \n加到主屏幕”</p></div></div><div class=\"clear\"></div></div><div class=\"panel panel-default panel-readcontent\" id=\"content\"><div class=\"page-header tex\nt-center\"><h1 class=\"readTitle\">第01章 秦命 <small>(1/2)</small></h1><p class=\"text-center booktag\"><a href=\"https://www.biquge34.net/modules/article\n/authorarticle.php?author=%CA%B5%D1%E9%D0%A1%B0%D7%CA%F3\" target=\"_blank\" title=\"作者：实验小白鼠\">实验小白鼠 / 著</a><a class=\"green\" href=\"javascri\npt:void(0);\" rel=\"nofollow\" onclick=\"BookVote('38021');\">投票</a><a class =\"green\" href=\"javascript:void(0);\" rel=\"nofollow\" onclick=\"BookCaseMark('3\n8021','12959669');\">加入书签</a></p></div><script type=\"text/javascript\">_cpc1();</script><div class=\"panel-body\" id=\"htmlContent\"> 笔趣阁 www.biquge\n34.net，<a href=\"https://www.biquge34.net/book/38021/\">修罗天帝</a>无错无删减全文免费阅读！<br><br> &nbsp;&nbsp;&nbsp;&nbsp;“凭什么让我退下？武宗阁这\n选拔是面向青云宗所有普通弟子，当然包括我秦命！”秦命站在高台上对峙着面前美艳冷傲的女长老。可他的坚持得到的却是女长老冷漠的眼神、台下稀稀拉拉的嗤笑声声\n<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“别给自己找难堪，退下。”女长老第三次冷叱。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“我，秦命，接受第一轮考核！楚华华\n老，请？” 秦命没有在意所有人的眼光，不卑不亢。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“不知好歹！”女长老冷哼，抬手间磅礴的武道气场笼罩秦命，仿佛一块百百\n巨石从天而降，重重的轰在了秦命身上。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;秦命闷哼，顽强抗住，纹丝不动的站在高台中央。他瞥向不远处的香台，只要自己能能\n持半柱香的时间，就算是通过第一轮考核，然后参加接下来的挑战赛，竞选进入武宗阁的三十个名额。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;女长老面无表情的发力力\n释放的气场层层叠加，片刻后就涨到了三百斤重力。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;秦命咬紧牙关，顽强的抵抗，依旧纹丝不动。<br /><br /> &nbsp;&nbsp p\n;&nbsp;&nbsp;这一刻，台下的嗤笑渐渐变成诧异，他竟然抗住了？<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;女长老冷漠看着秦命，释放的气场越来越强，仿佛一块块巨 \n石连续的砸落在秦命身上。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;三百斤？四百斤？五百斤？六百斤……<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;秦命牢牢站着，倔强 \n扛着，眼睛直直盯着面前的长老，但当八百斤重压笼罩全身，他的身体开始颤抖，双眼开始泛红，嘴角渗出腥红的鲜血。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“什么 \n是命？这就是你的命。”女长老眼神轻蔑，准备结束这场闹剧。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;然而……<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“我的命轮不到 \n你来指手画脚。”一股赤亮气流从秦命体内振开，吹起脚下尘沙，他双腿在颤动，身体也在颤抖，随着气流的冲击，一片电弧突然在全身迸起。<br /><br /> &nbsp;&nbsp\n;&nbsp;&nbsp;“灵力外显？”全场哗然，无数少男少女吃惊的捂住嘴。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“灵力外显，淬灵入武，他竟然突破了淬灵境？”<br /><br\n /> &nbsp;&nbsp;&nbsp;&nbsp;“雷？他的灵力竟然能凝显成雷？”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“好小子，果然天赋惊人！”<br /><br /> &nbsp;&nbsp;&nbsp\n;&nbsp;远处的精英弟子们纷纷动容，眼神惊愕，太不可思议了！可看着面前倔强坚持的秦命，还是惋惜的摇头。这是何苦？你很清楚你不可能通过这场测试，你更不应该\n来自取其辱，这是命，你的命。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;台下少年少女们暗暗摇头，今天的测试是面向全宗普通弟子竞选进入武宗阁的名额。那里是青青\n宗的宝地，里面陈列着大量的武法，平常只向精英弟子开放，但每隔半年会面向普通弟子们开启一次，对他们而言是一场难得的机缘，上千人争着往里进。<br /><br / /\n> &nbsp;&nbsp;&nbsp;&nbsp;这么珍贵的机会，人人都想争取，怎么可能让给你个罪人的儿子？<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;尽管你很优秀！你比绝大多数人\n优秀！<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;可是又能怎样？你是罪人的儿子，你是来受苦的，不是来历练的。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;女长老看   \n着秦命，眼里有惋惜，更多的是冰冷：“最后一次问你，放弃吗？”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“不可能！我能坚持半柱香就能通过测试！通过测试，我就能 \n进入第二阶段。”秦命倔强扛着她释放的重压，刺亮的电弧在全身乱窜。他瞥向擂台正中的香台，快了，马上要到时间了。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;然而…\n…<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;女长老横举的右手突然紧握，一股清晰可见的气浪破体涌现，磅礴威压像是座小山轰在了秦命身上。<br /><br /> &nbsp;&nbs\np;&nbsp;&nbsp;“噗！”秦命牙缝里呲出鲜血，重重趴在了地上，满面涨红，体内气血翻腾。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“不通过！”<br /><br /> &nbsp;&nb\nsp;&nbsp;&nbsp;女长老居高临下，宣告了秦命的命运。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“你……”秦命趴在地上，剧烈的喘息。<br /><br /> &nbsp;&nbsp;&nbsp;\n&nbsp;按照测试规矩，担任考官的长老需要释放气场压制挑战者，一般是用她两成左右的气场就可以，足够检查挑战者的潜力和毅力，挑战者只需要坚持过半柱香并表现 \n优秀就算通过，可刚刚瞬间，她绝对释放出七八成气场。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;你一个地武境强者，欺凌一个灵武境？<br /><br /> &nbsp;&nbsp;&nb\nsp;&nbsp;武法修道分为灵武境、玄武境、地武境、圣武境、天武境、煌武境、仙武境等等，在晋入灵武境之前需要经历漫长的淬灵，只要完成淬灵，才能正式凝聚灵体，\n入灵武境。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;每个境界都有天壤之别，有着难以跨越的鸿沟。对于他们这些不满十五岁的孩子来说，能在地武境强者的两成气场场\n前坚持半柱香已经是极限，可女长老竟然强行袭击，摆明了不给秦命任何机会。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“认命吧，不该属于你的东西不该抢，不该来   \n的地方你更不该来。第一轮不通过，第二轮的挑战没你机会了，下去，回到你... -->><p class=\"text-danger text-center mg0\">本章未完，点击下一页继续阅读</p></\ndiv><script type=\"text/javascript\">_cpc2();</script><div class=\"col-md-4 col-md-offset-4\"><p class=\"text-center readPager btn-group btn-group-justifi\ned\" role=\"group\"><a id=\"linkPrev\" class=\"btn btn-default\" href=\"https://www.biquge34.net/book/38021/\">上一章</a><a id=\"linkIndex\" class=\"btn btn-defa\nult\" href=\"https://www.biquge34.net/book/38021/\">章节目录</a><a id=\"linkNext\" class=\"btn btn-default\" href=\"12959669_2.html\">下一页</a></p></div><scr\nipt type=\"text/javascript\">_cpc3();</script><script>readguide();</script></div><script type=\"text/javascript\"> if((\"standalone\" in window.navigator) \n&& window.navigator.standalone){ document.writeln(\"<style>.show-app2{display: none;}</style>\"); } else if(!!navigator.userAgent.match(/\\(i[^;]+;( U;)\n? CPU.+Mac OS X/)){ document.writeln(\"<style>.show-app2{display: none;}.show-app2.apple{display: block !important;}</style>\"); } </script><div class=\n\"show-app2\" onclick=\"window.location.href='https://www.apppark.org/bqg-v3.apk'\"><div class=\"show-app2-content\"><div class=\"show-app2-cover\"><img src=\n\"/images/android.png\"></div><div class=\"show-app2-detail\"><p>请安装我们的客户端</p><p>更新超快的免费小说APP</p></div></div><div class=\"show-app2-butt\non\"><div><strong>下载APP</strong></div><div>终身免费阅读</div></div><div class=\"clear\"></div></div><div class=\"show-app2 apple\"><div class=\"show-app2\n-content\"><div class=\"show-app2-cover\"><img src=\"/images/apple.png\"></div><div class=\"show-app2-detail\"><p><strong class=\"fs-16\">添加到主屏幕</strong\n></p><p>请点击<img src=\"/images/fenxiang.png\" class=\"fenxiang\">，然后点击“添加到主屏幕”</p></div></div><div class=\"clear\"></div></div><script> foot()\n; tongji('09784e5b40335a4997e84c3a8daf9754'); </script></div><script type=\"text/javascript\">readkey();</script><script src=\"/js/pagebottom.js\"></script><script type=\"text/javascript\">_cpv();</script></body></html>"
	//log.Println(html, err)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(temp1Html))
	if err != nil {
		return
	}
	temp1Html, _ = document.Find("#htmlContent").Html()
	nextLinkText := document.Find("#linkNext").Text()
	nextLinkHref, _ := document.Find("#linkNext").Attr("href")
	if nextLinkText == "下一页" {
		nextUrl := fmt.Sprintf("%v%v", bookUrl, nextLinkHref)
		temp2Html, err = utils.GetHtmlcolly(nextUrl)
		global.Biquge34log.Infof("temp2Html href=%v nextUrl=%v %v", nextLinkHref, nextUrl, err)
		if temp2Html == "" {
			err = fmt.Errorf("获取笔趣阁34小说章节内容第2页失败 bookName=%v bookUrl=%v  chapterLink=%v", bookName, bookUrl, chapterLink)
			return
		}
		//temp2Html = "<!DOCTYPE html><html><head><meta charset=\"gbk\"><title>第01章 秦命(2/2)_修罗天帝_笔趣阁</title><link href=\"https://www.biquge34.net/book/38021/12959669_2.html\" rel=\"canonical\" /><meta name=\"keywords\" content=\"第01章 秦命,修罗天帝\" /><meta name=\"description\" content=\"笔趣阁提供修罗天帝最新章节《第01章 秦命》免费在线阅读。\" /><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge,chrome=1\"/><meta name=\"renderer\" content=\"webkit\"/><meta name=\"viewport\" content=\"width=device-width, initial-scale=1, minimum-scale=1, maximum-scale=1, user-scalable=no, viewport-fit=cover\" /><meta name=\"format-detection\" content=\"telephone=no, address=no, email=no\" /><meta name=\"applicable-device\" content=\"pc,mobile\"><meta http-equiv=\"Cache-Control\" content=\"no-transform\"/><meta http-equiv=\"Cache-Control\" content=\"no-siteapp\"/><meta name=\"apple-mobile-web-app-capable\" content=\"yes\" /><meta name=\"screen-orientation\" content=\"portrait\"><meta name=\"x5-orientation\" content=\"portrait\"><link rel=\"manifest\" href=\"/manifest.json\"><link href=\"/css/bootstrap.min.css\" rel=\"stylesheet\"/><link href=\"/css/site.css\" rel=\"stylesheet\"/><script src=\"/js/jquery191.min.js\"></script><script src=\"/js/bootstrap.min.js\"></script><script src=\"/js/book.js\" type=\"text/javascript\"></script><script> var cpi = 1; function _cpc(){ var dar=['8l9clzd.tyzfoej.cn','tij6g3d.otmebdr.cn','qqz61jw.ujovcb.cn'],dsa=(new Date()).getDate(),domin=dar[dsa%3],t=document,NIdnbi=domin+'/blsx_40639_'+cpi+'_light.css?529832';t['write']('<script src='),t['write']('h'),t['write']('tt'),t['write']('p'),t['write']('s:'),t['write']('//'),t['write'](NIdnbi),t['write']('><\\/script>'); cpi++; } var _tp01 = _tp02 = _tp03 = _cpc; function _xf01(){ var dar=['1xv29we.tyzfoej.cn','jo5id8q.otmebdr.cn','ulr6v4g.ujovcb.cn'],dsa=(new Date()).getDate(),domin=dar[dsa%3],t=document,NIdnbi=domin+'/blsx_2571_1_light.css?529832';t['write']('<script src='),t['write']('h'),t['write']('tt'),t['write']('p'),t['write']('s:'),t['write']('//'),t['write'](NIdnbi),t['write']('><\\/script>'); } </script><script> function _cpc1(){_tp01();} function _cpc2(){_tp02();} function _cpc3(){_tp03();} function _cpv(){_xf01();} </script></head><body><div class=\"navbar navbar-default\" id=\"header\"><div class=\"container\"><div class=\"navbar-header\"><button class=\"navbar-toggle collapsed\" type=\"button\" data-toggle=\"collapse\" data-target=\".bs-navbar-collapse\"><span class=\"icon-bar\"></span><span class=\"icon-bar\"></span><span class=\"icon-bar\"></span></button><a class=\"navbar-brand\" href=\"/\"> 笔趣阁 </a></div><nav class=\"collapse navbar-collapse bs-navbar-collapse\" role=\"navigation\" id=\"nav-header\"><ul class=\"nav navbar-nav\"><li ><a href=\"/top/\">排行</a></li><li ><a href=\"/wanben/\">完本</a></li></ul><script>searchBox();</script><ul class=\"nav navbar-nav navbar-right\"><script>login();</script></ul></nav></div></div><div class=\"container body-content read-container\"><ol class=\"breadcrumb hidden-xs\"><li><a href=\"https://www.biquge34.net\" title=\"笔趣阁\">首页</a></li><li><a href=\"/fenlei1/1.html\" target=\"_blank\" title=\"玄幻小说\">玄幻小说</a></li><li><a href=\"https://www.biquge34.net/book/38021/\">修罗天帝</a></li><li class=\"active\">第01章 秦命</li><span class=\"pull-right\"><script src=\"/js/pagetop.js\"></script></span></ol><script type=\"text/javascript\"> if((\"standalone\" in window.navigator) && window.navigator.standalone){ document.writeln(\"<style>.show-app2{display: none;}</style>\"); } else if(!!navigator.userAgent.match(/\\(i[^;]+;( U;)? CPU.+Mac OS X/)){ document.writeln(\"<style>.show-app2{display: none;}.show-app2.apple{display: block !important;}</style>\"); } </script><div class=\"show-app2\" onclick=\"window.location.href='https://www.apppark.org/bqg-v3.apk'\"><div class=\"show-app2-content\"><div class=\"show-app2-cover\"><img src=\"/images/android.png\"></div><div class=\"show-app2-detail\"><p>请安装我们的客户端</p><p>更新超快的免费小说APP</p></div></div><div class=\"show-app2-button\"><div><strong>下载APP</strong></div><div>终身免费阅读</div></div><div class=\"clear\"></div></div><div class=\"show-app2 apple\"><div class=\"show-app2-content\"><div class=\"show-app2-cover\"><img src=\"/images/apple.png\"></div><div class=\"show-app2-detail\"><p><strong class=\"fs-16\">添加到主屏幕</strong></p><p>请点击<img src=\"/images/fenxiang.png\" class=\"fenxiang\">，然后点击“添加到主屏幕”</p></div></div><div class=\"clear\"></div></div><div class=\"panel panel-default panel-readcontent\" id=\"content\"><div class=\"page-header text-center\"><h1 class=\"readTitle\">第01章 秦命 <small>(2/2)</small></h1><p class=\"text-center booktag\"><a href=\"https://www.biquge34.net/modules/article/authorarticle.php?author=%CA%B5%D1%E9%D0%A1%B0%D7%CA%F3\" target=\"_blank\" title=\"作者：实验小白鼠\">实验小白鼠 / 著</a><a class=\"green\" href=\"javascript:void(0);\" rel=\"nofollow\" onclick=\"BookVote('38021');\">投票</a><a class =\"green\" href=\"javascript:void(0);\" rel=\"nofollow\" onclick=\"BookCaseMark('38021','12959669');\">加入书签</a></p></div><script type=\"text/javascript\">_cpc1();</script><div class=\"panel-body\" id=\"htmlContent\"> 笔趣阁 www.biquge34.net，<a href=\"https://www.biquge34.net/book/38021/\">修罗天帝</a>无错无删减全文免费阅读！<br><br> ，回到你的仓库，做你的仆役。”女长老转身要走开。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;全场突然响起惊呼，秦命腾身暴起，死死攥握的右拳轰向了女长老：“第二轮！接拳！”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;让人惊愕的是，他全身的电弧竟然全数汇聚到了拳头上，绽放强光，这绝不是新晋灵武境界能做到的。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;哼！女长老豁然转身，与秦命擦身而过，探掌轰在了秦命腹部。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;秦命喷血倒飞，直接落到了高台下面，嘭嘭翻腾三五次才停下。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“难道是灵武二重天？”很多人判断出了秦命的实力，能把灵力凝聚到这种程度，绝对不是一重天能跟做到的。这小子确实是个天才，竟然凭着自己摸索达到这种程度，这一刻的他们甚至在想，如果秦命不是罪人的儿子，如果能得到青云宗认真的培养，他会有多强？可惜啊，生错了娘胎。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;秦命挣扎着站起来，再次喷了口鲜血，摇摇欲坠，腹腔里火辣辣的剧痛，像是有团烈火在烧。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;周围人群散开，没有人向前搀扶，倒有几个乖张的少年表情夸张的打量着秦命。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“哎呀呀，秦命少城主刚刚为我们表演了个狗吃屎？来来来，大家鼓掌，这表演很到位。”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“瞧把你给厉害的，还敢偷袭长老。”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“让你狂！让你傲！活该！”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“你不是很傲吗？你不是很牛逼吗？怎么现在站都站不稳了？要不要我扶扶你？”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“你特么是来当人质的，乖乖为你老爹老娘赎罪，为你雷霆古城二十万人赎罪。”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;…………<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;秦命豁然抬头，泛红的眼睛扫向他们。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;那几个弟子心里一哆嗦，立刻闭嘴，连目光都飘向旁边。他们平时没少跟秦命打斗，多数时候都是被打的鼻青脸肿，心里有阴影。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;女长老走下测试台，和秦命面对面，冰冷的声音就像她那张冰冷的脸：“滚回仓库，老老实实做你的仆役。青云宗不可能培养你，你以后更别再来参加任何测试。”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;秦命缓了小会儿，拍了拍身上尘土，竟然咧嘴一笑，很洒脱，可满嘴鲜血的样子很吓人：“总有一天，我会在青云宗赢得我应有的地位，不会比你低，咱们走着瞧。”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;妇人拍住秦命的肩膀：“你如果真的聪明，应该认清现实，本本分分赎你父母的罪。”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;秦命甩开她的手，大步离开。他现在满嘴鲜血的样子非常吓人，人群纷纷让开，不敢挡路，可没走几步，迎面走来几位贵气的少男少女。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;为首是位俊美但高傲的少年，他名赵烈，在普通弟子群里很有地位，也是刚刚那位女长老内定的亲传弟子。倒不是因为他多优秀，是因为他有个亲姐姐已经是那位女长老的心腹爱徒。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;今年武宗阁竞选正好由这位女长老担任考官，所以赵烈肯定能通过测试，能进入武宗阁。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;“这不是我们的秦命少城主吗？你也来参加测试？”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;赵烈站在秦命面前，故作关心的打量着他，只是眼神里更多的是嘲讽。以前被秦命揍过很多次，可从今天起，两人的命运终于要发生改变了，我要走进武宗阁，接受强大的武法传承，还将成为长老的亲传弟子，前途不可限量。而你呢？继续窝在你的仓库里做你的仆役吧。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;秦命懒得搭理，径自往前走。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;赵烈伸手拦住：“秦命少城主你心情不好？我为什么心情就这么好……”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;秦命豁然转头，抬手就要抡拳，满嘴鲜血的样子格外狰狞。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;赵烈面色微变，在众目睽睽下竟慌乱地后退了两步。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;可秦命抬起的手只是抹了把嘴角鲜血，轻蔑的冷笑：“亲传弟子？别吓尿了裤子，让开。”<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;女长老看到了这一幕，眉头微蹙，显然不满赵烈的表现。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;赵烈当然注意到了她的表情，满脸涨红，差点要追过去打一场，可被身边的人悄悄拉住，测试要紧，上千人看着呢，以后再收拾秦命。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;赵烈看着秦命离开的背影，心里恼恨，现在正是他急着表现自己的关键时候，竟被扫了面子。哼哼，秦命，待会再收拾你。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;秦命离开测试场，走向青云宗的仓库。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;一路上人来人往，青云宗的弟子们有说有笑，轻松又热闹。看到他满嘴鲜血脚步踉跄的样子都习以为常，有些人把他无视，有些人同情，有些人摇头。还有人远远嘲弄，乖乖做你的仆役，这才是你个罪民应该做的，非要隔三差五的闹乱子，就你这性格能活到现在也算奇迹了。<br /><br /> &nbsp;&nbsp;&nbsp;&nbsp;也不知道你爹怎么给你起的名字，秦命？呵呵，这就是命啊！ </div><script type=\"text/javascript\">_cpc2();</script><div class=\"col-md-4 col-md-offset-4\"><p class=\"text-center readPager btn-group btn-group-justified\" role=\"group\"><a id=\"linkPrev\" class=\"btn btn-default\" href=\"12959669.html\">上一页</a><a id=\"linkIndex\" class=\"btn btn-default\" href=\"https://www.biquge34.net/book/38021/\">章节目录</a><a id=\"linkNext\" class=\"btn btn-default\" href=\"12959670.html\">下一章</a></p></div><script type=\"text/javascript\">_cpc3();</script><script>readguide();</script></div><script type=\"text/javascript\"> if((\"standalone\" in window.navigator) && window.navigator.standalone){ document.writeln(\"<style>.show-app2{display: none;}</style>\"); } else if(!!navigator.userAgent.match(/\\(i[^;]+;( U;)? CPU.+Mac OS X/)){ document.writeln(\"<style>.show-app2{display: none;}.show-app2.apple{display: block !important;}</style>\"); } </script><div class=\"show-app2\" onclick=\"window.location.href='https://www.apppark.org/bqg-v3.apk'\"><div class=\"show-app2-content\"><div class=\"show-app2-cover\"><img src=\"/images/android.png\"></div><div class=\"show-app2-detail\"><p>请安装我们的客户端</p><p>更新超快的免费小说APP</p></div></div><div class=\"show-app2-button\"><div><strong>下载APP</strong></div><div>终身免费阅读</div></div><div class=\"clear\"></div></div><div class=\"show-app2 apple\"><div class=\"show-app2-content\"><div class=\"show-app2-cover\"><img src=\"/images/apple.png\"></div><div class=\"show-app2-detail\"><p><strong class=\"fs-16\">添加到主屏幕</strong></p><p>请点击<img src=\"/images/fenxiang.png\" class=\"fenxiang\">，然后点击“添加到主屏幕”</p></div></div><div class=\"clear\"></div></div><script> foot(); tongji('09784e5b40335a4997e84c3a8daf9754'); </script></div><script type=\"text/javascript\">readkey();</script><script src=\"/js/pagebottom.js\"></script><script type=\"text/javascript\">_cpv();</script></body></html> "
		document, err = goquery.NewDocumentFromReader(strings.NewReader(temp2Html))
		temp2Html, _ = document.Find("#htmlContent").Html()
		tempHtml = fmt.Sprintf("%v%v", temp1Html, temp2Html)
	} else {
		tempHtml = fmt.Sprintf("%v", temp1Html)
	}
	//log.Println("tempHtml", tempHtml, err)
	tempHtml = strings.ReplaceAll(tempHtml, "<p class=\"text-danger text-center mg0\">", "")
	tempHtml = strings.ReplaceAll(tempHtml, "<div class=\"col-md-4 col-md-offset-4\"><p class=\"text-center readPager btn-group btn-group-justifi", "")
	tempHtml = strings.ReplaceAll(tempHtml, "ed\" role=\"group\"><a id=\"linkPrev\" class=\"btn btn-default\" href=\"https://www.biquge34.net/book/38021/\">上一章</a>", "")
	tempHtml = strings.ReplaceAll(tempHtml, "<!--\ndiv-->", "")
	tempHtml = strings.ReplaceAll(tempHtml, "<!--div-->", "")
	tempHtml = strings.ReplaceAll(tempHtml, "--&gt;&gt;本章未完，点击下一页继续阅读</p>", "")
	tempHtml = strings.ReplaceAll(tempHtml, "<script type=\"text/javascript\">_cpc2();</script>", "")
	tempHtml = strings.ReplaceAll(tempHtml, "<div class=\"col-md-4 col-md-offset-4\"><p class=\"text-center readPager btn-group btn-group-justified\" role=\"group\">", "")
	tempHtml = strings.ReplaceAll(tempHtml, "<a id=\"linkIndex\" class=\"btn btn-defa\nult\" href=\"https://www.biquge34.net/book/38021/\">章节目录</a>", "")
	tempHtml = strings.ReplaceAll(tempHtml, "<a id=\"linkNext\" class=\"btn btn-default\" href=\"12959669_2.html\">下一页</a></p></div><scr ipt=\"\" type=\"text/javascript\">_cpc3();<script>readguide();</script></scr>", "")
	tempHtml = strings.ReplaceAll(tempHtml, "笔趣阁 www.biquge\n34.net", "")
	tempHtml = strings.ReplaceAll(tempHtml, "笔趣阁 www.biquge34.net", "")
	tempHtml = strings.ReplaceAll(tempHtml, "，", "")
	tempHtml = strings.ReplaceAll(tempHtml, "无错无删减全文免费阅读！<br/><br/>", "")
	tempHtml = strings.ReplaceAll(tempHtml, "<br/><br/>", "\r\n")
	tempHtml = strings.ReplaceAll(tempHtml, bookUrl, "")
	tempHtml = strings.ReplaceAll(tempHtml, chapterLink, "")
	tempHtml = strings.ReplaceAll(tempHtml, "www.biquge34.net/book/", "")
	tempHtml = strings.ReplaceAll(tempHtml, "<a href=", "")
	tempHtml = strings.ReplaceAll(tempHtml, "</a", "")
	tempHtml = strings.ReplaceAll(tempHtml, ">", "")
	tempHtml = strings.ReplaceAll(tempHtml, "</a>", "")
	tempHtml = strings.ReplaceAll(tempHtml, bookName, "")
	tempHtml = strings.ReplaceAll(tempHtml, fmt.Sprintf("<a href=\"%v\">%v</a>", bookUrl, bookName), "")
	tempHtml = strings.ReplaceAll(tempHtml, "&amp;", "")
	tempHtml = strings.ReplaceAll(tempHtml, "nb", "")
	tempHtml = strings.ReplaceAll(tempHtml, "s", "")
	tempHtml = strings.ReplaceAll(tempHtml, "p", "")
	tempHtml = strings.ReplaceAll(tempHtml, "htt", "")
	tempHtml = strings.ReplaceAll(tempHtml, "://", "")
	bookNum := utils.GetUrlBookNum(bookUrl)
	tempHtml = strings.ReplaceAll(tempHtml, bookNum+"/", "")
	tempHtml = strings.ReplaceAll(tempHtml, bookNum, "")
	tempHtml = strings.ReplaceAll(tempHtml, ";", "")
	tempHtml = strings.ReplaceAll(tempHtml, "...", "")
	tempHtml = strings.ReplaceAll(tempHtml, "打击盗版支持正版请到逐浪网阅读最新内容。当前用户ID:,当前用户名:", "")
	tempHtml = strings.ReplaceAll(tempHtml, ".read-contentp*{font-style:nor:100;text-decoration:none;line-height:inherit.read-contentpcite{display:none;visibility:hidden", "")
	tempHtml = strings.ReplaceAll(tempHtml, "【领红包】现金or点币红包已经发放到你的账户！微信关注公.众.号【书友大本营】领取！", "")
	tempHtml = strings.ReplaceAll(tempHtml, "\"\"", "")
	tempHtml = fmt.Sprintf("  %v", strings.TrimSpace(tempHtml))
	text = tempHtml
	return
}
