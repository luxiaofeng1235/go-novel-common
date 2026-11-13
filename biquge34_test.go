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
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"go-novel/app/models"
	"go-novel/app/service/collect/bqg349_service"
	"go-novel/app/service/common/chapter_service"
)

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
