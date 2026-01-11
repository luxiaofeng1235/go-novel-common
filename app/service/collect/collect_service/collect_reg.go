package collect_service

import (
	"encoding/json"
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/common/setting_service"
	"go-novel/global"
	"go-novel/utils"
	"regexp"
	"strings"
)

func FieldContent(collect *models.McCollect, bookUrl string) (collecBookInfoRes *models.CollecBookInfoRes, err error) {
	var html string
	collectId := collect.Id
	html, err = utils.GetHtml(bookUrl, collect.Charset, collect.UrlComplete, utils.SleepSecond)
	if html == "" {
		utils.GetS5()
		err = BookUrlUnLock(collectId)
		if err != nil {
			global.Collectlog.Errorf("获取小说详情页面失败1 bookUrl=%v err=%v", bookUrl, err.Error())
		}
		err = fmt.Errorf("获取小说详情页面失败2 bookUrl=%v", bookUrl)
		return
	}
	collecBookInfoRes = new(models.CollecBookInfoRes)
	collecBookInfoRes.BookUrl = bookUrl
	categoryNameReg := collect.CategoryNameReg
	if categoryNameReg == "" {
		err = fmt.Errorf("%v", "获取分类名称正则不能为空")
		return
	}
	matchCate := regexp.MustCompile(categoryNameReg).FindStringSubmatch(html)
	if len(matchCate) > 0 {
		categoryName := matchCate[1]

		categoryId := collect.CategoryFixed
		categoryWay := collect.CategoryWay
		if categoryWay <= 0 && categoryId <= 0 {
			var categorys []*models.CategoryReg
			err = json.Unmarshal([]byte(collect.Categorys), &categorys)
			if err != nil {
				global.Collectlog.Errorf("解析collect分类出错 collectId=%v err=%v", collect.Id, err.Error())
				err = fmt.Errorf("解析collect分类出错 collectId=%v err=%v", collect.Id, err.Error())
				return
			}
			collecBookInfoRes.ClassId = utils.CategoryEquiv(categorys, categoryName)
		} else {
			collecBookInfoRes.ClassId = categoryId
		}
		if collecBookInfoRes.ClassId > 0 {
			collecBookInfoRes.CategoryName = categoryName
		}
	} else {
		global.Collectlog.Errorf("获取小说分类失败 bookUrl=%v", bookUrl)
		//err = fmt.Errorf("获取小说分类失败 bookUrl=%v", bookUrl)
		//return
	}
	bookNameReg := collect.BookNameReg
	if bookNameReg == "" {
		err = fmt.Errorf("%v", "获取小说名称正则不能为空")
		return
	}
	matchBookName := regexp.MustCompile(bookNameReg).FindStringSubmatch(html)
	if len(matchBookName) > 0 {
		collecBookInfoRes.BookName = matchBookName[1]
	} else {
		err = RmmoveCollectPageBooks(collectId, bookUrl)
		if err != nil {
			return
		}
		err = fmt.Errorf("获取小说名称失败 bookUrl=%v", bookUrl)
		if err != nil {
			return
		}
		return
	}

	authorReg := collect.AuthorReg
	if authorReg == "" {
		err = fmt.Errorf("%v", "获取小说作者正则不能为空")
		return
	}
	matchAuthor := regexp.MustCompile(authorReg).FindStringSubmatch(html)
	if len(matchAuthor) > 0 {
		collecBookInfoRes.Author = matchAuthor[1]
	} else {
		err = RmmoveCollectPageBooks(collectId, bookUrl)
		if err != nil {
			return
		}
		err = fmt.Errorf("获取小说作者失败 bookUrl=%v", bookUrl)
		if err != nil {
			return
		}
		return
	}

	picReg := collect.PicReg
	if picReg == "" {
		err = fmt.Errorf("%v", "获取小说作者正则不能为空")
		return
	}
	matchPic := regexp.MustCompile(picReg).FindStringSubmatch(html)
	if len(matchPic) > 0 {
		collecBookInfoRes.Pic = matchPic[1]
	} else {
		err = RmmoveCollectPageBooks(collectId, bookUrl)
		if err != nil {
			return
		}
		err = fmt.Errorf("获取小说图片失败 bookUrl=%v", bookUrl)
		if err != nil {
			return
		}
		return
	}
	pic := collecBookInfoRes.Pic

	if collect.PicLocal > 0 && pic != "" {
		var uploadBookPicPath string
		uploadBookPicPath, err = setting_service.GetValueByName(utils.UploadBookPicPath)
		if err != nil {
			err = fmt.Errorf("获取小说上传目录失败 bookUrl=%v", bookUrl)
			return
		}
		var filePath string
		filePath, err = utils.DownImg(collecBookInfoRes.BookName, collecBookInfoRes.Author, pic, uploadBookPicPath)
		if err != nil {
			err = fmt.Errorf("下载小说图片失败 %v", err.Error())
			return
		}
		collecBookInfoRes.Pic = strings.TrimLeft(filePath, ".")
	}
	//descText := `<h3 class="bookinfo_intro"><strong>\s*《(.*?)》.*?<\/strong><br \/>\s*([\s\S]*?)<br\/>`
	descReg := collect.DescReg
	if descReg == "" {
		err = fmt.Errorf("%v", "获取小说简介正则不能为空")
		return
	}
	matchDesc := regexp.MustCompile(collect.DescReg).FindStringSubmatch(html)
	if len(matchDesc) > 1 {
		collecBookInfoRes.Desc = strings.TrimSpace(matchDesc[1])
	}
	if len(matchDesc) > 2 {
		collecBookInfoRes.Desc = strings.TrimSpace(matchDesc[2])
	}
	if collecBookInfoRes.Desc == "" {
		collecBookInfoRes.Desc = collecBookInfoRes.BookName
	}

	serializeReg := collect.SerializeReg
	if serializeReg == "" {
		err = fmt.Errorf("%v", "获取小说更新状态正则不能为空")
		return
	}
	matchSerialize := regexp.MustCompile(serializeReg).FindStringSubmatch(html)
	if len(matchSerialize) > 0 {
		collecBookInfoRes.Serialize = matchSerialize[1]
	} else {
		global.Collectlog.Errorf("获取小说连载状态失败 bookUrl=%v", bookUrl)
		collecBookInfoRes.Serialize = "连载中"
	}

	updateReg := collect.UpdateReg
	if updateReg == "" {
		err = fmt.Errorf("%v", "获取小说最近更新时间正则不能为空")
		return
	}
	matchUpdateTime := regexp.MustCompile(collect.UpdateReg).FindStringSubmatch(html)
	if len(matchUpdateTime) > 0 {
		collecBookInfoRes.UpdateTime = matchUpdateTime[1]
	} else {
		err = fmt.Errorf("%v", "获取小说最新更新时间失败")
		return
	}
	tagNameReg := collect.TagNameReg
	if tagNameReg == "" {
		err = fmt.Errorf("%v", "获取小说标签正则不能为空")
		return
	}
	matchTagName := regexp.MustCompile(tagNameReg).FindStringSubmatch(html)
	if len(matchTagName) > 0 {
		newStr := strings.Replace(matchTagName[1], "小说", "", -1)
		collecBookInfoRes.TagName = newStr
	} else {
		global.Collectlog.Errorf("获取小说标签失败 bookUrl=%v", bookUrl)
		//err = fmt.Errorf("获取小说标签失败 bookUrl=%v", bookUrl)
		//return
	}
	chapterSectionReg := collect.ChapterSectionReg
	if chapterSectionReg == "" {
		err = fmt.Errorf("%v", "获取小说章节正则不能为空")
		return
	}
	chapterUrlReg := collect.ChapterUrlReg
	if chapterUrlReg == "" {
		err = fmt.Errorf("%v", "获取小说章节链接正则不能为空")
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
		global.Collectlog.Errorf("获取小说章节区间失败 bookUrl=%v", bookUrl)
		err = fmt.Errorf("获取小说章节区间失败 bookUrl=%v", bookUrl)
		return
	}
	collecBookInfoRes.Chapters = chapters
	if len(chapters) <= 0 {
		err = RmmoveCollectPageBooks(collectId, bookUrl)
		if err != nil {
			return
		}
	}
	return
}
