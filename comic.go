package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-novel/app/models"
	"go-novel/db"
	"go-novel/global"
	"go-novel/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	db.InitZapLog()
	for {
		test()
	}
}

func test() {
	var comicList []*models.Comic
	var err error
	//list := redis_service.Get(fmt.Sprintf("%v", utils.ComicList))
	//if list != "" {
	//	err = json.Unmarshal([]byte(list), &comicList)
	//	if err != nil {
	//		global.Collectlog.Errorf("comicList 解析采集信息失败 err=%v", err.Error())
	//		return
	//	}
	//}
	if len(comicList) <= 0 {
		bookUrl := "https://www.bcloudmerge.com/bmergelists/9/全部/3/2.html"
		var html string
		html, err = utils.GetHtml(bookUrl, "utf-8", 1, 0)
		if html == "" {
			utils.GetS5()
		}
		// 使用 goquery 解析 HTML 字符串
		var doc *goquery.Document
		doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			log.Println(err)
			return
		}
		// 使用选择器定位目标元素并提取数据
		doc.Find(".cy_list_mh ul").Each(func(index int, ulElement *goquery.Selection) {
			// 提取每个 <ul> 元素中的信息
			var pic, comicName, comicHref, author, desc string

			titleElement := ulElement.Find("li.title a")
			comicName = titleElement.Text()
			comicHref, _ = titleElement.Attr("href")
			lastNumber := utils.GetLastNumber(comicHref)
			//fmt.Println(lastNumber) // 输出: 7111
			comicHref = fmt.Sprintf("https://www.bcloudmerge.com/menu/%v.html", lastNumber)
			//fmt.Println("标题:", comicName)
			//fmt.Println("链接地址:", comicHref)

			authorElement := ulElement.Find("li.zuozhe")
			author = strings.TrimPrefix(authorElement.Text(), "作者：")
			//fmt.Println("作者:", author)

			introElement := ulElement.Find("li.info")
			desc = strings.TrimPrefix(introElement.Text(), "简介：")
			//fmt.Println("简介:", desc)

			imageElement := ulElement.Find("li a.pic img")
			pic, _ = imageElement.Attr("src")
			//fmt.Println("图片地址:", pic)

			comic := &models.Comic{
				ComicName: comicName,
				ComicHref: comicHref,
				Author:    author,
				Pic:       pic,
				Desc:      desc,
			}
			comicList = append(comicList, comic)
		})

		//var jsonData []byte
		//jsonData, err = json.Marshal(comicList)
		//if err != nil {
		//	err = fmt.Errorf("转换json数据失败: %v", err.Error())
		//	return
		//}
		//err = redis_service.Set(fmt.Sprintf("%v", utils.ComicList), jsonData, 0)
		//if err != nil {
		//	global.Collectlog.Errorf("缓存采集当前漫画列表失败 %v", err.Error())
		//	return
		//}
	}

	var uploadBookPicPath string = "/mnt/comic/"
	for _, comic := range comicList {
		err = SaveDesc(comic.ComicName, comic.Author, comic.Desc, comic.ComicHref, comic.Pic)
		if err != nil {
			global.Errlog.Errorf("%v", err.Error())
			return
		}
		var filePath string
		uploadBookPicPath = fmt.Sprintf("%v%v/", uploadBookPicPath, comic.ComicName)
		filePath, err = DownImg(comic.ComicName, comic.Author, comic.Pic, uploadBookPicPath)
		if err != nil {
			global.Errlog.Errorf("下载pic图片失败%v %v %v", comic.ComicName, filePath, err.Error())
			utils.GetS5()
		} else {
			global.Errlog.Errorf("下载pic图片成功 %v %v", comic.ComicName, filePath)
		}
		var html string
		html, err = utils.GetHtml(comic.ComicHref, "utf-8", 1, 0)
		if html == "" {
			utils.GetS5()
		}
		if err != nil {
			global.Errlog.Errorf("%v", err.Error())
			return
		}
		var doc *goquery.Document
		doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			log.Println(err)
			return
		}

		var chapters []models.ComicChapter
		doc.Find(".cy_plist").Eq(1).Find("ul li a").Each(func(i int, s *goquery.Selection) {
			chapterHref, _ := s.Attr("href")
			chapterName := s.Text()
			chapter := models.ComicChapter{
				ChapterName: chapterName,
				ChapterHref: chapterHref,
			}
			chapters = append(chapters, chapter)
		})
		if len(chapters) > 0 {
			for _, chapter := range chapters {
				html, err = utils.GetHtml(chapter.ChapterHref, "utf-8", 1, 0)
				if html == "" {
					utils.GetS5()
				}
				if err != nil {
					global.Errlog.Errorf("%v", err.Error())
					return
				}
				doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
				if err != nil {
					global.Errlog.Errorf("获取章节链接出错 %v", err.Error())
					return
				}
				var imgs []string
				doc.Find(".mh_list a img").Each(func(i int, s *goquery.Selection) {
					src, _ := s.Attr("src")
					imgs = append(imgs, src)
				})
				for index, img := range imgs {
					uploadBookPicPath = fmt.Sprintf("/mnt/comic/%v/%v/", comic.ComicName, chapter.ChapterName)
					filePath, err = DownImg(fmt.Sprintf("%v", index), "", img, uploadBookPicPath)
					if err != nil {
						log.Println("下载图片失败", comic.ComicName, comic.ComicHref, uploadBookPicPath, filePath, err.Error())
						utils.GetS5()
					} else {
						log.Println("下载图片成功", comic.ComicName, comic.ComicHref, uploadBookPicPath, filePath)

					}
				}
			}

		}
		//log.Println(comic.ComicHref)
	}
}

func DownImg(name, author, picUrl, pathDir string) (filePath string, err error) {
	// 检查URL是否合法
	if !strings.HasPrefix(picUrl, "http://") && !strings.HasPrefix(picUrl, "https://") {
		err = fmt.Errorf("图片链接错误 picUrl=%v", picUrl)
		return
	}

	if name == "" {
		name = utils.GetFileName(picUrl)
	}

	// 获取文件扩展名
	fileExt := strings.ToLower(utils.GetExt(picUrl))
	fileName := fmt.Sprintf("%v.%v", name, fileExt)
	if author != "" {
		fileName = fmt.Sprintf("%v-%v.%v", name, author, fileExt)
	}

	uploadPath := fmt.Sprintf("%s", pathDir)

	err = utils.IsNotExistMkDir(uploadPath)
	if err != nil {
		err = fmt.Errorf("%v", "创建目录失败")
		return
	}

	// 文件完整路径
	filePath = fmt.Sprintf("%s%s", uploadPath, fileName)
	if utils.FileExist(filePath) {
		global.Collectlog.Errorf("%v 文件已存在", filePath)
		return
	}

	httpClient := &http.Client{}

	if utils.IsS5 {
		//httpTransport := utils.GetHttpTransport()
		//log.Println(utils.S5Domain, utils.S5Port, utils.S5Username, utils.S5Username)
		//httpClient = &http.Client{Transport: httpTransport}
	}
	// 发送HTTP请求获取文件内容
	resp, err := httpClient.Get(picUrl)
	if err != nil {
		err = fmt.Errorf("%v err=%v获取图片链接失败", picUrl, err.Error())
		utils.GetS5()
		return
	}
	defer resp.Body.Close()

	// 读取文件内容
	fileContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("%v 读取图片失败", picUrl)
		utils.GetS5()
		return
	}
	// 写入文件
	err = ioutil.WriteFile(filePath, fileContent, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("%v 保存图片失败", picUrl)
		return
	}
	return
}

func SaveDesc(comicName, author, desc, href, pic string) (err error) {
	text := ""
	text += fmt.Sprintf("漫画名称:%v\n", comicName)
	text += fmt.Sprintf("漫画作者:%v\n", author)
	text += fmt.Sprintf("漫画简介:%v\n", desc)
	text += fmt.Sprintf("漫画链接:%v\n", href)
	text += fmt.Sprintf("漫画图片:%v\n", pic)
	err = utils.WriteFile(fmt.Sprintf("/mnt/comic/%v/%v.txt", comicName, comicName), text)
	return
}
