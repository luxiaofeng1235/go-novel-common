/*
 * @Descripttion: 服务启动编排（脚手架最小链路）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 14:05:00
 */
package db

import (
	"flag"
	"fmt"
	"go-novel/config"
	"go-novel/routers/api_routes"
	"go-novel/routers/source_routes"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// 启动后端 -只包含MySQL和source
func StartAdminServer() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_ = config.GetString("server.env")
	if strings.TrimSpace(viper.GetString("jwt.secret")) == "" {
		log.Fatal("缺少 jwt.secret：请在 config.yml 中设置 jwt.secret（用于签发/校验登录 token）")
	}
	host, name, user, passwd := GetDB()
	InitMysql(host, name, user, passwd)
	sourceHost := viper.GetString("source.host")
	sourcePort := viper.GetString("source.port")
	if sourceHost == "" {
		sourceHost = "0.0.0.0"
	}
	if sourcePort == "" {
		sourcePort = "8007"
	}
	go source_routes.InitSourceRoutes(fmt.Sprintf("%s:%s", sourceHost, sourcePort))
}

// StartApiServer 启动 API 服务（脚手架最小启动链路）
func StartApiServer() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	// 触发配置加载（config 包 init 会读取根目录 config.yml）
	_ = config.GetString("server.env")
	if strings.TrimSpace(viper.GetString("jwt.secret")) == "" {
		log.Fatal("缺少 jwt.secret：请在 config.yml 中设置 jwt.secret（用于签发/校验登录 token）")
	}

	var apiHost, apiPort string
	// API 监听地址：统一从 server.host/server.port 读取（api.host/api.port 已废弃）
	defaultHost := strings.TrimSpace(viper.GetString("server.host"))
	defaultPort := viper.GetInt("server.port")
	if defaultHost == "" || defaultPort == 0 {
		log.Fatal("缺少 server.host/server.port：请在 config.yml 的 server 节点配置 API 监听地址（api.host/api.port 已废弃）")
	}
	flag.StringVar(&apiHost, "host", defaultHost, "api listen host")
	flag.StringVar(&apiPort, "port", fmt.Sprintf("%d", defaultPort), "api listen port")
	flag.Parse()

	if apiHost == "" {
		log.Fatal("host 不能为空：请通过 -host 或 config.yml 的 server.host 配置")
	}
	if apiPort == "" {
		log.Fatal("port 不能为空：请通过 -port 或 config.yml 的 server.port 配置")
	}

	sourceHost := viper.GetString("source.host")
	sourcePort := viper.GetString("source.port")
	if sourceHost == "" {
		sourceHost = "0.0.0.0"
	}
	if sourcePort == "" {
		sourcePort = "8007"
	}

	host, name, user, passwd := GetDB()
	InitMysql(host, name, user, passwd)

	redisAddr, redisPasswd, redisDB := GetRedis()
	fmt.Println("redis addr:", redisAddr, "redis db:", redisDB)
	InitRedis(redisAddr, redisPasswd, redisDB)

	// 静态资源服务依赖上传落盘目录（不依赖 DB），但启动顺序上放在 MySQL/Redis 成功之后
	go source_routes.InitSourceRoutes(fmt.Sprintf("%s:%s", sourceHost, sourcePort))

	InitZapLog()
	InitWs()

	api_routes.InitApiRoutes(fmt.Sprintf("%s:%s", apiHost, apiPort))
}
