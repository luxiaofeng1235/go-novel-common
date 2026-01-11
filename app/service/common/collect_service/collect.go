package collect_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"regexp"
)

func GetCollectById(id int64) (collect *models.McCollect, err error) {
	err = global.DB.Model(models.McCollect{}).Where("id = ?", id).Find(&collect).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

// 获取章节内容
func GetChapterContent(chapterText, html string) (chapterContent string, err error) {
	if chapterText == "" {
		global.Collectlog.Errorf("%v", "采集小说规则不能为空")
		err = fmt.Errorf("%v", "采集小说规则不能为空")
		return
	}
	if html == "" {
		global.Collectlog.Errorf("%v", "小说html页面不能为空")
		err = fmt.Errorf("%v", "小说html页面不能为空")
		return
	}
	match := regexp.MustCompile(chapterText).FindStringSubmatch(html)
	if len(match) > 1 {
		chapterContent = match[1]
		return
	} else {
		global.Collectlog.Errorf("采集小说内容失败 html=%v chapterText %v", html, chapterText)
		return
	}
	return
}
