/*
@Descripttion: 持续采集 bqg349 分类、书籍、章节，直接同步 mc_book 和章节文件
@Author: red
@Date: 2026-11-13
@LastEditors: red
@LastEditTime: 2026-11-13
*/
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"go-novel/app/models"
	"go-novel/app/service/collect/bqg349_service"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/db"
	"go-novel/global"
	"go-novel/utils"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	client, err := bqg349_service.NewClient(bqg349_service.DefaultSite, bqg349_service.DefaultAPI)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	db.InitDB() // 已同时初始化 MySQL 和 Redis，书籍和章节由本脚本直接同步。
	db.InitZapLog()
	err = runBqg349Collector(ctx, 5*time.Minute, func(ctx context.Context) error {
		return collectBqg349Site(ctx, client)
	})
	if ctx.Err() != nil {
		log.Printf("收到停止信号，采集已退出: %v", err)
		return
	}
	if err != nil {
		log.Fatal(err)
	}
}

// runBqg349Collector 每轮重新扫描，失败也继续轮询；等待期间可被停止信号唤醒。
func runBqg349Collector(ctx context.Context, interval time.Duration, collect func(context.Context) error) error {
	if interval <= 0 {
		return fmt.Errorf("采集间隔必须大于 0")
	}
	for round := 1; ; round++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		log.Printf("开始第 %d 轮采集 %s", round, bqg349_service.DefaultSite)
		err := collect(ctx)
		if ctx.Err() != nil {
			return errors.Join(err, ctx.Err())
		}
		if err != nil {
			log.Printf("本轮采集未全部完成，下一轮自动重试: %v", err)
		}
		log.Printf("第 %d 轮结束，%s 后继续采集新书和新章节", round, interval)
		timer := time.NewTimer(interval)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
	}
}

func collectBqg349Site(ctx context.Context, client *bqg349_service.Client) error {
	return walkBqg349Site(ctx, client, func(ctx context.Context, book *bqg349_service.Book) error {
		return syncBqg349Book(ctx, client, book)
	})
}

func walkBqg349Site(ctx context.Context, client *bqg349_service.Client, save func(context.Context, *bqg349_service.Book) error) error {
	seen := make(map[bqg349_service.ID]bool)
	var processed, failed int
	for _, category := range client.Categories() {
		if err := ctx.Err(); err != nil {
			return err
		}
		books, err := client.Books(ctx, category.URL)
		if err != nil {
			failed++
			log.Printf("获取分类失败 category=%s: %v", category.Name, err)
			continue
		}
		log.Printf("分类 %s 获取到 %d 本书", category.Name, len(books))
		for _, book := range books {
			if err := ctx.Err(); err != nil {
				return err
			}
			if seen[book.ID] {
				continue
			}
			seen[book.ID] = true
			processed++
			detail, err := client.Book(ctx, book.URL)
			if err == nil {
				err = save(ctx, detail)
			}
			if err != nil {
				if ctx.Err() != nil {
					return err
				}
				failed++
				log.Printf("采集失败 source-id=%d book=%q author=%q: %v", book.ID, book.Title, book.Author, err)
			}
		}
	}
	log.Printf("本轮结束 books=%d failures=%d", processed, failed)
	if failed > 0 {
		return fmt.Errorf("本轮有 %d 项采集失败，已完成进度保留", failed)
	}
	return nil
}

func syncBqg349Book(ctx context.Context, client *bqg349_service.Client, book *bqg349_service.Book) (err error) {
	local, err := saveBqg349Book(ctx, global.DB, book)
	if err != nil {
		return err
	}
	chapters, err := bqg349Chapters(ctx, client, book, false)
	if err != nil {
		return err
	}
	chapterFile, err := chapter_service.GetChapterFile(local.BookName, local.Author)
	if err != nil {
		return err
	}
	chapterPath := func(title string) (string, error) {
		_, path, err := book_service.GetChapterTxtFile(local.BookName, local.Author, title)
		return path, err
	}
	defer func() {
		// 采集失败或手动停止时，也把已经保存的章节进度同步回 mc_book。
		progressCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = errors.Join(err, saveBqg349BookProgress(progressCtx, global.DB, local, chapterFile))
	}()
	log.Printf("开始采集 local-id=%d source-id=%d book=%q chapters=%d", local.Id, book.ID, local.BookName, len(chapters))
	return saveBqg349Chapters(ctx, chapters, chapterFile, chapterPath, client.ChapterText, false)
}

// saveBqg349Book 按书名和作者查重；空库直接新增，后续更新同一条书籍记录。
func saveBqg349Book(ctx context.Context, database *gorm.DB, book *bqg349_service.Book) (*models.McBook, error) {
	if book == nil || book.ID <= 0 || book.Title == "" || book.Author == "" || book.URL == "" {
		return nil, fmt.Errorf("书籍信息不完整，无法入库")
	}
	var local models.McBook
	err := database.WithContext(ctx).Where("book_name = ? AND author = ?", book.Title, book.Author).Take(&local).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	serialize := 1 // 沿用原采集程序：1 连载，2 完结。
	if strings.Contains(book.Status, "完") {
		serialize = 2
	}
	now := time.Now().Unix()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		local = models.McBook{BookName: book.Title, Author: book.Author, BookType: 1, Addtime: now}
		if strings.Contains(book.Category, "女生") {
			local.BookType = 2
		}
		// 沿用本地分类配置，不写死分类 ID。
		var class models.McBookClass
		if classErr := database.WithContext(ctx).Where("class_name IN ?", []string{book.Category, book.Category + "小说"}).Take(&class).Error; classErr == nil {
			local.Cid, local.ClassName = class.Id, class.ClassName
			if class.BookType > 0 {
				local.BookType = class.BookType
			}
		}
		if err := database.WithContext(ctx).Create(&local).Error; err != nil {
			return nil, fmt.Errorf("新增 mc_book 失败: %w", err)
		}
		log.Printf("新书入库 local-id=%d source-id=%d book=%q author=%q", local.Id, book.ID, book.Title, book.Author)
	}
	data := map[string]any{
		"source_url": book.URL, "category_name": book.Category, "desc": book.Intro,
		"serialize": serialize, "last_chapter_title": book.LastChapter,
		"last_chapter_time": utils.DateToUnix(book.LastUpdate), "uptime": now,
	}
	if local.Pic == "" {
		data["pic"] = book.CoverURL
	}
	if err := database.WithContext(ctx).Model(&models.McBook{}).Where("id = ?", local.Id).Updates(data).Error; err != nil {
		return nil, fmt.Errorf("更新 mc_book 失败: %w", err)
	}
	return &local, nil
}

func saveBqg349BookProgress(ctx context.Context, database *gorm.DB, book *models.McBook, chapterFile string) error {
	chapters, err := chapter_service.GetChaptersByFile(chapterFile)
	if err != nil {
		return err
	}
	var last *models.McBookChapter
	var total int
	for _, chapter := range chapters {
		if chapter == nil {
			return fmt.Errorf("本地目录存在空章节，停止更新书籍进度")
		}
		total += chapter.TextNum
		if last == nil || chapter.Sort > last.Sort {
			last = chapter
		}
	}
	if last == nil {
		return nil
	}
	data := map[string]any{
		"chapter_num": len(chapters), "text_num": total,
		"update_chapter_id": last.Id, "update_chapter_title": last.ChapterName,
		"update_chapter_time": time.Now().Unix(), "uptime": time.Now().Unix(),
	}
	if book.ChapterNum == 0 {
		data["status"] = 1
	}
	return database.WithContext(ctx).Model(&models.McBook{}).Where("id = ?", book.Id).Updates(data).Error
}

func saveBqg349Chapters(ctx context.Context, chapters []bqg349_service.Chapter, chapterFile string, chapterPath func(string) (string, error), fetch func(context.Context, bqg349_service.Chapter) (string, error), overwrite bool) (err error) {
	// TXT 以章节标题的 MD5 命名，重复标题会覆盖正文，必须在写入前拒绝。
	names := make(map[string]bool, len(chapters))
	for _, chapter := range chapters {
		if strings.TrimSpace(chapter.Title) == "" || chapter.URL == "" {
			return fmt.Errorf("新站目录包含空标题或空链接")
		}
		if names[chapter.Title] {
			return fmt.Errorf("新站目录存在重复标题 %q，无法使用当前 TXT 存储格式", chapter.Title)
		}
		names[chapter.Title] = true
	}
	records, err := chapter_service.GetChaptersByFile(chapterFile)
	if err != nil {
		return err
	}
	sort.SliceStable(records, func(i, j int) bool {
		if records[i] == nil || records[j] == nil {
			return records[i] == nil && records[j] != nil
		}
		return records[i].Sort < records[j].Sort
	})
	index := make(map[string]*models.McBookChapter, len(records))
	var lastID int64
	for _, record := range records {
		if record == nil || record.ChapterName == "" || index[record.ChapterName] != nil {
			return fmt.Errorf("本地目录包含空记录或重复标题，请先检查 %s", chapterFile)
		}
		index[record.ChapterName] = record
		if record.Id > lastID {
			lastID = record.Id
		}
	}
	order, err := bqg349ChapterOrder(records, chapters)
	if err != nil {
		return err
	}
	// 一次读取目录、用 map 查重，每 25 章保存；退出时也保存已完成章节。
	var pending, downloaded, reused int
	flush := func() error {
		if pending == 0 {
			return nil
		}
		ordered := make([]*models.McBookChapter, 0, len(index))
		for _, title := range order {
			if record := index[title]; record != nil {
				record.Sort = len(ordered) + 1
				ordered = append(ordered, record)
			}
		}
		data, err := json.MarshalIndent(ordered, "", "  ")
		if err != nil {
			return err
		}
		if err := bqg349WriteFile(chapterFile, data); err != nil {
			return err
		}
		pending = 0
		return nil
	}
	defer func() { err = errors.Join(err, flush()) }()
	for _, chapter := range chapters {
		if err := ctx.Err(); err != nil {
			return err
		}
		txtFile, err := chapterPath(chapter.Title)
		if err != nil {
			return err
		}
		record := index[chapter.Title]
		text, readErr := os.ReadFile(txtFile)
		if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
			return fmt.Errorf("读取正文 %s 失败: %w", txtFile, readErr)
		}
		textNum := utf8.RuneCount(text)
		incomplete := record != nil && (record.IsLess != 0 || record.TextNum > textNum)
		if overwrite || incomplete || !utf8.Valid(text) || strings.TrimSpace(string(text)) == "" {
			body, err := fetch(ctx, chapter)
			if err != nil {
				// 停止本书，避免失败后继续追加造成目录顺序错乱；下次从文件续采。
				return fmt.Errorf("获取 %q 失败: %w", chapter.Title, err)
			}
			if strings.TrimSpace(body) == "" || !utf8.ValidString(body) {
				return fmt.Errorf("%q 返回空正文或无效 UTF-8", chapter.Title)
			}
			text = []byte(body)
			if err := bqg349WriteFile(txtFile, text); err != nil {
				return fmt.Errorf("保存 %q 失败: %w", chapter.Title, err)
			}
			downloaded++
		} else {
			reused++
		}
		textNum = utf8.RuneCount(text)
		if record == nil {
			lastID++
			record = &models.McBookChapter{Id: lastID, ChapterName: chapter.Title, Addtime: time.Now().Unix()}
			index[chapter.Title] = record
		} else if record.ChapterLink == chapter.URL && record.TextNum == textNum && record.IsLess == 0 {
			continue
		}
		record.ChapterLink = chapter.URL
		record.TextNum = textNum
		record.IsLess = 0
		pending++
		if pending >= 25 {
			if err := flush(); err != nil {
				return err
			}
			log.Printf("已保存 file=%s downloaded=%d reused=%d", chapterFile, downloaded, reused)
		}
	}
	if err := flush(); err != nil {
		return err
	}
	log.Printf("补采完成 file=%s downloaded=%d reused=%d", chapterFile, downloaded, reused)
	return nil
}

// bqg349ChapterOrder 将缺章插入新站的对应位置，保留旧章节 ID 和旧站独有章节。
func bqg349ChapterOrder(records []*models.McBookChapter, chapters []bqg349_service.Chapter) ([]string, error) {
	positions := make(map[string]int, len(records))
	for i, record := range records {
		positions[record.ChapterName] = i
	}
	order := make([]string, 0, len(records)+len(chapters))
	next := 0
	shared := 0
	for _, chapter := range chapters {
		position, exists := positions[chapter.Title]
		if !exists {
			order = append(order, chapter.Title)
			continue
		}
		if position < next {
			return nil, fmt.Errorf("本地与新站共有章节顺序不一致: %q，停止补采以免串章", chapter.Title)
		}
		shared++
		for next <= position {
			order = append(order, records[next].ChapterName)
			next++
		}
	}
	if len(records) > 0 && shared == 0 {
		return nil, fmt.Errorf("本地与新站目录没有共同标题，无法确定补采顺序")
	}
	for ; next < len(records); next++ {
		order = append(order, records[next].ChapterName)
	}
	return order, nil
}

type bqg349ChapterCache struct {
	DirectoryID   bqg349_service.ID
	LastChapterID bqg349_service.ID
	LastChapter   string
	LastUpdate    string
	Chapters      []bqg349_service.Chapter
}

func bqg349Chapters(ctx context.Context, client *bqg349_service.Client, book *bqg349_service.Book, refresh bool) ([]bqg349_service.Chapter, error) {
	key := fmt.Sprintf("%s:chapters:v1:%d", client.CachePrefix(), book.ID)
	if !refresh {
		data, err := global.Redis.Get(ctx, key).Bytes()
		if err == nil {
			var cached bqg349ChapterCache
			if json.Unmarshal(data, &cached) == nil && cached.DirectoryID == book.DirectoryID && cached.LastChapterID == book.LastChapterID && cached.LastChapter == book.LastChapter && cached.LastUpdate == book.LastUpdate && len(cached.Chapters) > 0 {
				return cached.Chapters, nil
			}
		} else if !errors.Is(err, redis.Nil) {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			log.Printf("目录缓存读取失败，将直接获取: %v", err)
		}
	}
	chapters, err := client.Chapters(ctx, book)
	if err != nil {
		return nil, err
	}
	cached := bqg349ChapterCache{DirectoryID: book.DirectoryID, LastChapterID: book.LastChapterID, LastChapter: book.LastChapter, LastUpdate: book.LastUpdate, Chapters: chapters}
	data, err := json.Marshal(cached)
	if err != nil {
		return nil, err
	}
	if err := global.Redis.Set(ctx, key, data, 30*time.Minute).Err(); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		log.Printf("目录缓存写入失败: %v", err)
	}
	return chapters, nil
}

// bqg349WriteFile 同目录临时文件写完再替换，避免中断时截断原目录或正文。
func bqg349WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.CreateTemp(dir, ".bqg349-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	mode := os.FileMode(0644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	} else if !errors.Is(err, os.ErrNotExist) {
		file.Close()
		return err
	}
	if err := file.Chmod(mode); err != nil {
		file.Close()
		return err
	}
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}
