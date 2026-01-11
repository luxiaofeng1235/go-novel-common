package models

type McBookClass struct {
	Id        int64  `gorm:"column:id" json:"id"`
	BookType  int    `gorm:"column:book_type" json:"book_type"`   // 类型 1-男生 2-女生 3-图书
	ClassName string `gorm:"column:class_name" json:"class_name"` // 分类名称
	ClassPic  string `gorm:"column:class_pic" json:"class_pic"`   // 分类展示图片
	BookCount int64  `gorm:"column:book_count" json:"book_count"` // 小说数量
	Sort      int    `gorm:"column:sort" json:"sort"`             // 排序ID
	TypeId    int64  `gorm:"column:type_id" json:"type_id"`       // 分类类型ID
	Status    int    `gorm:"column:status" json:"status"`         // 分类状态,0-不展示,1-展示
}

func (*McBookClass) TableName() string {
	return "mc_book_class"
}

type BookTypeReq struct {
	BookType int   `form:"book_type" json:"book_type"`
	UserId   int64 `form:"user_id" json:"user_id"`
}

type ClassListReq struct {
	BookType  string `form:"book_type"  json:"book_type"`
	ClassName string `form:"class_name"  json:"class_name"`
	TypeId    string `form:"type_id"  json:"type_id"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
}

type CreateClassReq struct {
	BookType  int    `form:"book_type"  json:"book_type"`
	ClassName string `form:"class_name"  json:"class_name"`
	ClassPic  string `form:"class_pic" json:"class_pic"`
	BookId    int64  `form:"book_id" json:"book_id"`
	TypeId    int64  `form:"type_id" json:"type_id"`
	Sort      int    `form:"sort" json:"sort"`
	Status    int    `form:"status" json:"status"`
}

type UpdateClassReq struct {
	ClassId int64 `form:"id"  json:"id"`
	CreateClassReq
}

type DeleteClassReq struct {
	ClassIds []int64 `json:"ids" form:"ids"`
}

type AssignClassReq struct {
	ClassIds []int64 `json:"ids" form:"ids"`
}
