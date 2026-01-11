package class_service

import (
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
)

func GetClassList(req *models.BookTypeReq) (typeRes []*models.ClassTypeRes, err error) {
	bookType := req.BookType
	var typeList []*models.McClassType
	typeList, err = GetClassTypes()
	if err != nil {
		return
	}
	if len(typeList) <= 0 {
		return
	}
	for _, val := range typeList {
		classTypeId := val.Id
		classListRes := getClassList(bookType, classTypeId)
		typeInfo := &models.ClassTypeRes{
			Id:           classTypeId,
			TypeName:     val.TypeName,
			Sort:         val.Sort,
			ClassListRes: classListRes,
		}
		typeRes = append(typeRes, typeInfo)
	}
	return
}

func getClassList(bookType int, classTypeId int64) (list []*models.ClassListRes) {
	var err error
	classList, err := GetClassByClassType(bookType, classTypeId)
	if err != nil {
		global.Errlog.Errorf("%v", err.Error())
		return
	}
	if len(classList) <= 0 {
		return
	}
	for _, val := range classList {
		var count int64
		if val.ClassPic != "" {
			count = GetBookCountByClassId(val.Id)
		}
		res := &models.ClassListRes{
			ClassId:   val.Id,
			ClassName: val.ClassName,
			Count:     count,
			Pic:       utils.GetFileUrl(val.ClassPic),
		}
		list = append(list, res)
	}
	return
}
