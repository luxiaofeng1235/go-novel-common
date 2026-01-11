package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mozillazg/go-pinyin"
	"go-novel/app/models"
	"go-novel/db"
	"go-novel/global"
	"go-novel/utils"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	db.InitZapLog()
	host, name, user, passwd := db.GetDB()
	db.InitDB(host, name, user, passwd)
	//log.Println(UpLoadFile("./test.jpg"))

	//var list []models.McComic
	//global.DB.Model(models.McComic{}).Where("id > 27").Find(&list)
	//for _, val := range list {
	//	pic, err := UpLoadFile(val.Pic)
	//	if err != nil {
	//		log.Println("失败", err.Error())
	//		return
	//	}
	//	data := make(map[string]interface{})
	//	data["pic"] = pic
	//	data["picx"] = pic
	//	data["cid"] = 1
	//	global.DB.Model(models.McComic{}).Where("id = ?", val.Id).Debug().Updates(data)
	//}
	//return

	//var list []models.McComicChapter
	//global.DB.Model(models.McComicChapter{}).Where("image like ?", "%E:%").Find(&list)
	//for _, val := range list {
	//	pic, err := UpLoadFile(val.Image)
	//	if err == nil {
	//		global.DB.Model(models.McComicChapter{}).Where("id = ?", val.Id).Debug().Update("image", pic)
	//	} else {
	//		log.Println("失败", err.Error())
	//	}
	//	//return
	//}
	//return

	var list []models.McComicPic
	//global.DB.Model(models.McComicPic{}).Where("id > 322").Find(&list)
	//for _, val := range list {
	//	pic, err := UpLoadFile(val.Img)
	//	if err == nil {
	//		data := make(map[string]interface{})
	//		data["img"] = pic
	//		data["md5"] = utils.Md5(pic)
	//		global.DB.Model(models.McComicPic{}).Where("id = ?", val.Id).Debug().Updates(data)
	//	} else {
	//		log.Println("失败", err.Error())
	//	}
	//}

	global.DB.Model(models.McComicPic{}).Order("id desc").Where("img like ?", "%E:%").Limit(250).Find(&list)
	for _, val := range list {
		pic, err := UpLoadFile(val.Img)
		if err == nil {
			data := make(map[string]interface{})
			data["img"] = pic
			data["md5"] = utils.Md5(pic)
			global.DB.Model(models.McComicPic{}).Where("id = ?", val.Id).Debug().Updates(data)
		} else {
			log.Println("失败", err.Error())
		}
	}
	return
	var start = time.Now()
	threadCount := 250
	listCount := len(list)
	if listCount > 0 {
		//threadsPerThread := chapterCount / threadCount
		threadsPerThread := int(float64(listCount)/float64(threadCount) + 0.5)
		var wg sync.WaitGroup
		wg.Add(threadCount)
		for i := 0; i < threadCount; i++ {
			startIndex := i * threadsPerThread
			endIndex := (i + 1) * threadsPerThread
			// 最后一个线程处理剩余的章节
			if i == threadCount-1 {
				endIndex = len(list)
			}
			tempCount := listCount - 1
			if endIndex > tempCount {
				startIndex = tempCount
				endIndex = tempCount
			}

			global.Collectlog.Errorf("正在采集 线程=%v 目标数量=%v 当前线程数量=%v startIndex=%v endIndex=%v", i, listCount, len(list[startIndex:endIndex]), startIndex, endIndex)
			go processChapterDownPic(list[startIndex:endIndex], &wg)
		}
		wg.Wait()
	}
	log.Fatalln(time.Since(start))
	return

	// 解密图片
	//xorNum := utils.ImgEncry
	//err := utils.DecodeImage("./xfile.jpg", "./test.png", xorNum)
	//if err != nil {
	//	log.Println(err.Error())
	//	return
	//}
	//var list []models.McComicChapter
	//global.DB.Model(models.McComicChapter{}).Where("image LIKE ?", "%-%").Find(&list)
	//
	//log.Println(list, len(list))
	//for _, chapter := range list {
	//	err := DeleteImageFile(chapter.Image)
	//	if err != nil {
	//		// 处理删除过程中的错误
	//		fmt.Printf("删除图片文件时出错: %v\n", err)
	//	} else {
	//		// 删除成功
	//		fmt.Println("图片文件已成功删除", chapter.Image)
	//	}
	//}
	//return
	rootDir := "E:\\mnt\\comic\\test"

	//imgs, err := utils.ScanDirPic(rootDir)
	//if err != nil {
	//	fmt.Printf("遍历目录时出错: %v\n", err)
	//}
	//imgs = imgs[:len(imgs)-1]
	//aa, err := json.Marshal(imgs)
	//log.Println(string(aa), err)
	//return
	rootDir = "E:\\mnt\\comic"
	baseDir, err := utils.ScanBaseDir(rootDir)
	if err != nil {
		return
	}
	if len(baseDir) <= 0 {
		return
	}
	//baseDir = []string{"test"}
	for _, comicName := range baseDir {
		descFile := fmt.Sprintf("E:\\mnt\\comic\\%v\\%v.txt", comicName, comicName)
		var conByte []byte
		conByte, err = ioutil.ReadFile(descFile)
		if err != nil {
			log.Println(descFile, err.Error())
			return
		}
		fileContent := string(conByte)
		author := getValueAfterColon(fileContent, "漫画作者")
		desc := getValueAfterColon(fileContent, "漫画简介")
		//pic := getValueAfterColon(fileContent, "漫画图片")
		uploadBookPicPath := fmt.Sprintf("/mnt/comic/%v/%v.jpg", comicName, comicName)
		//var filePath string
		//uploadBookPicPath := fmt.Sprintf("/mnt/comic/%v/", comicName)
		//filePath, err = DownPic(comicName, "", pic, uploadBookPicPath)
		//log.Println(pic, filePath, err)
		hits, hitsDay, hitsWeek, _, shits, score, _, _ := utils.GetRandNumBookHits()

		comic := models.McComic{
			Name:      comicName,
			Author:    author,
			Pic:       uploadBookPicPath,
			Yname:     strings.Join(pinyin.LazyPinyin(comicName, pinyin.NewArgs()), ""),
			Cid:       0,
			Serialize: "完结",
			Text:      desc,
			Content:   desc,
			Hits:      hits,
			Yhits:     hitsDay,
			Zhits:     hitsWeek,
			Rhits:     shits,
			Score:     score,
			Addtime:   utils.GetUnix(),
		}
		if err = global.DB.Create(&comic).Error; err != nil {
			global.Sqllog.Errorf("记录失败，稍后再试 err=%v", err.Error())
			return
		}
		rootDir = fmt.Sprintf("E:\\mnt\\comic\\%v", comicName)
		var chapters []models.DirPics
		chapters, err = utils.ScanDirPic(rootDir)
		if err != nil {
			fmt.Printf("遍历目录时出错: %v\n", err)
		}
		//return
		for index, img := range chapters {
			if strings.Contains(img.DirName, "章") {
				continue
			}
			if strings.TrimSpace(img.DirName) == strings.TrimSpace(comicName) {
				continue
			}
			//if strings.Contains(img.DirPics[0], comicName) {
			//	continue
			//}
			index++
			//log.Println(comicName, img.DirName, index)
			chapter := models.McComicChapter{
				Mid:     comic.Id,
				Xid:     index,
				Image:   img.DirPics[0],
				Name:    img.DirName,
				Pnum:    len(img.DirPics),
				Addtime: utils.GetUnix(),
			}
			if err = global.DB.Create(&chapter).Error; err != nil {
				global.Sqllog.Errorf("记录失败，稍后再试 err=%v", err.Error())
				return
			}
			for i, pic := range img.DirPics {
				comicPic := models.McComicPic{
					Cid: chapter.Id,
					Mid: comic.Id,
					Xid: i,
					Img: pic,
					Md5: utils.Md5(pic),
				}
				if err = global.DB.Create(&comicPic).Error; err != nil {
					global.Sqllog.Errorf("记录失败，稍后再试 err=%v", err.Error())
					return
				}
			}

		}
		//return
	}

}

func UpLoadFile(fileName string) (fileUrl string, err error) {
	// 打开要上传的图片文件
	file, err := os.Open(fileName)
	if err != nil {
		err = fmt.Errorf("无法打开图片文件 %v", err.Error())
		return
	}
	defer file.Close()

	// 创建一个缓冲区来存储表单数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加额外的参数
	writer.WriteField("facility", "windows_h5")
	writer.WriteField("deviceid", "17107462144859084108")
	writer.WriteField("apikey", "123123")
	writer.WriteField("user_id", "254514")
	writer.WriteField("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond)))

	// 创建一个表单字段，用于上传图片
	imageField, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		err = fmt.Errorf("创建表单字段失败:%v", err.Error())
		return
	}

	// 将图片内容复制到表单字段中
	_, err = io.Copy(imageField, file)
	if err != nil {
		err = fmt.Errorf("复制图片内容失败:%v", err.Error())
		return
	}

	// 完成表单写入
	writer.Close()

	// 创建一个POST请求用于上传图片
	request, err := http.NewRequest("POST", "http://103.36.90.182/v3.php/app/kuaikan_upload", body)
	if err != nil {
		err = fmt.Errorf("创建请求失败:%v", err.Error())
		return
	}

	// 设置请求头，指定multipart/form-data类型
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("发送请求失败:%v", err.Error())
		return
	}
	defer response.Body.Close()
	// 处理响应
	if response.StatusCode == http.StatusOK {
		// 读取文件内容
		var bodyBytes []byte
		bodyBytes, _ = ioutil.ReadAll(response.Body)
		var ciphertext models.UpLoadPicCiphertextRes
		if err = json.Unmarshal(bodyBytes, &ciphertext); err != nil {
			err = fmt.Errorf("解析返回参数错误 %v %v", err.Error(), string(bodyBytes))
			return
		}
		if ciphertext.Data == "" {
			err = fmt.Errorf("%v", "返回data为空")
			return
		}
		var plaintext string
		plaintext, _ = utils.AesDecryptByCFB(utils.ApiAesKey, ciphertext.Data)
		if plaintext == "" {
			err = fmt.Errorf("%v", "解密错误解密返回为空")
			return
		}
		var res models.UpLoadPicRes
		if err = json.Unmarshal([]byte(plaintext), &res); err != nil {
			err = fmt.Errorf("解析返回参数错误 %v %v", err.Error(), plaintext)
			return
		}
		if res.Code != 1 {
			err = fmt.Errorf("上传接口报错 err=%v", res.Msg)
			return
		}
		fileUrl = res.Url
		//{"code":1,"url":"upload\/kuaikan\/20240320\/66b156b48c05efe8750a3fcdfa14f54e_xfile.jpg","img":"https:\/\/t.shunfengs.com\/upload\/kuaikan\/20240320\/66b156b48c05efe8750a3fcdfa14f54e_xfile.jpg","msg":"\u56fe\u7247\u4e0a\u4f20\u5b8c\u6210"}
	} else {
		var bodyBytes []byte
		bodyBytes, _ = ioutil.ReadAll(response.Body)
		log.Println(string(bodyBytes))
		err = fmt.Errorf("图片上传失败 500错误 %v", response.StatusCode)
		return
	}
	return
}

func getValueAfterColon(fileContent string, key string) string {
	lines := strings.Split(fileContent, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, key) {
			value := strings.TrimPrefix(line, key+":")
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func DownPic(name, author, picUrl, pathDir string) (filePath string, err error) {
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
		httpTransport := utils.GetHttpTransport()
		log.Println(utils.S5Domain, utils.S5Port, utils.S5Username, utils.S5Username)
		httpClient = &http.Client{Transport: httpTransport}
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

func DeleteImageFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func processChapterDownPic(list []models.McComicPic, wg *sync.WaitGroup) {
	defer wg.Done()
	if len(list) <= 0 {
		return
	}
	for _, val := range list {
		pic, err := UpLoadFile(val.Img)
		if err == nil {
			data := make(map[string]interface{})
			data["img"] = pic
			data["md5"] = utils.Md5(pic)
			global.DB.Model(models.McComicPic{}).Where("id = ?", val.Id).Debug().Updates(data)
		} else {
			log.Println("失败", err.Error())
		}
	}
	return
}
