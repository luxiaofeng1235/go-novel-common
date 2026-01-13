package db

import (
	"go-novel/global"
	"go-novel/pkg/ws"
)

func InitWs() {
	hub := ws.NewHub()
	global.WsHub = hub
	go hub.Run()
}
