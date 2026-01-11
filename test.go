package main

import (
	"crypto/tls"
	"fmt"
	"go-novel/app/service/common/chapter_service"
	"go-novel/db"
	"go-novel/utils"
	"gopkg.in/gomail.v2"
	mailer "gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
)

func main() {
	host, name, user, passwd := db.GetDB()
	db.InitDB(host, name, user, passwd)
	db.InitZapLog()
	log.Println(utils.GetGuestName())
	log.Println(utils.RangeNum(1, 2))
	log.Println(utils.Md5("第73章 挖莲藕，摸菱角"))
	log.Println(utils.Md5("第73章 挖莲藕,摸菱角"))
	html, err1 := utils.GetHtmlcolly("https://www.27k.net/read/2834/1973655.html")
	log.Println(html, err1)
	return
	//fileName := strings.Join(pinyin.LazyPinyin("万相之王", pinyin.NewArgs()), "")
	//
	//log.Println(fileName)
	//author := utils.GetFirstLetter("天蚕土豆")
	//log.Println(fmt.Sprintf("%v-%v", fileName, author))
	//log.Println(utils.GetFileBase("test.jpg"))
	utils.AutoModel()
	return
	content, err1 := ioutil.ReadFile("/data/txt/bd1c8f228cdd256f4c3e770244108bb4/144fbc378cb33173932d7274ab7a0b78.txt")
	if err1 != nil {
		fmt.Println("无法读取文件:", err1)
		return
	}
	// 统计字节数
	byteCount := len(content)
	fmt.Println("文件字节数:", byteCount)

	// 将文件内容转换为字符串
	text1 := string(content)

	// 统计字符数
	charCount := len([]rune(text1))

	fmt.Println("文件字符数:", charCount)

	log.Println(chapter_service.GetTxtDir("逆天邪神", "火星引力"))
	return
	//var chapterTable string
	var err error
	//utils.AutoModel()
	//mdata, err := adver_service.GetAdverMap()
	//if err != nil {
	//	return
	//}
	//log.Println(mdata)
	//start := time.Now()
	//chapterTable, err = book_service.GetChapterTable(3)
	//if err != nil {
	//	err = fmt.Errorf("%v", "获取章节失败")
	//	return
	//}
	//log.Println(err)
	//chapters, err := book_service.GetChapterList(chapterTable, "desc")
	//if err != nil {
	//	return
	//}
	//log.Println(chapters, err)
	//log.Println(time.Since(start))
	return
	//checkin_service.CheckRemind()
	var registIds []string
	registIds = []string{"141fe1da9fad52d3c93"}
	_, err = utils.JpushMsg("您今天还未进行签到", registIds)
	log.Println(err)
	if err != nil {
		return
	}
	return
	log.Println(utils.GetWeekyUnix())
	pageNum := 1
	pageCount := 5
	bookNum := 5
	bookCount := 20

	totalPages := pageCount * bookCount
	readPages := (pageNum-1)*bookCount + bookNum
	progress := float64(readPages) / float64(totalPages) * 100
	log.Println(progress)
	//log.Println(utils.GetWanFormatted(20000000))
	//generate_id.InitId()
	//ids, err := generate_id.ThirHBOrder.NextID()
	//log.Println(ids, err)
	//log.Println(utils.GetRandomUsername())
	//utils.AutoModel()
	text := "  全本小说网 www.qb50.com，最快更新<a href=\"https://www.qb50.com/book_99149/\">重生都市仙帝</a>最新章节！<br><br>     同一时间，就在方灵给张逸风讲故事的时候，司马寒重伤而归的消息，传遍了天山各大门派。<br/><br/> >万妖圣祖</a>最新章节！<br><br>     少年手掌，抓握着那劈落下的刑刀，在刽子手，项权，周围无数民众震惊的眼神中缓缓站起了身。<br/><br/>     滴答！滴答……<br/><br/>     鲜血从他掌心之中流下，染红了刑刀，可是被一股力量保护，没有能劈断他的手。<br/><br/>     少年抬起了头，乱发遮盖下的暗金色眼眸中绽放出了冰冷的嗜血光芒。<br/><br/>     “怎，怎么可能？”<br/><br/>     刽子手震惊的望着少年，自己也是九重体魄境界的强者，这一刀蕴含千斤之力，哪怕是钢铁都可以劈开，竟然，竟然被他用手抓握挡住了。<br/><br/>     “小子，撒手！”<br/><br/>     刽子手怒吼，想抽刀，可是却是抽不出，那只手掌仿佛一双大铁钳子死死抓握住了。<br/><br/>     “你，将是我涅槃重生所杀第一人！”<br/><br/>     项尘露出了冰冷的笑意，另一手成拳，一拳狂暴轰杀而出。<br/><br/>     “吼！”<br/><br/>     拳风呼啸，竟然发出野兽咆哮之声，拳头绽放金色光芒，蕴含恐怖的力道轰击在了震惊的刽子手身躯。<br/><br/>     嘭！<br/><br/>     “噗嗤……”<br/><br/>     刽子手惨叫，双目怒瞪滚圆快要裂眶而出，身躯仿佛被一柄万斤巨锤轰中，一声惨叫，被轰击崩飞而退十多米远，鲜血狂喷，鲜血中还蕴含破碎的内脏碎沫，胸膛凹陷一大片胸骨破碎。<br/><br/>     人还没有落地就已经断气。<br/><br/>     力量，一股强大无匹，前所未有感觉过的恐怖力量充斥满了这具身体。<br/><br/>     寂静，全场瞬间寂静下来，鸦雀无声，所有人震惊的望向了少年，眼眸中全是不可思议。<br/><br/>     怎么，前一刻还是待宰的羔羊，下一秒瞬间怎么变成了那噬人的猛虎！<br/><br/>     无数人目瞪口呆的望着那死去断气的刽子手，随后一声尖叫声打破安静。<br/><br/>     “杀人啦！”<br/><br/>     民众尖叫，四散而逃。<br/><br/>     项权也是瞪大眼睛望向了死去的刽子手，随后暴怒出声。<br/><br/>     “你们还愣着干什么，还不快杀了他！”<br/><br/>     项权尖声叫道，身躯连连后退，再望向这少年的身躯中，对方的软弱气势瞬间变了，变得凌厉，霸道，仿佛一头远古凶兽散发可怕的危险气息。<br/><br/>     “啊，是！”<br/><br/>     周围的十多名带刀护卫也是瞬间回过神来，一个个抽出了四尺腰刀围了上来，围住了项尘。<br/><br/>     项尘身子微低，口中发出了低沉的咆哮之声，口中犬齿外露。<br/><br/>     “小畜生，去死！”<br/><br/>     一名护卫低吼壮胆，脚掌一踏，身躯冲向项尘，其他人也纷纷动了，十多人抽刀劈杀向了项尘。<br/><br/>     唰！<br/><br/>     这一刀寒光劈来，刀上缭绕着一股红色的真气蕴含炙热温度，刀锋破空呼呼做响，这一刀直取项尘头颅。<br/><br/>     这些人，都是大王妃麾下的养的家仆，修为最低都有体魄巅峰境界修为，一个人可以击杀寻常二三十个普通大汉，对常人来说都是精英。<br/><br/>     嘭！<br/><br/>     项尘脚步一踏，微屈的身子宛如一张蓄力的长弓瞬间爆发力量，速度快如一颗出膛的炮弹，瞬间冲到了对方面前，一拳重击轰在对方胸膛。<br/><br/>     “咔嚓！”<br/><br/>     胸骨断碎，这名护卫惨叫一声，被一拳轰飞十多米，口吐鲜血直接没有了生气，体内心脏都被一拳狂暴击碎。<br/><br/>     第二人，死！<br/><br/>     “唰！”<br/><br/>     周围好几道真气刀劲劈来，花岗岩石地面被划出了一道道口气撕裂而来，三名护卫的攻击联手杀至。<br/><br/>     嘭！<br/><br/>     项尘身躯一踏冲天而起，竟然飞跃起了十多米高，脱离了几道刀气的劈杀，他身躯扑落，手掌之中，一道道淡金色利爪竟然指骨之中生长而出，一爪撕裂向了一名护卫。<br/><br/>     噗嗤！<br/><br/>     锋利的爪子划过，竟然有五道金色爪劲如同光刃撕裂而出，那护卫脖子被撕裂过，头颅抛飞而起鲜血井喷，无头的尸体借着惯性又冲了几步！<br/><br/>     第三人，杀！<br/><br/>     真武修者会武技，那是一种能提升自己真气战斗力的武功，不过，项尘不会，他只是凭借自己体内强大的爆发力量和自己的本能杀人。<br/><br/>     “嘶吼……！”<br/><br/>     落地后的项尘一声低吼，双手撑地，整个人如同一头下山猛虎扑出，身躯又弹射而出，瞬间冲至一人身前，手指利爪刺出插入了对方的胸膛。<br/><br/>     “噗嗤！”<br/><br/>     手掌掏出之时，那人胸膛被贯穿出血洞，那人张大嘴巴后退，身子一扬，嘭的一声倒在地上，死不瞑目。<br/><br/>     第四人，陨！<br/><br/>     项尘宛如脱缰野马，下山猛虎，扑冲，撕杀，重拳击打，利爪撕裂，都是最简单的招式，可是却是蕴含可怕的霸道力量没有人能抗过他一击，速度又是快得惊人。<br/><br/>     “怎么可能，这，这小畜生不是没有半点灵力修为吗？怎么会变得如此强大？”<br/><br/>     项权见项尘杀人宛如猛虎撕羊，眼眸中全是惊骇，这些护卫可不是普通人啊，都是拥有真气的修真武者。<br/><br/>     嘭！<br/><br/>     终于，最后一人被他一拳打爆了头颅，十多人，数息之间被他一人击杀！<br/><br/>     “啊，这，这……”<br/><br/>     项权更是吓得一屁股坐在了地上，面色苍白。<br/><br/>     滴答，滴答……<br/><br/>     浑身流淌敌人鲜血的项尘一步步走向了项权，项权吓得不停后爬。<br/><br/>     “不要过来，不要过来！”<br/><br/>     项权惊恐大吼，起身爬起，转身奔逃。<br/><br/>     呼……！<br/><br/>     项尘双手撑地，后腿发力一扑，身躯如同一头猛虎跨越了十数米距离挡在了项权前方。<br/><br/>     “啊！”<br/><br/>     项权吓得又一屁股坐在地上，面色苍白，他的修为和那些护卫差不多。<br/><br/>     “不要杀我，不要杀我！”<br/><br/>     项权惊恐嘶吼，恐惧望着染血的少年，随后他竟然跪在了地上，惊恐求饶说道：“少爷饶命，少爷饶命，不是我要杀您，是大王妃要杀您啊。”<br/><br/>     项权磕头求饶，一把鼻涕一把泪。<br/><br/>     “少爷……”<br/><br/>     项尘冰冷自嘲笑了，：“你不是一直叫我小畜生小野种吗？原来，有实力就可以是少爷？原来，这世间的道理都需要实力去悍卫。”<br/><br/>     他自言自语自嘲的笑着。<br/><br/>     他走到了项权身前，冰冷道：“四年前，你大儿子患重病无钱可医治，母亲怜你赏你灵药。<br/><br/>     而你被妖兽重创，也是我母亲所救，却是投靠了大王妃来害我！大恩换来大仇，你放心，我，早晚会让大王妃那个贱人下去陪你，让你继续侍奉她的，你的儿子得到的前程，我也会亲手撕碎！<br/><br/>     忘恩负义者，杀！”<br/><br/>     “不！”<br/><br/>     项权一听，一股死亡危机瞬间降临心头。<br/><br/>     嘭！<br/><br/>     项尘一脚踢在了他的下巴上，可怕的力量竟然直接把下巴踢碎，脖子向后而断扭曲一百八十度瞬间惨死当场。<br/><br/>     呼呼……<br/><br/>     北风卷过，吹来了一缕寒意，灰尘四起。<br/><br/>     整个刑场空空荡荡，血腥味弥漫，只留下来十多具尸体还有一名染血的少年。<br/><br/>     少年走向了两颗头颅，明叔的，红袖的，他抱在了怀中，无声嘶吼着。<br/><br/>     不过这一次，他没有再流出泪。<br/><br/>     泪干了，心已经痛麻木了。<br/><br/>     “明叔，红袖，我项尘发誓，我会活着，把你们的那一份也活下去，哪怕天弃我，神弃我，人弃我，还是妖魔鬼怪都弃我项尘，我也会活下去！”<br/><br/>     项尘对两人首级三跪，哽咽低沉道，发下誓言，随后把首级装入布袋里背在背上。<br/><br/>     “妹妹！”<br/><br/>     项尘却是没有停留，狂奔离开这里，奔向了王家所在，眼眸中全是担忧，今日是妹妹叶柔被逼嫁给王家那个花心纨绔少爷的日子，他绝对不能让妹妹被那个纨绔少爷为妾毁了以后余生。<br/><br/>     纵身死，也得杀入王家，救出妹妹！ "
	log.Println(utils.ReplaceText(text))

	//text = `重生都市仙帝</a>最新章节！<br><br> 张逸风拖着行李，穿过了几条小道，刚要抵达宿舍。一道声音忽然叫住了他。<br/><br/> “张逸风？”<br/><br/> 声音有些冷`
	//
	//// 使用正则表达式将 <br> 替换为换行符
	//re := regexp.MustCompile(`<br\s*/?>`)
	//text = re.ReplaceAllString(text, "\n")
	//
	//fmt.Println(text)
	return
	//startTime, endTime := utils.GetWeekDayRange(1)
	//log.Println(1, startTime, endTime)
	//log.Println(utils.GetAgoDayUnix(7 - 1))

	//log.Println(utils.DaysSince(1700148514))
	//log.Println(utils.GetUnix(), utils.GetAgoDayUnix(1))

	//banners, err := book_service.GetBookClassList()
	//if err != nil {
	//	return
	//}
	//log.Println(banners)
	//arr := book_service.GetAppData(banners)
	//
	//log.Printf("%+v", arr)
	//for i, v := range arr {
	//	log.Println(i, v)
	//}
	// 调用 getRank 函数
	//log.Println(book_service.GetRank(global.DB, 1))
	//log.Println(book_service.GetRank(1, "cion"))
	//var chapters []*models.McBookChapter
	//chapters, err := book_service.GetChapterList("mc_book_chapter_1", 1, 5)
	//if err != nil {
	//	return
	//}
	//for _, i2 := range chapters {
	//	log.Printf("%+v", i2)
	//}
	//chapter, err := book_service.GetChapterByChapterId("mc_book_chapter_1", 111)
	//log.Println(chapter, err, chapter == nil)
	//log.Println(utils.GetCion(50, 1))
	//log.Println(book_service.GetBookTxt(1, 1, ""))
}

func SendEmail1(to, code string) (err error) {
	var host string = "smtp.qiye.aliyun.com"
	var port int = 465
	var username string = "admin@yaodudushu.site"
	var passwd string = "rens123.123"

	// 1. 首先构建一个 Message 对象，也就是邮件对象
	msg := mailer.NewMessage()
	// 2. 填充 From，注意第一个字母要大写
	msg.SetHeader("From", username)
	// 3. 填充 To
	msg.SetHeader("To", to)
	// 5. 设置邮件标题
	msg.SetHeader("Subject", "go-novel")
	// 6. 设置要发送的邮件正文
	// 第一个参数是类型，第二个参数是内容
	// 如果是 html，第一个参数则是 `text/html`
	msg.SetBody("text/html", code)
	// 7. 添加附件，注意，这个附件是完整路径
	//msg.Attach("/Users/yufei/Downloads/1.jpg")
	// 8. 创建 smtp 实例
	// 如果你的阿里云企业邮箱则是密码，否则一般情况下国内国外使用的都是授权码
	dialer := mailer.NewDialer(host, port, username, passwd)
	// 9. 发送邮件，连接邮件服务器，发送完就关闭
	return dialer.DialAndSend(msg)
}

func SendEmail2(to, code string) (err error) {
	// 1. 首先构建一个 Message 对象，也就是邮件对象
	msg := mailer.NewMessage()
	// 2. 填充 From，注意第一个字母要大写
	msg.SetHeader("From", "admin@yaodudushu.site")
	// 3. 填充 To
	msg.SetHeader("To", to)
	// 5. 设置邮件标题
	msg.SetHeader("Subject", "code")
	// 6. 设置要发送的邮件正文
	// 第一个参数是类型，第二个参数是内容
	// 如果是 html，第一个参数则是 `text/html`
	msg.SetBody("text/html", code)
	// 7. 添加附件，注意，这个附件是完整路径
	//msg.Attach("/Users/yufei/Downloads/1.jpg")
	// 8. 创建 smtp 实例
	// 如果你的阿里云企业邮箱则是密码，否则一般情况下国内国外使用的都是授权码
	dialer := mailer.NewDialer("smtpdm.aliyun.com", 465, "admin@yaodudushu.site", "rens123.123")
	// 9. 发送邮件，连接邮件服务器，发送完就关闭
	return dialer.DialAndSend(msg)
}

func SendEmail(to, code string) (err error) {
	//username := "admin@yaodudushu.site"
	//password := "rens123.123"

	username := "services@yaodudushu.site"
	password := "rens456.456"

	//username := "support@yaodudushu.site"
	//password := "rens678.678"
	// 1. 首先构建一个 Message 对象，也就是邮件对象
	msg := gomail.NewMessage()
	// 2. 填充 From，注意第一个字母要大写
	msg.SetHeader("From", username)
	// 3. 填充 To
	msg.SetHeader("To", to)
	// 5. 设置邮件标题
	msg.SetHeader("Subject", "code")
	// 6. 设置要发送的邮件正文
	// 第一个参数是类型，第二个参数是内容
	// 如果是 html，第一个参数则是 `text/html`
	msg.SetBody("text/html", code)
	// 7. 添加附件，注意，这个附件是完整路径
	//msg.Attach("/Users/yufei/Downloads/1.jpg")
	// 8. 创建 Dialer
	dialer := gomail.NewDialer("smtp.sg.aliyun.com", 465, username, password)
	// 9. 禁用 SSL
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// 10. 发送邮件
	return dialer.DialAndSend(msg)
}
