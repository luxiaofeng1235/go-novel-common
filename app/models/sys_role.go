package models

type SysRole struct {
	Id         int64   `gorm:"column:id" json:"id"`
	Status     int     `gorm:"column:status" json:"status"`           // 状态;0:禁用;1:正常
	CreateTime int64   `gorm:"column:create_time" json:"create_time"` // 创建时间
	UpdateTime int64   `gorm:"column:update_time" json:"update_time"` // 更新时间
	Sort       float64 `gorm:"column:sort" json:"sort"`               // 排序
	Name       string  `gorm:"column:name" json:"name"`               // 角色名称
	Remark     string  `gorm:"column:remark" json:"remark"`           // 备注
}

func (*SysRole) TableName() string {
	return "sys_role"
}

// 获取角色列表结构体
type RoleListReq struct {
	Name      string `form:"name" json:"name"`
	Status    string `form:"status" json:"status"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	PageNum   uint   `form:"pageNum" json:"pageNum"`
	PageSize  uint   `form:"pageSize" json:"pageSize"`
}

// 创建角色结构体
type CreateRoleReq struct {
	Name    string  `form:"name"    json:"name"`   // 角色名称
	Remark  string  `form:"remark"  json:"remark"` // 备注
	Status  int     `form:"status"  json:"status"` // 状态;0:禁用;1:正常
	Sort    float64 `form:"sort"    json:"sort"`   // 排序
	MenuIds []int   `form:"menuIds" json:"menuIds"`
}
type UpdateRoleReq struct {
	RoleId int64 `json:"roleId"`
	CreateRoleReq
}

// 批量删除角色结构体
type DeleteRoleReq struct {
	RoleIds []int64 `json:"roleIds" form:"roleIds"`
}

// 设置用户状态参数
type RoleStatusReq struct {
	RoleId int64 `form:"roleId" json:"roleId"`
	Status int   `form:"status" json:"status"`
}
