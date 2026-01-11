package models

type McAdverPackageDetail struct {
	Id           int64  `gorm:"column:id" json:"id"`                         //自增主键ID
	Pid          int    `gorm:"column:pid" json:"pid"`                       //关联mc_adver_package中的ID
	ThirdValueId string `gorm:"column:third_value_id" json:"third_value_id"` //第三方广告ID
	AddTime      int64  `gorm:"column:addtime" json:"addtime"`               //添加时间
	Uptime       int64  `gorm:"column:uptime" json:"uptime"`                 //更新时间
}

// 自动添加表名称
func (McAdverPackageDetail) TableName() string {
	return "mc_adver_package_detail"
}

// 搜索列表数据
type AdverPackageDetailReq struct {
	ThirdValueId string `form:"project_name"  json:"project_name"`
	Pid          int    `form:"pid" json:"pid"` //通过主键来进行关联搜索
}

// 创建广告的包关联的ID和信息
type CreateAdverPackageDetailReq struct {
	Pid          int    `form:"pid"  json:"pid"`
	ThirdValueId string `form:"third_value_id" json:"third_value_id"`
}

// 更新广告的包关联的ID信息
type UpdateAdverPackageDetailReq struct {
	DetailId int64 `form:"id"  json:"id"`
	CreateAdverPackageDetailReq
}

// 删除包明细中的数据
type DeleteAdverPackageDetailReq struct {
	DetailIds []int64 `json:"ids" form:"ids"`
}
