package main

import (
	"go-novel/app/models"
	"go-novel/db"
	"go-novel/global"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	host, name, user, passwd := db.GetDB()
	db.InitDB(host, name, user, passwd)
	addr, passwd, defaultdb := db.GetRedis()
	db.InitRedis(addr, passwd, defaultdb)
	db.InitZapLog()
	db.InitKeyLock()

	var books []*models.McBook
	global.DB.Model(models.McBook{}).Debug().Order("id desc").Where("source_url like ?", "%"+"biquge34"+"%").Find(&books)
	if len(books) <= 0 {
		return
	}
	for _, book := range books {

	}
}
