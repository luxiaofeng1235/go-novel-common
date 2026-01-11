package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/search_service"
	"go-novel/utils"
)

type Search struct{}

func (search *Search) SearchHistory(c *gin.Context) {
	var req models.SearchHistoryReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	searchs, err := search_service.SearchHistory(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"searchs": searchs,
	}
	utils.SuccessEncrypt(c, res, "获取数据成功")
}

func (search *Search) HotSearchRank(c *gin.Context) {
	var req models.SearchHotReq
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	req.Ip = utils.RemoteIp(c)
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark
	list, err := search_service.GetHotSearchRank(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"list": list,
	}
	utils.SuccessEncrypt(c, res, "获取数据成功")
}
