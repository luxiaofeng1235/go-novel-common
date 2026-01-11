package admin

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/adver_service"
	"go-novel/utils"
	"io/ioutil"
	"log"
	"strconv"
)

type Adver struct{}

// 广告列表
func (adver *Adver) AdverList(c *gin.Context) {
	//初始化需要绑定的参数方便进行表单接收
	var req models.AdverListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	//遍历数据的对象列表
	list, total, err := adver_service.AdverListSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取列表失败")
		return
	}
	//返回定义的数组对象信息，方便进行组装
	res := gin.H{
		"total":       total,
		"currentPage": req.PageNum,
		"list":        list,
	}
	utils.Success(c, res, "ok")
}

// 广告项目列表
func (adver *Adver) AdverProjectList(c *gin.Context) {
	var req models.AdverProjectReq
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定错误")
		return
	}
	list, total, err := adver_service.AdverProjectListSearch(&req)
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
	return
}

// 获取类型列表
func (adver *Adver) ProjectTypeList(c *gin.Context) {
	var req models.AdverProjectReq
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	//获取查询的所有类型.只返回类型
	list, err := adver_service.GetALlProjectTypeList()
	if err != nil {
		utils.Fail(c, err, "获取失败")
		return
	}
	res := gin.H{
		"list": list,
	}
	utils.Success(c, res, "ok")
	return
}

// 保存广告项目信息
func (adver *Adver) AddAdverProject(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateAdverProjectReq
		if err := c.ShouldBind(&req); err != nil {
			utils.Fail(c, err, "绑定参数错误")
			return
		}
		//保存广告设置信息
		isUpdate, err := adver_service.SaveAdverProject(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}
		utils.Success(c, "", "ok")
		return
	}
	projectId, _ := strconv.Atoi(c.Query("id"))
	projectInfo, err := adver_service.GetAdverProjectInfo(int64(projectId))
	if err != nil {
		utils.Fail(c, err, "获取数据失败")
		return
	}
	res := gin.H{
		"info": projectInfo,
	}
	utils.Success(c, res, "ok")
}

// 获取所有的广告包信息
func (adver *Adver) GetAllPackageList(c *gin.Context) {
	packageList, err := adver_service.GetALlAdminPackageList()
	if err != nil {
		utils.Fail(c, err, "获取失败")
		return
	}
	res := gin.H{
		"packageList": packageList,
	}
	utils.Success(c, res, "ok")
}

// 广告包列表
func (adver *Adver) AdverPackageList(c *gin.Context) {
	var req models.AdverPackageReq
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "绑定参数错误")
		return
	}
	//获取广告包管理列表
	list, total, err := adver_service.AdvertPackageListSearch(&req)
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

// 添加广告包信息
func (adver *Adver) CreateAdverPackage(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateAdverPackageReq
		//拉取对应的数据信息
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatalf("read body failed at Before,err:%s", err.Error())
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		//log.Println("read body: ", string(body))

		if err := c.ShouldBind(&req); err != nil {
			utils.Fail(c, err, "绑定参数失败")
			return
		}
		username, _ := c.Get("username")
		req.CreateUser = username.(string)
		//添加广告包的信息
		InsertId, err := adver_service.CreateAdverPackage(&req)
		if err != nil {
			utils.Fail(c, err, "创建失败")
			return
		}
		//成功返回对应的状态信息
		utils.Success(c, InsertId, "ok")
		return
	}
	utils.Success(c, "", "ok")

}

// 删除广告包设置信息
func (adver *Adver) DeleteAdverPackage(c *gin.Context) {

	var req models.DeleteAdverPackageReq

	//参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "绑定参数失败")
		return
	}

	isDelete, err := adver_service.DeleteAdverPackage(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}

// 更新包信息
func (adver *Adver) UpdateAdverPackage(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateAdverPackageReq
		//拉取对应的数据信息
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatalf("read body failed at Before,err:%s", err.Error())
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		//log.Println("read body: ", string(body))

		//log.Println(req.PushData)
		if err := c.ShouldBind(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		//更新列表数据
		isUpdate, err := adver_service.UpdateAdverPackage(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}
		utils.Success(c, "", "ok")
		return
	}
	packageId, _ := strconv.Atoi(c.Query("id"))
	packageRes, err := adver_service.GetAdvertPackageInfo(int64(packageId))
	if err != nil {
		utils.Fail(c, err, "获取数据失败")
		return
	}
	//查询关联的类目
	projectList, err := adver_service.GetAllProjectList(int64(packageId))
	if err != nil {
	}
	packageRes.Extra = map[string]interface{}{
		"details": projectList,
	}
	//if packageRes.Id == 0 {
	//	packageRes = []
	//}
	res := gin.H{
		"packageInfo": packageRes,
	}
	utils.Success(c, res, "ok")
}

// 创建广告
func (adver *Adver) CreateAdver(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateAdverReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		//请求创建广告信息
		InsertId, err := adver_service.CreateAdver(&req)
		if err != nil {
			utils.Fail(c, err, "创建失败")
			return
		}
		//成功返回对应的状态信息
		utils.Success(c, InsertId, "ok")
		return
	}

	utils.Success(c, "", "ok")
}

// 更新广告设置
func (adver *Adver) UpdateAdver(c *gin.Context) {
	//判断是否为POST请求设置
	if c.Request.Method == "POST" {
		var req models.UpdateAdverReq

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatalf("read body failed at Before,err:%s", err.Error())
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		//log.Println("read body: ", string(body))

		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := adver_service.UpdateAdver(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}
	//这里是查看获取当前的ID并发送GET请求处理方式
	//字符转数字
	adverId, _ := strconv.Atoi(c.Query("id"))
	//根据当前的ID获取广告的基础信息设置
	adverInfo, err := adver_service.GetAdverById(int64(adverId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}
	//判断图片是否为空
	if adverInfo.Pic != "" {
		adverInfo.Pic = utils.GetAdminFileUrl(adverInfo.Pic)
	}
	//返回结果集
	res := gin.H{
		"adverInfo": adverInfo,
	}
	utils.Success(c, res, "ok")
}

// 删除广告配置--这里需要根据对应的选项来删除广告
func (adver *Adver) DelAdver(c *gin.Context) {
	var req models.DeleteAdverReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	//调用删除配置的信息
	isDelete, err := adver_service.DeleteAdver(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}
