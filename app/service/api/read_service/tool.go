package read_service

import (
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/utils"
)

func getReadBook(bookId, userId int64) (read *models.BookReadListRes, err error) {
	book, err := book_service.GetBookById(bookId)
	if err != nil {
		return
	}
	newsChapterId, _ := chapter_service.GetBookNewChapterId(book.BookName, book.Author)
	var readChapterId int64
	var readChapterName string
	readChapterId, readChapterName = book_service.GetReadChapterIdByBookId(bookId, userId)
	if readChapterId > 0 && readChapterName == "" {
		readChapterName = chapter_service.GetChapterNameByChapterId(book.BookName, book.Author, readChapterId)
	}
	if readChapterName == "" {
		readChapterName = "未阅读"
	}
	var isShelf int64
	var count int64
	count = book_service.GetShelfCountByBookId(bookId, userId)
	if count > 0 {
		isShelf = 1
	}

	read = &models.BookReadListRes{
		Author:          book.Author,
		Bid:             book.Id,
		BookName:        book.BookName,
		NewsChapterId:   newsChapterId,
		ChapterNum:      book.ChapterNum,
		Serialize:       book.Serialize,
		ReadChapterName: readChapterName,
		ReadChapterId:   readChapterId,
		Pic:             utils.GetFileUrl(book.Pic),
		IsShelf:         isShelf,
	}
	return
}
