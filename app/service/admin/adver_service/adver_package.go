package adver_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/admin/setting_service"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

/*
* @note 获取包的基本信息
* @param id int64 包ID
* @return object ,total , err
 */
func GetAdvertPackageInfo(id int64) (adver *models.AdvertProjectInfoRes, err error) {
	err = global.DB.Model(&models.McAdverPackage{}).Where("id", id).Debug().Find(&adver).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

// 根据对应字段来获取关联的字段信息
func GetPackageByCondition(project_name, package_name string) (myString []int64, err error) {
	if project_name == "" && package_name == "" {
		return
	}
	var list []*models.McAdverPackage
	////查看当前满足条件的关联信息
	db := global.DB.Model(&models.McAdverPackage{})
	if project_name != "" {
		db = db.Where("project_name LIKE ?", "%"+project_name+"%")
	}
	if package_name != "" {
		db = db.Where("package_name LIKE ?", "%"+package_name+"%")
	}
	//查询关联的数据信息
	err = db.Debug().Find(&list).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	//重新定义切片，直接返回
	var slice []int64
	for _, v := range list {
		slice = append(slice, v.Id)
	}
	//处理相关数据的返回IDS
	return slice, nil
}

/*
* @note 获取所有的项目列表
* @param list array
* @return object ,total , err
 */
func GetALlAdminPackageList() (list []*models.McAdverPackage, err error) {
	db := global.DB.Model(&models.McAdverPackage{}).Order("id asc")
	err = db.Debug().Find(&list).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

/*
* @note 查询广告项目列表数据
* @param req object 项目名称
* @return object ,total , err
 */
func AdverProjectListSearch(req *models.AdverProjectReq) (list []*models.McAdverProject, total int64, err error) {
	db := global.DB.Model(&models.McAdverProject{}).Order("id desc")
	pageNum := req.PageNum
	pageSize := req.PageSize
	packageId := req.PackageId
	if packageId > 0 {
		db = db.Where("package_id = ?", packageId)
	}
	//统计总条数
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Debug().Find(&list).Error
	} else {
		err = db.Offset(pageNum).Debug().Find(&list).Error
	}
	if len(list) <= 0 {
		return
	}
	return list, total, err
}

/*
* @note 广告包管理搜索
* @param req object 搜索参数
* @return object ,total , err
 */
func AdvertSettingPackageList(req *models.AdverPackageReq) (list []*models.AdverSettingRes, total int64, err error) {
	db := global.DB.Model(&models.McAdverPackage{}).Order("id desc")
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Debug().Find(&list).Error
	} else {
		err = db.Offset(pageNum).Debug().Find(&list).Error
	}
	if len(list) <= 0 {
		return
	}
	for _, packageInfo := range list {
		agreementInfo, _ := setting_service.GetAgreementInfo(packageInfo.Id)
		if agreementInfo != nil {
			packageInfo.UserAgreementUrl = agreementInfo.UserAgreementUrl //用户协议
			packageInfo.PrivacyUrl = agreementInfo.PrivacyUrl             //隐私协议
		}
	}
	return list, total, err
}

/*
* @note 获取广告的包的基本信息
* @param id int64 包ID
* @return object ,total , err
 */
func GetAdverProjectInfo(id int64) (adver *models.McAdverProject, err error) {
	err = global.DB.Model(&models.McAdverProject{}).Where("id", id).Find(&adver).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

/*
* @note 获取项目列表数据
* @param id int64 包ID
* @return object ,total , err
 */
func GetAllProjectList(packge_id int64) (list []*models.McAdverProject, err error) {
	if packge_id <= 0 {
		return
	}
	err = global.DB.Model(&models.McAdverProject{}).Where("package_id", packge_id).Debug().Find(&list).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

/*
* @note 获取所有的类型列表
* @param id int64 包ID
* @return object ,total , err
 */
func GetALlProjectTypeList() (list []*models.AdverProjectTypeListRes, err error) {
	err = global.DB.Model(&models.McAdver{}).Order("id asc").Debug().Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

/*
* @note 保存项目名称信息
* @param req object 搜索参数
* @return object ,total , err
 */
func SaveAdverProject(req *models.UpdateAdverProjectReq) (res bool, err error) {
	id := req.ProjectId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不能为空")
		return
	}
	//广告类型
	adverType := strings.TrimSpace(req.AdverType)
	if adverType == "" {
		err = fmt.Errorf("%v", "类型不能为空")
		return
	}
	//广告对应的值
	adverValueString := strings.TrimSpace(req.AdverValueString)
	if adverValueString == "" {
		err = fmt.Errorf("%v", "请填写广告ID")
		return
	}
	var mapData = make(map[string]interface{})
	mapData["adver_type"] = adverType
	mapData["adver_value_string"] = adverValueString
	mapData["uptime"] = utils.GetUnix()
	if err = global.DB.Model(models.McAdverProject{}).Where("id", id).Debug().Updates(mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

/*
* @note 广告包管理搜索
* @param req object 搜索参数
* @return object ,total , err
 */
func AdvertPackageListSearch(req *models.AdverPackageReq) (list []*models.McAdverPackage, total int64, err error) {
	db := global.DB.Model(&models.McAdverPackage{}).Order("id desc")
	projectName := strings.TrimSpace(req.ProjectName)
	//项目名称搜索
	if projectName != "" {
		db = db.Where("project_name  = ?", projectName)
	}
	//appid搜索
	appId := strings.TrimSpace(req.AppId)
	if appId != "" {
		db = db.Where("app_id = ?", appId)
	}
	//包名搜索
	packageName := strings.TrimSpace(req.PackageName)
	if packageName != "" {
		db = db.Where("package_name  = ?", packageName)
	}
	deviceType := strings.TrimSpace(req.DeviceType) //设备类型搜索
	if deviceType != "" {
		db = db.Where("device_type = ?", deviceType)
	}
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Debug().Find(&list).Error
	} else {
		err = db.Offset(pageNum).Debug().Find(&list).Error
	}
	if len(list) > 0 {
		return
	}
	return list, total, err
}

/*
* @note 获取所有的广告包列表-版本更新进行关联
* @return list,err
 */
func GetAllAdverPackage() (list []*models.McAdverPackage, err error) {
	db := global.DB.Model(&models.McAdverPackage{}).Debug().Order("id desc")
	err = db.Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, err
}

/*
* @note 创建包的基本信息
* @param req object 创建参数
* @return InsertId,err
 */
func CreateAdverPackage(req *models.CreateAdverPackageReq) (InsertId int64, err error) {
	projectName := strings.TrimSpace(req.ProjectName) //项目名称
	appId := strings.TrimSpace(req.AppId)             //appid
	packageName := strings.TrimSpace(req.PackageName) //包的名称
	deviceType := req.DeviceType                      //设备类型
	mark := req.Mark                                  //渠道号
	if projectName == "" {
		err = fmt.Errorf("%v", "项目名称")
		return
	}
	if appId == "" {
		err = fmt.Errorf("%v", "appid不能为空")
		return
	}
	if packageName == "" {
		err = fmt.Errorf("%v", "包名不能为空")
		return
	}
	if deviceType == "" {
		err = fmt.Errorf("%v", "渠道不为空")
		return
	} else {
		if deviceType != "android" && deviceType != "ios" {
			err = fmt.Errorf("%v", "设备类型只能为ios或android")
			return
		}
	}

	//待添加的数据
	adverPackge := models.McAdverPackage{
		ProjectName: projectName,
		AppId:       appId,
		CreateUser:  req.CreateUser,
		PackageName: packageName,
		DeviceType:  deviceType,
		Mark:        mark,
		AddTime:     utils.GetUnix(),
	}
	if err = global.DB.Create(&adverPackge).Error; err != nil {
		return 0, err
	}

	Id := adverPackge.Id

	///////////////////////////////////////关联处理先关的数据信息
	//处理数据信息
	pushData := req.PushData
	//处理同步广告子元素信息
	err = SyncProjectData(pushData, Id)

	return Id, nil
}

/*
* @note 删除项目表的基础信息关联
* @param package_id int 包ID
* @return res 结果集合 ,err 错误信息
 */
func DeleteProjectList(package_id int64) (res bool, err error) {
	if package_id <= 0 {
		return
	}
	db := global.DB.Debug().Model(models.McAdverProject{}).Where("package_id", package_id)
	var total int64

	err = db.Count(&total).Error
	if total > 0 {
		fmt.Println("***************************删除关联无用子表单数据**************************")
		err = db.Delete(&models.McAdverProject{}).Error
		if err != nil {
			return res, err
		}
		return res, nil
	} else {
		return res, nil
	}
}

// 同步包的基础信息
// type []MyAjaxModelsLis
func SyncProjectData(mList []models.MyAjaxModelsList, package_id int64) (err error) {
	if package_id < 0 || len(mList) == 0 {
		return nil
	}
	////删除关联的操作数据信息
	_, err = DeleteProjectList(package_id) //删除关联的操作数据
	insertData := make([]map[string]interface{}, 0, len(mList))

	for _, value := range mList {
		//重新定义切片来进行循环插入相关数据
		insertData = append(insertData, map[string]interface{}{
			"package_id":         package_id,
			"adver_type":         value.AdverType,
			"adver_type_name":    value.AdverTypeName,
			"adver_value_string": value.AdverValueString,
			"addtime":            utils.GetUnix(),
		})
	}
	err = global.DB.Debug().Model(&models.McAdverProject{}).Debug().CreateInBatches(&insertData, 100).Error
	if err != nil {
		return err
	}
	return
}

/*
* @note 更新广告信息设置
* @param req object 传入的参数
* @return res 结果集合 ,err 错误信息
 */

func UpdateAdverPackage(req *models.UpdateAdverPackageReq) (res bool, err error) {
	Id := req.PackageId
	if Id <= 0 {
		err = fmt.Errorf("%v", "id不能为空")
		return
	}
	projectName := strings.TrimSpace(req.ProjectName)
	appId := strings.TrimSpace(req.AppId)
	packageName := strings.TrimSpace(req.PackageName)
	deviceType := req.DeviceType
	mark := req.Mark //渠道号

	///////////////////////////////////////关联处理先关的数据信息
	//处理数据信息
	pushData := req.PushData
	//处理同步广告子元素信息
	err = SyncProjectData(pushData, Id)
	//fmt.Println(pushData)
	//for _, val := range pushData {
	//	log.Println(val.AdverValueString, val.AdverType, val.AdverTypeName)
	//}

	if projectName == "" {
		err = fmt.Errorf("%v", "项目名称不能为空")
		return
	}
	if appId == "" {
		err = fmt.Errorf("%v", "appid不能为空")
		return
	}
	if packageName == "" {
		err = fmt.Errorf("%v", "包名不能为空")
		return
	}
	if deviceType == "" {
		err = fmt.Errorf("%v", "渠道不为空")
		return
	} else {
		if deviceType != "android" && deviceType != "ios" {
			err = fmt.Errorf("%v", "设备类型只能为ios或android")
			return
		}
	}
	var mapData = make(map[string]interface{})
	mapData["project_name"] = projectName
	mapData["app_id"] = appId
	mapData["package_name"] = packageName
	mapData["device_type"] = deviceType
	mapData["mark"] = mark
	mapData["uptime"] = utils.GetUnix()
	if err = global.DB.Model(models.McAdverPackage{}).Where("id", Id).Debug().Updates(&mapData).Error; err != nil {
		global.Sqllog.Error(err)
		return false, err
	}
	return true, nil
}

/*
* @note 删除包的相关信息
* @param req object 传入的参数
* @return res 结果集合 ,err 错误信息
 */

func DeleteAdverPackage(req *models.DeleteAdverPackageReq) (res bool, err error) {
	ids := req.PackageIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Debug().Delete(&models.McAdverPackage{}).Error
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return
	}
	return true, nil
}
