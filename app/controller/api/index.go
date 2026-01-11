package api

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/gojsonq/v2"
	"go-novel/utils"
	"log"
	"time"
)

type Index struct{}

func (index *Index) IndexGet(c *gin.Context) {
	gq := gojsonq.New().File("./data.json")
	start := time.Now()

	pageNum := 3
	pageSize := 1000
	//total := gq.Count()
	list := gq.Offset((pageNum - 1) * pageSize).Limit(pageSize).Get()
	log.Println(gq.Get())

	//list := gq.Get()
	//log.Println(list)
	log.Println(time.Since(start).Milliseconds())
	//c.Writer.Write([]byte("首页测试"))
	utils.SuccessEncrypt(c, list, "获取列表成功")
}
