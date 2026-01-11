package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initNoticeRoutes(r *gin.RouterGroup) gin.IRoutes {
	noticeAdmin := new(admin.Notice)
	notice := r.Group("/notice").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		notice.GET("/noticeList", noticeAdmin.NoticeList)
		notice.GET("/addNotice", noticeAdmin.CreateNotice)
		notice.POST("/addNotice", noticeAdmin.CreateNotice)
		notice.GET("/editNotice", noticeAdmin.UpdateNotice)
		notice.POST("/editNotice", noticeAdmin.UpdateNotice)
		notice.POST("/deleteNotice", noticeAdmin.DeleteNotice)
	}

	return r
}
