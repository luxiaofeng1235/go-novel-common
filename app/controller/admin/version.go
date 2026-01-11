package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/version_service"
	"go-novel/utils"
	"strconv"
)

type Version struct{}

func (version *Version) AppVersionList(c *gin.Context) {
	var req models.AppVersionListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := version_service.VersionListSearch(&req)
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

// 创建接口
func (version *Version) CreateAppVersion(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateAppVersionReq
		if err := c.ShouldBind(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		insertId, err := version_service.CreateVersion(&req)
		if err != nil {
			utils.Fail(c, err, "保存失败")
			return
		}
		utils.Success(c, gin.H{
			"insertId": insertId,
		}, "ok")
		return
	}
	utils.Success(c, "", "ok")
}

// 删除渠道
func (version *Version) DeleteAppVersion(c *gin.Context) {
	var req models.DeleteVersionReq
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	//删除
	isDelete, err := version_service.DeleteVersionRes(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除失败")
		return
	}
	utils.Success(c, "", "删除成功")
	return

}

func (version *Version) UpdateAppVersion(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateAppVersionReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := version_service.UpdateVersion(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	versionId, _ := strconv.Atoi(c.Query("id"))
	versionInfo, err := version_service.GetVersionById(int64(versionId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"versionInfo": versionInfo,
	}
	utils.Success(c, res, "ok")
}
