package utils

import (
	"go-novel/app/models"
)

const (
	New = "new"

	CollectChapterUrlTemp = "collect_chapter_url_temp" //章节列表

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
	RoundDot              = "・"
	UploadBookChapterPath = "uploadBookChapterPath"
	Rank                  = "rank"
	DefaultPic            = "/general/nocover.jpg" //默认封面路径
	REPLACEFOLDER         = "/data/pic"
)

const (
	DBNAME    = "novel"
	MYSQLUSER = "root"
)

var (
	HighScoreBookPage = make(map[string]int)
)

const (
	CollectThreadCount = 5
)

const (
	ApiAesKey = "WB0nMZHXlxNndORe"
)

const (
	//极光推送
	JKey    = "17a6b6188486a31196873ff8"
	Jsecret = "6a6263dcd3f2c252bd36821e"
)

var (
	//开启s5代理
	IsS5               = false
	S5Type             = "rank"
	S5RankUrl2         = "http://api.yilian.top/v2/proxy/proxies?token=Lc7Qk7BnqQe9DFIITKwpesxysbr3xxyX&pull_num=1&format=json&protocol=1&separator=1&auto_shield=true"
	S5Username         = ""
	S5Passwd           = ""
	S5Domain           = ""
	S5Port             = ""
	S5ExpireTime int64 = 0
	S5Proxys     []models.Socket5Proxy
)

const (
	BucketName        = "tutututu"
	BucketDomain      = "https://t.shunfengs.com"
	ACCOUNT_ID        = "93d07f3af38ec5fc90eca35db0e8c20b"
	ACCESS_KEY_ID     = "7d7ee826341b3e361adaf70e16b2bfcd"
	ACCESS_KEY_SECRET = "28c9df055d5ae870cdf50bb7cab8bd137095c2bc659b9d15701d354c3a0a092f"
	ImgEncry          = byte(136)
)
