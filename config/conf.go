package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
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
