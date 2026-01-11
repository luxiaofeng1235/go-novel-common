package main

import (
	"go-novel/app/service/collect/collect_service"
	"go-novel/db"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	host, name, user, passwd := db.GetDB()
	db.InitDB(host, name, user, passwd)
	addr, passwd, defaultdb := db.GetRedis()
	db.InitRedis(addr, passwd, defaultdb)
	db.InitZapLog()
	db.InitNsqProducer()
	db.InitNsqConsumer()
	db.InitKeyLock()
	for {
		log.Println(collect_service.StartCollect(2, true))
		time.Sleep(time.Second * 1)
	}
	//collect_service.UpdateCollect()
	//log.Println(collect_service.StartCollect(1, true))

	//for {
	//	collect_service.StartCollect(4, true)
	//	time.Sleep(time.Second * 1)
	//}
	//log.Println(collect_service.StartCollect(1, true))
}
