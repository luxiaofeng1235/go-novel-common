package global

import (
	"github.com/allegro/bigcache/v3"
	"github.com/go-redis/redis/v8"
	"github.com/nsqio/go-nsq"
	"github.com/olahol/melody"
	"github.com/oschwald/geoip2-golang"
	go_keylock "github.com/sjy3/go-keylock"
	"github.com/thedevsaddam/gojsonq/v2"
	"go-novel/pkg/ws"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

/*
全局变量
*/
var (
	DB            *gorm.DB
	Redis         *redis.Client
	NsqPro        *nsq.Producer
	KeyLock       *go_keylock.KeyLock
	Errlog        *zap.SugaredLogger
	Zssqlog       *zap.SugaredLogger
	Sqllog        *zap.SugaredLogger
	Paylog        *zap.SugaredLogger
	Wslog         *zap.SugaredLogger
	Collectlog    *zap.SugaredLogger
	Nsqlog        *zap.SugaredLogger
	Updatelog     *zap.SugaredLogger
	Jsonq         *zap.SugaredLogger
	Biquge34log   *zap.SugaredLogger
	Paoshu8log    *zap.SugaredLogger
	Xswlog        *zap.SugaredLogger
	Lydlog        *zap.SugaredLogger
	Bqg24log      *zap.SugaredLogger
	Siluke520log  *zap.SugaredLogger
	VivoClicklog  *zap.SugaredLogger
	SmClicklog    *zap.SugaredLogger
	BaiduClicklog *zap.SugaredLogger
	Requestlog    *zap.SugaredLogger
	GeoCityReader *geoip2.Reader
	GeoAsnReader  *geoip2.Reader
	Ws            *melody.Melody
	WsHub         *ws.Hub
	Gojsonq       *gojsonq.JSONQ
	Bigcache      *bigcache.BigCache
)
