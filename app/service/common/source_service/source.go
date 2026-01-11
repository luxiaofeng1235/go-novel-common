package source_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/utils"
	"regexp"
)

func GetSourceLastChapter(bookUrl string, collect *models.McCollect) (sourceName, updateTime, chapterTitle string, err error) {
	sourceName = collect.Title
	updateTimeReg := collect.UpdateReg
	chapterSectionReg := collect.ChapterSectionReg
	chapterUrlReg := collect.ChapterUrlReg
	var html string
	html, err = utils.GetHtml(bookUrl, "utf-8", 1, 0)
	if html == "" {
		err = fmt.Errorf("获取小说详情页面失败 bookUrl=%v", bookUrl)
		return
	}
	matchUpdateTime := regexp.MustCompile(updateTimeReg).FindStringSubmatch(html)
	if len(matchUpdateTime) > 0 {
		updateTime = matchUpdateTime[1]
	} else {
		err = fmt.Errorf("%v", "获取小说最新更新时间失败")
		return
	}

	var chapters []*models.CollectChapterInfo
	matchChapterSection := regexp.MustCompile(chapterSectionReg).FindStringSubmatch(html)
	if len(matchChapterSection) > 0 {
		ulContent := matchChapterSection[0]
		liMatches := regexp.MustCompile(chapterUrlReg).FindAllStringSubmatch(ulContent, -1)
		for _, liMatch := range liMatches {
			if len(liMatch) > 2 {
				href := liMatch[1]
				chapterName := liMatch[2]
				chapter := &models.CollectChapterInfo{
					ChapterTitle: chapterName,
					ChapterLink:  href,
				}
				chapters = append(chapters, chapter)
			}
		}
	} else {
		err = fmt.Errorf("%v", "获取小说章节区间失败")
		return
	}
	chapter := chapters[len(chapters)-1]
	chapterTitle = chapter.ChapterTitle
	return
}
