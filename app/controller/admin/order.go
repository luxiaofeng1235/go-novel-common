package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/order_service"
	"go-novel/utils"
	"strconv"
)

type Order struct{}

func (order *Order) OrderList(c *gin.Context) {
	var req models.OrderListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := order_service.OrderListSearch(&req)
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

func (order *Order) UpdateOrder(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateOrderReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := order_service.UpdateOrder(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	orderId, _ := strconv.Atoi(c.Query("id"))
	orderInfo, err := order_service.GetOrderById(int64(orderId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"orderInfo": orderInfo,
	}
	utils.Success(c, res, "ok")
}
