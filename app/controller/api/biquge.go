package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"strings"
)

// 笔趣阁对外的关联的解密程序单独用接口进行访问配置
type Biquge struct{}

func (biquge *Biquge) GeneralDecrypt(c *gin.Context) {
	var req utils.ReqGeneralDecrypt
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	content := req.Content
	content = strings.TrimSpace(content)
	//对相关数据进行解码操作
	decoded, err := utils.ComicApiDecryptV1(content)
	if err != nil {
		utils.FailEncrypt(c, err, "解码失败")
		return
	}
	log.Printf("当前的解码对应的值为：%v\n", decoded)
	res := gin.H{
		"detail": decoded,
	}
	utils.SuccessEncrypt(c, res, "ok")
}

// 笔趣阁解锁章节内容
func (biquge *Biquge) ContentDecrypt(c *gin.Context) {
	var req utils.ReqBqgContentDecrypt
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	content := strings.TrimSpace(req.Content)
	//处理解密字符串中的base64最后包含的=或者==有问题的数据
	newContent := utils.RemoveEqualSigns(content)
	//对内容单独进行解码
	origString, err := utils.PswDecrypt(newContent)
	if err != nil {
		utils.FailEncrypt(c, err, "章节内容解码失败")
		return
	}
	res := gin.H{
		"content": origString,
	}
	utils.SuccessEncrypt(c, res, "章节内容")
	return
}

// 关联章节内容的解密
func (biquge *Biquge) ChapterDecrypt(c *gin.Context) {
	var req utils.ReqChapterDecrypt
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	//获取接收的对象参数信息
	domainUrl := strings.TrimSpace(req.DomainUrl)
	fmt.Println(domainUrl)
	//异步请求的章节目录信息
	//封面的资源url：res.ycukhv.com
	//catalog.lmmwlkj.com //使用这个域名，这个域名数据比较全面上面的坑爹呀，好多章节都是空的
	//url := "https://catalog.ycukhv.com/"
	path := strings.TrimSpace(req.Path)

	url1 := "https://catalog.lmmwlkj.com/"    //系统默认的找的另外一个域名
	url2 := "https://catalog.ycukhv.com/"     //系统默认的的APP里的一个域名
	url3 := "https://chapter.jhkhmgj.com/"    //系统默认找到的的备用一个域名
	apiUrl := fmt.Sprintf("%v%v", url1, path) //拼接请求的解码路径
	jsonData := utils.GetBaiduResponse(apiUrl)
	if jsonData == "" {
		//这个地方用CDN里的数据进行重新请求一次
		apiUrl = fmt.Sprintf("%v%v", url2, path) //获取新的重新请求一次
		jsonData = utils.GetBaiduResponse(apiUrl)
		if jsonData == "" {
			apiUrl = fmt.Sprintf("%v%v", url3, path) //用url3的域名的再去获取一次
			jsonData = utils.GetBaiduResponse(apiUrl)
			if jsonData == "" {
				utils.FailEncrypt(c, nil, "获取章节目录为空")
				return
			}
		}
	}
	//定义解析的章节的结构体
	var response utils.BiqugeChapterListItem
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		global.Errlog.Infof("Error parsing JSON results = %+v", err)
		utils.FailEncrypt(c, err, "转换json失败")
		return
	}
	if response.Code == 1 {
		// 创建新对象数组
		var newDataItems []utils.BiqugeChapterDataItem
		for _, val := range response.Data {
			//使用替换后的值
			replaceStr := utils.RemoveEqualSigns(strings.TrimSpace(val.Name))
			var chapterName string
			//解密相关的数据操作
			chapterName, err = utils.PswDecrypt(replaceStr)
			if err != nil {
				log.Fatalf("Error decrypt: %s", err)
			}
			//追加数组对象
			newDataItem := utils.BiqugeChapterDataItem{
				Name:      chapterName, // 新的解密后的章节名称
				URL:       val.URL,
				IsContent: val.IsContent,
				Path:      fmt.Sprintf("%v%v", domainUrl, val.Path), // 路径的拼接规则： xxx域名（POST过来） + path路径信息
				UpdatedAt: val.UpdatedAt,
			}
			newDataItems = append(newDataItems, newDataItem)
		}
		//log.Println(newDataItems)
		res := gin.H{
			"list": newDataItems,
		}
		log.Printf("请求的章节目录的url=%v\n", apiUrl)
		utils.SuccessEncrypt(c, res, "ok")
		return
	} else {
		utils.FailEncrypt(c, err, "获取章节目录失败")
		return
	}

}
