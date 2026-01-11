package db

import (
	"github.com/nsqio/go-nsq"
	"go-novel/app/service/common/nsq_service"
	"go-novel/global"
	"go-novel/utils"
	"log"
)

func InitNsqConsumer() {
	//nsq_service.InitConsumer(utils.UpdateChapter, utils.Default, utils.NsqConsumerIP)
	//nsq_service.InitConsumer(utils.UpdateChapterText, utils.Default, utils.NsqConsumerIP)
	//nsq_service.InitConsumer(utils.ChapterText, utils.Default, utils.NsqConsumerIP)
	nsq_service.InitConsumer(utils.UpdateBook, utils.Default, utils.NsqConsumerIP)
	nsq_service.InitConsumer(utils.SourceUpdateLastChapter, utils.Default, utils.NsqConsumerIP)
	nsq_service.InitConsumer(utils.UpdateComic, utils.Default, utils.NsqConsumerIP)
}

func InitNsqProducer() {
	p, err := nsq.NewProducer(utils.NsqProducerIP, nsq.NewConfig()) // 新建生产者
	if err != nil {
		log.Fatal("err", err.Error())
		return
	}
	global.NsqPro = p
}
