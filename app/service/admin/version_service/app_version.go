package version_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/admin/adver_service"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetVersionById(id int64) (version *models.McAppVersion, err error) {
	err = global.DB.Model(models.McAppVersion{}).Where("id", id).First(&version).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	//处理图片问题
	if version.ForceUrl != "" {
		//只有没有发现里面有http才进行替换拼接，不然会有问题
		if strings.Contains(version.ForceUrl, "https") == false || strings.Contains(version.ForceUrl, "http") == false {
			version.ForceUrl = utils.GetApkFileUrl(version.ForceUrl)
		}
	}
	return
}

// 版本列表哦诉讼
func VersionListSearch(req *models.AppVersionListReq) (list []*models.AppVersionAdminListRes, total int64, err error) {
	db := global.DB.Model(&models.McAppVersion{}).Order("id desc").Debug()

	device := strings.TrimSpace(req.Device)
	if device != "" {
		db = db.Where("device = ?", device)
	}
	packageName := strings.TrimSpace(req.PackageName)
	projectName := strings.TrimSpace(req.ProjectName)
	if packageName != "" || projectName != "" {
		var ids []int64
		ids, err = adver_service.GetPackageByCondition(projectName, packageName)
		if err != nil {
			global.Sqllog.Errorf("VersionListSearch GetPackageByCondition err:%s", err.Error())
		}
		//如果有匹配到就直接按照这个查询即可
		if ids != nil {
			//如果不为空，按照数组的形式来进行分割，转换
			db = db.Where("package_id in(?)", ids)
		} else {
			db = db.Where("package_id in (-1)") //默认没查到给一个空的，防止所有数据都出来
		}
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	if len(list) <= 0 {
		return
	}
	for _, version := range list {
		packageId := version.PackageId
		if packageId <= 0 {
			continue
		}
		//获取包关联的相关信息
		packageInfo, _ := adver_service.GetAdvertPackageInfo(int64(packageId))
		if packageInfo != nil {
			version.ProjectName = packageInfo.ProjectName //项目名称
			version.PackageName = packageInfo.PackageName //包名称
			version.AppId = packageInfo.AppId             //appid
			version.Mark = packageInfo.Mark               //渠道号
		} else {
			version.ProjectName = ""
			version.ProjectName = ""
			version.AppId = ""
			version.Mark = ""
		}
	}
	return list, total, err
}

// 创建接口
func CreateVersion(req *models.CreateAppVersionReq) (insertId int64, err error) {
	device := req.Device
	version := req.Version
	packageId := req.PackageId
	downUrl := strings.TrimSpace(req.DownUrl)
	isForce := req.IsForce
	updateText := strings.TrimSpace(req.UpdateText)
	commentStatus := req.CommentStatus
	copyrightStatus := req.CopyrightStatus
	forceType := req.ForceType
	forceFileSize := req.ForceFileSize
	forceUrl := strings.TrimSpace(req.ForceUrl)

	//判断当前是否有添加的重复项目
	IsCount, err := GetVersionByPackageId(int64(packageId))
	if err != nil {
		global.Sqllog.Info("%v", err.Error())
	}
	if IsCount != 0 {
		err = fmt.Errorf("%v", "已添加当前项目，请更换其他项目添加~")
		return
	}
	//创建接口数据信息
	appData := models.McAppVersion{
		Device:          device,
		Version:         version,
		UpdateText:      updateText,
		PackageId:       packageId,
		IsForce:         isForce,
		CommentStatus:   commentStatus,
		CopyrightStatus: copyrightStatus,
		ForceType:       forceType,
		ForceUrl:        forceUrl,
		DownUrl:         downUrl,
		ForceFileSize:   forceFileSize,
		Addtime:         utils.GetUnix(),
	}
	err = global.DB.Create(&appData).Error
	if err != nil {
		return 0, err
	}
	return appData.Id, nil
}

/*
* @note 删除包的相关信息
* @param req object 传入的参数
* @return res 结果集合 ,err 错误信息
 */

func DeleteVersionRes(req *models.DeleteVersionReq) (res bool, err error) {
	ids := req.VersionIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Debug().Delete(&models.McAppVersion{}).Error
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return
	}
	return true, nil
}

func UpdateVersion(req *models.UpdateAppVersionReq) (res bool, err error) {
	id := req.VersionId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	version := strings.TrimSpace(req.Version)
	downUrl := strings.TrimSpace(req.DownUrl)
	isForce := req.IsForce
	commentStatus := req.CommentStatus
	packageId := req.PackageId //应用ID
	updateText := strings.TrimSpace(req.UpdateText)
	copyrightStatus := req.CopyrightStatus
	forceType := req.ForceType
	forceUrl := req.ForceUrl
	forceFileSize := req.ForceFileSize

	var mapData = make(map[string]interface{})
	mapData["version"] = version
	mapData["down_url"] = downUrl
	mapData["is_force"] = isForce
	mapData["update_text"] = updateText
	mapData["uptime"] = utils.GetUnix()
	mapData["comment_status"] = commentStatus
	mapData["copyright_status"] = copyrightStatus //版全面状态
	mapData["package_id"] = packageId
	mapData["force_type"] = forceType
	mapData["force_url"] = forceUrl
	mapData["force_file_size"] = forceFileSize
	if err = global.DB.Model(models.McAppVersion{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

// 获取有当前添加的版本号
func GetVersionByPackageId(package_id int64) (total int64, err error) {
	db := global.DB.Model(models.McAppVersion{}).Debug()
	db = db.Where("package_id = ?", package_id)
	err = db.Count(&total).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return 0, err
	}
	return total, nil
}
