/*
 * @Descripttion: Source 静态资源路由入口
 * @Author: red
 * @Date: 2026-01-12 11:58:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 13:45:00
 */
package source_routes

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-novel/config"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func InitSourceRoutes(addr string) {
	if viper.GetBool("source.debug") || viper.GetBool("server.debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r = sourceRouter(r)
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

func getHttpString(addr string) string {
	if addr != "" {
		return addr
	}
	host := viper.GetString("source.host")
	port := viper.GetString("source.port")
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "8007"
	}
	return fmt.Sprintf("%s:%s", host, port)
}

func sourceRouter(r *gin.Engine) *gin.Engine {
	// 触发配置加载（config 包 init 会读取根目录 config.yml）
	_ = config.GetString("server.env")

	// 按目录存在性挂载，避免误导（目录不存在时 gin.Static 也会挂载，但访问 404）
	mountStatic(r, "/public", "./public")
	mountStatic(r, "/resource", "./public/resource")
	mountStatic(r, "/static", "./public/static")
	mountStatic(r, "/dist", "./public/dist")
	mountStatic(r, "/dist_gold", "./public/dist_gold")

	// 兼容常见 favicon
	if fileExists("./public/favicon.ico") {
		r.StaticFile("/favicon.ico", "./public/favicon.ico")
	}

	log.Println("初始化 source 静态路由完成！")
	return r
}

func mountStatic(r *gin.Engine, urlPrefix, localDir string) {
	if dirExists(localDir) {
		r.Static(urlPrefix, localDir)
	}
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
