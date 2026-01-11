package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-novel/global"
	"go-novel/utils/zaplog"
	"io"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinBodyLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 获取所有header参数
		headers := c.Request.Header
		headerMap := map[string]interface{}{
			"mark":          headers.Get("mark"),
			"PACKAGE":       headers.Get("PACKAGE"),
			"imei":          headers.Get("imei"),
			"androidid":     headers.Get("androidid"),
			"oaid":          headers.Get("oaid"),
			"ua":            headers.Get("ua"),
			"OS":            headers.Get("OS"),
			"uuid":          headers.Get("uuid"),
			"model":         headers.Get("model"),
			"BRAND":         headers.Get("BRAND"),
			"Authorization": headers.Get("Authorization"),
			"CLIENTVERSION": headers.Get("CLIENTVERSION"),
		}
		c.Set("headerMap", headerMap)
		// 获取body内容
		body := c.Request.Body
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			global.Errlog.Error("request get body error:", err.Error())
			c.AbortWithStatusJSON(500, gin.H{"error": "unable to read request body"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重置body以供后续使用
		c.Next()
		// 解析body中的JSON内容
		bodyMap := map[string]interface{}{}
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &bodyMap); err != nil {
				global.Errlog.Error("unmarshal body error:", err.Error())
			}
		}

		// 获取query和form参数
		query := c.Request.URL.Query()
		queryMap := make(map[string]interface{})
		for k, v := range query {
			queryMap[k] = v[0]
		}

		form := c.Request.PostForm
		formMap := make(map[string]interface{})
		for k, v := range form {
			formMap[k] = v[0]
		}

		// 合并所有参数
		params := map[string]interface{}{
			"headers": headerMap,
			"query":   queryMap,
			"form":    formMap,
			"body":    bodyMap,
		}
		responseTime := time.Since(startTime)

		info := logrus.Fields{
			"status": c.Writer.Status(),
			"ip":     c.ClientIP(),
			"method": c.Request.Method,
			//"response": string(blw.body.Bytes()),
			"response_time": responseTime.String(), // Add response time
		}
		for k, v := range params {
			info[k] = v
		}
		//转为json
		//jsonInfo, err := json.Marshal(info)
		//过滤一些特殊接口 book/chapter 返回内容太多不需要记载
		whitelist := []string{
			"/api/book/read",
			"/api/adver/getAdverMapList",
			"/api/vip/vipBookStore",
			"/api/book/chapter",
			"/api/book/todayUpdateBooks",
			"/api/book/getSectionForYouRec",
			"/api/class/list",
			"/api/book/rankList",
			"/api/book/getSectionEnd",
		}
		skipLogging := false
		for _, path := range whitelist {
			if c.Request.URL.Path == path {
				skipLogging = true
				break
			}
		}
		response := string(blw.body.Bytes())
		//if global.Requestlog != nil {
		//	if skipLogging {
		//
		//		global.Requestlog.Info("path:", c.Request.URL.String(), "  request:", string(jsonInfo), "  response:太多了不记录")
		//	} else {
		//		global.Requestlog.Info("path:", c.Request.URL.String(), "  request:", string(jsonInfo), "  response:", response)
		//	}
		//}

		go func() {
			// 将日志发送到 ZincSearch
			responseTime := time.Since(startTime).Milliseconds() // 转换为整数
			if !skipLogging {
				logData := map[string]interface{}{
					"status":        c.Writer.Status(),
					"ip":            c.ClientIP(),
					"method":        c.Request.Method,
					"path":          c.Request.URL.String(),
					"response_time": responseTime,
					"params":        params,
					"timestamp":     time.Now().Format(time.DateTime),
					"response":      response,
				}

				zaplog.SendLogToZincSearch(logData, "request")
			}
		}()
		//处理请求
		c.Next()
	}
}

func MergeMap(mObj ...map[string]interface{}) map[string]interface{} {
	newObj := map[string]interface{}{}
	for _, m := range mObj {
		for k, v := range m {
			newObj[k] = v
		}
	}
	return newObj
}
