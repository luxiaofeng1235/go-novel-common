package rank_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetRankById(id int64) (tag *models.McRank, err error) {
	err = global.DB.Model(models.McRank{}).Where("id", id).First(&tag).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func CheckRankNameUnique(tagName string, id int64) bool {
	var count int64
	model := global.DB.Model(models.McRank{}).Where("rank_name = ?", tagName)
	if id != 0 {
		model = model.Where("id!=?", id)
	}
	model.Count(&count)
	return count == 0
}

func CheckRankCodeUnique(tagCode string, id int64) bool {
	var count int64
	model := global.DB.Model(models.McRank{}).Where("rank_code = ?", tagCode)
	if id != 0 {
		model = model.Where("id!=?", id)
	}
	model.Count(&count)
	return count == 0
}

func RankListSearch(req *models.RankListReq) (list []*models.McRank, total int64, err error) {
	db := global.DB.Model(&models.McRank{}).Order("id desc")

	rankName := strings.TrimSpace(req.RankName)
	if rankName != "" {
		db = db.Where("rank_name = ?", rankName)
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

func CreateRank(req *models.CreateRankReq) (InsertId int64, err error) {
	rankName := strings.TrimSpace(req.RankName)
	if rankName == "" {
		err = fmt.Errorf("%v", "排行榜名称不能为空")
		return
	}

	rankCode := strings.TrimSpace(req.RankCode)
	if rankCode == "" {
		err = fmt.Errorf("%v", "排行榜标识不能为空")
		return
	}

	if !CheckRankNameUnique(rankName, 0) {
		err = fmt.Errorf("%v", "排行榜名称已经存在")
		return
	}

	if !CheckRankCodeUnique(rankName, 0) {
		err = fmt.Errorf("%v", "排行榜标识已经存在")
		return
	}
	sort := req.Sort
	status := req.Status
	rank := models.McRank{
		RankName: rankName,
		RankCode: rankCode,
		Sort:     sort,
		Status:   status,
		Addtime:  utils.GetUnix(),
	}

	if err = global.DB.Create(&rank).Error; err != nil {
		return 0, err
	}

	return rank.Id, nil
}

func UpdateRank(req *models.UpdateRankReq) (res bool, err error) {
	id := req.RankId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	rankName := strings.TrimSpace(req.RankName)
	if rankName == "" {
		err = fmt.Errorf("%v", "排行榜名称不能为空")
		return
	}
	if !CheckRankNameUnique(rankName, id) {
		err = fmt.Errorf("%v", "排行榜名称已经存在")
		return
	}

	rankCode := strings.TrimSpace(req.RankCode)
	if rankCode == "" {
		err = fmt.Errorf("%v", "排行榜标识不能为空")
		return
	}
	if !CheckRankCodeUnique(rankCode, id) {
		err = fmt.Errorf("%v", "排行榜标识已经存在")
		return
	}
	var mapData = make(map[string]interface{})
	mapData["rank_name"] = rankName
	mapData["rank_code"] = rankCode
	mapData["sort"] = req.Sort
	mapData["status"] = req.Status
	mapData["uptime"] = utils.GetUnix()
	if err = global.DB.Model(models.McRank{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DeleteRank(req *models.DeleteRankReq) (res bool, err error) {
	id := req.RankId
	if id <= 0 {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return

	}
	err = global.DB.Where("id = ?", id).Delete(&models.McRank{}).Error
	if err != nil {
		return
	}
	return true, nil
}
