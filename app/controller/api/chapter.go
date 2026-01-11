package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/chapter_service"
	"go-novel/utils"
)

type Chapter struct{}

func (chapter *Chapter) FeedbackAdd(c *gin.Context) {
	var req models.ChapterFeedBackAddReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	req.Ip = utils.RemoteIp(c)
	err := chapter_service.FeedBackAdd(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, "", "反馈成功")
}
