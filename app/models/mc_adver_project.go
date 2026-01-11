package models

type McAdverProject struct {
	Id               int64  `gorm:"column:id" json:"id"`                           //自增主键ID
	AdverType        int    `gorm:"column:adver_type" json:"adver_type"`           //广告类型 和广告那边的对应
	PackageId        int64  `gorm:"column:package_id" json:"package_id"`           //关联包ID
	AdverTypeName    string `gorm:"column:adver_type_name" json:"adver_type_name"` //广告类型名称，方便进行查询显示利用
	AdverValueString string `gorm:"column:adver_value_string" json:"adver_value_string"`
	AddTime          int64  `gorm:"column:addtime" json:"addtime"` //添加时间
	Uptime           int64  `gorm:"column:uptime" json:"uptime"`   //更新时间
}

func (McAdverProject) TableName() string {
	return "mc_adver_project"
}

// 搜索包管理需要用到的参数
type AdverProjectReq struct {
	AdverType        string `form:"adver_type"  json:"adver_type"`
	AdverTypeName    string `form:"adver_type_name" json:"adver_type_name"`
	PackageId        int64  `form:"package_id" json:"package_id"`
	AdverValueString string `form:"adver_value_string" json:"adver_value_string"`
	PageNum          int    `form:"pageNum" json:"pageNum"`
	PageSize         int    `form:"pageSize" json:"pageSize"`
}

// 创建广告包的管理名称
type CreateAdverProjectReq struct {
	AdverType        string `form:"adver_type"  json:"adver_type"`
	PackageId        string `form:"package_id" json:"package_id"`
	AdverValueString string `form:"adver_value_string" json:"adver_value_string"`
}

// 更新包需要提交的数据
type UpdateAdverProjectReq struct {
	ProjectId int64 `form:"id"  json:"id"`
	CreateAdverProjectReq
}
