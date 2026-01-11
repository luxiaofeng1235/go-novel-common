package models

type McTag struct {
	Id         int64  `gorm:"column:id" json:"id"`
	BookType   int    `gorm:"column:book_type" json:"book_type"`     // 标签类型 1-男生 2-女生
	ColumnType int    `gorm:"column:column_type" json:"column_type"` // 栏目类型 1-推荐 2-热门 3-经典
	TagName    string `gorm:"column:tag_name" json:"tag_name"`       // 标签名称
	IsNew      int    `gorm:"column:is_new" json:"is_new"`           //是否为新书
	Sort       int    `gorm:"column:sort" json:"sort"`               // 排序ID
	Status     int    `gorm:"column:status" json:"status"`           //状态
	BookCount  int64  `gorm:"column:book_count" json:"book_count"`   //小说数量
	Addtime    int64  `gorm:"column:addtime" json:"addtime"`         // 添加时间
	Uptime     int64  `gorm:"column:uptime" json:"uptime"`           // 更新时间
}

func (*McTag) TableName() string {
	return "mc_tag"
}

type TagListReq struct {
	TagName    string `form:"tag_name" json:"tag_name"`
	BookType   string `form:"book_type" json:"book_type"`
	ColumnType string `form:"column_type" json:"column_type"`
	IsNew      string `form:"is_new" json:"is_new"`
	PageNum    int    `form:"pageNum" json:"pageNum"`
	PageSize   int    `form:"pageSize" json:"pageSize"`
}

type CreateTagReq struct {
	BookType   int    `form:"book_type" json:"book_type"`
	ColumnType int    `form:"column_type" json:"column_type"`
	TagName    string `form:"tag_name"  json:"tag_name"`
	Sort       int    `form:"sort" json:"sort"`
	IsNew      int    `form:"is_new" json:"is_new"`
	Status     int    `form:"status" json:"status"`
}

type UpdateTagReq struct {
	TagId int64 `form:"id"  json:"id"`
	CreateTagReq
}

type DeleteTagReq struct {
	TagId int64 `json:"id" form:"id"`
}

type AssignTagReq struct {
	TagIds []int64 `json:"ids" form:"ids"`
}
