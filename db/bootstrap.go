/*
 * @Descripttion: 服务启动编排（脚手架最小链路）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 11:45:00
 */
package db

import (
	"go-novel/routers/api_routes"
	"log"
)

// StartApiServer 启动 API 服务（脚手架最小启动链路）
func StartApiServer() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	host, name, user, passwd := GetDB()
	InitMysql(host, name, user, passwd)
	InitZapLog()
	InitWs()
	api_routes.InitApiRoutes()
}
