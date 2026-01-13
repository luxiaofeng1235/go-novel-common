/*
 * @Descripttion: 配置加载（读取根目录 config.yml，并可选合并 config/upload.yml）
 * @Author: red
 * @Date: 2026-01-13 09:20:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-13 09:20:00
 */
package config

import (
	"bytes"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	initConfig()
}

func initConfig() {
	dir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(strings.Replace(dir, "\\", "/", -1))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// 业务配置拆分：如存在 config/upload.yml，则合并进主配置（不影响线上只用 config.yml 的场景）
	uploadConfigPath := filepath.Join(dir, "config", "upload.yml")
	if b, err := os.ReadFile(uploadConfigPath); err == nil {
		// viper.MergeConfig 使用当前的 config type 解析
		if err := viper.MergeConfig(bytes.NewReader(b)); err != nil {
			log.Fatalln(err)
		}
	}
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}
