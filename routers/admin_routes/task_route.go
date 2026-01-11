package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initTaskRoutes(r *gin.RouterGroup) gin.IRoutes {
	checkinAdmin := new(admin.Checkin)
	checkin := r.Group("/checkin").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		checkin.GET("/rewardList", checkinAdmin.RewardList)
		checkin.GET("/addReward", checkinAdmin.CreateReward)
		checkin.POST("/addReward", checkinAdmin.CreateReward)
		checkin.GET("/editReward", checkinAdmin.UpdateReward)
		checkin.POST("/updateReward", checkinAdmin.UpdateReward)
		checkin.POST("/deleteReward", checkinAdmin.DelReward)
		checkin.GET("/checkinList", checkinAdmin.CheckinList)
	}

	taskAdmin := new(admin.Task)
	task := r.Group("/task").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		task.GET("/taskList", taskAdmin.TaskList)
		task.GET("/addTask", taskAdmin.CreateTask)
		task.POST("/addTask", taskAdmin.CreateTask)
		task.GET("/editTask", taskAdmin.UpdateTask)
		task.POST("/updateTask", taskAdmin.UpdateTask)
		task.POST("/deleteTask", taskAdmin.DelTask)
		task.GET("/recordList", taskAdmin.RecordList)
	}
	return r
}
