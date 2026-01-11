package admin

import (
	"go-novel/app/models"
	"go-novel/app/service/admin/admin_service"
	"go-novel/app/service/admin/role_service"
	"go-novel/utils"
	"strconv"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

type Admin struct{}

// 获取管理员列表
func (admin *Admin) AdminList(c *gin.Context) {
	var req models.AdminListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	// 获取管理员列表
	list, total, err := admin_service.AdminListSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取管理员列表失败")
		return
	}

	admins := make([]map[string]interface{}, len(list))
	//获取所有角色信息
	allRoles, err := role_service.GetRoleList()
	if err != nil {
		utils.Fail(c, err, "获取管理员角色数据失败")
		return
	}

	for k, u := range list {
		admins[k] = structs.Map(u)
		roles, err := admin_service.GetAdminRole(u.Id, allRoles)
		if err != nil {
			utils.Fail(c, err, "获取管理员角色数据失败")
			return
		}
		roleInfo := make([]map[string]interface{}, 0, len(roles))
		for _, r := range roles {
			roleInfo = append(roleInfo, map[string]interface{}{"roleId": r.Id, "name": r.Name})
		}
		admins[k]["roleInfo"] = roleInfo
	}

	res := gin.H{
		"total":       total,
		"currentPage": req.PageNum,
		"list":        admins,
	}
	utils.Success(c, res, "ok")
}

// 创建管理员
func (admin *Admin) CreateAdmin(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateAdminReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		// 创建管理员
		InsertId, err := admin_service.CreateAdmin(&req)
		if err != nil {
			utils.Fail(c, err, "创建管理员失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	//获取角色信息
	roleListEntities, err := role_service.GetRoleList()
	if err != nil {
		utils.Fail(c, err, "获取角色数据失败")
		return
	}

	res := gin.H{
		"roleList": roleListEntities,
	}
	utils.Success(c, res, "ok")
}

// 更新管理员
func (admin *Admin) UpdateAdmin(c *gin.Context) {
	if c.Request.Method == "POST" {
		username, _ := c.Get("username")

		var req models.UpdateAdminReq
		// 参数绑定
		if err := c.ShouldBind(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		//获取path中的userId
		if req.AdminId <= 0 {
			utils.Fail(c, nil, "用户ID不正确")
			return
		}

		//不能禁用自己
		if req.Username == username && req.Status == 0 {
			utils.Fail(c, nil, "不能禁用自己")
			return
		}

		isUpdate, err := admin_service.UpdateAdmin(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改管理员信息失败")
			return
		}
		utils.Success(c, "", "修改管理员信息成功")
		return
	}

	//管理员用户信息
	adminId, _ := strconv.Atoi(c.Query("id"))
	adminEntity, err := admin_service.GetAdminById(int64(adminId))
	if err != nil {
		utils.Fail(c, nil, "获取管理员数据失败")
		return
	}

	//获取角色信息
	roleListEntities, err := role_service.GetRoleList()
	if err != nil {
		utils.Fail(c, nil, "获取角色数据失败")
		return
	}

	//获取已选择的角色信息
	checkedRoleIds, err := admin_service.GetAdminRoleIds(int64(adminId))
	if err != nil {
		utils.Fail(c, nil, "获取用户角色数据失败")
	}
	if checkedRoleIds == nil {
		checkedRoleIds = []uint{}
	}

	res := gin.H{
		"roleList":       roleListEntities,
		"adminInfo":      adminEntity,
		"checkedRoleIds": checkedRoleIds,
	}
	utils.Success(c, res, "ok")
}

// 删除管理员
func (admin *Admin) DeleteAdmin(c *gin.Context) {
	var req models.DeleteAdminReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	username, _ := c.Get("username")
	adminEntity, err := admin_service.GetAdminByUsername(username.(string))
	if err != nil {
		utils.Fail(c, nil, "获取管理员数据失败")
		return
	}
	// 前端传来的用户ID
	reqUserIds := req.AdminIds
	for k, _ := range reqUserIds {
		if reqUserIds[k] == adminEntity.Id {
			utils.Fail(c, nil, "用户不能删除自己")
			return
		}
	}

	isDelete, err := admin_service.DeleteAdmin(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除管理员失败")
		return
	}

	utils.Success(c, "", "删除管理员信息成功")
	return
}

// 重置管理员密码
func (admin *Admin) ResetAdminPwd(c *gin.Context) {
	var req models.ResetPwdReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	if err := admin_service.ResetAdminPwd(&req); err != nil {
		utils.Fail(c, err, "用户密码重置失败")
		return
	}

	utils.Success(c, "", "用户密码重置成功")
}

// 修改管理员状态
func (admin *Admin) ChangeAdminStatus(c *gin.Context) {
	var req models.AdminStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	if err := admin_service.ChangeAdminStatus(&req); err != nil {
		utils.Fail(c, err, "修改管理员状态失败")
		return
	}

	utils.Success(c, "", "用户状态设置成功")
}
