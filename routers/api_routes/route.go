/*
 * @Descripttion: API 路由入口（脚手架：仅 user + ws）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 13:45:00
 */
package api_routes

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-novel/middleware"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func InitApiRoutes(addr string) {
	if viper.GetBool("api.debug") || viper.GetBool("server.debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.Cors())
	r = apiRouter(r)
	s := &http.Server{
		Addr:         getHttpString(addr),
		Handler:      r,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}

	// 强制使用 IPv4 监听，避免 Windows 侧仅出现 [::1] 导致 127.0.0.1 无法访问
	ln, err := net.Listen("tcp4", s.Addr)
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

// 获取启动地址
func getHttpString(addr string) string {
	if addr != "" {
		return addr
	}
	host := strings.TrimSpace(viper.GetString("server.host"))
	port := viper.GetInt("server.port")
	if host == "" || port == 0 {
		// 兜底：按脚手架约定必须在 server 配置监听地址，这里避免默默用默认值造成联调混乱
		log.Fatal("缺少 server.host/server.port：请在 config.yml 的 server 节点配置 API 监听地址")
	}
	return fmt.Sprintf("%s:%d", host, port)
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
