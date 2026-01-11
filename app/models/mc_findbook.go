package models

type McFindbook struct {
	Id         int64  `gorm:"column:id" json:"id"`
	BookName   string `gorm:"column:book_name" json:"book_name"`     // 小说书名
	Author     string `gorm:"column:author" json:"author"`           // 小说作者
	SourceName string `gorm:"column:source_name" json:"source_name"` // 来源网站名称
	Status     int    `gorm:"column:status" json:"status"`           //状态
	Uid        int64  `gorm:"column:uid" json:"uid"`                 // 用户ID
	BookTimes  int    `gorm:"column:book_times" json:"book_times"`   //求书次数
	Addtime    int64  `gorm:"column:addtime" json:"addtime"`         // 添加时间
	Uptime     int64  `gorm:"column:uptime" json:"uptime"`           // 更新时间
}

func (*McFindbook) TableName() string {
	return "mc_findbook"
}

type FindBookListReq struct {
	BookName  string `form:"book_name" json:"book_name"`
	Status    string `form:"status" json:"status"`
	UserId    string `form:"user_id" json:"user_id"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
}

type CreateFindBookReq struct {
	BookName   string `form:"book_name" json:"book_name"`
	Author     string `form:"author" json:"author"`
	SourceName string `form:"source_name" json:"source_name"`
	Status     int    `form:"status" json:"status"`
	BookTimes  int    `form:"book_times" json:"book_times"`
	UserId     int64  `form:"user_id" json:"user_id"`
}

type UpdateFindBookReq struct {
	FindBookId int64 `form:"id" json:"id"`
	CreateFindBookReq
}

type DelFindBookReq struct {
	FindBookIds []int64 `json:"ids" form:"ids"`
}

type ApiFindbookListReq struct {
	Page   int   `form:"page" json:"page"`
	Size   int   `form:"size" json:"size"`
	UserId int64 `form:"user_id" json:"user_id"`
}
