package adver_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

/*
* @note 根据当前ID获取广告信息
* @param id integer 广告ID
* @return object,err
 */
func GetAdverById(id int64) (adver *models.McAdver, err error) {
	err = global.DB.Model(models.McAdver{}).Where("id", id).First(&adver).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

/*
* @note 检查广告的唯一性
* @param advertName string 广告名称
* @param id int 广告ID
* @return interger
 */
func CheckAdverNameUnique(adverName string, id int64) bool {
	var count int64
	model := global.DB.Model(models.McAdver{}).Where("adver_name = ?", adverName)
	if id != 0 {
		model = model.Where("id!=?", id)
	}
	model.Count(&count)
	return count == 0
}

/*
* @note 广告列表搜索
* @param req object 搜索参数
* @return object ,total , err
 */
func AdverListSearch(req *models.AdverListReq) (list []*models.McAdver, total int64, err error) {
	db := global.DB.Model(&models.McAdver{}).Order("id desc")

	adverName := strings.TrimSpace(req.AdverName)
	if adverName != "" {
		db = db.Where("adver_name = ?", adverName)
	}

	adverCode := strings.TrimSpace(req.AdverCode)
	if adverCode != "" {
		db = db.Where("adver_code = ?", adverCode)
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
	for _, adver := range list {
		adver.Pic = utils.GetAdminFileUrl(adver.Pic)
	}
	return list, total, err
}

/*
* @note 判断是否在数组中的元素信息
* @param items object 切片对象
* @param item integer 元素值
* @return bool
 */
func IsContainInt(items []int, item int) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

/*
* @note 创建广告设置
* @param req object 传入的参数
* @return InsertId ,err
 */
func CreateAdver(req *models.CreateAdverReq) (InsertId int64, err error) {
	adverType := req.AdverType
	isLocal := req.IsLocal
	weight := req.Weight
	AdverPosition := req.AdverPosition //广告浏览设置
	AdverNum := req.AdverNum           //设置广告的时间或总数
	ErrorNum := req.ErrorNum           //错误次数
	AdverTime := req.AdverTime         //免广告时间
	adverName := strings.TrimSpace(req.AdverName)
	EveryShowNum := req.EveryShowNum                 //每天显示的激励广告最大次数
	ReadTurningSHowTimes := req.ReadTurningSHowTimes //阅读中翻页的的次数
	ReadRollTime := req.ReadRollTime                 //阅读中滚动的时间
	AdverValueIos := req.AdverValueIos               //ios广告
	if adverName == "" {
		err = fmt.Errorf("%v", "广告名称不能为空")
		return
	}

	//判断是否在对象中
	targetArr := []int{1, 2, 3, 4, 5, 9}
	index := IsContainInt(targetArr, adverType)
	if index {
		//if !(weight > 1 && weight <= 100) {
		//	err = fmt.Errorf("%v", "权重取值范围为1-100")
		//	return
		//}
	}

	pic := strings.TrimSpace(req.Pic)
	adverLink := strings.TrimSpace(req.AdverLink)
	if req.IsLocal == 1 {
		if pic == "" {
			err = fmt.Errorf("%v", "广告图片不能为空")
			return
		}
		if adverLink == "" {
			err = fmt.Errorf("%v", "广告链接不能为空")
			return
		}
	}
	advertValue := strings.TrimSpace(req.AdverValue)

	if !CheckAdverNameUnique(adverName, 0) {
		err = fmt.Errorf("%v", "广告名称已经存在")
		return
	}

	status := req.Status
	adver := models.McAdver{
		AdverType:            adverType,
		IsLocal:              isLocal,
		AdverName:            adverName,
		Pic:                  pic,
		Weight:               weight,
		AdverLink:            adverLink,
		AdverValue:           advertValue,
		AdverValueIos:        AdverValueIos,
		Status:               status,
		AdverPosition:        AdverPosition,
		AdverNum:             AdverNum,
		AdverTime:            AdverTime,
		ErrorNum:             ErrorNum,
		EveryShowNum:         EveryShowNum,
		ReadTurningSHowTimes: ReadTurningSHowTimes,
		ReadRollTime:         ReadRollTime,
		Addtime:              utils.GetUnix(),
	}

	if err = global.DB.Create(&adver).Error; err != nil {
		return 0, err
	}

	return adver.Id, nil
}

/*
* @note 更新广告信息设置
* @param req object 传入的参数
* @return res 结果集合 ,err 错误信息
 */
func UpdateAdver(req *models.UpdateAdverReq) (res bool, err error) {
	id := req.AdverId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	weight := req.Weight

	//判断判断本地广告是否在对象中
	targetArr := []int{1, 2, 3, 4, 5, 9}
	index := IsContainInt(targetArr, req.AdverType)
	if index {
		//if !(weight > 1 && weight <= 100) {
		//	err = fmt.Errorf("%v", "权重取值范围为1-100")
		//	return
		//}
	}
	adverName := strings.TrimSpace(req.AdverName)
	pic := strings.TrimSpace(req.Pic)
	adverLink := strings.TrimSpace(req.AdverLink)
	if !CheckAdverNameUnique(adverName, id) {
		err = fmt.Errorf("%v", "广告名称已经存在")
		return
	}

	if req.IsLocal == 1 {
		if pic == "" {
			err = fmt.Errorf("%v", "广告图片不能为空")
			return
		}
		if adverLink == "" {
			err = fmt.Errorf("%v", "广告链接不能为空")
			return
		}
	}

	//校验是否可以选择其他类型的 ,如果已经选择了就不能选择其他的啦
	var total int64
	err = global.DB.Model(models.McAdver{}).Debug().Where("adver_type=? and id!=?", req.AdverType, req.AdverId).Count(&total).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	if total > 0 {
		err = fmt.Errorf("%v", "该类型已存在，请勿重复选择")
		return
	}

	adverValue := strings.TrimSpace(req.AdverValue)
	var mapData = make(map[string]interface{})
	mapData["adver_name"] = adverName
	if pic != "" {
		mapData["pic"] = pic
	}
	mapData["weight"] = weight
	mapData["adver_link"] = adverLink
	mapData["adver_value"] = adverValue
	mapData["adver_value_ios"] = req.AdverValueIos
	mapData["status"] = req.Status
	mapData["adver_position"] = req.AdverPosition
	mapData["error_num"] = req.ErrorNum
	mapData["adver_num"] = req.AdverNum
	mapData["adver_time"] = req.AdverTime
	mapData["read_turning_show_times"] = req.ReadTurningSHowTimes
	mapData["read_roll_time"] = req.ReadRollTime
	mapData["uptime"] = utils.GetUnix()
	mapData["every_show_num"] = req.EveryShowNum
	if err = global.DB.Model(models.McAdver{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

/*
* @note 删除广告信息设置
* @param req object 传入的参数
* @return res 结果集合 ,err 错误信息
 */
func DeleteAdver(req *models.DeleteAdverReq) (res bool, err error) {
	id := req.AdverId
	if id <= 0 {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return

	}
	err = global.DB.Where("id = ?", id).Delete(&models.McAdver{}).Error
	if err != nil {
		return
	}
	return true, nil
}
