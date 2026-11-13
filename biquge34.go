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
	"sync"
	"sync/atomic"
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

const (
	bqg349BookWorkers    = 5
	bqg349ChapterWorkers = 20
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
		// 多本书共享 20 个正文下载额度。
		fetch := limitBqg349ChapterFetch(client.ChapterText, bqg349ChapterWorkers)
		log.Printf("批量采集并发配置 book-workers=%d chapter-workers-total=%d", bqg349BookWorkers, bqg349ChapterWorkers)
		return walkBqg349Site(ctx, client, func(ctx context.Context, book *bqg349_service.Book) error {
			return syncBqg349Book(ctx, client, book, fetch)
		})
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

type bqg349SiteClient interface {
	Categories() []bqg349_service.Category
	Books(context.Context, string) ([]bqg349_service.Book, error)
	Book(context.Context, string) (*bqg349_service.Book, error)
}

func walkBqg349Site(ctx context.Context, client bqg349SiteClient, save func(context.Context, *bqg349_service.Book) error) error {
	seen := make(map[bqg349_service.ID]bool)
	var processed int
	var failed atomic.Int64
	jobs := make(chan bqg349_service.Book, bqg349BookWorkers)
	var wg sync.WaitGroup
	var identityMu sync.Mutex
	// 不同源 ID 可能对应同一本书，必须在入库和写文件前按书名、作者再次去重。
	identities := make(map[[2]string]bool)
	for worker := 0; worker < bqg349BookWorkers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for book := range jobs {
				if ctx.Err() != nil {
					return
				}
				detail, err := client.Book(ctx, book.URL)
				if err == nil {
					identity := [2]string{detail.Title, detail.Author}
					identityMu.Lock()
					duplicate := identities[identity]
					identities[identity] = true
					identityMu.Unlock()
					if duplicate {
						log.Printf("跳过重复书籍 source-id=%d book=%q author=%q", book.ID, detail.Title, detail.Author)
						continue
					}
					err = save(ctx, detail)
				}
				if err != nil {
					if ctx.Err() != nil {
						return
					}
					failed.Add(1)
					log.Printf("采集失败 source-id=%d book=%q author=%q: %v", book.ID, book.Title, book.Author, err)
				}
			}
		}()
	}
categories:
	for _, category := range client.Categories() {
		if ctx.Err() != nil {
			break
		}
		books, err := client.Books(ctx, category.URL)
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			failed.Add(1)
			log.Printf("获取分类失败 category=%s: %v", category.Name, err)
			continue
		}
		log.Printf("分类 %s 获取到 %d 本书", category.Name, len(books))
		for _, book := range books {
			if seen[book.ID] {
				continue
			}
			seen[book.ID] = true
			processed++
			select {
			case jobs <- book:
			case <-ctx.Done():
				break categories
			}
		}
	}
	close(jobs)
	// 等待所有在采书籍保存各自进度后再结束本轮或响应停止信号。
	wg.Wait()
	if err := ctx.Err(); err != nil {
		return err
	}
	log.Printf("本轮结束 books=%d failures=%d", processed, failed.Load())
	if failures := failed.Load(); failures > 0 {
		return fmt.Errorf("本轮有 %d 项采集失败，已完成进度保留", failures)
	}
	return nil
}

func limitBqg349ChapterFetch(fetch func(context.Context, bqg349_service.Chapter) (string, error), concurrency int) func(context.Context, bqg349_service.Chapter) (string, error) {
	slots := make(chan struct{}, concurrency)
	return func(ctx context.Context, chapter bqg349_service.Chapter) (string, error) {
		select {
		case slots <- struct{}{}:
			defer func() { <-slots }()
		case <-ctx.Done():
			return "", ctx.Err()
		}
		if err := ctx.Err(); err != nil {
			return "", err
		}
		return fetch(ctx, chapter)
	}
}

func syncBqg349Book(ctx context.Context, client *bqg349_service.Client, book *bqg349_service.Book, fetch func(context.Context, bqg349_service.Chapter) (string, error)) (err error) {
	local, err := saveBqg349Book(ctx, global.DB, book)
	if err != nil {
		return err
	}
	chapters, err := bqg349Chapters(ctx, client, book)
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
		// 失败或手动停止时，也把已保存的目录和正文统计同步回 mc_book。
		progressCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = errors.Join(err, saveBqg349BookProgress(progressCtx, global.DB, local, chapterFile))
		if err == nil {
			log.Printf("整本同步完成 local-id=%d book=%q chapters=%d", local.Id, local.BookName, len(chapters))
		}
	}()
	complete, err := bqg349BookComplete(ctx, book.Status, chapters, chapterFile, chapterPath)
	if err != nil {
		return err
	}
	if complete {
		log.Printf("已完结且采集完整，跳过采集 local-id=%d book=%q local=%d remote=%d", local.Id, local.BookName, len(chapters), len(chapters))
		return nil
	}
	log.Printf("开始处理书籍 local-id=%d source-id=%d book=%q cid=%d class=%q chapters=%d workers=%d", local.Id, book.ID, local.BookName, local.Cid, local.ClassName, len(chapters), bqg349ChapterWorkers)
	return saveBqg349Chapters(ctx, local.BookName, chapters, chapterFile, chapterPath, fetch)
}

func bqg349IsFinished(status string) bool {
	switch strings.TrimSpace(status) {
	case "完结", "已完结", "完本", "已完本", "完成", "已完成":
		return true
	default:
		return false
	}
}

// bqg349BookComplete 复用落盘的章节完成标记；只读目录和文件属性，不读取正文。
func bqg349BookComplete(ctx context.Context, status string, chapters []bqg349_service.Chapter, chapterFile string, chapterPath func(string) (string, error)) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if !bqg349IsFinished(status) || len(chapters) == 0 {
		return false, nil
	}
	records, err := chapter_service.GetChaptersByFile(chapterFile)
	if err != nil {
		return false, err
	}
	if len(records) != len(chapters) {
		return false, nil
	}
	seen := make(map[string]bool, len(chapters))
	for i, chapter := range chapters {
		if err := ctx.Err(); err != nil {
			return false, err
		}
		record := records[i]
		if chapter.Title == "" || chapter.URL == "" || seen[chapter.Title] || record == nil ||
			record.ChapterName != chapter.Title || record.ChapterLink != chapter.URL ||
			record.Sort != i+1 || record.IsLess != 0 || record.TextNum <= 0 {
			return false, nil
		}
		seen[chapter.Title] = true
		path, err := chapterPath(chapter.Title)
		if err != nil {
			return false, err
		}
		info, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if !info.Mode().IsRegular() || info.Size() < int64(record.TextNum) {
			return false, nil
		}
	}
	return true, nil
}

// bqg349ClassCandidates 兼容书源分类简称和本地展示名称，不写死数据库 ID。
func bqg349ClassCandidates(category string) []string {
	name := strings.TrimSuffix(strings.TrimSpace(category), "小说")
	if strings.HasPrefix(name, "女生") && name != "女生" {
		name = strings.TrimPrefix(name, "女生")
	}
	aliases := map[string][]string{
		"玄幻": {"玄幻奇幻", "玄幻言情"}, "奇幻": {"玄幻奇幻", "玄幻言情"},
		"武侠": {"传统武侠", "武侠仙侠"}, "仙侠": {"古典仙侠", "武侠仙侠", "仙侠奇缘"},
		"修真": {"古典仙侠", "仙侠", "武侠仙侠", "仙侠奇缘"},
		"都市": {"都市言情", "现代言情"}, "都市言情": {"都市", "现代言情"},
		"历史": {"历史题材", "历史军事"},
		"军事": {"军旅战争", "历史军事"}, "现实": {"现实题材"},
		"网游": {"游戏异界", "游戏", "游戏竞技"}, "游戏": {"网游", "游戏异界", "游戏竞技"},
		"科幻": {"未来科幻", "科幻空间"}, "悬疑": {"悬疑灵异", "悬疑推理"},
		"灵异": {"悬疑灵异", "悬疑", "悬疑推理"}, "体育": {"体育竞技"},
		"女生": {"现代言情", "古代言情"}, "女频": {"女生", "现代言情"},
		"言情": {"现代言情", "古代言情"},
	}
	names := []string{strings.TrimSpace(category), name}
	names = append(names, aliases[name]...)
	seen := make(map[string]bool)
	result := make([]string, 0, len(names))
	for _, candidate := range names {
		if candidate != "" && !seen[candidate] {
			result = append(result, candidate)
			seen[candidate] = true
		}
	}
	return result
}

func matchBqg349BookClass(ctx context.Context, database *gorm.DB, category string, local *models.McBook) (*models.McBookClass, error) {
	// 保留有效的本地分类，包括后台手工指定的分类；同步其最新名称。
	if local.Cid > 0 {
		var class models.McBookClass
		err := database.WithContext(ctx).Where("id = ?", local.Cid).Take(&class).Error
		if err == nil {
			return &class, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	bookType := local.BookType
	if bookType != 2 {
		bookType = 1
	}
	if strings.Contains(category, "女生") || strings.Contains(category, "女频") {
		bookType = 2
	}
	names := bqg349ClassCandidates(category)
	var matches []models.McBookClass
	if err := database.WithContext(ctx).Where("book_type IN ? AND status = 1 AND class_name IN ?", []int{1, 2}, names).Order("sort ASC, id ASC").Find(&matches).Error; err != nil {
		return nil, err
	}
	for _, name := range names {
		for i := range matches {
			if matches[i].BookType == bookType && matches[i].ClassName == name {
				return &matches[i], nil
			}
		}
	}
	// “古代言情”等明确的女频分类可由实际分类记录识别，避免默认归到男生。
	if bookType == 1 {
		for _, name := range names {
			for i := range matches {
				if matches[i].ClassName == name {
					return &matches[i], nil
				}
			}
		}
	}
	var fallback models.McBookClass
	if err := database.WithContext(ctx).Where("book_type = ? AND class_name = ? AND status = 1", bookType, "其他").Order("id ASC").Take(&fallback).Error; err != nil {
		return nil, fmt.Errorf("来源分类 %q 无可用本地分类（book_type=%d）: %w", category, bookType, err)
	}
	return &fallback, nil
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
	if bqg349IsFinished(book.Status) {
		serialize = 2
	}
	now := time.Now().Unix()
	isNew := errors.Is(err, gorm.ErrRecordNotFound)
	if isNew {
		local = models.McBook{BookName: book.Title, Author: book.Author, BookType: 1, ClassName: book.Category, Status: 1, Addtime: now, Uptime: now}
	}
	class, err := matchBqg349BookClass(ctx, database, book.Category, &local)
	if err != nil {
		return nil, err
	}
	local.Cid, local.ClassName, local.BookType = class.Id, class.ClassName, class.BookType
	if isNew {
		if err := database.WithContext(ctx).Create(&local).Error; err != nil {
			return nil, fmt.Errorf("新增 mc_book 失败: %w", err)
		}
		log.Printf("新书入库 local-id=%d source-id=%d book=%q author=%q", local.Id, book.ID, book.Title, book.Author)
	}
	data := map[string]any{
		"cid": local.Cid, "class_name": local.ClassName, "book_type": local.BookType,
		"source_url": book.URL, "category_name": book.Category, "desc": book.Intro,
		"serialize": serialize, "last_chapter_title": book.LastChapter,
		"last_chapter_time": utils.DateToUnix(book.LastUpdate), "uptime": now,
	}
	// 沿用原采集器的随机热度范围，只补空值/零值，后续轮询不覆盖已有统计。
	hits, day, week, month, shits, score, readCount, searchCount := utils.GetRandNumBookHits()
	for _, counter := range []struct {
		column string
		old    int64
		value  int64
	}{
		{"hits", local.Hits, hits}, {"hits_day", local.HitsDay, day},
		{"hits_week", local.HitsWeek, week}, {"hits_month", local.HitsMonth, month},
		{"shits", local.Shits, shits}, {"read_count", local.ReadCount, readCount},
		{"search_count", local.SearchCount, searchCount},
	} {
		if counter.old <= 0 {
			data[counter.column] = counter.value
		}
	}
	if local.SearchNum <= 0 {
		if local.SearchCount > 0 {
			searchCount = local.SearchCount
		}
		data["search_num"] = searchCount + int64(utils.RangeNum(0, 3000))
	}
	if local.Score <= 0 {
		data["score"] = score
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
		if chapter.IsLess != 0 || chapter.TextNum <= 0 {
			continue
		}
		total += chapter.TextNum
		if last == nil || chapter.Sort > last.Sort {
			last = chapter
		}
	}
	if len(chapters) == 0 {
		return nil
	}
	data := map[string]any{
		"chapter_num": len(chapters), "text_num": total,
		"uptime": time.Now().Unix(),
	}
	if last != nil {
		data["update_chapter_id"] = last.Id
		data["update_chapter_title"] = last.ChapterName
		data["update_chapter_time"] = time.Now().Unix()
	}
	if book.ChapterNum == 0 {
		data["status"] = 1
	}
	return database.WithContext(ctx).Model(&models.McBook{}).Where("id = ?", book.Id).Updates(data).Error
}

// saveBqg349Chapters 先保存完整目录，再按目录 URL 并发获取正文。
func saveBqg349Chapters(ctx context.Context, bookName string, chapters []bqg349_service.Chapter, chapterFile string, chapterPath func(string) (string, error), fetch func(context.Context, bqg349_service.Chapter) (string, error)) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}
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
	// 全量目录先落盘，尚未取得正文的新章节用 is_less=1 标记。
	var pending, added, downloaded, reused int
	for _, chapter := range chapters {
		record := index[chapter.Title]
		if record == nil {
			lastID++
			record = &models.McBookChapter{Id: lastID, ChapterName: chapter.Title, IsLess: 1, Addtime: time.Now().Unix()}
			index[chapter.Title] = record
			added++
			pending++
		}
		if record.ChapterLink != chapter.URL {
			record.ChapterLink = chapter.URL
			pending++
		}
	}
	ordered := make([]*models.McBookChapter, 0, len(index))
	for _, title := range order {
		record := index[title]
		if record.Sort != len(ordered)+1 {
			record.Sort = len(ordered) + 1
			pending++
		}
		ordered = append(ordered, record)
	}
	flush := func() error {
		if pending == 0 {
			return nil
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
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := flush(); err != nil {
		return err
	}
	log.Printf("目录保存完成 book=%q total=%d added=%d file=%s", bookName, len(ordered), added, chapterFile)
	defer func() { err = errors.Join(err, flush()) }()
	for start := 0; start < len(chapters); start += 25 {
		if err := ctx.Err(); err != nil {
			return err
		}
		end := start + 25
		if end > len(chapters) {
			end = len(chapters)
		}
		batch := chapters[start:end]
		previousDownloaded := downloaded
		results := downloadBqg349Batch(ctx, batch, index, chapterPath, fetch)
		var batchErr error
		for i, chapter := range batch {
			result := results[i]
			if result.err != nil {
				batchErr = errors.Join(batchErr, result.err)
				continue
			}
			if result.downloaded {
				downloaded++
			} else {
				reused++
			}
			record := index[chapter.Title]
			if record.TextNum != result.textNum || record.IsLess != 0 {
				record.TextNum = result.textNum
				record.IsLess = 0
				pending++
			}
		}
		if err := flush(); err != nil {
			return errors.Join(batchErr, err)
		}
		if downloaded > previousDownloaded || batchErr != nil {
			log.Printf("正文下载进度 book=%q processed=%d/%d downloaded=%d reused=%d", bookName, end, len(chapters), downloaded, reused)
		}
		if batchErr != nil {
			return batchErr
		}
	}
	log.Printf("正文采集完成 book=%q total=%d downloaded=%d reused=%d", bookName, len(chapters), downloaded, reused)
	return nil
}

type bqg349DownloadResult struct {
	textNum    int
	downloaded bool
	err        error
}

// downloadBqg349Batch 最多使用 20 个协程；批内目录只读，JSON 由主协程保存。
func downloadBqg349Batch(ctx context.Context, chapters []bqg349_service.Chapter, index map[string]*models.McBookChapter, chapterPath func(string) (string, error), fetch func(context.Context, bqg349_service.Chapter) (string, error)) []bqg349DownloadResult {
	results := make([]bqg349DownloadResult, len(chapters))
	jobs := make(chan int)
	workers := bqg349ChapterWorkers
	if len(chapters) < workers {
		workers = len(chapters)
	}
	var wg sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				chapter := chapters[i]
				results[i] = downloadBqg349Chapter(ctx, chapter, index[chapter.Title], chapterPath, fetch)
			}
		}()
	}
	for i := range chapters {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return results
}

func downloadBqg349Chapter(ctx context.Context, chapter bqg349_service.Chapter, record *models.McBookChapter, chapterPath func(string) (string, error), fetch func(context.Context, bqg349_service.Chapter) (string, error)) bqg349DownloadResult {
	if err := ctx.Err(); err != nil {
		return bqg349DownloadResult{err: err}
	}
	txtFile, err := chapterPath(chapter.Title)
	if err != nil {
		return bqg349DownloadResult{err: err}
	}
	text, err := os.ReadFile(txtFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return bqg349DownloadResult{err: fmt.Errorf("读取正文 %s 失败: %w", txtFile, err)}
	}
	textNum := utf8.RuneCount(text)
	// 新目录里的待下载标记不应触发重复下载已完整落盘、尚未更新 JSON 的 TXT。
	incomplete := record != nil && (record.TextNum > textNum || (record.IsLess != 0 && record.TextNum > 0))
	if !incomplete && utf8.Valid(text) && strings.TrimSpace(string(text)) != "" {
		return bqg349DownloadResult{textNum: textNum}
	}
	body, err := fetch(ctx, chapter)
	if err != nil {
		return bqg349DownloadResult{err: fmt.Errorf("获取 %q 失败: %w", chapter.Title, err)}
	}
	if strings.TrimSpace(body) == "" || !utf8.ValidString(body) {
		return bqg349DownloadResult{err: fmt.Errorf("%q 返回空正文或无效 UTF-8", chapter.Title)}
	}
	if err := bqg349WriteFile(txtFile, []byte(body)); err != nil {
		return bqg349DownloadResult{err: fmt.Errorf("保存 %q 失败: %w", chapter.Title, err)}
	}
	return bqg349DownloadResult{textNum: utf8.RuneCountInString(body), downloaded: true}
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

func bqg349Chapters(ctx context.Context, client *bqg349_service.Client, book *bqg349_service.Book) ([]bqg349_service.Chapter, error) {
	key := fmt.Sprintf("%s:chapters:v1:%d", client.CachePrefix(), book.ID)
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
	chapters, err := client.Chapters(ctx, book)
	if err != nil {
		return nil, err
	}
	cached := bqg349ChapterCache{DirectoryID: book.DirectoryID, LastChapterID: book.LastChapterID, LastChapter: book.LastChapter, LastUpdate: book.LastUpdate, Chapters: chapters}
	data, err = json.Marshal(cached)
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
