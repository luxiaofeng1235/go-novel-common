package book_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetChapterCountByBookId(tableName string) (count int64) {
	var err error
	err = global.DB.Table(tableName).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetChapterIdByName(tableName, chapterName string) (chapterId int64) {
	var err error
	err = global.DB.Table(tableName).Select("id").Debug().Where("chapter_name = ?", chapterName).Scan(&chapterId).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetChapterByChapterId(tableName string, chapterId int64) (chapter *models.McBookChapter, err error) {
	err = global.DB.Table(tableName).Where("id = ?", chapterId).Find(&chapter).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetFirstChapter(tableName string) (chapter *models.McBookChapter, err error) {
	err = global.DB.Table(tableName).Order("sort ASC").Find(&chapter).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getChapterPrev(tableName string, bookId int64, sort int) (chapterId int64) {
	var err error
	db := global.DB.Table(tableName).Order("sort DESC").Select("id")
	err = db.Where("sort < ?", sort).Limit(1).Find(&chapterId).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getChapterNext(tableName string, bookId int64, sort int) (chapterId int64) {
	var err error
	db := global.DB.Table(tableName).Order("sort ASC").Select("id")
	err = db.Where("sort > ?", sort).Limit(1).Find(&chapterId).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
