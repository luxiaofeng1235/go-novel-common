package models

type McBookCommnetInfo struct {
	Id           int64   `gorm:"column:id" json:"id"`                       //书评ID
	BookId       string  `gorm:"column:book_id" json:"book_id"`             //网站书籍ID
	Title        string  `gorm:"column:title" json:"title"`                 //标题
	Author       string  `gorm:"column:author" json:"author"`               //作者
	BookUrl      string  `gorm:"column:book_url" json:"book_url"`           //书籍URL
	CommentCount int     `gorm:"column:comment_count" json:"comment_count"` //评分总人数
	CoverLogo    string  `gorm:"column:cover_logo" json:"cover_logo"`       //书籍封面
	Category     string  `gorm:"column:category" json:"category"`           //分类信息
	Score        float32 `gorm:"column:score" json:"score"`                 //书籍评分
	NearyTime    string  `gorm:"column:neary_time" json:"neary_time"`       //书的更新时间（目标源同步过来）
	Addtime      int64   `gorm:"column:addtime" json:"addtime"`             //添加时间
	Uptime       int64   `gorm:"column:uptime" json:"uptime"`               //更新时间
}

func (*McBookCommnetInfo) TableName() string {
	return "mc_book_comment_info"
}

// 请求的书籍列表信息
type BookCommentListReq struct {
	Id       int    `form:"id" json:"id"`
	BookId   int    `form:"book_id" json:"book_id"`
	Title    string `form:"title" json:"title"`
	UserId   int    `form:"user_id" json:"user_id"`
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

// 访问书评信息列表
type BookCommentInfoReq struct {
	BookId int `form:"book_id" json:"book_id"`
}
