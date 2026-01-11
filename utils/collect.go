package utils

import (
	"encoding/json"
	"fmt"
	"go-novel/app/models"
	"regexp"
	"strconv"
	"strings"
)

//获取当前的列表信息
func GetListSourceURL(listPageUrls string) (urls []string, err error) {
	listPages := []models.ListPageReg{}
	err = json.Unmarshal([]byte(listPageUrls), &listPages)
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
		return
	}

	for _, value := range listPages {
		switch value.Type {
		case 1:
			pageStart := value.PageStart
			pageEnd := value.PageEnd
			pageInc := value.PageInc
			pageDesc := value.PageDesc

			for i := pageStart; i <= pageEnd; i += pageInc {
				urls = append(urls, strings.ReplaceAll(value.Url, "[page]", strconv.Itoa(i)))
			}

			if pageDesc {
				urls = reverseSlice(urls)
			}
		default:
			urls = append(urls, value.Url)
		}
	}

	return
}

func reverseSlice(urls []string) []string {
	length := len(urls)
	reversed := make([]string, length)
	for i, url := range urls {
		reversed[length-i-1] = url
	}
	return reversed
}

//url完成的配置信息
func UrlComplete(html, baseURL string) string {
	hrefRegex := regexp.MustCompile(`(?i)\bhref=[\'\"]([^\'\"]*)[\'\"]`)
	srcRegex := regexp.MustCompile(`(?i)\bsrc=[\'\"]([^\'\"]*)[\'\"]`)

	html = hrefRegex.ReplaceAllStringFunc(html, func(match string) string {
		hrefValue := extractURL(match)
		return fmt.Sprintf(`href="%s"`, createURL(hrefValue, baseURL))
	})

	html = srcRegex.ReplaceAllStringFunc(html, func(match string) string {
		srcValue := extractURL(match)
		return fmt.Sprintf(`src="%s"`, createURL(srcValue, baseURL))
	})

	return html
}

func extractURL(tag string) string {
	urlRegex := regexp.MustCompile(`(?i)[\'\"]([^\'\"]*)[\'\"]`)
	matches := urlRegex.FindStringSubmatch(tag)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func createURL(url, baseURL string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}

	if strings.HasPrefix(url, "/") {
		baseURLParts := strings.Split(baseURL, "/")
		baseURL = strings.Join(baseURLParts[:3], "/")
	}

	return fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(url, "/"))
}
