package tag_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetTagById(id int64) (tag *models.McTag, err error) {
	err = global.DB.Model(models.McTag{}).Where("id", id).First(&tag).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetTagList() (list []*models.McTag, err error) {
	err = global.DB.Model(models.McTag{}).Where("status = 1").Order("sort asc").Find(&list).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetTagBySex() (tagList []*models.McTag, err error) {
	tagList, err = GetTagList()
	if err != nil {
		err = fmt.Errorf("获取标签列表失败 err=%v", err.Error())
		return
	}
	if len(tagList) > 0 {
		for _, val := range tagList {
			text := val.TagName
			if val.BookType == 1 {
				text = fmt.Sprintf("男生-%v", text)
			} else if val.BookType == 2 {
				text = fmt.Sprintf("女生-%v", text)
			}
			val.TagName = text
		}
	}
	return
}

func CheckTagNameUnique(tagName string, id int64) bool {
	var count int64
	model := global.DB.Model(models.McTag{}).Where("tag_name = ?", tagName).Debug()
	if id != 0 {
		model = model.Where("id!=?", id)
	}
	model.Count(&count)
	return count == 0
}

func TagListSearch(req *models.TagListReq) (list []*models.McTag, total int64, err error) {
	db := global.DB.Model(&models.McTag{}).Order("id desc").Debug()

	tagName := strings.TrimSpace(req.TagName)
	if tagName != "" {
		//用模糊搜索
		likeName := "%" + tagName + "%"
		db = db.Where("tag_name like ?", likeName)
	}

	bookType := strings.TrimSpace(req.BookType)
	if bookType != "" {
		db = db.Where("book_type = ?", bookType)
	}

	columnType := strings.TrimSpace(req.ColumnType)
	if columnType != "" {
		db = db.Where("column_type = ?", columnType)
	}

	isNew := strings.TrimSpace(req.IsNew)
	if isNew != "" {
		db = db.Where("is_new = ?", isNew)
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

	if len(list) > 0 {
		for _, val := range list {
			val.BookCount = GetBookCountByTagId(val.ColumnType, val.Id)
		}
	}
	return list, total, err
}

func CreateTag(req *models.CreateTagReq) (InsertId int64, err error) {
	tagName := strings.TrimSpace(req.TagName)
	if tagName == "" {
		err = fmt.Errorf("%v", "标签名称不能为空")
		return
	}

	if !CheckTagNameUnique(tagName, 0) {
		err = fmt.Errorf("%v", "标签已经存在")
		return
	}

	var isNew int
	if strings.Contains(tagName, "新书") {
		isNew = 1 //标记有新书就是1
	} else {
		isNew = 0 //否则就是0
	}

	bookType := req.BookType
	columnType := req.ColumnType
	sort := req.Sort
	status := req.Status
	tag := models.McTag{
		BookType:   bookType,
		ColumnType: columnType,
		TagName:    tagName,
		Sort:       sort,
		IsNew:      isNew, //是否为新书标签
		Status:     status,
		Addtime:    utils.GetUnix(),
	}

	if err = global.DB.Create(&tag).Error; err != nil {
		return 0, err
	}

	return tag.Id, nil
}

func UpdateTag(req *models.UpdateTagReq) (res bool, err error) {
	id := req.TagId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	tagName := strings.TrimSpace(req.TagName)
	if tagName == "" {
		err = fmt.Errorf("%v", "标签名称不能为空")
		return
	}
	if !CheckTagNameUnique(tagName, id) {
		err = fmt.Errorf("%v", "标签已经存在")
		return
	}
	var isNew int
	if strings.Contains(tagName, "新书") {
		isNew = 1 //标记有新书就是1
	} else {
		isNew = 0 //否则就是0
	}
	var mapData = make(map[string]interface{})
	mapData["book_type"] = req.BookType
	mapData["column_type"] = req.ColumnType
	mapData["tag_name"] = tagName
	mapData["is_new"] = isNew //是否为新书标签
	mapData["sort"] = req.Sort
	mapData["status"] = req.Status
	mapData["uptime"] = utils.GetUnix()
	if err = global.DB.Model(models.McTag{}).Debug().Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DeleteTag(req *models.DeleteTagReq) (res bool, err error) {
	id := req.TagId
	if id <= 0 {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return

	}
	err = global.DB.Where("id = ?", id).Delete(&models.McTag{}).Error
	if err != nil {
		return
	}
	return true, nil
}

func AssignTag(req *models.AssignTagReq) (res bool, err error) {
	ids := req.TagIds
	if len(ids) > 0 {
		var tagList []*models.McTag
		err = global.DB.Where("id in(?)", ids).Find(&tagList).Error
		if err != nil {
			return
		}
		for _, tag := range tagList {
			tagName := utils.TrimDotTag(tag.TagName)
			data := make(map[string]interface{})
			data["book_type"] = tag.BookType
			if tag.ColumnType == 1 {
				data["is_rec"] = 1
			} else if tag.ColumnType == 2 {
				data["is_hot"] = 1
			} else if tag.ColumnType == 3 {
				data["is_classic"] = 1
			}
			//判断是否为新书属性
			if tag.IsNew == 1 {
				data["is_new"] = 1
			}
			data["tid"] = tag.Id
			data["tag_name"] = tag.TagName
			likeName := "%" + tagName + "%"
			global.DB.Model(models.McBook{}).Where("book_name like ? or tags like ?", likeName, likeName).Debug().Updates(data)
		}
	} else {
		err = fmt.Errorf("%v", "归类失败，参数错误")
		return
	}
	return true, nil
}
