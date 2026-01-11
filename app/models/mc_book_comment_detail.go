package models

type McBookCommentDetail struct {
	Id            int64   `gorm:"column:id" json:"id"`                           //评论ID
	BookId        int     `gorm:"column:book_id" json:"book_id"`                 //书籍关联ID
	UserId        int     `gorm:"column:user_id" json:"user_id"`                 //用户ID
	AvtarId       string  `gorm:"column:avtar_id" json:"avtar_id"`               //avtar_id关联用户的唯一标识
	Username      string  `gorm:"column:username" json:"username"`               //用户名
	AvtarUrl      string  `gorm:"column:avtar_url" json:"avtar_url"`             //用户头像
	Score         float32 `gorm:"column:score" json:"score"`                     //评分
	SynUpdateTime string  `gorm:"column:syn_update_time" json:"syn_update_time"` //评论同步时间（网站拉取）
	Content       string  `gorm:"column:content" json:"content"`                 //评论内容
	AddTime       int64   `gorm:"column:add_time" json:"add_time"`               //添加时间
}

func (*McBookCommentDetail) TableName() string {
	return "mc_book_comment_detail"
}

// 请求的评论列表信息的默认请求参数信息
type BookCommentDetailReq struct {
	Id     int    `form:"weight" json:"weight"`
	BookId string `form:"column:book_id" json:"book_id"`
}

type BookCommnetUserReq struct {
	PageSize int `form:"pageSize" json:"pageSize"`
}

// 根据单个用户获取书评
type BookCommentSingeUidReq struct {
	UserId int `form:"user_id" json:"user_id"`
}

// 指定返回的数据信息
type BookCommentListRes struct {
	Id            int64                   `form:"id" json:"id"`
	BookId        int                     `form:"book_id" json:"book_id"`
	UserId        int                     `form:"user_id" json:"user_id"`
	AvtarId       string                  `form:"avtar_id" json:"avtar_id"`
	Username      string                  `form:"username" json:"username"`
	AvtarUrl      string                  `form:"avtar_url" json:"avtar_url"`
	Score         float32                 `form:"score" json:"score"`
	SynUpdateTime string                  `form:"syn_update_time" json:"syn_update_time"`
	Content       string                  `form:"content" json:"content"`
	addTime       string                  `form:"add_time" json:"add_time"`
	BookList      []*BookCommentDetailRes `form:"book_info" json:"book_info"`
}

// 定义书的基本详情信息
type BookCommentDetailRes struct {
	Title     string `form:"column:title" json:"title"`
	Author    string `form:"column:author" json:"author"`
	CoverLogo string `form:"column:cover_logo" json:"cover_logo"`
}

// 用户评论列表数据
type BookCommentUserListRes struct {
	Num      int    `form:"column:num" json:"num"`
	UserId   int    `form:"column:user_id" json:"user_id"`
	Username string `form:"column:username" json:"username"`
	AvtarUrl string `form:"column:avtar_url" json:"avtar_url"`
}
