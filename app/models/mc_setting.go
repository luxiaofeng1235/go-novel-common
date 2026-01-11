package models

type McSetting struct {
	Id    int64  `gorm:"column:id" json:"id"`       // ID
	Label string `gorm:"column:label" json:"label"` // 配置名称
	Name  string `gorm:"column:name" json:"name"`   // 配置名称
	Value string `gorm:"column:value" json:"value"` // 配置值
}

func (*McSetting) TableName() string {
	return "mc_setting"
}

type GetValueReq struct {
	Name string `form:"name" json:"name"`
}

type GetRanksRes struct {
	Name  string `form:"name" json:"name"`
	Label string `form:"label" json:"label"`
}

type GetVipPriceRes struct {
	Name  string `form:"name" json:"name"`
	Value string `form:"value" json:"value"`
}

// updateReq 用于存储页面更新(新增、修改)网址的信息
type SettingUpdateReq struct {
	WebContent map[string]interface{} `form:"webContent" json:"webContent"` // 站点信息
}

type SettingUpdateOneReq struct {
	Key   string `form:"key" json:"key"`
	Value string `form:"value" json:"value"`
}
