package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/casbin_adapter_service"
	"go-novel/app/service/admin/menu_service"
	"go-novel/app/service/admin/role_service"
	"go-novel/utils"
	"strconv"
)

type Role struct{}

// 获取角色列表
func (role *Role) RoleList(c *gin.Context) {
	var req models.RoleListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	// 获取
	list, total, err := role_service.RoleListSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取角色列表失败")
		return
	}
	utils.Success(c, gin.H{"list": list, "total": total, "currentPage": req.PageNum}, "ok")
}

// 创建角色
func (role *Role) CreateRole(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateRoleReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		// 创建角色
		InsertId, err := role_service.CreateRole(&req)
		if err != nil {
			utils.Fail(c, err, "创建角色失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	//获取菜单信息
	mListEntities, err := menu_service.GetMenuList()
	if err != nil {
		utils.Fail(c, err, "获取菜单数据失败")
		return
	}
	var mList []map[string]interface{}
	for _, entity := range mListEntities {
		m := map[string]interface{}{
			"id":    entity.Id,
			"pid":   entity.Pid,
			"label": entity.Title,
		}
		mList = append(mList, m)
	}

	mList = menu_service.PushSonToParent(mList, 0, "pid", "children", false)
	utils.Success(c, gin.H{"menuList": mList}, "ok")
}

// 修改角色
func (role *Role) UpdateRole(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateRoleReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		// 修改角色
		isUpdate, err := role_service.UpdateRole(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改角色失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	//获取角色信息
	roleId, _ := strconv.Atoi(c.Query("id"))
	roleInfo, err := role_service.GetRoleById(int64(roleId))
	if err != nil {
		utils.Fail(c, nil, "获取角色数据失败")
		return
	}

	//获取菜单信息
	mListEntities, err := menu_service.GetMenuList()
	if err != nil {
		utils.Fail(c, nil, "获取菜单数据失败")
		return
	}

	//获取角色关联的菜单规则
	gpSlice, err := casbin_adapter_service.GetRoleIdRule(roleId)
	if err != nil {
		utils.Fail(c, nil, "获取权限处理器失败")
		return
	}

	var mList []map[string]interface{}
	for _, entity := range mListEntities {
		m := gin.H{
			"id":    entity.Id,
			"pid":   entity.Pid,
			"label": entity.Title,
		}
		mList = append(mList, m)
	}

	mList = menu_service.PushSonToParent(mList, 0, "pid", "children", true)
	res := gin.H{
		"menuList":     mList,
		"role":         roleInfo,
		"checkedRules": gpSlice,
	}
	utils.Success(c, res, "ok")
}

// 删除角色
func (role *Role) DeleteRole(c *gin.Context) {
	var req models.DeleteRoleReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := role_service.DeleteRole(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除角色失败")
		return
	}

	utils.Success(c, "", "删除角色信息成功")
	return
}

// 修改角色状态
func (role *Role) ChangeRoleStatus(c *gin.Context) {
	var req models.RoleStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	if err := role_service.ChangeRoleStatus(&req); err != nil {
		utils.Fail(c, err, "修改角色状态失败")
		return
	}

	utils.Success(c, "", "角色状态设置成功")
}
