/*
@Descripttion: bqg349 真实协议样例与异常响应测试，不依赖外网或数据库
@Author: red
@Date: 2026-11-13
@LastEditors: red
@LastEditTime: 2026-11-13
*/
package bqg349_service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// 预期值来自 Playwright 调用站点 enaes/get_api 的结果。
func TestBrowserTokens(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value any
		want  string
	}{
		{"book", struct {
			ID int `json:"id"`
		}{2530}, "IZzVX5XJA1oySI5jozY8vQ=="},
		{"chapter", struct {
			ID        int `json:"id"`
			ChapterID int `json:"chapterid"`
		}{2530, 1}, "b+vXnT3wjuXQsxBmZh033ZjqwezLEinKfOakcVaiDx0="},
	} {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tokenQuery(tc.value)
			if err != nil || query.Get("token") != tc.want {
				t.Fatalf("token = %q, err = %v; want %q", query.Get("token"), err, tc.want)
			}
		})
	}
}

func TestSourceFlowWithAliasDirectory(t *testing.T) {
	var directoryID ID
	directoryQuery, _ := tokenQuery(struct {
		ID ID `json:"id"`
	}{99})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Referer") != DefaultSite+"/" {
			t.Error("缺少站点 Referer")
		}
		switch r.URL.Path {
		case "/api/sort":
			if r.URL.Query().Get("sort") != "xuanhuan" {
				t.Error("分类参数错误")
			}
			fmt.Fprint(w, `{"data":[{"id":"2530","title":" 测试书 ","author":"\n 作者 \n"}]}`)
		case "/api/book":
			fmt.Fprint(w, `{"id":"2530","dirid":99,"title":"测试书","author":"作者","full":"完结","sortname":"玄幻"}`)
		case "/api/booklist":
			if r.URL.Query().Get("token") != directoryQuery.Get("token") {
				t.Error("目录应使用 dirid")
			}
			directoryID = 99
			// 短篇、番外和公告也必须保留，链接编号由原始数组位置决定。
			fmt.Fprint(w, `{"list":["第一章 开始","番外","完本感言"]}`)
		case "/api/chapter":
			fmt.Fprint(w, `{"id":2530,"chapterid":"1","chaptername":"第一章 开始","txt":"正文 http 示例，保留英文 words 和标点。\r\n第二段。"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client, err := NewClient(DefaultSite+"/#/", server.URL)
	if err != nil {
		t.Fatal(err)
	}
	client.interval = 0
	ctx := context.Background()
	books, err := client.Books(ctx, client.Categories()[0].URL)
	if err != nil || len(books) != 1 {
		t.Fatalf("books = %v, err = %v", books, err)
	}
	if books[0].Title != "测试书" || books[0].Author != "作者" {
		t.Fatalf("未清理元数据空白: %+v", books[0])
	}
	book, err := client.Book(ctx, books[0].URL)
	if err != nil {
		t.Fatal(err)
	}
	if book.CoverURL != DefaultSite+"/bookimg/2/2530.jpg" {
		t.Fatalf("封面地址: %s", book.CoverURL)
	}
	chapters, err := client.Chapters(ctx, book)
	if err != nil || len(chapters) != 3 || directoryID != 99 {
		t.Fatalf("chapters = %v, err = %v", chapters, err)
	}
	if chapters[2].URL != DefaultSite+"/#/book/2530/3.html" {
		t.Fatalf("章节序号错误: %+v", chapters)
	}
	text, err := client.Text(ctx, chapters[0].URL)
	if err != nil || text != "正文 http 示例，保留英文 words 和标点。\n第二段。" {
		t.Fatalf("text = %q, err = %v", text, err)
	}
}

func TestInvalidResponses(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{"HTTP error", 503, `unavailable`},
		{"HTML challenge", 200, `<html>challenge</html>`},
		{"JSON error", 200, `{"error":"not found"}`},
		{"empty text", 200, `{"id":2530,"chapterid":1,"chaptername":"第一章","txt":" "}`},
		{"wrong chapter", 200, `{"id":2530,"chapterid":2,"chaptername":"第二章","txt":"其他正文"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				fmt.Fprint(w, tc.body)
			}))
			defer server.Close()
			client, _ := NewClient(DefaultSite, server.URL)
			if text, err := client.Text(context.Background(), DefaultSite+"/#/book/2530/1.html"); err == nil || text != "" {
				t.Fatalf("错误响应被当成正文: %q, %v", text, err)
			}
		})
	}
}

func TestInvalidLinksAndCancellation(t *testing.T) {
	client, _ := NewClient(DefaultSite, DefaultAPI)
	for _, link := range []string{"https://other.example/#/book/2530/1.html", DefaultSite + "/#/book/0/1.html", DefaultSite + "/#/book/2530/1_2.html", DefaultSite + "/#/book/2530/"} {
		if _, err := client.Text(context.Background(), link); err == nil {
			t.Errorf("接受了无效章节地址: %s", link)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := client.Book(ctx, DefaultSite+"/#/book/2530/"); err == nil || !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("取消请求未生效: %v", err)
	}
}

func TestFindBookRequiresUniqueIdentity(t *testing.T) {
	for _, tc := range []struct {
		name   string
		books  string
		detail string
		wantOK bool
	}{
		{"exact match", `[{"id":2,"title":"测试书番外","author":"作者"},{"id":3,"title":"测试书","author":"另一作者"},{"id":2530,"title":" 测试书 ","author":"作者"}]`, "", true},
		{"duplicate result", `[{"id":2530,"title":"测试书","author":"作者"},{"id":2530,"title":"测试书","author":"作者"}]`, "", true},
		{"ambiguous", `[{"id":2530,"title":"测试书","author":"作者"},{"id":2531,"title":"测试书","author":"作者"}]`, "", false},
		{"different author", `[{"id":2530,"title":"测试书","author":"另一作者"}]`, "", false},
		{"missing", `[]`, "", false},
		{"changed detail", `[{"id":2530,"title":"测试书","author":"作者"}]`, "另一作者", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/search":
					if r.URL.Query().Get("q") != "测试书" {
						t.Error("搜索关键词未正确传递")
					}
					fmt.Fprintf(w, `{"data":%s}`, tc.books)
				case "/api/book":
					author := tc.detail
					if author == "" {
						author = "作者"
					}
					fmt.Fprintf(w, `{"id":2530,"dirid":99,"title":"测试书","author":%q}`, author)
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()
			client, _ := NewClient(DefaultSite, server.URL)
			client.interval = 0
			book, err := client.FindBook(context.Background(), " 测试书 ", "作者")
			if tc.wantOK {
				if err != nil || book == nil || book.ID != 2530 || book.DirectoryID != 99 {
					t.Fatalf("book = %+v, err = %v", book, err)
				}
			} else if err == nil || book != nil {
				t.Fatalf("不明确的书籍被接受: %+v, %v", book, err)
			}
		})
	}
}

func TestChapterTitleValidation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":2530,"chapterid":1,"chaptername":"新标题","txt":"短篇，http words...；\r\n下一行\r末行"}`)
	}))
	defer server.Close()
	client, _ := NewClient(DefaultSite, server.URL)
	client.interval = 0
	chapter := Chapter{Title: "旧标题", URL: DefaultSite + "/#/book/2530/1.html"}
	if text, err := client.ChapterText(context.Background(), chapter); err == nil || text != "" {
		t.Fatalf("目录标题已变化仍返回正文: %q, %v", text, err)
	}
	chapter.Title = "新标题"
	if text, err := client.ChapterText(context.Background(), chapter); err != nil || text != "短篇，http words...；\n下一行\n末行" {
		t.Fatalf("正文被误清洗: %q, %v", text, err)
	}
}

func TestRequestRetryPolicy(t *testing.T) {
	for _, tc := range []struct {
		name       string
		status     int
		retryAfter string
		attempts   int32
		wantOK     bool
	}{
		{"temporary failure", 503, "", 2, true},
		{"rate limit", 429, "1", 2, true},
		{"missing book", 404, "", 1, false},
		{"long retry after", 429, "120", 1, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if attempts.Add(1) == 1 {
					w.Header().Set("Retry-After", tc.retryAfter)
					w.WriteHeader(tc.status)
					return
				}
				fmt.Fprint(w, `{"id":2530,"chapterid":1,"chaptername":"第一章","txt":"正文"}`)
			}))
			defer server.Close()
			client, _ := NewClient(DefaultSite, server.URL)
			client.interval = 0
			text, err := client.Text(context.Background(), DefaultSite+"/#/book/2530/1.html")
			if attempts.Load() != tc.attempts || (err == nil) != tc.wantOK {
				t.Fatalf("attempts = %d, text = %q, err = %v", attempts.Load(), text, err)
			}
		})
	}
}

func TestCancelWhileWaiting(t *testing.T) {
	t.Run("rate limiter", func(t *testing.T) {
		client, _ := NewClient(DefaultSite, DefaultAPI)
		client.next = time.Now().Add(time.Hour)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()
		if err := client.wait(ctx); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("等待限速时未取消: %v", err)
		}
	})
	t.Run("retry backoff", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var attempts atomic.Int32
		requested := make(chan struct{})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.Header().Set("Retry-After", "30")
			w.WriteHeader(http.StatusTooManyRequests)
			close(requested)
		}))
		defer server.Close()
		client, _ := NewClient(DefaultSite, server.URL)
		response := make(chan error, 1)
		go func() {
			_, err := client.Text(ctx, DefaultSite+"/#/book/2530/1.html")
			response <- err
		}()
		deadline := time.NewTimer(2 * time.Second)
		defer deadline.Stop()
		select {
		case <-requested:
		case <-deadline.C:
			t.Fatal("未发起请求")
		}
		cancel()
		select {
		case err := <-response:
			if !errors.Is(err, context.Canceled) || attempts.Load() != 1 {
				t.Fatalf("取消后仍在重试: %v, attempts=%d", err, attempts.Load())
			}
		case <-deadline.C:
			t.Fatal("取消未及时结束重试")
		}
	})
}
