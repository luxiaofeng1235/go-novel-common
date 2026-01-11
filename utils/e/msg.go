//Package e ...
/*
 * @Descripttion:
 * @Author: congz
 * @Date: 2020-07-15 14:53:24
 * @LastEditors: congz
 * @LastEditTime: 2020-08-18 21:46:46
 */
package e

// MsgFlags 状态码map
var MsgFlags = map[int]string{
	Success:                           "ok",
	Error:                             "fail",
	NotFound:                          "找不到该页面了",
	NotPermission:                     "访问被拒绝",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[Error]
}
