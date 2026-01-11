package middleware

import (
	"github.com/gocolly/colly"
	"go-novel/utils"
	"strings"
	"time"
)

// RetryMiddleware 定义一个自定义中间件
type RetryMiddleware struct {
	Retries  int64
	Interval time.Duration
}

// Apply 实现中间件接口
func (m *RetryMiddleware) Apply(c *colly.Collector) {
	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("retries", 0)
	})

	c.OnError(func(r *colly.Response, err error) {
		if err != nil {
			if strings.Contains(err.Error(), "Timeout") {
				retries := r.Ctx.Get("retries")
				retriesInt := utils.FormatInt64(retries)
				if retriesInt < m.Retries {
					r.Ctx.Put("retries", retriesInt+1)
					c.Visit(r.Request.URL.String())
					time.Sleep(m.Interval)
				}
			}
		}
	})
}
