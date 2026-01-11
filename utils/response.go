package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"go-novel/global"
	"go-novel/utils/e"
	"net/http"
)

// Response 基础序列化器
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Key   string      `json:"key,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty"`
}

type PageResult struct {
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

func Fail(c *gin.Context, err error, msg string) {
	if err != nil {
		msg = fmt.Sprintf("%v", err.Error())
	}
	response := Response{
		Code: e.Error,
		Msg:  fmt.Sprintf("%v", msg),
	}
	c.Set("apiReturnRes", JSONString(response))
	c.JSON(http.StatusOK, response)
}

// 成功信息
func SuccessEncrypt(c *gin.Context, data interface{}, msg string) {
	response := Response{
		Code: e.Success,
		Data: data,
		Msg:  msg,
	}
	isEncrypt := GetApiEncrypt()
	if isEncrypt {
		encrypt, _ := AesEncryptByCFB(ApiAesKey, JSONString(response))
		res := gin.H{
			"data": encrypt,
		}
		c.JSON(http.StatusOK, res)
		return
	}
	c.JSON(http.StatusOK, response)
}

func FailEncrypt(c *gin.Context, err error, msg string) {
	if err != nil {
		msg = fmt.Sprintf("%v", err.Error())
	}
	response := Response{
		Code: e.Error,
		Msg:  fmt.Sprintf("%v", msg),
	}

	isEncrypt := GetApiEncrypt()
	if isEncrypt {
		encrypt, _ := AesEncryptByCFB(ApiAesKey, JSONString(response))
		res := gin.H{
			"data": encrypt,
		}
		c.JSON(http.StatusOK, res)
		return
	}
	c.JSON(http.StatusOK, response)
}

// 成功信息
func Success(c *gin.Context, data interface{}, msg string) {
	response := Response{
		Code: e.Success,
		Data: data,
		Msg:  msg,
	}
	c.Set("apiReturnRes", JSONString(response))
	c.JSON(http.StatusOK, response)
}

// 处理200的返回值
func SuccessBaidu(c *gin.Context, data interface{}, msg string) {
	response := Response{
		Code: e.SUCCESS200,
		Data: data,
		Msg:  msg,
	}
	c.Set("apiReturnRes", JSONString(response))
	c.JSON(http.StatusOK, response)
}

func WsFail(s *melody.Session, err error, protocol, key, msg string) {
	if err != nil {
		msg = fmt.Sprintf("%v", err.Error())
	}
	response := Response{
		Code: e.Error,
		Key:  fmt.Sprintf("%v", key),
		Msg:  fmt.Sprintf("%v", msg),
	}
	err = s.Write([]byte(fmt.Sprintf("%v%v", protocol, JSONString(response))))
	if err != nil {
		if global.Errlog != nil {
			global.Errlog.Error(err.Error())
		}
		_ = s.Write([]byte(fmt.Sprintf("%v%v", protocol, err.Error())))
	}
}

// 成功信息
func WsSuccess(s *melody.Session, protocol, key string, data interface{}, msg string) {
	response := Response{
		Code: e.Success,
		Data: data,
		Key:  key,
		Msg:  msg,
	}
	var err error
	err = s.Write([]byte(fmt.Sprintf("%v%v", protocol, JSONString(response))))
	if err != nil {
		if global.Errlog != nil {
			global.Errlog.Error(err.Error())
		}
		_ = s.Write([]byte(fmt.Sprintf("%v%v", protocol, err.Error())))
	}
}

// TokenData 带有token的Data结构
type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}
