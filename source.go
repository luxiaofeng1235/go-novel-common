package main

import (
	"go-novel/db"
	"go-novel/routers/source_routes"
)

func main() {
	host, name, user, passwd := db.GetDB()
	db.InitDB(host, name, user, passwd)
	db.InitZapLog()
	db.InitNsqProducer()
	db.InitNsqConsumer()
	source_routes.InitSourceRoutes()
}
