package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initCommentRoutes(r *gin.RouterGroup) gin.IRoutes {
	commentAdmin := new(admin.Comment)
	feedback := r.Group("/comment").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		feedback.GET("/commentList", commentAdmin.CommentList)
		feedback.GET("/sonComments", commentAdmin.SonComments)
		feedback.GET("/editComment", commentAdmin.UpdateComment)
		feedback.POST("/updateComment", commentAdmin.UpdateComment)
		feedback.POST("/deleteComment", commentAdmin.DelComment)
	}

	return r
}
