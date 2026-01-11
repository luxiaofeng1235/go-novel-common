package source_routes

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-novel/middleware"
	"log"
	"net/http"
	"time"
)

func InitSourceRoutes() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Use(sessions.Sessions("go-admin", cookie.NewStore([]byte("2616411af10baad87e8d3a356feb0a5c"))))
	//路由设置
	r = sourceRouter(r)
	s := &http.Server{
		Addr:         getHttpString(),
		Handler:      r,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	_ = s.ListenAndServe()
}

// 获取启动地址
func getHttpString() string {
	var port string
	flag.StringVar(&port, "port", viper.GetString("source.port"), "default :port")
	flag.Parse()
	return fmt.Sprintf(":%s", port)
}

// api路由
func sourceRouter(r *gin.Engine) *gin.Engine {
	//加载public目录
	r.Static("/public", "./public")
	r.Static("/attachment", "./attachment")
	system := r.Group("source")
	// 注册路由组
	initCommonRoutes(system)
	initBookRoutes(system)
	log.Println("初始化路由完成！")
	return r
}
