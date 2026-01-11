package models

type McClassType struct {
	Id       int64  `gorm:"column:id" json:"id"`
	TypeName string `gorm:"column:type_name" json:"type_name"` // 分类名称
	Sort     int    `gorm:"column:sort" json:"sort"`           // 排序ID
	Status   int    `gorm:"column:status" json:"status"`       // 分类状态,0-不展示,1-展示
}

func (*McClassType) TableName() string {
	return "mc_class_type"
}

type ClassTypeRes struct {
	Id           int64           `form:"id" json:"id"`
	TypeName     string          `form:"type_name" json:"type_name"`
	Sort         int             `form:"sort" json:"sort"`
	ClassListRes []*ClassListRes `form:"class_list" json:"class_list"`
}

type ClassListRes struct {
	ClassId   int64  `form:"class_id" json:"class_id"`
	ClassName string `form:"class_name" json:"class_name"`
	Count     int64  `form:"count" json:"count,omitempty"`
	Pic       string `form:"pic" json:"pic,omitempty"`
}

type TypeListReq struct {
	TypeName string `form:"type_name" json:"type_name"`
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

type CreateTypeReq struct {
	TypeName string `form:"type_name"  json:"type_name"`
	Sort     int    `form:"sort" json:"sort"`
	Status   int    `form:"status" json:"status"`
}

type UpdateTypeReq struct {
	TypeId int64 `form:"id"  json:"id"`
	CreateTypeReq
}

type DeleteTypeReq struct {
	TypeIds []int64 `json:"ids" form:"ids"`
}
