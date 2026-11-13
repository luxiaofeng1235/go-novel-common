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
	"strings"
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
			err = walkBqg349Site(context.Background(), client, func(ctx context.Context, book *bqg349_service.Book) error {
				saved = append(saved, book.ID)
				if failFirst && book.ID == 2530 {
					return errors.New("模拟入库失败")
				}
				return nil
			})
			if (err != nil) != failFirst || !reflect.DeepEqual(saved, []bqg349_service.ID{2530, 2531}) {
				t.Fatalf("空库发现新书、分类去重或失败后继续异常: saved=%v err=%v", saved, err)
			}
		})
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
	remote.Status = "已完成"
	updated, err := saveBqg349Book(ctx, tx, remote)
	if err != nil || updated.Id != local.Id {
		t.Fatalf("重复采集未更新原记录: local=%v err=%v", updated, err)
	}
	catalog := filepath.Join(t.TempDir(), "chapters.json")
	data, _ := json.Marshal([]*models.McBookChapter{{Id: 1, Sort: 1, ChapterName: "第一章", TextNum: 5}, {Id: 2, Sort: 2, ChapterName: "第二章", TextNum: 3}})
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
	if count != 1 || saved.Serialize != 2 || saved.Status != 1 || saved.ChapterNum != 2 || saved.TextNum != 8 || saved.UpdateChapterId != 2 || saved.UpdateChapterTitle != "第二章" || saved.SourceUrl != remote.URL || saved.Desc != "简介" {
		t.Fatalf("书籍或章节进度同步有误: count=%d saved=%+v", count, saved)
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
	if err := saveBqg349Chapters(context.Background(), chapters, catalog, path, fetch, false); err != nil {
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
	if err := saveBqg349Chapters(context.Background(), chapters, catalog, path, fetch, false); err != nil {
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
		if chapter.Title == "第三章" {
			t.Error("失败后不应继续下载后面的章节")
		}
		return "第一章正文", nil
	}
	if err := saveBqg349Chapters(context.Background(), chapters, catalog, path, fetch, false); !errors.Is(err, failed) {
		t.Fatalf("未传递错误: %v", err)
	}
	saved, err := chapter_service.GetChaptersByFile(catalog)
	if err != nil || len(saved) != 1 || saved[0].ChapterName != "第一章" {
		t.Fatalf("失败前完成的章节未保存: %v, %v", saved, err)
	}
	fetch = func(ctx context.Context, chapter bqg349_service.Chapter) (string, error) {
		if chapter.Title == "第一章" {
			t.Error("续采不应重复下载第一章")
		}
		return "短公告", nil
	}
	if err := saveBqg349Chapters(context.Background(), chapters, catalog, path, fetch, false); err != nil {
		t.Fatal(err)
	}
	saved, err = chapter_service.GetChaptersByFile(catalog)
	if err != nil || len(saved) != 3 || saved[1].TextNum != 3 {
		t.Fatalf("续采或短篇保留失败: %v, %v", saved, err)
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
	if err := saveBqg349Chapters(context.Background(), chapters, catalog, path, fetch, true); err == nil {
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
