package utils

import (
	"github.com/olahol/melody"
	"go-novel/app/models"
)

const (
	Nickname              = "nickname"
	Email                 = "email"
	Sex                   = "sex"
	Pic                   = "pic"
	Tel                   = "tel"
	BookType              = "book_type"
	Invite                = "invite"
	Follow                = "follow"
	Passwd                = "passwd"
	Rec                   = "rec"
	Hot                   = "hot"
	Hits                  = "hits"
	Search                = "search"
	Serialize             = "serialize"
	Score                 = "score"
	New                   = "new"
	Choice                = "choice"
	Classic               = "classic"
	Classicrec            = "classicrec"
	Today                 = "today"
	Yesterday             = "yesterday"
	Agoday                = "agoday"
	Notice                = "notice"
	Praise                = "praise"
	Comment               = "comment"
	CollectInfo           = "collect_info"             //采集信息
	CollectSourceUrl      = "collect_source_url"       //分页链接
	CollectPageUrl        = "collect_page_url"         //分页链接
	CollectPageBookUrl    = "collect_page_book_url"    //分页小说列表
	CollectChapterUrl     = "collect_chapter_url"      //章节列表
	CollectChapterUrlTemp = "collect_chapter_url_temp" //章节列表
	CollectList           = "collect_list"
	CollectChapter        = "collect_chapter"
	CollectLog            = "collect_log"
	ZssqCategory          = "zssq_category"     //分类列表
	ZssqBooks             = "zssq_books"        //小说列表
	Biquge34Chapters      = "biquge34_chapters" //分类列表
	PaoshubaChapters      = "paoshuba_chapters" //分类列表
	XswCategory           = "xsw_category"
	XswBooks              = "xsw_books"
	LydCategory           = "lyd_category"
	LydBooks              = "lyd_books"
	Bqg24Category         = "bqg24_category" //分类列表
	Bqg24Books            = "bqg24_books"
	Siluke520Category     = "siluke520_category" //分类列表
	Siluke520Books        = "siluke520_books"
	Chapter               = "chapter"
	DefaultSource         = "默认书源"
	RoundDot              = "・"
	RoundDotBig           = "●"
	UploadBookPicPath     = "uploadBookPicPath"
	UploadBookChapterPath = "uploadBookChapterPath"
	UploadBookTextPath    = "uploadBookTextPath"
	Rank                  = "rank"
	ComicList             = "comicList"
	DefaultPic            = "/general/nocover.jpg" //默认封面路径
	REPLACEFOLDER         = "/data/pic"
	LOCALUPLOAD           = "/data/upload/"                     //本地路径替换
	REPLACEAPK            = "/www/wwwroot/down.mnjkfup.cn"      //需要替换的下载显示域名
	UPLOADAPK             = "/www/wwwroot/down.mnjkfup.cn/apk/" //上传的apk路径信息12
	EncryptVer0           = 0                                   //笔趣阁加密默认值
	EncryptVer1           = 1                                   //笔趣阁解密默认值
	AdminHostIp           = "103.36.91.36"                      //web机器的IP地址
)

var (
	HighScoreBookPage = make(map[string]int)
)

const (
	Pay_Cion_Name      = "金币"
	Author_Fc_Book     = 10
	Pay_Rmb_Cion       = 10
	Mc_Book_Key        = "Cbj2s3kV5hHzBD7my8oE"
	Pl_Time            = 30
	Pl_Add_Num         = 30
	User_Reg_Vip_Day   = 0
	CollectThreadCount = 5
	SleepSecond        = 0
)

var (
	UserList = make(map[string]int64)
	NodeList = make(map[string]*melody.Session)
)

const (
	ApiAesKey = "WB0nMZHXlxNndORe"
)
const (
	Zero    = "0x00"
	Zone    = "0x01"
	Ztwo    = "0x02"
	Zthree  = "0x03"
	Zfour   = "0x04"
	Zfive   = "0x05"
	Zsix    = "0x06"
	Zseven  = "0x07"
	Zeight  = "0x08"
	Znine   = "0x09"
	Zten    = "0x10"
	Zeleven = "0x11"
	ApiWs   = "/api/message/ws"
)

const (
	//1 专门用来测试的 22%
	//2 201 支付宝超级快手 26%
	//3 271 支付宝YY 23%
	//4 333 微信扫码 20%
	Appid           = "M1701969066"
	AppSecret       = "34686d1c044c4fb4996d3c0fae0131d0"
	ChannelCode     = 201
	AliChannelCode  = 271
	WxChannelCode   = 333
	UnifiedOrderUrl = "https://lantianpayq8lw8h.zzbbm.xyz/api/pay/unifiedorder"
	QueryOrderUrl   = "https://lantianpayq8lw8h.zzbbm.xyz/api/pay/query"
	ReturnUrl       = "http://103.36.90.182:8005/api/order/returnUrl"
	NotifyUrl       = "http://103.36.90.182:8005/api/order/notifyUrl"
)

const (
	//极光推送
	JKey    = "17a6b6188486a31196873ff8"
	Jsecret = "6a6263dcd3f2c252bd36821e"
)

const (
	NsqProducerIP           = "127.0.0.1:4150"
	NsqConsumerIP           = "127.0.0.1:4161"
	Default                 = "default"
	ChapterText             = "chapter_text"
	UpdateBook              = "update_book"
	UpdateChapter           = "update_chapter"
	UpdateChapterText       = "update_chapter_text"
	SourceUpdateLastChapter = "source_update_last_chapter"
	UpdateComic             = "update_comic"
)

var (
	//开启s5代理
	IsS5 = false
	//S5Type = "rank"
	S5Type = "rank"
	//Token1 = "56edbb1f-6b97-4897-9006-751b78b6e085"
	//Token2     = "061eb411-c713-4e1a-8a37-20e885ca1e50"
	//Token3     = "4a6d8c06-78ce-4205-b4cc-d2b4b361a942"
	//S5RankUrl2         = "https://tj.xiaobaibox.com/goldprod/ippool/list?&country=CN&loop=1&postal=110000"
	S5RankUrl          = "http://webapi.http.zhimacangku.com/getip?neek=321a408a&num=1&type=2&pro=0&city=0&yys=0&port=2&pack=341942&ts=1&ys=1&cs=1&lb=1&sb=&pb=4&mr=1&regions=110000,140000,210000,220000,230000,360000,370000,410000,510000,520000,500000,340000,440000,610000,130000,320000,430000,420000,120000,310000,460000,530000,350000,620000"
	S5RankUrl2         = "http://api.yilian.top/v2/proxy/proxies?token=Lc7Qk7BnqQe9DFIITKwpesxysbr3xxyX&pull_num=1&format=json&protocol=1&separator=1&auto_shield=true"
	S5Username         = ""
	S5Passwd           = ""
	S5Domain           = ""
	S5Port             = ""
	S5ExpireTime int64 = 0
	S5Proxys     []models.Socket5Proxy
	S5ProxyIndex int = 0
)

const (
	BucketName        = "tutututu"
	BucketDomain      = "https://t.shunfengs.com"
	ACCOUNT_ID        = "93d07f3af38ec5fc90eca35db0e8c20b"
	ACCESS_KEY_ID     = "7d7ee826341b3e361adaf70e16b2bfcd"
	ACCESS_KEY_SECRET = "28c9df055d5ae870cdf50bb7cab8bd137095c2bc659b9d15701d354c3a0a092f"
	ImgEncry          = byte(136)
)
