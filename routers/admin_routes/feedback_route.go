package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initFeedbackRoutes(r *gin.RouterGroup) gin.IRoutes {
	feedbackAdmin := new(admin.Feedback)
	feedback := r.Group("/feedback").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		feedback.GET("/helpList", feedbackAdmin.HelpList)
		feedback.GET("/addHelp", feedbackAdmin.CreateHelp)
		feedback.POST("/addHelp", feedbackAdmin.CreateHelp)
		feedback.GET("/editHelp", feedbackAdmin.UpdateHelp)
		feedback.POST("/editHelp", feedbackAdmin.UpdateHelp)
		feedback.POST("/deleteHelp", feedbackAdmin.DeleteHelp)
		feedback.GET("/feedbackList", feedbackAdmin.FeedbackList)
		feedback.POST("/feedbackReply", feedbackAdmin.FeedbackReply)
		feedback.GET("/feedbackBookList", feedbackAdmin.FeedbackBookList)
		feedback.GET("/editFeedbackBook", feedbackAdmin.UpdateFeedbackBook)
		feedback.POST("/updateFeedbackBook", feedbackAdmin.UpdateFeedbackBook)
		feedback.POST("/delFeedbackBook", feedbackAdmin.DelFeedbackBook)
	}

	return r
}
