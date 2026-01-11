package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/vip_service"
	"go-novel/utils"
	"strconv"
)

type Vip struct{}

func (vip *Vip) CardList(c *gin.Context) {
	var req models.CardListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := vip_service.CardListSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取列表失败")
		return
	}

	res := gin.H{
		"total":       total,
		"currentPage": req.PageNum,
		"list":        list,
	}
	utils.Success(c, res, "ok")
}

func (vip *Vip) CreateCard(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateCardReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := vip_service.CreateCard(&req)
		if err != nil {
			utils.Fail(c, err, "创建失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	utils.Success(c, "", "ok")
}

func (vip *Vip) UpdateCard(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateCardReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := vip_service.UpdateVipCard(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	cardId, _ := strconv.Atoi(c.Query("id"))
	cardInfo, err := vip_service.GetCardById(int64(cardId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"cardInfo": cardInfo,
	}
	utils.Success(c, res, "ok")
}

func (vip *Vip) DelVipCard(c *gin.Context) {
	var req models.DeleteCardReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := vip_service.DeleteCard(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}
