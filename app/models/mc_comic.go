package models

type McComic struct {
	Id                int64   `gorm:"id" json:"id"`
	SourceId          int64   `gorm:"source_id" json:"source_id"`                     // 来源ID
	Name              string  `gorm:"name" json:"name"`                               // 标题
	Yname             string  `gorm:"yname" json:"yname"`                             // 英文别名
	Pic               string  `gorm:"pic" json:"pic"`                                 // 竖版封面
	Picx              string  `gorm:"picx" json:"picx"`                               // 横版封面
	Cid               int64   `gorm:"cid" json:"cid"`                                 // 分类ID
	Tid               int     `gorm:"tid" json:"tid"`                                 // 1推荐，0未推
	Ttid              int     `gorm:"ttid" json:"ttid"`                               // 人工推荐(banner),0:不推荐,1:推荐
	Serialize         string  `gorm:"serialize" json:"serialize"`                     // 状态
	Author            string  `gorm:"author" json:"author"`                           // 漫画作者
	Uid               int64   `gorm:"uid" json:"uid"`                                 // 用户ID
	Notice            string  `gorm:"notice" json:"notice"`                           // 公告
	PicAuthor         string  `gorm:"pic_author" json:"pic_author"`                   // 图作者
	TxtAuthor         string  `gorm:"txt_author" json:"txt_author"`                   // 文作者
	Text              string  `gorm:"text" json:"text"`                               // 一句话简介
	Content           string  `gorm:"content" json:"content"`                         // 介绍
	Hits              int64   `gorm:"hits" json:"hits"`                               // 总点击
	Yhits             int64   `gorm:"yhits" json:"yhits"`                             // 月点击
	Zhits             int64   `gorm:"zhits" json:"zhits"`                             // 周点击
	Rhits             int64   `gorm:"rhits" json:"rhits"`                             // 日点击
	HitsUptime        int64   `gorm:"hits_uptime" json:"hits_uptime"`                 // 统计更新时间
	Shits             int     `gorm:"shits" json:"shits"`                             // 收藏人气
	Pay               int     `gorm:"pay" json:"pay"`                                 // 是否收费1金币，2VIP
	Cion              int     `gorm:"cion" json:"cion"`                               // 打赏总额
	Ticket            int     `gorm:"ticket" json:"ticket"`                           // 月票总额
	Sid               int     `gorm:"sid" json:"sid"`                                 // 0正常1锁定
	Nums              int     `gorm:"nums" json:"nums"`                               // 章节总数
	LatestChapterId   int     `gorm:"latest_chapter_id" json:"latest_chapter_id"`     // 最新章节ID
	LatestChapterName string  `gorm:"latest_chapter_name" json:"latest_chapter_name"` // 最新章节名称
	Score             float64 `gorm:"score" json:"score"`                             // 总得分
	Did               int     `gorm:"did" json:"did"`                                 // 采集资源ID
	Ly                string  `gorm:"ly" json:"ly"`                                   // 采集来源标识
	Yid               int     `gorm:"yid" json:"yid"`                                 // 0正常，1待审核
	Msg               string  `gorm:"msg" json:"msg"`                                 // 未审核原因
	Addtime           int64   `gorm:"addtime" json:"addtime"`                         // 入库时间
	Uptime            int64   `gorm:"uptime" json:"uptime"`                           // 更新时间
}

func (*McComic) TableName() string {
	return "mc_comic"
}

type Comic struct {
	ComicName string `form:"comic_name" json:"comic_name"`
	Author    string `form:"author" json:"author"`
	ComicHref string `form:"comic_href" json:"comic_href"`
	Pic       string `form:"pic" json:"pic"`
	Desc      string `form:"desc" json:"desc"`
}

type ComicChapter struct {
	ChapterName string `form:"chapter_name" json:"chapter_name"`
	ChapterHref string `form:"chapter_href" json:"chapter_href"`
}

type UpLoadPicCiphertextRes struct {
	Code int    `form:"code"    json:"code"`
	Msg  string `form:"msg"    json:"msg"`
	Data string `form:"data"    json:"data"`
}

type UpLoadPicRes struct {
	Code int    `form:"code"    json:"code"`
	Msg  string `form:"msg"    json:"msg"`
	Url  string `form:"url"    json:"url"`
	Img  string `form:"img"    json:"img"`
}

type DirPics struct {
	DirPath string   `form:"dir_path"    json:"dir_path"`
	DirName string   `form:"dir_name"    json:"dir_name"`
	DirPics []string `form:"dir_pics"    json:"dir_pics"`
}
