/*
@Descripttion: bqg349 动态书源接口适配，规则来自站点 read.js/common.js
@Author: red
@Date: 2026-11-13
@LastEditors: red
@LastEditTime: 2026-11-13
*/
package bqg349_service

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DefaultSite     = "https://www.bqg349.xyz"
	DefaultAPI      = "https://apibi.cc"
	maxResponseSize = 16 << 20
)

// ID 兼容接口中数字和字符串两种编号格式。
type ID int64

func (id *ID) UnmarshalJSON(data []byte) error {
	n, err := strconv.ParseInt(strings.Trim(string(data), "\""), 10, 64)
	if err != nil || n < 0 {
		return fmt.Errorf("无效的书源编号 %s", data)
	}
	*id = ID(n)
	return nil
}

type Category struct {
	Key  string
	Name string
	URL  string
}

type Book struct {
	ID            ID     `json:"id"`
	DirectoryID   ID     `json:"dirid"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Intro         string `json:"intro"`
	Category      string `json:"sortname"`
	Status        string `json:"full"`
	LastChapterID ID     `json:"lastchapterid"`
	LastChapter   string `json:"lastchapter"`
	LastUpdate    string `json:"lastupdate"`
	URL           string `json:"-"`
	CoverURL      string `json:"-"`
}

type Chapter struct {
	Title string
	URL   string
}

type Client struct {
	site     string
	api      string
	http     *http.Client
	mu       sync.Mutex
	next     time.Time
	interval time.Duration
}

func NewClient(site, api string) (*Client, error) {
	normalize := func(raw string) (string, error) {
		u, err := url.Parse(strings.TrimSpace(raw))
		if err != nil || u == nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" || u.User != nil || (u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.ForceQuery || (u.Fragment != "" && u.Fragment != "/") {
			return "", fmt.Errorf("无效的书源地址 %q", raw)
		}
		return u.Scheme + "://" + u.Host, nil
	}
	site, err := normalize(site)
	if err != nil {
		return nil, err
	}
	api, err = normalize(api)
	if err != nil {
		return nil, err
	}
	return &Client{site: site, api: api, http: &http.Client{Timeout: 30 * time.Second}, interval: 250 * time.Millisecond}, nil
}

// CachePrefix 隔离不同站点/接口的进度，避免读取旧书源缓存。
func (c *Client) CachePrefix() string {
	return fmt.Sprintf("bqg349:%x", md5.Sum([]byte(c.site+"|"+c.api)))
}

func (c *Client) Categories() []Category {
	categories := []Category{
		{Key: "xuanhuan", Name: "玄幻"}, {Key: "wuxia", Name: "武侠"},
		{Key: "dushi", Name: "都市"}, {Key: "lishi", Name: "历史"},
		{Key: "wangyou", Name: "网游"}, {Key: "kehuan", Name: "科幻"},
		{Key: "mm", Name: "女生"}, {Key: "finish", Name: "完本"},
	}
	for i := range categories {
		categories[i].URL = c.site + "/#/" + categories[i].Key
	}
	return categories
}

// Books 与网页一致：sort 接口一次返回整份列表，滚动加载没有服务端翻页。
func (c *Client) Books(ctx context.Context, categoryURL string) ([]Book, error) {
	var key string
	for _, category := range c.Categories() {
		if categoryURL == category.URL {
			key = category.Key
			break
		}
	}
	if key == "" {
		return nil, fmt.Errorf("未知分类地址 %q", categoryURL)
	}
	return c.books(ctx, "sort", url.Values{"sort": {key}})
}

// Search 使用站点公开搜索接口；旧站编号不能直接当作新站编号。
func (c *Client) Search(ctx context.Context, keyword string) ([]Book, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, fmt.Errorf("搜索关键词为空")
	}
	return c.books(ctx, "search", url.Values{"q": {keyword}})
}

// FindBook 只接受书名和作者均一致的唯一结果，避免同名书或别名书串书。
func (c *Client) FindBook(ctx context.Context, title, author string) (*Book, error) {
	title, author = strings.TrimSpace(title), strings.TrimSpace(author)
	if title == "" || author == "" {
		return nil, fmt.Errorf("书名或作者为空")
	}
	books, err := c.Search(ctx, title)
	if err != nil {
		return nil, err
	}
	var match *Book
	for i := range books {
		book := &books[i]
		if book.Title != title || book.Author != author {
			continue
		}
		if match != nil && match.ID != book.ID {
			return nil, fmt.Errorf("书籍 %q 作者 %q 在新站有多个匹配，请指定书籍链接", title, author)
		}
		match = book
	}
	if match == nil {
		return nil, fmt.Errorf("新站未找到书名和作者完全匹配的书籍: %q / %q", title, author)
	}
	book, err := c.Book(ctx, match.URL)
	if err != nil {
		return nil, err
	}
	if book.Title != title || book.Author != author {
		return nil, fmt.Errorf("书籍详情与搜索结果不一致: %q / %q", book.Title, book.Author)
	}
	return book, nil
}

func (c *Client) books(ctx context.Context, endpoint string, query url.Values) ([]Book, error) {
	var result struct {
		Data *[]Book `json:"data"`
	}
	if err := c.get(ctx, endpoint, query, &result); err != nil {
		return nil, err
	}
	if result.Data == nil {
		return nil, fmt.Errorf("%s 接口缺少 data 列表", endpoint)
	}
	for i := range *result.Data {
		book := &(*result.Data)[i]
		book.Title = strings.TrimSpace(book.Title)
		book.Author = strings.TrimSpace(book.Author)
		if book.ID <= 0 || book.Title == "" || book.Author == "" {
			return nil, fmt.Errorf("%s 接口第 %d 本书的信息不完整", endpoint, i+1)
		}
		book.URL = c.bookURL(book.ID)
		book.CoverURL = c.coverURL(book.ID)
	}
	return *result.Data, nil
}

func (c *Client) Book(ctx context.Context, link string) (*Book, error) {
	id, chapterID, err := c.parseLink(link)
	if err != nil || chapterID != 0 {
		return nil, fmt.Errorf("无效的书籍地址 %q", link)
	}
	query, err := tokenQuery(struct {
		ID ID `json:"id"`
	}{id})
	if err != nil {
		return nil, err
	}
	var book Book
	if err := c.get(ctx, "book", query, &book); err != nil {
		return nil, err
	}
	book.Title = strings.TrimSpace(book.Title)
	book.Author = strings.TrimSpace(book.Author)
	if book.ID != id || book.DirectoryID <= 0 || book.Title == "" || book.Author == "" {
		return nil, fmt.Errorf("书籍 %d 的详情不完整或编号不匹配", id)
	}
	book.URL = c.bookURL(id)
	book.CoverURL = c.coverURL(id)
	return &book, nil
}

func (c *Client) Chapters(ctx context.Context, book *Book) ([]Chapter, error) {
	if book == nil || book.ID <= 0 || book.DirectoryID <= 0 {
		return nil, fmt.Errorf("书籍或目录编号为空")
	}
	// 部分别名书的 dirid 与书籍 id 不同；目录用 dirid，阅读链接仍用 id。
	query, err := tokenQuery(struct {
		ID ID `json:"id"`
	}{book.DirectoryID})
	if err != nil {
		return nil, err
	}
	var result struct {
		List []string `json:"list"`
	}
	if err := c.get(ctx, "booklist", query, &result); err != nil {
		return nil, err
	}
	if len(result.List) == 0 {
		return nil, fmt.Errorf("书籍 %d 的目录为空", book.ID)
	}
	chapters := make([]Chapter, 0, len(result.List))
	for i, title := range result.List {
		title = strings.TrimSpace(title)
		if title == "" {
			return nil, fmt.Errorf("书籍 %d 第 %d 条目录标题为空", book.ID, i+1)
		}
		chapters = append(chapters, Chapter{Title: title, URL: fmt.Sprintf("%s%d.html", c.bookURL(book.ID), i+1)})
	}
	return chapters, nil
}

func (c *Client) Text(ctx context.Context, link string) (string, error) {
	return c.text(ctx, link, "")
}

// ChapterText 校验目录标题，避免目录缓存过期后把其他章节正文写入旧标题。
func (c *Client) ChapterText(ctx context.Context, chapter Chapter) (string, error) {
	if strings.TrimSpace(chapter.Title) == "" {
		return "", fmt.Errorf("章节标题为空")
	}
	return c.text(ctx, chapter.URL, strings.TrimSpace(chapter.Title))
}

func (c *Client) text(ctx context.Context, link, expectedTitle string) (string, error) {
	id, chapterID, err := c.parseLink(link)
	if err != nil || chapterID <= 0 {
		return "", fmt.Errorf("无效的章节地址 %q", link)
	}
	query, err := tokenQuery(struct {
		ID        ID `json:"id"`
		ChapterID ID `json:"chapterid"`
	}{id, chapterID})
	if err != nil {
		return "", err
	}
	var result struct {
		ID        ID     `json:"id"`
		ChapterID ID     `json:"chapterid"`
		Title     string `json:"chaptername"`
		Text      string `json:"txt"`
	}
	if err := c.get(ctx, "chapter", query, &result); err != nil {
		return "", err
	}
	text := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(result.Text, "\r\n", "\n"), "\r", "\n"))
	if result.ID != id || result.ChapterID != chapterID || strings.TrimSpace(result.Title) == "" || text == "" {
		return "", fmt.Errorf("书籍 %d 第 %d 章返回空正文或编号不匹配", id, chapterID)
	}
	if expectedTitle != "" && strings.TrimSpace(result.Title) != expectedTitle {
		return "", fmt.Errorf("书籍 %d 第 %d 章标题已变化，目录为 %q，正文为 %q，请刷新目录", id, chapterID, expectedTitle, result.Title)
	}
	return text, nil
}

func (c *Client) coverURL(id ID) string {
	return fmt.Sprintf("%s/bookimg/%d/%d.jpg", c.site, id/1000, id)
}

func (c *Client) bookURL(id ID) string {
	return fmt.Sprintf("%s/#/book/%d/", c.site, id)
}

func (c *Client) parseLink(link string) (ID, ID, error) {
	u, err := url.Parse(strings.TrimSpace(link))
	if err != nil || u == nil || u.User != nil || u.Scheme+"://"+u.Host != c.site || u.Path != "/" || u.RawQuery != "" || u.ForceQuery {
		return 0, 0, fmt.Errorf("无效地址")
	}
	parts := strings.Split(strings.Trim(u.Fragment, "/"), "/")
	if len(parts) < 2 || len(parts) > 3 || parts[0] != "book" {
		return 0, 0, fmt.Errorf("无效书籍路径")
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || id <= 0 {
		return 0, 0, fmt.Errorf("无效书籍编号")
	}
	var chapterID int64
	if len(parts) == 3 {
		if !strings.HasSuffix(parts[2], ".html") {
			return 0, 0, fmt.Errorf("无效章节路径")
		}
		chapterID, err = strconv.ParseInt(strings.TrimSuffix(parts[2], ".html"), 10, 64)
		if err != nil || chapterID <= 0 {
			return 0, 0, fmt.Errorf("无效章节编号")
		}
	}
	return ID(id), ID(chapterID), nil
}

// tokenQuery 复现网页公开的 enaes：MD5 十六进制字符串的前/后 16 字节
// 分别作为 IV/key，使用 AES-CBC、PKCS7 和标准 Base64。此 token 不涉及用户登录。
func tokenQuery(value any) (url.Values, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	digest := md5.Sum([]byte("book@token.html"))
	code := hex.EncodeToString(digest[:])
	block, err := aes.NewCipher([]byte(code[16:]))
	if err != nil {
		return nil, err
	}
	padding := aes.BlockSize - len(data)%aes.BlockSize
	data = append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
	cipher.NewCBCEncrypter(block, []byte(code[:16])).CryptBlocks(data, data)
	return url.Values{"token": {base64.StdEncoding.EncodeToString(data)}}, nil
}

func (c *Client) get(ctx context.Context, endpoint string, query url.Values, output any) error {
	for attempt := 0; ; attempt++ {
		if err := c.wait(ctx); err != nil {
			return err
		}
		err := c.getOnce(ctx, endpoint, query, output)
		if err == nil || ctx.Err() != nil || attempt >= 2 {
			return err
		}
		delay := time.Duration(1<<attempt) * 500 * time.Millisecond
		var statusErr *httpError
		var networkErr net.Error
		if errors.As(err, &statusErr) {
			if statusErr.status != http.StatusTooManyRequests && statusErr.status != http.StatusInternalServerError && statusErr.status != http.StatusBadGateway && statusErr.status != http.StatusServiceUnavailable && statusErr.status != http.StatusGatewayTimeout {
				return err
			}
			if statusErr.retryAfter > delay {
				delay = statusErr.retryAfter
			}
		} else if !errors.As(err, &networkErr) && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
			return err
		}
		if err := waitContext(ctx, delay); err != nil {
			return err
		}
	}
}

// wait 在等待限速时释放锁，让其他等待中的请求也能及时响应取消。
func (c *Client) wait(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		c.mu.Lock()
		delay := time.Until(c.next)
		if delay <= 0 {
			c.next = time.Now().Add(c.interval)
			c.mu.Unlock()
			return nil
		}
		c.mu.Unlock()
		if err := waitContext(ctx, delay); err != nil {
			return err
		}
	}
}

func waitContext(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type httpError struct {
	endpoint   string
	status     int
	retryAfter time.Duration
}

func (e *httpError) Error() string {
	return fmt.Sprintf("请求 %s: HTTP %d", e.endpoint, e.status)
}

func (c *Client) getOnce(ctx context.Context, endpoint string, query url.Values, output any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.api+"/api/"+endpoint+"?"+query.Encode(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Referer", c.site+"/")
	req.Header.Set("Origin", c.site)
	response, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("请求 %s 失败: %w", endpoint, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		statusErr := &httpError{endpoint: endpoint, status: response.StatusCode}
		if value := response.Header.Get("Retry-After"); value != "" {
			if seconds, err := strconv.ParseInt(value, 10, 32); err == nil && seconds > 0 {
				statusErr.retryAfter = time.Duration(seconds) * time.Second
			} else if until, err := http.ParseTime(value); err == nil {
				statusErr.retryAfter = time.Until(until)
			}
		}
		// 超过等待预算时交给下次运行重试，避免提前重发请求。
		if statusErr.retryAfter > 30*time.Second {
			return fmt.Errorf("接口要求等待 %s: %s", statusErr.retryAfter, statusErr)
		}
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 4<<10))
		return statusErr
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, maxResponseSize+1))
	if err != nil {
		return fmt.Errorf("读取 %s 响应失败: %w", endpoint, err)
	}
	if len(data) > maxResponseSize {
		return fmt.Errorf("%s 响应超过 16 MiB", endpoint)
	}
	if err := json.Unmarshal(data, output); err != nil {
		return fmt.Errorf("%s 接口未返回有效 JSON: %w", endpoint, err)
	}
	return nil
}
