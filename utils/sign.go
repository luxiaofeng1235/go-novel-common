package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
)

var (
	StoreNonce sync.Map
)

// 签名算法
func PaySign(params map[string]interface{}, token string) string {
	var keys []string
	//如果该key对应的value为空，则不参与签名
	for k, v := range params {
		if v != "" {
			keys = append(keys, k)
		}
	}
	// 对参数按照参数名进行升序排列
	sort.Strings(keys)

	// 将参数和值进行拼接，用"&"连接
	var s []string
	for _, k := range keys {
		v := params[k]
		s = append(s, fmt.Sprintf("%v=%v", k, v))
	}
	signStr := strings.Join(s, "&")
	signStr = fmt.Sprintf("%v&key=%v", signStr, token)
	signStr = Md5(signStr)
	return signStr
}

// 签名算法
func ApiSign(params map[string]string, secret string) string {
	var keys []string
	//如果该key对应的value为空，则不参与签名
	for k, v := range params {
		if v != "" {
			keys = append(keys, k)
		}
	}
	// 对参数按照参数名进行升序排列
	sort.Strings(keys)

	// 将参数和值进行拼接，用"&"连接
	var s []string
	for _, k := range keys {
		v := params[k]
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}
	signStr := strings.Join(s, "&")

	signStr = fmt.Sprintf("%v&%v", signStr, Md5(secret))

	signStr = strings.ToUpper(Md5(signStr))

	return signStr
}

func SendApiPostSign(strUrl string, mapParams map[string]string, headers map[string]string) string {
	httpClient := &http.Client{}

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Language", "zh-cn")
	if len(headers) > 0 {
		for key, val := range headers {
			request.Header.Add(key, val)
		}
	}

	response, err := httpClient.Do(request)
	defer response.Body.Close()
	if nil != err {
		return err.Error()
	}

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}
