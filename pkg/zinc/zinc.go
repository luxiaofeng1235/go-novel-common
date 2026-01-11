package zinc

import (
	"github.com/cockroachdb/errors"
	"github.com/go-resty/resty/v2"
	"go-novel/config"
	"go-novel/utils"
	"net/http"
	"time"
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

// NewConfig initializes the Zinc configuration
func NewConfig() (string, string, string) {
	env := config.GetString("server.env")

	if env == utils.Local {
		return "http://103.36.91.96:4080", "admin", "SInR5cCI6IkpXV#25"
	} else if env == utils.Dev {
		return "http://127.0.0.1:4080", "admin", "SInR5cCI6IkpXV#25"
	} else {
		return "http://127.0.0.1:4080", "admin", "SInR5cCI6IkpXV#25"
	}
}
