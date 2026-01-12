/*
 * @Descripttion: source 静态资源服务入口
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 11:58:00
 */
package main

import (
	"go-novel/db"
	"go-novel/routers/source_routes"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	db.InitZapLog()
	source_routes.InitSourceRoutes()
}
