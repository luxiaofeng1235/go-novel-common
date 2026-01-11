package models

type SysAuthRule struct {
	Id         int64          `gorm:"column:id" json:"id"`
	Pid        int64          `gorm:"column:pid" json:"pid"`                 // 父ID
	Name       string         `gorm:"column:name" json:"name"`               // 规则名称
	Title      string         `gorm:"column:title" json:"title"`             // 规则名称
	Icon       string         `gorm:"column:icon" json:"icon"`               // 图标
	Condition  string         `gorm:"column:condition" json:"condition"`     // 条件
	Remark     string         `gorm:"column:remark" json:"remark"`           // 备注
	MenuType   int            `gorm:"column:menu_type" json:"menu_type"`     // 类型 0目录 1菜单 2按钮
	CreateTime int64          `gorm:"column:create_time" json:"create_time"` // 创建时间
	UpdateTime int64          `gorm:"column:update_time" json:"update_time"` // 更新时间
	Weigh      int            `gorm:"column:weigh" json:"weigh"`             // 权重
	Status     int            `gorm:"column:status" json:"status"`           // 状态
	AlwaysShow int            `gorm:"column:always_show" json:"always_show"` // 显示状态
	Path       string         `gorm:"column:path" json:"path"`               // 路由地址
	Component  string         `gorm:"column:component" json:"component"`     // 组件路径
	IsFrame    int            `gorm:"column:is_frame" json:"is_frame"`       // 是否外链 1是 0否
	Children   []*SysAuthRule `gorm:"-" json:"children,omitty"`              // 子菜单集合
}

func (*SysAuthRule) TableName() string {
	return "sys_auth_rule"
}

type SysAuthRuleReqSearch struct {
	Status string `json:"status" `
	Title  string `json:"menuName" `
}

// 获取管理员列表结构体
type MenuListReq struct {
	Name     string `form:"name" json:"name"`
	Title    string `form:"title" json:"title"`
	MenuType int    `form:"menu_type" json:"menu_type"`
	Status   string `form:"status" json:"status" `
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

// 创建菜单结构体
type CreateMenuReq struct {
	Pid        int64  `form:"pid" json:"pid"`                 // 父ID
	Name       string `form:"name" json:"name"`               // 规则名称
	Title      string `form:"title" json:"title"`             // 规则名称
	Icon       string `form:"icon" json:"icon"`               // 图标
	Condition  string `form:"condition" json:"condition"`     // 条件
	MenuType   int    `form:"menu_type" json:"menu_type"`     // 类型 0目录 1菜单 2按钮
	Weigh      int    `form:"weigh" json:"weigh"`             // 权重
	Status     int    `form:"status" json:"status"`           // 状态
	AlwaysShow int    `form:"always_show" json:"always_show"` // 显示状态
	Path       string `form:"path" json:"path"`               // 路由地址
	Component  string `form:"component" json:"component"`     // 组件路径
	IsFrame    int    `form:"is_frame" json:"is_frame"`       // 是否外链 1是 0否
}

type UpdateMenuReq struct {
	MenuId int64 `json:"menuId"`
	CreateMenuReq
}

// 批量删除菜单结构体
type DeleteMenuReq struct {
	MenuIds []int64 ` form:"menuIds" json:"menuIds"`
}
