package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go-novel/global"
	"io/ioutil"
	"net/http"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyBytes []byte // 我们需要的body内容
		//// 从原有Request.Body读取
		var err error
		bodyBytes, err = ioutil.ReadAll(c.Request.Body)
		if err != nil {
			global.Paylog.Infof("read request body failed,err =%s", err)
			return
		}
		// 新建缓冲区并替换原有Request.body
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		// PUT或者POST方法时打印body
		if c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPost {
			data := string(bodyBytes)

			//var str bytes.Buffer
			//_ = json.Indent(&str, []byte(data), "", "    ")
			//fmt.Println("formated: ", str.String())
			//global.Jsonlog.Infof("接受数据iuser=%s", str.String())
			global.Paylog.Infof("接受数据 %s", data)

		}
		// 处理请求
		c.Next()
	}
}
