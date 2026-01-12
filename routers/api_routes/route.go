/*
 * @Descripttion: API 路由入口（脚手架：仅 user + ws）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 13:25:00
 */
package api_routes

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-novel/middleware"
	"log"
	"net/http"
	"time"
)

func InitApiRoutes() {
	if viper.GetBool("api.debug") || viper.GetBool("server.debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.Cors())
	r = apiRouter(r)
	s := &http.Server{
		Addr:         getHttpString(),
		Handler:      r,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

// 获取启动地址
func getHttpString() string {
	var host string
	var port string
	flag.StringVar(&host, "host", viper.GetString("api.host"), "default host")
	flag.StringVar(&port, "port", viper.GetString("api.port"), "default :port")
	flag.Parse()
	if host == "" {
		host = "0.0.0.0"
	}
	return fmt.Sprintf("%s:%s", host, port)
}

// api路由
func apiRouter(r *gin.Engine) *gin.Engine {
	system := r.Group("api")
	initUserRoutes(system)
	initCommonRoutes(system)
	initWsRoutes(system)
	log.Println("初始化路由完成！")
	return r
}
