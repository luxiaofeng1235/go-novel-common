package db

import (
	"github.com/olahol/melody"
	"go-novel/global"
	"net/http"
	"time"
)

func InitWs() {
	m := melody.New()
	m.Config = &melody.Config{
		WriteWait:         m.Config.WriteWait,
		PongWait:          m.Config.PongWait,
		PingPeriod:        time.Second * 1,
		MaxMessageSize:    m.Config.MaxMessageSize,
		MessageBufferSize: m.Upgrader.ReadBufferSize,
	}
	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	global.Ws = m
}
