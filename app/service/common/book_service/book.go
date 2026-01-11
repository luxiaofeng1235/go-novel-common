package book_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"gorm.io/gorm"
)

func GetBookById(id int64) (book *models.McBook, err error) {
	err = global.DB.Model(models.McBook{}).Where("id = ?", id).Find(&book).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetBookByBookName(bookName, author string) (book *models.McBook, err error) {
	err = global.DB.Model(models.McBook{}).Where("book_name = ? and author = ?", bookName, author).Find(&book).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetBookIdByBookNameAndAuthor(bookName, author string) (bookId int64) {
	var err error
	err = global.DB.Model(models.McBook{}).Select("id").Where("book_name = ? and author = ?", bookName, author).Scan(&bookId).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdateBookInfoByName(bookName, author string, data map[string]interface{}) (err error) {
	err = global.DB.Model(models.McBook{}).Where("book_name = ? and author = ?", bookName, author).Updates(&data).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdateHitsByBookName(bookName, author string) (err error) {
	data := make(map[string]interface{})
	data["hits"] = gorm.Expr("hits + 1")
	data["hits_month"] = gorm.Expr("hits_month + 1")
	data["hits_week"] = gorm.Expr("hits_week + 1")
	data["hits_day"] = gorm.Expr("hits_day + 1")
	err = global.DB.Model(models.McBook{}).Where("book_name = ? and author = ?", bookName, author).Updates(data).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetBookNewChapterId(bookId int64) (newChapterId int64) {
	//章节表
	chapterTable, err := GetChapterTable(bookId)
	if err != nil {
		err = fmt.Errorf("%v", "获取章节失败")
		return
	}
	err = global.DB.Table(chapterTable).Order("sort desc,id desc").Select("id").Limit(1).Scan(&newChapterId).Error
	if err != nil {
		return
	}
	return
}

func GetChapterNameByChapterId(bookId, chapterId int64) (chapterName string) {
	//章节表
	chapterTable, err := GetChapterTable(bookId)
	if err != nil {
		err = fmt.Errorf("%v", "获取章节失败")
		return
	}
	err = global.DB.Table(chapterTable).Select("chapter_name").Where("id = ?", chapterId).Limit(1).Scan(&chapterName).Error
	if err != nil {
		return
	}
	return
}

func GetSearchBookCountByName(bookName string) (count int64) {
	err := global.DB.Model(models.McBook{}).Where("book_name LIKE ?", "%"+bookName+"%").Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetBookCountByName(bookName, author string) (count int64) {
	err := global.DB.Model(models.McBook{}).Where("book_name = ? and author = ?", bookName, author).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdateSearchNumByName(bookName string) (err error) {
	err = global.DB.Model(models.McBook{}).Where("book_name LIKE ?", "%"+bookName+"%").Debug().Update("search_count", gorm.Expr("search_count + ?", 1)).Error
	if err != nil {
		global.Sqllog.Error(err.Error())
		return
	}
	return
}

func GetBookNameById(id int64) (bookName string) {
	var book models.McBook
	var err error
	err = global.DB.Model(models.McBook{}).Select("name,uid").Where("id = ? and yid = 0", id).Find(&book).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	bookName = book.BookName
	return
}

func GetBookCountById(bookId int64) (count int64) {
	var err error
	err = global.DB.Model(models.McBook{}).Where("id = ?", bookId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
