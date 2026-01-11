package class_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"strings"
)

func GetTypeById(id int64) (classType *models.McClassType, err error) {
	err = global.DB.Model(models.McClassType{}).Where("id", id).First(&classType).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetTypeList() (list []*models.McClassType, err error) {
	err = global.DB.Model(models.McClassType{}).Where("status = 1").Order("sort desc").Find(&list).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func TypeListSearch(req *models.TypeListReq) (list []*models.McClassType, total int64, err error) {
	db := global.DB.Model(&models.McClassType{}).Order("id desc")

	typeName := strings.TrimSpace(req.TypeName)
	if typeName != "" {
		db = db.Where("type_name = ?", typeName)
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
	return list, total, err
}

func CreateType(req *models.CreateTypeReq) (InsertId int64, err error) {
	if req.TypeName == "" {
		err = fmt.Errorf("%v", "类型名称不能为空")
		return
	}

	classType := models.McClassType{
		TypeName: req.TypeName,
		Status:   req.Status,
		Sort:     req.Sort,
	}

	if err = global.DB.Create(&classType).Error; err != nil {
		return 0, err
	}

	return classType.Id, nil
}

func UpdateType(req *models.UpdateTypeReq) (res bool, err error) {
	id := req.TypeId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	if req.TypeName == "" {
		err = fmt.Errorf("%v", "类型名称不能为空")
		return
	}
	var mapData = make(map[string]interface{})
	mapData["type_name"] = req.TypeName
	mapData["status"] = req.Status
	mapData["sort"] = req.Sort
	if err = global.DB.Model(models.McClassType{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DeleteType(req *models.DeleteTypeReq) (res bool, err error) {
	ids := req.TypeIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.McClassType{}).Error
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return
	}
	return true, nil
}
