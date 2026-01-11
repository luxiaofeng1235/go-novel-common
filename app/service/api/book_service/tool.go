package book_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/api/comment_service"
	"go-novel/app/service/api/user_service"
	"go-novel/global"
	"go-novel/utils"
	"regexp"
	"sort"
	"strings"
)

func getCommentByBookId(bookId, userId int64) (commentList []*models.CommentListRes) {
	comments, err := comment_service.GetCommentsByBookId(bookId, 1)
	if err != nil {
		return
	}
	if len(comments) <= 0 {
		return
	}
	for _, comment := range comments {
		user, _ := user_service.GetUserById(comment.Uid)
		com := models.CommentListRes{
			Id:          comment.Id,
			UserId:      comment.Uid,
			Nickname:    user.Nickname,
			Pic:         utils.GetFileUrl(user.Pic),
			Text:        comment.Text,
			ReplyNum:    comment.ReplyNum,
			Addtime:     comment.Addtime,
			Score:       comment.Score,
			Ip:          comment.Ip,
			IsPraise:    0,
			PraiseCount: comment.PraiseCount,
		}
		if userId > 0 {
			count := comment_service.GetIsPraiseByUserId(comment.Id, userId)
			if count > 0 {
				com.IsPraise = 1
			}
		}
		commentList = append(commentList, &com)
	}
	return
}

func getUserCountById(id int64) (count int64) {
	err := global.DB.Model(models.McUser{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetChapterContent(content string) (result string) {
	// 替换所有的"[n]"为换行符
	replacedText := strings.Replace(content, "[n]", "\n", -1)
	// 截取字符串，并在截取后替换"\n"为"[n]"
	result = strings.Replace(replacedText, "\n", "[n]", 250)
	return
}

func getNewBookTag(bookType, columnType int) (tags []*models.McTag, err error) {
	db := global.DB.Model(models.McTag{}).Debug()
	db = db.Where("status = 1 and is_new = 1")
	//查询book_type是男生还是说女生，1：男生 2：女生
	if bookType > 0 {
		db = db.Where("book_type = ? ", bookType)
	}
	//现在不查column_type
	if columnType > 0 {
		//db = db.Where("column_type = ?", columnType)
	}
	err = db.Find(&tags).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getCollectById(id int64) (collect *models.McCollect, err error) {
	err = global.DB.Model(models.McCollect{}).Where("id = ?", id).Find(&collect).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getBookUrlByBookId(sourceId, bookId int64) (sourceUrl string) {
	var err error
	err = global.DB.Model(models.McBookSource{}).Select("source_url").Where("sid = ? and bid = ?", sourceId, bookId).Scan(&sourceUrl).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetSourceChapters(collectId int64, bookUrl, sortStatus string) (chapters []*models.McBookChapter, err error) {
	var collect *models.McCollect
	collect, err = getCollectById(collectId)
	if err != nil {
		return
	}
	chapterSectionReg := collect.ChapterSectionReg
	chapterUrlReg := collect.ChapterUrlReg
	var html string
	html, err = utils.GetHtml(bookUrl, "utf-8", 1, 0)
	if html == "" {
		err = fmt.Errorf("获取小说详情页面失败 bookUrl=%v", bookUrl)
		return
	}

	matchChapterSection := regexp.MustCompile(chapterSectionReg).FindStringSubmatch(html)
	if len(matchChapterSection) > 0 {
		ulContent := matchChapterSection[0]
		liMatches := regexp.MustCompile(chapterUrlReg).FindAllStringSubmatch(ulContent, -1)
		var sortInt int
		for _, liMatch := range liMatches {
			sortInt++
			if len(liMatch) > 2 {
				href := liMatch[1]
				chapterName := liMatch[2]
				chapter := &models.McBookChapter{
					ChapterName: chapterName,
					ChapterLink: href,
					Sort:        sortInt,
					Vip:         0,
					Cion:        0,
					TextNum:     2000,
					Addtime:     utils.GetUnix(),
				}
				chapters = append(chapters, chapter)
			}
		}

		if len(chapters) <= 0 {
			return
		}

		var lastIndex int = len(chapters) - 1
		chapters[0].IsFirst = 1
		chapters[lastIndex].IsLast = 1

		if sortStatus == "desc" {
			sort.Slice(chapters, func(i, j int) bool {
				return chapters[i].Sort > chapters[j].Sort
			})
			chapters[0].IsLast = 1
			chapters[lastIndex].IsFirst = 1
		}
	} else {
		err = fmt.Errorf("%v", "获取小说章节区间失败")
		return
	}
	return
}
