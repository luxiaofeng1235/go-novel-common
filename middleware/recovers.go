package middleware

import (
	"fmt"
	"go-novel/utils/zaplog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"runtime"
)

func RecoverWithZincSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 跳过记录标志（根据需要设置）
		skipLogging := false

		// 请求参数
		params := ""
		if c.Request.Method == http.MethodGet {
			params = c.Request.URL.RawQuery
		} else {
			// 解析表单参数（POST、PUT 等）
			if err := c.Request.ParseForm(); err == nil {
				params = c.Request.Form.Encode()
			}
		}

		// 延迟处理函数，用于捕获 panic
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, false)]

				// 格式化错误信息
				errorMessage := fmt.Sprintf("Panic recovered: %v\nStack trace: %s", err, stack)

				// 获取请求相关信息
				requestMethod := c.Request.Method
				requestURI := c.Request.RequestURI
				clientIP := c.ClientIP()
				userAgent := c.Request.UserAgent()
				queryParams := c.Request.URL.RawQuery
				postParams := params

				// 格式化并打印错误信息到控制台（红色）
				fmt.Printf("\033[1;31;40m%s\033[0m\n", "系统错误: "+fmt.Sprint(err))
				fmt.Printf("\033[1;31;40mmethod: %s\033[0m\n", requestMethod)
				fmt.Printf("\033[1;31;40muri: %s\033[0m\n", requestURI)
				fmt.Printf("\033[1;31;40mclient_ip: %s\033[0m\n", clientIP)
				fmt.Printf("\033[1;31;40muser_agent: %s\033[0m\n", userAgent)
				fmt.Printf("\033[1;31;40mquery_params: %s\033[0m\n", queryParams)
				fmt.Printf("\033[1;31;40mpost_params: %s\033[0m\n", postParams)
				fmt.Printf("\033[1;31;40merror: %s\033[0m\n", errorMessage)

				// 异步发送日志到 ZincSearch
				go func() {
					// 计算响应时间
					responseTime := time.Since(startTime).Milliseconds()

					// 检查是否需要跳过记录
					if !skipLogging {
						logData := map[string]interface{}{
							"status":        c.Writer.Status(),
							"ip":            clientIP,
							"method":        requestMethod,
							"path":          requestURI,
							"response_time": responseTime,
							"params":        params,
							"timestamp":     time.Now().Format(time.RFC3339),
							"error":         fmt.Sprint(err),
							"stack":         string(stack),
						}
						zaplog.SendLogToZincSearch(logData, "error")
					}
				}()

				// 返回自定义的错误响应
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error",
				})
			}
		}()

		// 继续处理请求
		c.Next()

		// 在请求完成后，可以记录响应信息（可选）
		// 例如，如果需要记录响应体，可以在这里获取
		// 这里假设 response 已经被记录
	}
}
