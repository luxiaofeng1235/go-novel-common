package models

import (
	"go-novel/global"
)

type Config struct {
	Id        int64  `gorm:"column:id" json:"id"`
	Key       string `gorm:"column:key" json:"key"`
	Name      string `gorm:"column:name" json:"name"`
	Value     string `gorm:"column:value" json:"value"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}
type configValue struct {
	Name  string
	Value string
}

func (*Config) TableName() string {
	return "mc_config"
}

// 查询key = "key"的配置
func GetConfigByKey(key string) ([]*Config, error) {
	var config []*Config
	query := global.DB.Model(config)
	query = query.Where("`key` = ?", key)
	err := query.Find(&config).Error
	if err != nil {
		global.Errlog.Errorf("GetConfigBykey error: %v", err)
		return config, err
	}
	return config, nil
}
