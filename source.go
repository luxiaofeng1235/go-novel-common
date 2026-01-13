/*
 * @Descripttion: source 静态资源服务入口
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 13:45:00
 */
package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	_ "go-novel/config"
	"go-novel/db"
	"go-novel/routers/source_routes"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	var host, port string
	flag.StringVar(&host, "host", viper.GetString("source.host"), "source listen host")
	flag.StringVar(&port, "port", viper.GetString("source.port"), "source listen port")
	flag.Parse()

	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "8007"
	}

	db.InitZapLog()
	source_routes.InitSourceRoutes(fmt.Sprintf("%s:%s", host, port))
}
