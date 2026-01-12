package zinc

import (
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-resty/resty/v2"
)

type ZincClient struct {
	ZincHost     string
	ZincUser     string
	ZincPassword string
}

// NewClient returns a new instance of ZincClient
func NewClient(host, user, passwd string) *ZincClient {
	return &ZincClient{
		ZincHost:     host,
		ZincUser:     user,
		ZincPassword: passwd,
	}
}

// SendLog sends log data to ZincSearch
func (c *ZincClient) SendLog(logData map[string]interface{}) error {
	logData["@timestamp"] = time.Now().Format(time.RFC3339)
	resp, err := c.request().SetBody(logData).Post("/api/request/_doc") // 假设 ZincSearch 接口为 /api/logs
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.Status())
	}
	return nil
}

func (c *ZincClient) request() *resty.Request {
	client := resty.New()
	client.DisableWarn = true
	client.SetBaseURL(c.ZincHost)
	client.SetBasicAuth(c.ZincUser, c.ZincPassword)
	return client.R()
}
