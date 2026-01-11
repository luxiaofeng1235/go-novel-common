package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/class_service"
	"go-novel/utils"
	"strconv"
)

type Class struct{}

func (class *Class) TypeList(c *gin.Context) {
	var req models.TypeListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := class_service.TypeListSearch(&req)
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

func (class *Class) CreateType(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateTypeReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := class_service.CreateType(&req)
		if err != nil {
			utils.Fail(c, err, "创建失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	utils.Success(c, "", "ok")
}

func (class *Class) UpdateType(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateTypeReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := class_service.UpdateType(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	typeId, _ := strconv.Atoi(c.Query("id"))
	typeInfo, err := class_service.GetTypeById(int64(typeId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"typeInfo": typeInfo,
	}
	utils.Success(c, res, "ok")
}

func (class *Class) DelType(c *gin.Context) {
	var req models.DeleteTypeReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := class_service.DeleteType(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}

func (class *Class) ClassList(c *gin.Context) {
	var req models.ClassListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := class_service.ClassListSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取列表失败")
		return
	}

	typeList, err := class_service.GetTypeList()
	if err != nil {
		utils.Fail(c, err, "获取分类列表失败")
		return
	}

	res := gin.H{
		"total":       total,
		"currentPage": req.PageNum,
		"list":        list,
		"typeList":    typeList,
	}
	utils.Success(c, res, "ok")
}

func (class *Class) CreateClass(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateClassReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := class_service.CreateClass(&req)
		if err != nil {
			utils.Fail(c, err, "创建失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	typeList, err := class_service.GetTypeList()
	if err != nil {
		utils.Fail(c, err, "获取类型列表失败")
		return
	}

	res := gin.H{
		"typeList": typeList,
	}

	utils.Success(c, res, "ok")
}

func (class *Class) UpdateClass(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateClassReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := class_service.UpdateClass(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	classId, _ := strconv.Atoi(c.Query("id"))
	classInfo, err := class_service.GetClassById(int64(classId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}
	classInfo.ClassPic = utils.GetAdminFileUrl(classInfo.ClassPic)
	typeList, err := class_service.GetTypeList()
	if err != nil {
		utils.Fail(c, err, "获取类型列表失败")
		return
	}

	res := gin.H{
		"classInfo": classInfo,
		"typeList":  typeList,
	}
	utils.Success(c, res, "ok")
}

func (class *Class) BookList(c *gin.Context) {
	var req models.BookListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := class_service.BookListSearch(&req)
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

func (class *Class) DelClass(c *gin.Context) {
	var req models.DeleteClassReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := class_service.DeleteClass(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}

func (class *Class) AssignClass(c *gin.Context) {
	var req models.AssignClassReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := class_service.AssignClass(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "归类信息失败")
		return
	}

	utils.Success(c, "", "归类信息成功")
	return
}
