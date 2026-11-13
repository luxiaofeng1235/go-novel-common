/*
@Descripttion: bqg349 补采目录排序、断点续采和原子写入回归测试
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
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-novel/app/models"
	"go-novel/app/service/collect/bqg349_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/db"

	drivermysql "github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestBqg349DiscoversSourceBooksWithoutLocalRows(t *testing.T) {
	for _, failFirst := range []bool{false, true} {
		t.Run(fmt.Sprintf("first_fails_%t", failFirst), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/sort":
					if r.URL.Query().Get("sort") == "xuanhuan" || r.URL.Query().Get("sort") == "finish" {
						fmt.Fprint(w, `{"data":[{"id":2530,"title":"书一","author":"作者"},{"id":2531,"title":"书二","author":"作者"}]}`)
					} else {
						fmt.Fprint(w, `{"data":[]}`)
					}
				case "/api/book":
					// 目录 token 与书籍 token 相同，借助已验证的样例区分编号。
					if r.URL.Query().Get("token") == "IZzVX5XJA1oySI5jozY8vQ==" {
						fmt.Fprint(w, `{"id":2530,"dirid":2530,"title":"书一","author":"作者"}`)
					} else {
						fmt.Fprint(w, `{"id":2531,"dirid":2531,"title":"书二","author":"作者"}`)
					}
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()
			client, err := bqg349_service.NewClient(bqg349_service.DefaultSite, server.URL)
			if err != nil {
				t.Fatal(err)
			}
			var saved []bqg349_service.ID
			var mu sync.Mutex
			err = walkBqg349Site(context.Background(), client, func(ctx context.Context, book *bqg349_service.Book) error {
				mu.Lock()
				saved = append(saved, book.ID)
				mu.Unlock()
				if failFirst && book.ID == 2530 {
					return errors.New("模拟入库失败")
				}
				return nil
			})
			sort.Slice(saved, func(i, j int) bool { return saved[i] < saved[j] })
			if (err != nil) != failFirst || !reflect.DeepEqual(saved, []bqg349_service.ID{2530, 2531}) {
				t.Fatalf("空库发现新书、分类去重或失败后继续异常: saved=%v err=%v", saved, err)
			}
		})
	}
}

type bqg349TestSite struct {
	books []bqg349_service.Book
}

func newBqg349TestSite(count int) *bqg349TestSite {
	site := &bqg349TestSite{}
	for i := 1; i <= count; i++ {
		site.books = append(site.books, bqg349_service.Book{ID: bqg349_service.ID(i), Title: fmt.Sprintf("测试书%d", i), Author: "作者", URL: fmt.Sprintf("book-%d", i)})
	}
	return site
}

func (site *bqg349TestSite) Categories() []bqg349_service.Category {
	return []bqg349_service.Category{{Name: "分类一", URL: "one"}, {Name: "分类二", URL: "two"}}
}

func (site *bqg349TestSite) Books(ctx context.Context, _ string) ([]bqg349_service.Book, error) {
	return site.books, ctx.Err()
}

func (site *bqg349TestSite) Book(ctx context.Context, url string) (*bqg349_service.Book, error) {
	for _, book := range site.books {
		if book.URL == url {
			return &book, ctx.Err()
		}
	}
	return nil, errors.New("没有此书")
}

func TestBqg349ProcessesMultipleBooksAndContinuesAfterFailure(t *testing.T) {
	site := newBqg349TestSite(bqg349BookWorkers + 3)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	started := make(chan struct{}, len(site.books))
	release := make(chan struct{})
	var active, peak atomic.Int32
	counts := make(map[bqg349_service.ID]int)
	var mu sync.Mutex
	done := make(chan error, 1)
	go func() {
		done <- walkBqg349Site(ctx, site, func(ctx context.Context, book *bqg349_service.Book) error {
			current := active.Add(1)
			defer active.Add(-1)
			for previous := peak.Load(); current > previous; previous = peak.Load() {
				if peak.CompareAndSwap(previous, current) {
					break
				}
			}
			mu.Lock()
			counts[book.ID]++
			mu.Unlock()
			started <- struct{}{}
			select {
			case <-release:
			case <-ctx.Done():
				return ctx.Err()
			}
			if book.ID == 1 {
				return errors.New("一本书暂时失败")
			}
			return nil
		})
	}()
	for i := 0; i < bqg349BookWorkers; i++ {
		select {
		case <-started:
		case <-ctx.Done():
			t.Fatal("其他书籍仍在等待第一本完成")
		}
	}
	close(release)
	select {
	case err := <-done:
		if err == nil || peak.Load() != bqg349BookWorkers || active.Load() != 0 {
			t.Fatalf("书籍并发上限、错误汇总或退出有误: peak=%d active=%d err=%v", peak.Load(), active.Load(), err)
		}
	case <-ctx.Done():
		t.Fatal("多本书采集未结束")
	}
	if len(counts) != len(site.books) {
		t.Fatalf("失败后未继续其他书籍: %v", counts)
	}
	for id, count := range counts {
		if count != 1 {
			t.Fatalf("同一源 ID 跨分类重复处理: id=%d count=%d", id, count)
		}
	}
}

func TestBqg349DeduplicatesBookIdentityAcrossSourceIDs(t *testing.T) {
	site := newBqg349TestSite(2)
	site.books[1].Title, site.books[1].Author = site.books[0].Title, site.books[0].Author
	var saves atomic.Int32
	err := walkBqg349Site(context.Background(), site, func(context.Context, *bqg349_service.Book) error {
		saves.Add(1)
		return nil
	})
	if err != nil || saves.Load() != 1 {
		t.Fatalf("同名同作者的别名书重复入库或写文件: calls=%d err=%v", saves.Load(), err)
	}
}

func TestBqg349CancelConcurrentBooks(t *testing.T) {
	site := newBqg349TestSite(bqg349BookWorkers * 4)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started := make(chan struct{}, len(site.books))
	var active, saves atomic.Int32
	done := make(chan error, 1)
	go func() {
		done <- walkBqg349Site(ctx, site, func(ctx context.Context, book *bqg349_service.Book) error {
			active.Add(1)
			defer active.Add(-1)
			saves.Add(1)
			started <- struct{}{}
			<-ctx.Done()
			return ctx.Err()
		})
	}()
	for i := 0; i < bqg349BookWorkers; i++ {
		select {
		case <-started:
		case <-time.After(2 * time.Second):
			t.Fatal("多本书采集未启动")
		}
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) || active.Load() != 0 || saves.Load() != bqg349BookWorkers {
			t.Fatalf("停止后仍处理书籍或未等待退出: active=%d calls=%d err=%v", active.Load(), saves.Load(), err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("停止信号未取消排队和在采书籍")
	}
}

func TestBqg349BooksShareChapterDownloadLimit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dir := t.TempDir()
	started := make(chan struct{}, bqg349ChapterWorkers*4)
	attempted := make(chan struct{}, bqg349ChapterWorkers*4)
	release := make(chan struct{})
	var active, peak atomic.Int32
	fetch := limitBqg349ChapterFetch(func(ctx context.Context, chapter bqg349_service.Chapter) (string, error) {
		current := active.Add(1)
		defer active.Add(-1)
		for previous := peak.Load(); current > previous; previous = peak.Load() {
			if peak.CompareAndSwap(previous, current) {
				break
			}
		}
		started <- struct{}{}
		select {
		case <-release:
			return chapter.URL + "正文", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}, bqg349ChapterWorkers)
	done := make(chan error, 2)
	for book := 1; book <= 2; book++ {
		bookDir := filepath.Join(dir, fmt.Sprint(book))
		if err := os.MkdirAll(bookDir, 0755); err != nil {
			t.Fatal(err)
		}
		catalog := filepath.Join(bookDir, "chapters.json")
		if err := os.WriteFile(catalog, []byte("[]"), 0644); err != nil {
			t.Fatal(err)
		}
		var chapters []bqg349_service.Chapter
		for i := 1; i <= bqg349ChapterWorkers+3; i++ {
			chapters = append(chapters, bqg349_service.Chapter{Title: fmt.Sprintf("第%d章", i), URL: fmt.Sprintf("book-%d-chapter-%d", book, i)})
		}
		go func() {
			path := func(title string) (string, error) { return filepath.Join(bookDir, title+".txt"), nil }
			done <- saveBqg349Chapters(ctx, bookDir, chapters, catalog, path, func(ctx context.Context, chapter bqg349_service.Chapter) (string, error) {
				attempted <- struct{}{}
				return fetch(ctx, chapter)
			}, false)
		}()
	}
	for i := 0; i < bqg349ChapterWorkers*2; i++ {
		select {
		case <-attempted:
		case <-ctx.Done():
			t.Fatal("两本书未同时提交正文下载任务")
		}
	}
	for i := 0; i < bqg349ChapterWorkers; i++ {
		select {
		case <-started:
		case <-ctx.Done():
			t.Fatal("正文下载未填满全局并发额度")
		}
	}
	close(release)
	for i := 0; i < 2; i++ {
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
		case <-ctx.Done():
			t.Fatal("共享下载额度造成阻塞")
		}
	}
	if peak.Load() != bqg349ChapterWorkers || active.Load() != 0 {
		t.Fatalf("正文并发随书籍数量叠加或协程未退出: peak=%d active=%d", peak.Load(), active.Load())
	}
	for book := 1; book <= 2; book++ {
		catalog := filepath.Join(dir, fmt.Sprint(book), "chapters.json")
		saved, err := chapter_service.GetChaptersByFile(catalog)
		if err != nil || len(saved) != bqg349ChapterWorkers+3 {
			t.Fatalf("多书章节目录不完整: book=%d count=%d err=%v", book, len(saved), err)
		}
		for i, record := range saved {
			body, err := os.ReadFile(filepath.Join(dir, fmt.Sprint(book), record.ChapterName+".txt"))
			if record.IsLess != 0 || record.Sort != i+1 || record.Id != int64(i+1) || err != nil || string(body) != fmt.Sprintf("book-%d-chapter-%d正文", book, i+1) {
				t.Fatalf("多本书并发时串章或顺序错误: book=%d record=%+v body=%q err=%v", book, record, body, err)
			}
		}
	}
}

func TestBqg349CancelWaitingForSharedDownloadSlot(t *testing.T) {
	var calls atomic.Int32
	started := make(chan struct{}, 1)
	fetch := limitBqg349ChapterFetch(func(ctx context.Context, chapter bqg349_service.Chapter) (string, error) {
		calls.Add(1)
		if chapter.Title == "之后的请求" {
			return "正文", nil
		}
		started <- struct{}{}
		<-ctx.Done()
		return "", ctx.Err()
	}, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	first := make(chan error, 1)
	go func() {
		_, err := fetch(ctx, bqg349_service.Chapter{Title: "占用额度"})
		first <- err
	}()
	select {
	case <-started:
	case <-ctx.Done():
		t.Fatal("下载额度未被占用")
	}
	waitCtx, stopWaiting := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer stopWaiting()
	if _, err := fetch(waitCtx, bqg349_service.Chapter{Title: "排队中"}); !errors.Is(err, context.DeadlineExceeded) || calls.Load() != 1 {
		t.Fatalf("等待全局下载额度时无法取消: calls=%d err=%v", calls.Load(), err)
	}
	cancel()
	select {
	case err := <-first:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("占用额度的请求未取消: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("下载协程未释放额度")
	}
	nextCtx, stopNext := context.WithTimeout(context.Background(), time.Second)
	defer stopNext()
	if text, err := fetch(nextCtx, bqg349_service.Chapter{Title: "之后的请求"}); err != nil || text != "正文" || calls.Load() != 2 {
		t.Fatalf("取消后泄漏了下载额度: text=%q calls=%d err=%v", text, calls.Load(), err)
	}
}

// 使用当前连接上的临时表验证实际 MySQL 写入，不修改正式 mc_book 数据。
func TestBqg349MySQLBookAndChapterSync(t *testing.T) {
	if os.Getenv("BQG349_TEST_MYSQL") != "1" {
		t.Skip("设置 BQG349_TEST_MYSQL=1 执行 MySQL 临时表集成测试")
	}
	host, name, user, password := db.GetDB()
	config := drivermysql.NewConfig()
	config.Net, config.Addr, config.DBName = "tcp", host, name
	config.User, config.Passwd, config.Timeout = user, password, 4*time.Second
	database, err := gorm.Open(gormmysql.Open(config.FormatDSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	tx := database.Begin()
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	defer tx.Rollback()
	for _, table := range []string{"mc_book", "mc_book_class"} {
		var tableName, ddl string
		if err := tx.Raw("SHOW CREATE TABLE "+table).Row().Scan(&tableName, &ddl); err != nil {
			t.Fatal(err)
		}
		if err := tx.Exec(strings.Replace(ddl, "CREATE TABLE", "CREATE TEMPORARY TABLE", 1)).Error; err != nil {
			t.Fatal(err)
		}
		defer tx.Exec("DROP TEMPORARY TABLE IF EXISTS " + table)
	}
	ctx := context.Background()
	remote := &bqg349_service.Book{ID: 2530, Title: "采集测试书", Author: "测试作者", Category: "玄幻", Status: "连载", Intro: "简介", URL: bqg349_service.DefaultSite + "/#/book/2530/", CoverURL: "cover", LastChapter: "第二章", LastUpdate: "2026-11-13"}
	local, err := saveBqg349Book(ctx, tx, remote)
	if err != nil || local == nil || local.Id <= 0 {
		t.Fatalf("空库未新增书籍: local=%v err=%v", local, err)
	}
	var initial models.McBook
	if err := tx.Where("id = ?", local.Id).Take(&initial).Error; err != nil {
		t.Fatal(err)
	}
	if initial.ReadCount < 1 || initial.ReadCount >= 10000 || initial.SearchCount < 1 || initial.SearchCount >= 1000 || initial.SearchNum < initial.SearchCount || initial.SearchNum >= initial.SearchCount+3000 || initial.Hits < 5000 || initial.Hits >= 9999 || initial.Score < 6 || initial.Score >= 9 || initial.ClassName != remote.Category || initial.Status != 1 || initial.Addtime <= 0 || initial.Uptime < initial.Addtime {
		t.Fatalf("新书随机热度、分类或时间有误: %+v", initial)
	}
	remote.Status = "已完成"
	updated, err := saveBqg349Book(ctx, tx, remote)
	if err != nil || updated.Id != local.Id {
		t.Fatalf("重复采集未更新原记录: local=%v err=%v", updated, err)
	}
	catalog := filepath.Join(t.TempDir(), "chapters.json")
	data, _ := json.Marshal([]*models.McBookChapter{{Id: 1, Sort: 1, ChapterName: "第一章", TextNum: 5}, {Id: 2, Sort: 2, ChapterName: "第二章", TextNum: 3}, {Id: 3, Sort: 3, ChapterName: "第三章", IsLess: 1}})
	if err := os.WriteFile(catalog, data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := saveBqg349BookProgress(ctx, tx, updated, catalog); err != nil {
		t.Fatal(err)
	}
	var saved models.McBook
	if err := tx.Where("id = ?", local.Id).Take(&saved).Error; err != nil {
		t.Fatal(err)
	}
	var count int64
	if err := tx.Model(&models.McBook{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 1 || saved.Serialize != 2 || saved.Status != 1 || saved.ChapterNum != 3 || saved.TextNum != 8 || saved.UpdateChapterId != 2 || saved.UpdateChapterTitle != "第二章" || saved.SourceUrl != remote.URL || saved.Desc != "简介" {
		t.Fatalf("书籍或章节进度同步有误: count=%d saved=%+v", count, saved)
	}
	stats := func(book models.McBook) []any {
		return []any{book.Hits, book.HitsDay, book.HitsWeek, book.HitsMonth, book.Shits, book.Score, book.ReadCount, book.SearchCount, book.SearchNum, book.ClassName, book.Addtime}
	}
	if !reflect.DeepEqual(stats(initial), stats(saved)) {
		t.Fatalf("重复采集覆盖了已有热度或入库时间: before=%v after=%v", stats(initial), stats(saved))
	}
	// 已采集的零值补齐，有效统计和手动隐藏状态不覆盖。
	if err := tx.Model(&models.McBook{}).Where("id = ?", local.Id).Updates(map[string]any{"read_count": 0, "search_count": 0, "search_num": 0, "status": 0}).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := saveBqg349Book(ctx, tx, remote); err != nil {
		t.Fatal(err)
	}
	var filled models.McBook
	if err := tx.Where("id = ?", local.Id).Take(&filled).Error; err != nil {
		t.Fatal(err)
	}
	if filled.ReadCount <= 0 || filled.SearchCount <= 0 || filled.SearchNum < filled.SearchCount || filled.Hits != saved.Hits || filled.Status != 0 {
		t.Fatalf("补零值时覆盖了已有统计或隐藏状态: %+v", filled)
	}
}

func TestBqg349CollectorContinuesAfterSuccessAndFailure(t *testing.T) {
	for _, failFirst := range []bool{false, true} {
		name := "success"
		if failFirst {
			name = "temporary failure"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			calls := 0
			err := runBqg349Collector(ctx, time.Millisecond, func(ctx context.Context) error {
				calls++
				if calls == 2 {
					cancel()
				}
				if failFirst && calls == 1 {
					return errors.New("暂时失败")
				}
				return nil
			})
			if calls != 2 || !errors.Is(err, context.Canceled) {
				t.Fatalf("未持续采集或无法停止: calls=%d, err=%v", calls, err)
			}
		})
	}
}

func TestBqg349CollectorStopsDuringWait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started := make(chan struct{})
	result := make(chan error, 1)
	go func() {
		result <- runBqg349Collector(ctx, time.Hour, func(context.Context) error {
			close(started)
			return nil
		})
	}()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("未立即开始首轮采集")
	}
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("停止错误: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("停止信号未及时唤醒轮询")
	}
}

func TestBqg349SaveResumeAndOrder(t *testing.T) {
	dir := t.TempDir()
	catalog := filepath.Join(dir, "chapters.json")
	records := []*models.McBookChapter{
		{Id: 10, Sort: 1, ChapterName: "第一章"},
		{Id: 30, Sort: 3, ChapterName: "第三章"},
	}
	data, _ := json.Marshal(records)
	if err := os.WriteFile(catalog, data, 0644); err != nil {
		t.Fatal(err)
	}
	path := func(title string) (string, error) { return filepath.Join(dir, title+".txt"), nil }
	for _, record := range records {
		file, _ := path(record.ChapterName)
		if err := os.WriteFile(file, []byte("已有正文"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	chapters := []bqg349_service.Chapter{
		{Title: "第一章", URL: "first"}, {Title: "第二章", URL: "second"}, {Title: "第三章", URL: "third"},
	}
	var fetched []string
	body := "正文，http words...；\n第二段"
	fetch := func(ctx context.Context, chapter bqg349_service.Chapter) (string, error) {
		fetched = append(fetched, chapter.Title)
		return body, nil
	}
	if err := saveBqg349Chapters(context.Background(), "测试书", chapters, catalog, path, fetch, false); err != nil {
		t.Fatal(err)
	}
	saved, err := chapter_service.GetChaptersByFile(catalog)
	if err != nil || len(saved) != 3 {
		t.Fatalf("saved = %v, err = %v", saved, err)
	}
	wantIDs := []int64{10, 31, 30}
	for i, record := range saved {
		if record.Id != wantIDs[i] || record.Sort != i+1 || record.ChapterName != chapters[i].Title || record.ChapterLink != chapters[i].URL {
			t.Fatalf("章节 ID、顺序或链接有误: %+v", record)
		}
	}
	file, _ := path("第二章")
	content, err := os.ReadFile(file)
	if err != nil || string(content) != body || !reflect.DeepEqual(fetched, []string{"第二章"}) {
		t.Fatalf("正文或下载列表有误: %q, %v, %v", content, fetched, err)
	}
	if err := saveBqg349Chapters(context.Background(), "测试书", chapters, catalog, path, fetch, false); err != nil {
		t.Fatal(err)
	}
	if len(fetched) != 1 {
		t.Fatalf("续采重复下载已有正文: %v", fetched)
	}
}

func TestBqg349SaveCheckpointOnFailure(t *testing.T) {
	dir := t.TempDir()
	catalog := filepath.Join(dir, "chapters.json")
	if err := os.WriteFile(catalog, []byte("[]"), 0644); err != nil {
		t.Fatal(err)
	}
	path := func(title string) (string, error) { return filepath.Join(dir, title+".txt"), nil }
	chapters := []bqg349_service.Chapter{{Title: "第一章", URL: "first"}, {Title: "公告", URL: "notice"}, {Title: "第三章", URL: "third"}}
	failed := errors.New("临时网络错误")
	fetch := func(ctx context.Context, chapter bqg349_service.Chapter) (string, error) {
		if chapter.Title == "公告" {
			return "", failed
		}
		return "第一章正文", nil
	}
	if err := saveBqg349Chapters(context.Background(), "测试书", chapters, catalog, path, fetch, false); !errors.Is(err, failed) {
		t.Fatalf("未传递错误: %v", err)
	}
	saved, err := chapter_service.GetChaptersByFile(catalog)
	if err != nil || len(saved) != 3 || saved[0].ChapterName != "第一章" || saved[0].IsLess != 0 || saved[1].IsLess != 1 || saved[2].IsLess != 0 {
		t.Fatalf("完整目录或正文完成状态未保存: %v, %v", saved, err)
	}
	third, _ := path("第三章")
	if _, err := os.Stat(third); err != nil {
		t.Fatalf("并发预下载的正文未保留供续采复用: %v", err)
	}
	fetch = func(ctx context.Context, chapter bqg349_service.Chapter) (string, error) {
		if chapter.Title != "公告" {
			t.Errorf("续采不应重复下载已保存正文: %s", chapter.Title)
		}
		return "短公告", nil
	}
	if err := saveBqg349Chapters(context.Background(), "测试书", chapters, catalog, path, fetch, false); err != nil {
		t.Fatal(err)
	}
	saved, err = chapter_service.GetChaptersByFile(catalog)
	if err != nil || len(saved) != 3 || saved[1].TextNum != 3 {
		t.Fatalf("续采或短篇保留失败: %v, %v", saved, err)
	}
}

func TestBqg349ConcurrentDownloadOrderAndResume(t *testing.T) {
	dir := t.TempDir()
	catalog := filepath.Join(dir, "chapters.json")
	if err := os.WriteFile(catalog, []byte("[]"), 0644); err != nil {
		t.Fatal(err)
	}
	path := func(title string) (string, error) { return filepath.Join(dir, title+".txt"), nil }
	var chapters []bqg349_service.Chapter
	for i := 1; i <= 58; i++ {
		chapters = append(chapters, bqg349_service.Chapter{Title: fmt.Sprintf("第%d章", i), URL: fmt.Sprintf("chapter-%d", i)})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	started := make(chan struct{}, len(chapters))
	release := make(chan struct{})
	second := make(chan struct{})
	var active, peak atomic.Int32
	fetch := func(ctx context.Context, chapter bqg349_service.Chapter) (string, error) {
		catalogRecords, err := chapter_service.GetChaptersByFile(catalog)
		if err != nil || len(catalogRecords) != len(chapters) {
			return "", fmt.Errorf("正文下载之前未保存完整目录: count=%d err=%v", len(catalogRecords), err)
		}
		current := active.Add(1)
		defer active.Add(-1)
		for previous := peak.Load(); current > previous; previous = peak.Load() {
			if peak.CompareAndSwap(previous, current) {
				break
			}
		}
		started <- struct{}{}
		select {
		case <-release:
		case <-ctx.Done():
			return "", ctx.Err()
		}
		// 强制第二章先于第一章返回，验证完成顺序不会改变目录和 ID。
		if chapter.Title == "第2章" {
			close(second)
		}
		if chapter.Title == "第1章" {
			select {
			case <-second:
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
		return chapter.Title + "正文", nil
	}
	done := make(chan error, 1)
	go func() { done <- saveBqg349Chapters(ctx, "测试书", chapters, catalog, path, fetch, false) }()
	for i := 0; i < bqg349ChapterWorkers; i++ {
		select {
		case <-started:
		case <-ctx.Done():
			t.Fatalf("未启动 %d 个并发下载协程", bqg349ChapterWorkers)
		}
	}
	close(release)
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("并发下载未完成")
	}
	if peak.Load() != bqg349ChapterWorkers || active.Load() != 0 {
		t.Fatalf("并发数未受限或协程未退出: peak=%d active=%d", peak.Load(), active.Load())
	}
	saved, err := chapter_service.GetChaptersByFile(catalog)
	if err != nil || len(saved) != len(chapters) {
		t.Fatalf("目录未完整保存: count=%d err=%v", len(saved), err)
	}
	for i, record := range saved {
		file, _ := path(chapters[i].Title)
		body, err := os.ReadFile(file)
		if record.Id != int64(i+1) || record.Sort != i+1 || record.ChapterName != chapters[i].Title || record.ChapterLink != chapters[i].URL || err != nil || string(body) != chapters[i].Title+"正文" {
			t.Fatalf("并发下载造成章节错序或串文: index=%d record=%+v body=%q err=%v", i, record, body, err)
		}
	}
	if err := saveBqg349Chapters(ctx, "测试书", chapters, catalog, path, func(context.Context, bqg349_service.Chapter) (string, error) {
		t.Error("完整正文不应重复下载")
		return "", errors.New("重复下载")
	}, false); err != nil {
		t.Fatal(err)
	}
}

func TestBqg349CancelConcurrentDownloads(t *testing.T) {
	dir := t.TempDir()
	catalog := filepath.Join(dir, "chapters.json")
	if err := os.WriteFile(catalog, []byte("[]"), 0644); err != nil {
		t.Fatal(err)
	}
	path := func(title string) (string, error) { return filepath.Join(dir, title+".txt"), nil }
	var chapters []bqg349_service.Chapter
	for i := 1; i <= 25; i++ {
		chapters = append(chapters, bqg349_service.Chapter{Title: fmt.Sprintf("第%d章", i), URL: fmt.Sprintf("chapter-%d", i)})
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started := make(chan struct{}, len(chapters))
	var active atomic.Int32
	fetch := func(ctx context.Context, chapter bqg349_service.Chapter) (string, error) {
		active.Add(1)
		defer active.Add(-1)
		started <- struct{}{}
		<-ctx.Done()
		return "", ctx.Err()
	}
	done := make(chan error, 1)
	go func() { done <- saveBqg349Chapters(ctx, "测试书", chapters, catalog, path, fetch, false) }()
	for i := 0; i < bqg349ChapterWorkers; i++ {
		select {
		case <-started:
		case <-time.After(2 * time.Second):
			t.Fatal("下载协程未启动")
		}
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) || active.Load() != 0 {
			t.Fatalf("停止后未取消全部协程: active=%d err=%v", active.Load(), err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("并发下载停止后未及时退出")
	}
	saved, err := chapter_service.GetChaptersByFile(catalog)
	if err != nil || len(saved) != len(chapters) {
		t.Fatalf("停止后丢失了完整目录: saved=%v err=%v", saved, err)
	}
	for _, chapter := range saved {
		if chapter.IsLess != 1 || chapter.TextNum != 0 {
			t.Fatalf("取消时将未下载正文错误标记为完成: %+v", chapter)
		}
	}
}

func TestBqg349ChapterOrder(t *testing.T) {
	for _, tc := range []struct {
		name    string
		local   []string
		remote  []string
		want    []string
		wantErr bool
	}{
		{"middle gap", []string{"一", "三"}, []string{"一", "二", "三"}, []string{"一", "二", "三"}, false},
		{"local exclusive", []string{"一", "旧站番外", "三"}, []string{"一", "二", "三", "四"}, []string{"一", "二", "旧站番外", "三", "四"}, false},
		{"new book", nil, []string{"一", "二"}, []string{"一", "二"}, false},
		{"conflicting order", []string{"一", "二"}, []string{"二", "一"}, nil, true},
		{"no shared titles", []string{"第一章", "第二章"}, []string{"一", "二"}, nil, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var local []*models.McBookChapter
			for _, title := range tc.local {
				local = append(local, &models.McBookChapter{ChapterName: title})
			}
			var remote []bqg349_service.Chapter
			for _, title := range tc.remote {
				remote = append(remote, bqg349_service.Chapter{Title: title})
			}
			got, err := bqg349ChapterOrder(local, remote)
			if (err != nil) != tc.wantErr || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("order = %v, err = %v; want %v", got, err, tc.want)
			}
		})
	}
}

func TestBqg349FailedOverwritePreservesText(t *testing.T) {
	dir := t.TempDir()
	catalog, textFile := filepath.Join(dir, "chapters.json"), filepath.Join(dir, "chapter.txt")
	original := []byte("原正文，保留 English 和标点。")
	if err := os.WriteFile(catalog, []byte("[]"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(textFile, original, 0640); err != nil {
		t.Fatal(err)
	}
	chapters := []bqg349_service.Chapter{{Title: "第一章", URL: "first"}}
	path := func(string) (string, error) { return textFile, nil }
	fetch := func(context.Context, bqg349_service.Chapter) (string, error) { return "", errors.New("下载失败") }
	if err := saveBqg349Chapters(context.Background(), "测试书", chapters, catalog, path, fetch, true); err == nil {
		t.Fatal("覆盖失败未返回错误")
	}
	content, err := os.ReadFile(textFile)
	if err != nil || string(content) != string(original) {
		t.Fatalf("下载失败破坏了原正文: %q, %v", content, err)
	}
	if err := bqg349WriteFile(textFile, []byte("新正文")); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(textFile)
	if err != nil || info.Mode().Perm() != 0640 {
		t.Fatalf("替换未保留原权限: %v, %v", info, err)
	}
}
