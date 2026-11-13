/*
@Descripttion: 使用 bqg349 动态接口补采已有小说，支持旧书源精确匹配和断点续采
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
	"flag"
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

	"github.com/go-redis/redis/v8"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	bookID := flag.Int64("book-id", 3196, "本地小说 ID；0 表示补采所有 biquge34/bqg349 来源的小说")
	bookURL := flag.String("book-url", "", "指定新站书籍链接，必须与本地书名、作者一致，仅用于单本补采")
	refresh := flag.Bool("refresh", false, "忽略目录缓存，重新获取目录")
	overwrite := flag.Bool("overwrite", false, "重新下载已有正文，用于修复旧采集内容")
	flag.Parse()
	if *bookID < 0 || (*bookID == 0 && *bookURL != "") {
		log.Fatal("book-id 必须大于等于 0；book-url 必须同时指定单本 book-id")
	}
	client, err := bqg349_service.NewClient(bqg349_service.DefaultSite, bqg349_service.DefaultAPI)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	db.InitDB() // 已同时初始化 MySQL 和 Redis；补采不需要启动 NSQ 消费者。
	db.InitZapLog()
	if err := collectBqg349Books(ctx, client, *bookID, *bookURL, *refresh, *overwrite); err != nil {
		log.Fatal(err)
	}
}

func collectBqg349Books(ctx context.Context, client *bqg349_service.Client, bookID int64, bookURL string, refresh, overwrite bool) error {
	var lastID int64
	var processed, failed int
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		query := global.DB.WithContext(ctx).Model(&models.McBook{}).Order("id ASC").Limit(100)
		if bookID > 0 {
			query = query.Where("id = ?", bookID)
		} else {
			query = query.Where("id > ? AND (source_url LIKE ? OR source_url LIKE ?)", lastID, "%biquge34%", "%bqg349%")
		}
		var books []*models.McBook
		if err := query.Find(&books).Error; err != nil {
			return fmt.Errorf("查询待补采小说失败: %w", err)
		}
		if len(books) == 0 {
			break
		}
		for _, book := range books {
			processed++
			lastID = book.Id
			if err := collectBqg349Book(ctx, client, book, bookURL, refresh, overwrite); err != nil {
				if ctx.Err() != nil {
					return err
				}
				failed++
				log.Printf("补采失败 id=%d book=%q author=%q: %v", book.Id, book.BookName, book.Author, err)
			}
		}
		if bookID > 0 {
			break
		}
	}
	if processed == 0 {
		return fmt.Errorf("未找到待补采小说 book-id=%d", bookID)
	}
	log.Printf("补采结束 total=%d success=%d failed=%d", processed, processed-failed, failed)
	if failed > 0 {
		return fmt.Errorf("%d 本小说补采失败，可重新运行续采", failed)
	}
	return nil
}

func collectBqg349Book(ctx context.Context, client *bqg349_service.Client, local *models.McBook, explicitURL string, refresh, overwrite bool) (err error) {
	var book *bqg349_service.Book
	if explicitURL != "" {
		book, err = client.Book(ctx, explicitURL)
	} else if strings.HasPrefix(local.SourceUrl, bqg349_service.DefaultSite+"/#/book/") {
		book, err = client.Book(ctx, local.SourceUrl)
	} else {
		// 老站编号与新站不同，禁止只替换域名拼接 URL。
		book, err = client.FindBook(ctx, local.BookName, local.Author)
	}
	if err != nil {
		return err
	}
	if book.Title != strings.TrimSpace(local.BookName) || book.Author != strings.TrimSpace(local.Author) {
		return fmt.Errorf("新站书名或作者不匹配: %q / %q", book.Title, book.Author)
	}
	chapters, err := bqg349Chapters(ctx, client, book, refresh)
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
	log.Printf("开始补采 id=%d book=%q url=%s chapters=%d", local.Id, local.BookName, book.URL, len(chapters))
	return saveBqg349Chapters(ctx, chapters, chapterFile, chapterPath, client.ChapterText, overwrite)
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
