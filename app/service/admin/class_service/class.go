package class_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetClassById(id int64) (bookClass *models.McBookClass, err error) {
	err = global.DB.Model(models.McBookClass{}).Where("id", id).First(&bookClass).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetBookCountByClassId(classId int64) (count int64) {
	err := global.DB.Model(models.McBook{}).Where("cid = ?", classId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetClassList() (list []*models.McBookClass, err error) {
	err = global.DB.Model(models.McBookClass{}).Where("status = 1").Order("sort desc").Find(&list).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetClassNameById(classId int64) (className string) {
	var err error
	err = global.DB.Model(models.McBookClass{}).Select("class_name").Where("id", classId).Scan(&className).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func ClassListSearch(req *models.ClassListReq) (list []*models.McBookClass, total int64, err error) {
	db := global.DB.Model(&models.McBookClass{}).Order("id desc")

	className := strings.TrimSpace(req.ClassName)
	if className != "" {
		db = db.Where("class_name = ?", className)
	}

	bookType := strings.TrimSpace(req.BookType)
	if bookType != "" {
		db = db.Where("book_type = ?", bookType)
	}

	typeId := strings.TrimSpace(req.TypeId)
	if typeId != "" {
		db = db.Where("type_id = ?", typeId)
	}

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

	for _, val := range list {
		val.ClassPic = utils.GetAdminFileUrl(val.ClassPic)
		val.BookCount = GetBookCountByClassId(val.Id)
	}
	return list, total, err
}

func CreateClass(req *models.CreateClassReq) (InsertId int64, err error) {
	bookType := req.BookType
	if bookType <= 0 {
		err = fmt.Errorf("%v", "阅读类型不能为空")
		return
	}
	className := strings.TrimSpace(req.ClassName)
	if className == "" {
		err = fmt.Errorf("%v", "分类名称不能为空")
		return
	}
	typeId := req.TypeId
	if typeId <= 0 {
		err = fmt.Errorf("%v", "分类类型不能为空")
		return
	}
	classPic := strings.TrimSpace(req.ClassPic)
	bookId := req.BookId
	if bookId > 0 {
		classPic = strings.TrimSpace(getBookPicById(bookId))
	}
	status := req.Status
	class := models.McBookClass{
		ClassName: req.ClassName,
		BookType:  req.BookType,
		ClassPic:  classPic,
		Sort:      req.Sort,
		TypeId:    req.TypeId,
		Status:    status,
	}

	if err = global.DB.Create(&class).Error; err != nil {
		return 0, err
	}

	return class.Id, nil
}

func UpdateClass(req *models.UpdateClassReq) (res bool, err error) {
	id := req.ClassId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	bookType := req.BookType
	if bookType <= 0 {
		err = fmt.Errorf("%v", "阅读类型不能为空")
		return
	}
	className := strings.TrimSpace(req.ClassName)
	if className == "" {
		err = fmt.Errorf("%v", "分类名称不能为空")
		return
	}
	typeId := req.TypeId
	if typeId <= 0 {
		err = fmt.Errorf("%v", "分类类型不能为空")
		return
	}
	status := req.Status
	classPic := strings.TrimSpace(req.ClassPic)
	bookId := req.BookId
	if bookId > 0 {
		classPic = strings.TrimSpace(getBookPicById(bookId))
	}
	var mapData = make(map[string]interface{})
	mapData["book_type"] = bookType
	mapData["class_name"] = className
	mapData["status"] = status
	if classPic != "" {
		mapData["class_pic"] = classPic
	}
	mapData["type_id"] = typeId
	mapData["sort"] = req.Sort
	if err = global.DB.Model(models.McBookClass{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DeleteClass(req *models.DeleteClassReq) (res bool, err error) {
	ids := req.ClassIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.McBookClass{}).Error
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return
	}
	return true, nil
}

func AssignClass(req *models.AssignClassReq) (res bool, err error) {
	ids := req.ClassIds
	if len(ids) > 0 {
		var classList []*models.McBookClass
		err = global.DB.Where("id in(?)", ids).Find(&classList).Error
		if err != nil {
			return
		}
		for _, class := range classList {
			class.ClassName = strings.TrimSpace(class.ClassName)
			data := make(map[string]interface{})
			data["book_type"] = class.BookType
			data["cid"] = class.Id
			data["class_name"] = class.ClassName
			likeName := "%" + class.ClassName + "%"
			global.DB.Model(models.McBook{}).Where("book_name like ? or tags like ?", likeName, likeName).Updates(data)
		}
	} else {
		err = fmt.Errorf("%v", "归类失败，参数错误")
		return
	}
	return true, nil
}
