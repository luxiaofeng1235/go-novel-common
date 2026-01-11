package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/utils"
	"io/ioutil"
	"log"
	"net/http"
)

func ApiAes1() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取原始请求 Body 数据
		var err error
		//bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		//if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "读取请求 Body 失败"})
		//	c.Abort()
		//	return
		//}
		//log.Println("bodyBytes", string(bodyBytes))
		//
		//// 恢复原始请求 Body
		//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// 解析原始 JSON 数据
		//var req models.ApiAesReq
		//if err = json.Unmarshal(bodyBytes, &req); err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": "参数绑定失败"})
		//	c.Abort()
		//	return
		//}

		// 参数绑定
		//if err := c.ShouldBind(&req); err != nil {
		//	utils.Fail(c, err, "参数绑定失败")
		//	return
		//}
		//log.Println("req", req.Data)

		var req models.ApiAesReq
		// 参数绑定
		if err := c.ShouldBind(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}

		var apiJson string
		apiJson, err = utils.AesDecryptByCFB(utils.ApiAesKey, req.Data)
		if err != nil {
			utils.Fail(c, err, "参数解密失败")
			return
		}
		log.Println("req", apiJson)

		//newJSONData1, err := json.Marshal(apiJson)
		//if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON 序列化失败"})
		//	c.Abort()
		//	return
		//}
		//
		//aa := models.TodayUpdateBooksReq{
		//	Page: 100,
		//	Size: 200,
		//}

		// 将修改后的 JSON 数据重新序列化
		//newJSONData, err := json.Marshal(aa)
		//if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON 序列化失败"})
		//	c.Abort()
		//	return
		//}
		//log.Println("newJSONData1", apiJson)
		//log.Println("newJSONData", string(newJSONData))
		//bb := "{\n    //页数\n    \"page\": 11,\n    //每页条数\n    \"size\": 101\n}"

		// 定义 JSON 字符串
		//jsonStr := `{
		//		//页数
		//		"page": 11,
		//		//每页条数
		//		"size": 101
		//	}`
		c.Request.Body = ioutil.NopCloser(bytes.NewReader([]byte(apiJson)))
		c.Request.ContentLength = int64(len(apiJson))
		//var jsonData models.TodayUpdateBooksReq
		//// 参数绑定
		//if err = c.ShouldBind(&jsonData); err != nil {
		//	log.Println("err", err.Error())
		//	utils.Fail(c, err, "参数绑定失败")
		//	return
		//}
		//log.Println("jsonData", jsonData)
		// 将字符串解析为 JSON 对象
		//var jsonData models.TodayUpdateBooksReq
		//err = json.Unmarshal([]byte(jsonStr), &jsonData)
		//if err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": "解析 JSON 失败"})
		//	c.Abort()
		//	return
		//}

		// 将 JSON 对象重新转换为字符串
		//newJSONStr, err := json.Marshal(jsonData)
		//if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON 序列化失败"})
		//	c.Abort()
		//	return
		//}
		// 更新请求的 Body
		//c.Request.Body = ioutil.NopCloser(bytes.NewReader([]byte(newJSONStr)))
		//c.Request.ContentLength = int64(len(newJSONStr))

		// 重置请求体为修改后的 JSON 数据
		//c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(len(newJSONData)))
		//c.Request.ContentLength = int64(len(newJSONData))

		//var bodyBytes []byte // 我们需要的body内容
		////// 从原有Request.Body读取
		//bodyBytes, err = ioutil.ReadAll(c.Request.Body)
		//if err != nil {
		//	global.Paylog.Infof("read request body failed,err =%s", err)
		//	return
		//}

		// 新建缓冲区并替换原有Request.body
		//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(newJSONData))
		// 更新请求的 Body
		//c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(len(newJSONData)))
		//c.Request.ContentLength = int64(len(newJSONData))
		// PUT或者POST方法时打印body
		if c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPost {
			//data := string(bodyBytes)

			//var str bytes.Buffer
			//_ = json.Indent(&str, []byte(data), "", "    ")
			//fmt.Println("formated: ", str.String())
			//global.Jsonlog.Infof("接受数据iuser=%s", str.String())
			//global.Paylog.Infof("接受数据 %s", newJSONData)

		}
		// 处理请求
		c.Next()
	}
}
