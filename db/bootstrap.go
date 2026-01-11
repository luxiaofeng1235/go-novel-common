package db

import (
	"go-novel/routers/admin_routes"
	"go-novel/routers/api_routes"
	"log"
)

func StartServer() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	//log.Println(host, name, user, passwd)
	InitDB()
	InitZapLog()
	//InitZinc()
	InitWs()
	InitGeoReadre()
	InitBigcache()
	api_routes.InitApiRoutes()
}

func StartBiqugeServer() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	InitZapLog()
	api_routes.InitBiqugeRoutes()
}

// 设置后台启动的服务
func StartAdminServer() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	//log.Println(host, name, user, passwd)
	InitDB()
	InitZapLog()
	InitWs()
	InitGeoReadre()
	InitBigcache()
	admin_routes.InitRoutes()
}
