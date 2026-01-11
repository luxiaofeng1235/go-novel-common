package models

type McRules struct {
	Id       int64  `gorm:"id" json:"id"`
	RuleName string `gorm:"rule_name" json:"rule_name"`    // 采集名称
	RuleCode string `gorm:"rule_code" json:"rule_code"`    // 采集标识
	Addtime  int64  `gorm:"column:addtime" json:"addtime"` // 创建时间
	Uptime   int64  `gorm:"column:uptime" json:"uptime"`   // 更新时间
}

func (*McRules) TableName() string {
	return "mc_rules"
}
