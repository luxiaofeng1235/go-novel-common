package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/menu_service"
	"go-novel/utils"
	"strconv"
)

type Menu struct{}

// 获取菜单列表
func (menu *Menu) MenuList(c *gin.Context) {
	var req models.MenuListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	// 获取
	list, total, err := menu_service.MenuListSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取菜单列表失败")
		return
	}
	utils.Success(c, gin.H{"list": list, "total": total, "currentPage": req.PageNum}, "ok")
}

// 创建菜单
func (menu *Menu) CreateMenu(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateMenuReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		// 创建菜单
		InsertId, err := menu_service.CreateMenu(&req)
		if err != nil {
			utils.Fail(c, err, "添加菜单失败")
			return
		}
		utils.Success(c, InsertId, "ok")
		return
	}

	//获取父级菜单信息
	list, err := menu_service.GetIsMenuList()
	if err != nil {
		utils.Fail(c, err, "获取数据失败")
		return
	}
	utils.Success(c, gin.H{"parentList": list}, "ok")
}

// 更新菜单
func (menu *Menu) UpdateMenu(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateMenuReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		// 创建菜单
		isUpdate, err := menu_service.UpdateMenu(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改菜单失败")
			return
		}
		utils.Success(c, "", "ok")
		return
	}

	//菜单信息
	menuId, _ := strconv.Atoi(c.Query("id"))
	menuInfo, err := menu_service.GetMenuById(int64(menuId))
	if err != nil {
		utils.Fail(c, nil, "获取菜单信息失败")
		return
	}

	//获取父级菜单信息
	menuList, err := menu_service.GetIsMenuList()
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}
	utils.Success(c, gin.H{"parentList": menuList, "menu": menuInfo}, "ok")
}

// 删除菜单
func (menu *Menu) DeleteMenu(c *gin.Context) {
	var req models.DeleteMenuReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := menu_service.DeleteMenu(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除菜单失败")
		return
	}
	utils.Success(c, "", "删除菜单信息成功")
	return
}
