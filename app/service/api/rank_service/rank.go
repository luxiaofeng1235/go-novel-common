package rank_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetRankList() (ranks []*models.McRank, err error) {
	err = global.DB.Model(models.McRank{}).Order("sort asc").Where("status = 1").Find(&ranks).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
