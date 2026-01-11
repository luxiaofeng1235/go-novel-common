package admin_routes

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-novel/config"
	"go-novel/middleware"
	"go-novel/utils"
	"log"
)

func InitRoutes() {
	if viper.GetBool("server.debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	//路由设置
	r = router(r)

	_ = r.Run(getHttpString())
}

// 设置路由
func router(r *gin.Engine) *gin.Engine {
	//加载public目录
	if config.GetString("server.env") == utils.Local { //本地图片资源判断
		r.Static("/mnt", "E:\\mnt")
	} else { //线上资源判断
		r.Static("/mnt", "/mnt")
	}
	r.Static("/public", "./public")
	//配置dist静态资源
	r.Static("/static", "./public/dist/static/")
	r.Static("/js", "./public/dist/js")
	r.StaticFile("/favicon.ico", "./public/dist/favicon.ico")
	//两种加载index.html写法都可以
	r.StaticFile("/", "./public/dist/index.html")
	//r.Use(static.Serve("/", static.LocalFile("./public/dist/", true)))
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(sessions.Sessions("go-admin", cookie.NewStore([]byte("2616411af10baad87e8d3a356feb0a5c"))))
	r = adminRouter(r)
	return r
}

// 后台路由
func adminRouter(r *gin.Engine) *gin.Engine {
	system := r.Group("system")
	// 注册路由组
	initIndexRoutes(system)
	initAuthRoutes(system)
	initCommonRoutes(system)
	initSettingRoutes(system)
	initUserRoutes(system)
	initFeedbackRoutes(system)
	initTaskRoutes(system)
	initCommentRoutes(system)
	initOrderRoutes(system)
	initClassRoutes(system)
	initBookRoutes(system)
	initAdverRoutes(system)
	initTagRoutes(system)
	initFindBookRoutes(system)
	initNoticeRoutes(system)
	initVersionRoutes(system)
	initCollectRoutes(system)
	initRankRoutes(system)
	initVipRoutes(system)
	log.Println("初始化路由完成！")
	return r
}

// 获取启动地址
func getHttpString() string {
	var port string
	flag.StringVar(&port, "port", viper.GetString("server.port"), "default :port")
	flag.Parse()
	return fmt.Sprintf(":%s", port)
}
