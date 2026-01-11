package models

type McCollect struct {
	Id                int64  `gorm:"column:id" json:"id"`
	Title             string `gorm:"column:title" json:"title"`
	Link              string `gorm:"column:link" json:"link"`
	Charset           string `gorm:"column:charset" json:"charset"`                         // 网站编码
	UrlComplete       int    `gorm:"column:url_complete" json:"url_complete"`               // 章节网址补全
	UrlReverse        int    `gorm:"column:url_reverse" json:"url_reverse"`                 // 倒序采集
	PicLocal          int    `gorm:"column:pic_local" json:"pic_local"`                     // 图片本地化 1-是 0-否
	Categorys         string `gorm:"column:categorys" json:"categorys"`                     // 栏目转换
	CategoryWay       int    `gorm:"column:category_way" json:"category_way"`               // 入库方式 0-对应分类 1-固定分类
	CategoryFixed     int64  `gorm:"column:category_fixed" json:"category_fixed"`           // 固定分类
	ListPageReg       string `gorm:"column:list_page_reg" json:"list_page_reg"`             // 列表地址
	ListSectionReg    string `gorm:"column:list_section_reg" json:"list_section_reg"`       // 列表区间正则 获取列表页盒子
	ListUrlReg        string `gorm:"column:list_url_reg" json:"list_url_reg"`               // 网址正则  获取列表页小说标题和链接
	ChapterSectionReg string `gorm:"column:chapter_section_reg" json:"chapter_section_reg"` // 章节区间正则 获取详情页章节盒子
	ChapterUrlReg     string `gorm:"column:chapter_url_reg" json:"chapter_url_reg"`         // 网址正则  获取详情页章节标题和链接
	ChapterTextReg    string `gorm:"column:chapter_text_reg" json:"chapter_text_reg"`       // 网址正则  获取详情内容
	CategoryNameReg   string `gorm:"column:category_name_reg" json:"category_name_reg"`     // 获取分类名称正则
	BookNameReg       string `gorm:"column:book_name_reg" json:"book_name_reg"`             // 获取小说名称正则
	DescReg           string `gorm:"column:desc_reg" json:"desc_reg"`                       // 获取小说简介正则
	PicReg            string `gorm:"column:pic_reg" json:"pic_reg"`                         // 获取小说图片正则
	AuthorReg         string `gorm:"column:author_reg" json:"author_reg"`                   // 获取作者正则
	Status            int    `gorm:"column:status" json:"status"`                           // 数据状态
	SerializeReg      string `gorm:"column:serialize_reg" json:"serialize_reg"`             // 获取小说正则
	TagNameReg        string `gorm:"column:tag_name_reg" json:"tag_name_reg"`               // 获取标签正则
	UpdateReg         string `gorm:"column:update_reg" json:"update_reg"`                   // 最新时间
	DescReplaceReg    string `gorm:"column:desc_replace_reg" json:"desc_replace_reg"`       // 简介替换规则
	TextReplaceReg    string `gorm:"column:text_replace_reg" json:"text_replace_reg"`       // 内容替换规则
	CollectTime       int64  `gorm:"column:collect_time" json:"collect_time"`               // 采集时间
	ChapterMode       string `gorm:"column:chapter_mode" json:"chapter_mode"`               // 获取章节所有章节 app换源使用
	Addtime           int64  `gorm:"column:addtime" json:"addtime"`                         // 创建时间
	Uptime            int64  `gorm:"column:uptime" json:"uptime"`                           // 更新时间
}

func (*McCollect) TableName() string {
	return "mc_collect"
}

type CollectRes struct {
	PageNum int               `form:"pageNum" json:"pageNum"`
	Count   int               `form:"count" json:"count"`
	Data    []*CollectDataRes `form:"data" json:"data"`
}

type CollectDataRes struct {
	Url  string `form:"url" json:"url"`
	Lock int    `form:"lock" json:"lock"`
}

type CollectListData struct {
	Url     string `form:"url" json:"url"`
	Lock    int    `form:"lock" json:"lock"`
	Count   int    `form:"count" json:"count"`
	PageNum int    `form:"pageNum" json:"pageNum"`
}

type CollecBookInfoRes struct {
	BookUrl      string                `form:"book_url" json:"book_url"`
	CategoryName string                `form:"category_name" json:"category_name"`
	BookName     string                `form:"book_name" json:"book_name"`
	ClassId      int64                 `form:"class_id" json:"class_id"`
	SourceUrl    string                `form:"source_url" json:"source_url"`
	Author       string                `form:"author" json:"author"`
	Serialize    string                `form:"serialize" json:"serialize"`
	Pic          string                `form:"pic" json:"pic"`
	Desc         string                `form:"desc" json:"desc"`
	TagName      string                `form:"tag_name" json:"tag_name"`
	UpdateTime   string                `form:"update_time" json:"update_time"`
	Chapters     []*CollectChapterInfo `form:"chapters" json:"chapters"`
}

type CollectChapterInfo struct {
	ChapterTitle string `form:"chapter_title" json:"chapter_title"`
	ChapterLink  string `form:"chapter_link" json:"chapter_link"`
}

type CategoryReg struct {
	Target string `form:"target" json:"target"`
	Local  int64  `form:"local" json:"local"`
}

type TextReplace struct {
	Find    string `form:"find" json:"find"`
	Replace string `form:"replace" json:"replace"`
}

type PageLinkReq struct {
	Collect  *McCollect `form:"collect"  json:"collect"`
	PageLink string     `form:"page_link"  json:"page_link"`
}

type BookChapterReq struct {
	Collect  *McCollect `form:"collect"  json:"collect"`
	BookLink string     `form:"book_link"  json:"book_link"`
}

type ChapterTextReq struct {
	Collect      *McCollect `form:"collect"  json:"collect"`
	BookId       int64      `form:"book_id" json:"book_id"`
	BookName     string     `form:"book_name" json:"book_name"`
	Author       string     `form:"author" json:"author"`
	ChapterTitle string     `form:"chapter_title" json:"chapter_title"`
	ChapterLink  string     `form:"chapter_link" json:"chapter_link"`
}

type CollectPageUrlRes struct {
	PageNum   int      `form:"pageNum" json:"pageNum"`
	PageCount int      `form:"pageCount" json:"pageCount"`
	PageUrls  []string `form:"pageUrls" json:"pageUrls"`
}

type CollectBookUrl struct {
	Url  string `form:"url" json:"url"`
	Lock int    `form:"lock" json:"lock"`
}

type CollectPageBook struct {
	PageNum   int               `form:"pageNum" json:"pageNum"`
	PageCount int               `form:"pageCount" json:"pageCount"`
	BookCount int               `form:"bookCount" json:"bookCount"`
	BookUrls  []*CollectBookUrl `form:"bookUrls" json:"bookUrls"`
}

type CollectPageBookChapter struct {
	Collect  *McCollect            `form:"collect"  json:"collect"`
	BookName string                `form:"book_name" json:"book_name"`
	Author   string                `form:"author" json:"author"`
	Chapters []*CollectChapterInfo `form:"chapters" json:"chapters"`
}

type CollectListReq struct {
	Title    string `form:"title" json:"title"`
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

type CreateCollectReq struct {
	Title             string            `form:"title" json:"title"`
	Link              string            `form:"link" json:"link"`
	Charset           string            `form:"charset" json:"charset"`
	UrlComplete       int               `form:"url_complete" json:"url_complete"`
	UrlReverse        int               `form:"url_reverse" json:"url_reverse"`
	PicLocal          int               `form:"pic_local" json:"pic_local"`
	Categorys         string            `form:"categorys" json:"categorys"`
	CategoryFixed     int64             `form:"category_fixed" json:"category_fixed"`
	ListPageReg       string            `form:"list_page_reg" json:"list_page_reg"`
	ListSectionReg    string            `form:"list_section_reg" json:"list_section_reg"`
	ListUrlReg        string            `form:"list_url_reg" json:"list_url_reg"`
	ChapterSectionReg string            `form:"chapter_section_reg" json:"chapter_section_reg"`
	ChapterUrlReg     string            `form:"chapter_url_reg" json:"chapter_url_reg"`
	ChapterTextReg    string            `form:"chapter_text_reg" json:"chapter_text_reg"`
	CategoryNameReg   string            `form:"category_name_reg" json:"category_name_reg"`
	BookNameReg       string            `form:"book_name_reg" json:"book_name_reg"`
	DescReg           string            `form:"desc_reg" json:"desc_reg"`
	PicReg            string            `form:"pic_reg" json:"pic_reg"`
	AuthorReg         string            `form:"author_reg" json:"author_reg"`
	Status            int               `form:"status" json:"status"`
	CategoryWay       int               `form:"category_way" json:"category_way"`
	SerializeReg      string            `form:"serialize_reg" json:"serialize_reg"`
	UpdateReg         string            `form:"update_reg" json:"update_reg"`
	TagNameReg        string            `form:"tag_name_reg" json:"tag_name_reg"`
	DescReplaceReg    string            `form:"desc_replace_reg" json:"desc_replace_reg"`
	TextReplaceReg    string            `form:"text_replace_reg" json:"text_replace_reg"`
	ListPageArr       []*ListPageReg    `form:"listPageArr"  json:"listPageArr"`
	CategoryArr       []*CategoryReg    `form:"categoryArr"  json:"categoryArr"`
	DescReplaceArr    []*TextReplaceReg `form:"descReplaceArr" json:"descReplaceArr"`
	TextReplaceArr    []*TextReplaceReg `form:"textReplaceArr" json:"textReplaceArr"`
}

type UpdateCollectReq struct {
	CollectId int64 `form:"id"  json:"id"`
	CreateCollectReq
}

type DeleteCollectReq struct {
	CollectIds []int64 `form:"ids" json:"ids"`
}

type GetCollectRes struct {
	Id                int64             `gorm:"id" json:"id"`
	Title             string            `form:"title" json:"title"`
	Link              string            `form:"link" json:"link"`
	Charset           string            `form:"charset" json:"charset"`
	UrlComplete       int               `form:"url_complete" json:"url_complete"`
	UrlReverse        int               `form:"url_reverse" json:"url_reverse"`
	PicLocal          int               `form:"pic_local" json:"pic_local"`
	CategoryWay       int               `form:"category_way" json:"category_way"`
	CategoryFixed     int64             `form:"category_fixed" json:"category_fixed"`
	ListSectionReg    string            `form:"list_section_reg" json:"list_section_reg"`
	ListUrlReg        string            `form:"list_url_reg" json:"list_url_reg"`
	ChapterSectionReg string            `form:"chapter_section_reg" json:"chapter_section_reg"`
	ChapterUrlReg     string            `form:"chapter_url_reg" json:"chapter_url_reg"`
	ChapterTextReg    string            `form:"chapter_text_reg" json:"chapter_text_reg"`
	CategoryNameReg   string            `form:"category_name_reg" json:"category_name_reg"`
	BookNameReg       string            `form:"book_name_reg" json:"book_name_reg"`
	DescReg           string            `form:"desc_reg" json:"desc_reg"`
	PicReg            string            `form:"pic_reg" json:"pic_reg"`
	AuthorReg         string            `form:"author_reg" json:"author_reg"`
	Status            int               `form:"status" json:"status"`
	SerializeReg      string            `form:"serialize_reg" json:"serialize_reg"`
	TagNameReg        string            `form:"tag_name_reg" json:"tag_name_reg"`
	UpdateReg         string            `form:"update_reg" json:"update_reg"`
	CategoryArr       []*CategoryReg    `form:"categoryArr" json:"categoryArr"`
	ListPageArr       []*ListPageReg    `form:"listPageArr" json:"listPageArr"`
	DescReplaceArr    []*TextReplaceReg `form:"descReplaceArr" json:"descReplaceArr"`
	TextReplaceArr    []*TextReplaceReg `form:"textReplaceArr" json:"textReplaceArr"`
}

type ListPageReg struct {
	Url       string `form:"url" json:"url"`
	Type      int    `form:"type" json:"type"`
	PageStart int    `form:"pageStart" json:"pageStart"`
	PageEnd   int    `form:"pageEnd" json:"pageEnd"`
	PageInc   int    `form:"pageInc" json:"pageInc"`
	PageDesc  bool   `form:"pageDesc" json:"pageDesc"`
}

type TextReplaceReg struct {
	Find     string `form:"find" json:"find"`
	Replaces string `form:"replaces" json:"replaces"`
}

type NsqCollectBookPush struct {
	BookName           string `form:"book_name" json:"book_name"`
	Author             string `form:"author" json:"author"`
	Pic                string `form:"pic" json:"pic"`
	Desc               string `form:"desc" json:"desc"`
	ClassId            int64  `form:"class_id" json:"class_id"`
	CategoryName       string `form:"category_name" json:"category_name"`
	Tags               string `form:"tags" json:"tags"`
	Serialize          int    `form:"serialize" json:"serialize"`
	TextNum            int    `form:"text_num" json:"text_num"`
	ChapterNum         int    `form:"chapter_num" json:"chapter_num"`
	SourceId           int64  `form:"source_id" json:"source_id"`
	SourceUrl          string `form:"source_url" json:"source_url"`
	BookType           int    `form:"book_type" json:"book_type"`
	IsClassic          int    `form:"is_classic" json:"is_classic"`
	LastChapterTitle   string `form:"last_chapter_title" json:"last_chapter_title"`
	LastChapterTime    string `form:"last_chapter_time" json:"last_chapter_time"`
	UpdateChapterId    int64  `form:"update_chapter_id" json:"update_chapter_id"`
	UpdateChapterTitle string `form:"update_chapter_title" json:"update_chapter_title"`
	UpdateChapterTime  int64  `form:"update_chapter_time" json:"update_chapter_time"`
}

type Socket5Res struct {
	Code int `form:"code" json:"code"`
	Data struct {
		List []*Socket5Proxy `form:"list" json:"list"`
	} `form:"data" json:"data"`
}

type Socket5Proxy struct {
	Ip         string `form:"ip" json:"ip"`
	Port       string `form:"port" json:"port"`
	Username   string `form:"username" json:"username"`
	Password   string `form:"password" json:"password"`
	OnlineDate string `form:"online_date" json:"online_date"`
	Node       string `form:"node" json:"node"`
}

type ZhimaSocket5Res struct {
	Code int                  `form:"code" json:"code"`
	Msg  string               `form:"msg" json:"msg"`
	Data []*ZhimaSocket5Proxy `form:"data" json:"data"`
}

type ZhimaSocket5Proxy struct {
	Ip         string `form:"ip" json:"ip"`
	Port       int    `form:"port" json:"port"`
	City       string `form:"city" json:"city"`
	Isp        string `form:"isp" json:"isp"`
	ExpireTime string `form:"expire_time" json:"expire_time"`
}

type YilianSocket5Res struct {
	Errcode int                   `form:"errcode" json:"errcode"`
	Errmsg  string                `form:"errmsg" json:"errmsg"`
	Data    []*YilianSocket5Proxy `form:"data" json:"data"`
}

type YilianSocket5Proxy struct {
	Ip              string `form:"ip" json:"ip"`
	Username        string `form:"proxy_user" json:"proxy_user"`
	Passwd          string `form:"proxy_pass" json:"proxy_pass"`
	Port            int    `form:"proxy_port" json:"proxy_port"`
	ExpireTime      string `form:"expire_time" json:"expire_time"`
	ExpireTimestamp int64  `form:"expire_timestamp" json:"expire_timestamp"`
}

type XswCategory struct {
	CategoryName string `form:"category_name" json:"category_name"`
	CategoryKey  string `form:"category_key" json:"category_key"`
	CategoryHref string `form:"category_href" json:"category_href"`
}

type XswPageBooks struct {
	PageLink     string         `form:"page_link" json:"page_link"`
	NextPageLink string         `form:"next_page_link" json:"next_page_link"`
	TotalPage    string         `form:"total_page" json:"total_page"`
	Books        []*XswPageBook `form:"books" json:"books"`
}

type XswPageBook struct {
	BookName string `form:"book_name" json:"book_name"`
	Author   string `form:"author"   json:"author"`
	Pic      string `form:"pic" json:"pic"`
	Link     string `form:"link" json:"link"`
	Use      int    `form:"use" json:"use"`
}

type XswBookDesc struct {
	BookName         string `form:"book_name" json:"book_name"`
	Author           string `form:"author"   json:"author"`
	Pic              string `form:"pic" json:"pic"`
	Desc             string `form:"desc" json:"desc"`
	Serialize        int    `form:"serialize"   json:"serialize"`
	CategoryName     string `form:"category_name" json:"category_name"`
	ClassId          int64  `form:"class_id" json:"class_id"`
	TextNum          int    `form:"text_num"   json:"text_num"`
	LastChapterTitle string `form:"last_chapter_title" json:"last_chapter_title"`
	LastChapterTime  string `form:"last_chapter_time" json:"last_chapter_time"`
	ChapterNum       int    `form:"chapter_num"   json:"chapter_num"`
	IsClassic        int    `form:"is_classic"   json:"is_classic"`
	SourceUrl        string `form:"source_url" json:"source_url"`
	Use              int    `form:"use" json:"use"`
}

type LydCategory struct {
	CategoryName string `form:"category_name" json:"category_name"`
	CategoryKey  string `form:"category_key" json:"category_key"`
	CategoryHref string `form:"category_href" json:"category_href"`
	Use          int    `form:"use" json:"use"`
}

type LydPageBooks struct {
	PageLink  string         `form:"page_link" json:"page_link"`
	Page      int            `form:"page" json:"page"`
	TotalPage int            `form:"total_page" json:"total_page"`
	BookNum   int            `form:"book_num" json:"book_num"`
	Books     []*LydPageBook `form:"books" json:"books"`
}

type LydPageBook struct {
	BookName string `form:"book_name" json:"book_name"`
	Author   string `form:"author"   json:"author"`
	Pic      string `form:"pic" json:"pic"`
	BookLink string `form:"book_link" json:"book_link"`
	Use      int    `form:"use" json:"use"`
}

type LydBookDesc struct {
	BookName         string `form:"book_name" json:"book_name"`
	Author           string `form:"author"   json:"author"`
	Pic              string `form:"pic" json:"pic"`
	Desc             string `form:"desc" json:"desc"`
	Serialize        int    `form:"serialize"   json:"serialize"`
	CategoryName     string `form:"category_name" json:"category_name"`
	ClassId          int64  `form:"class_id" json:"class_id"`
	TextNum          int    `form:"text_num"   json:"text_num"`
	LastChapterTitle string `form:"last_chapter_title" json:"last_chapter_title"`
	LastChapterTime  string `form:"last_chapter_time" json:"last_chapter_time"`
	ChapterNum       int    `form:"chapter_num"   json:"chapter_num"`
	IsClassic        int    `form:"is_classic"   json:"is_classic"`
	BookLink         string `form:"book_link" json:"book_link"`
	SourceUrl        string `form:"source_url" json:"source_url"`
	Use              int    `form:"use" json:"use"`
}

type Bqg24Category struct {
	CategoryName string `form:"category_name" json:"category_name"`
	CategoryKey  string `form:"category_key" json:"category_key"`
	CategoryHref string `form:"category_href" json:"category_href"`
	Use          int    `form:"use" json:"use"`
}

type Bqg24PageBooks struct {
	PageLink  string           `form:"page_link" json:"page_link"`
	NextLink  string           `form:"next_link" json:"next_link"`
	CateNum   int              `form:"cate_num" json:"cate_num"`
	PageNum   int              `form:"page_num" json:"page_num"`
	TotalPage int              `form:"total_page" json:"total_page"`
	BookNum   int              `form:"book_num" json:"book_num"`
	Books     []*Bqg24PageBook `form:"books" json:"books"`
}

type Bqg24PageBook struct {
	BookName string `form:"book_name" json:"book_name"`
	Author   string `form:"author"   json:"author"`
	BookLink string `form:"book_link" json:"book_link"`
	Use      int    `form:"use" json:"use"`
}

type Bqg24BookDesc struct {
	BookName         string `form:"book_name" json:"book_name"`
	Author           string `form:"author"   json:"author"`
	Pic              string `form:"pic" json:"pic"`
	Desc             string `form:"desc" json:"desc"`
	Serialize        int    `form:"serialize"   json:"serialize"`
	CategoryName     string `form:"category_name" json:"category_name"`
	ClassId          int64  `form:"class_id" json:"class_id"`
	TextNum          int    `form:"text_num"   json:"text_num"`
	LastChapterTitle string `form:"last_chapter_title" json:"last_chapter_title"`
	LastChapterTime  string `form:"last_chapter_time" json:"last_chapter_time"`
	ChapterNum       int    `form:"chapter_num"   json:"chapter_num"`
	IsClassic        int    `form:"is_classic"   json:"is_classic"`
	SourceUrl        string `form:"source_url" json:"source_url"`
	Use              int    `form:"use" json:"use"`
}

type Siluke520Category struct {
	CategoryName string `form:"category_name" json:"category_name"`
	CategoryKey  string `form:"category_key" json:"category_key"`
	CategoryHref string `form:"category_href" json:"category_href"`
	Use          int    `form:"use" json:"use"`
}

type Siluke520PageBooks struct {
	PageLink  string               `form:"page_link" json:"page_link"`
	NextLink  string               `form:"next_link" json:"next_link"`
	CateNum   int                  `form:"cate_num" json:"cate_num"`
	PageNum   int                  `form:"page_num" json:"page_num"`
	TotalPage int                  `form:"total_page" json:"total_page"`
	BookNum   int                  `form:"book_num" json:"book_num"`
	Books     []*Siluke520PageBook `form:"books" json:"books"`
}

type Siluke520PageBook struct {
	BookName string `form:"book_name" json:"book_name"`
	Author   string `form:"author"   json:"author"`
	BookLink string `form:"book_link" json:"book_link"`
	Use      int    `form:"use" json:"use"`
}

type Siluke520BookDesc struct {
	BookName         string `form:"book_name" json:"book_name"`
	Author           string `form:"author"   json:"author"`
	Pic              string `form:"pic" json:"pic"`
	Desc             string `form:"desc" json:"desc"`
	Serialize        int    `form:"serialize"   json:"serialize"`
	CategoryName     string `form:"category_name" json:"category_name"`
	ClassId          int64  `form:"class_id" json:"class_id"`
	TextNum          int    `form:"text_num"   json:"text_num"`
	LastChapterTitle string `form:"last_chapter_title" json:"last_chapter_title"`
	LastChapterTime  string `form:"last_chapter_time" json:"last_chapter_time"`
	ChapterNum       int    `form:"chapter_num"   json:"chapter_num"`
	IsClassic        int    `form:"is_classic"   json:"is_classic"`
	SourceUrl        string `form:"source_url" json:"source_url"`
	Use              int    `form:"use" json:"use"`
}
