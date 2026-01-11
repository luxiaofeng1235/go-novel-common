package api_routes

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-novel/app/controller/api"
	"go-novel/app/service/api/checkin_service"
	"go-novel/config"
	"go-novel/middleware"
	"go-novel/utils"
	"log"
	"net/http"
	"time"
)

func InitApiRoutes() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(sessions.Sessions("go-admin", cookie.NewStore([]byte("2616411af10baad87e8d3a356feb0a5c"))))
	go api.WsHandleMsg()
	go checkin_service.CheckRemind()
	//路由设置
	r = apiRouter(r)
	s := &http.Server{
		Addr:         getHttpString(),
		Handler:      r,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	_ = s.ListenAndServe()
}

// 笔趣阁的端口监听
func InitBiqugeRoutes() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(sessions.Sessions("go-admin", cookie.NewStore([]byte("2616411af10baad87e8d3a356feb0a5c"))))
	go checkin_service.CheckRemind()
	//路由设置
	r = apiBiqugeRouter(r)
	s := &http.Server{
		Addr:         getBiqugeHttpString(),
		Handler:      r,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	_ = s.ListenAndServe()
}

// 获取启动地址
func getHttpString() string {
	var port string
	flag.StringVar(&port, "port", viper.GetString("api.port"), "default :port")
	flag.Parse()
	return fmt.Sprintf(":%s", port)
}

// 获取笔趣阁的端口监听
func getBiqugeHttpString() string {
	var port string
	flag.StringVar(&port, "port", viper.GetString("biquge.port"), "default :port")
	flag.Parse()
	return fmt.Sprintf(":%s", port)
}

// 笔趣阁的路由配置
func apiBiqugeRouter(r *gin.Engine) *gin.Engine {
	system := r.Group("api")
	//注册路由组件
	initBiqugeRoutes(system)
	log.Println("初始化路由完成！")
	return r
}

// api路由
func apiRouter(r *gin.Engine) *gin.Engine {
	//加载public目录
	if config.GetString("server.env") == utils.Local { //本地图片资源判断
		r.Static("/mnt", "E:\\mnt")
	} else { //线上资源判断
		r.Static("/mnt", "/mnt")
	}
	r.Static("/public", "./public")
	////配置dist静态资源
	r.Static("/static", "./public/dist_gold/static/")
	r.StaticFile("/favicon.ico", "./public/dist_gold/favicon.ico")
	////两种加载index.html写法都可以
	//r.StaticFile("/", "./public/dist/index.html")
	r.Use(gin.Logger(), middleware.RecoverWithZincSearch(), middleware.GinBodyLogMiddleware())
	system := r.Group("api")
	// 注册路由组
	initCommonRoutes(system)

	initIndexRoutes(system)
	initUserRoutes(system)
	initBookRoutes(system)
	initCommentRoutes(system)
	initBookShelfRoutes(system)
	initChapterRoutes(system)
	initReadRoutes(system)
	initTaskRoutes(system)
	initClassRoutes(system)
	initFeedbackRoutes(system)
	initMessageRoutes(system)
	initSearchRoutes(system)
	initSettingRoutes(system)
	initPayRoutes(system)
	initVipRoutes(system)
	initWithdrawRoutes(system)
	initOrderRoutes(system)
	initSourceRoutes(system)
	initAdverRoutes(system)
	initFindBookRoutes(system)
	initVersionRoutes(system)
	initRankRoutes(system)
	log.Println("初始化路由完成！")
	return r
}
