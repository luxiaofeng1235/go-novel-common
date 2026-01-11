package models

type McBookFeedback struct {
	Id           int64  `gorm:"column:id" json:"id"`
	Uid          int64  `gorm:"column:uid" json:"uid"`                     // 反馈人UID
	Ip           string `gorm:"column:ip" json:"ip"`                       // 反馈人IP
	Username     string `gorm:"column:username" json:"username"`           // 用户名
	BookName     string `gorm:"column:book_name" json:"book_name"`         // 小说名称
	Author       string `gorm:"column:author" json:"author"`               // 小说作者
	ChapterName  string `gorm:"column:chapter_name" json:"chapter_name"`   // 章节名称
	Status       int    `gorm:"column:status" json:"status"`               // 反馈处理状态 0-未处理 1-处理中 2-已处理
	FeedbackType int    `gorm:"column:feedback_type" json:"feedback_type"` // 1-错字漏字 2-内容排版错乱 3-章节顺序错乱 4-章节缺失或重复 5-更新延迟或断更
	Bid          int64  `gorm:"column:bid" json:"bid"`                     // 小说id
	Cid          int64  `gorm:"column:cid" json:"cid"`                     // 章节id
	Handtime     int64  `gorm:"column:handtime" json:"handtime"`           // 处理时间
	Addtime      int64  `gorm:"column:addtime" json:"addtime"`             // 添加时间
}

func (*McBookFeedback) TableName() string {
	return "mc_book_feedback"
}

type FeedbackBookListSearchReq struct {
	BookName     string `form:"book_name"  json:"book_name"`
	Author       string `form:"author"  json:"author"`
	ChapterName  string `form:"chapter_name"  json:"chapter_name"`
	UserId       string `form:"user_id"  json:"user_id"`
	Status       string `form:"status"  json:"status"`
	FeedbackType string `form:"feedback_type"  json:"feedback_type"`
	BeginTime    string `form:"beginTime" json:"beginTime"`
	EndTime      string `form:"endTime" json:"endTime"`
	PageNum      int    `form:"pageNum" json:"pageNum"`
	PageSize     int    `form:"pageSize" json:"pageSize"`
}

type UpdateFeedbackBookReq struct {
	FeedbackBookId int64  `form:"id" json:"id"`
	BookName       string `form:"book_name"  json:"book_name"`
	Author         string `form:"author"  json:"author"`
	ChapterName    string `form:"chapter_name"  json:"chapter_name"`
	Status         int    `form:"status" json:"status"`
	FeedbackType   int    `form:"feedback_type" json:"feedback_type"`
}

type DelFeedbackBookReq struct {
	FeedbackBookIds []int64 `json:"ids" form:"ids"`
}
