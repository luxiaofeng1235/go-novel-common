/*
 * @Descripttion: API 服务入口
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 13:55:00
 */
package main

import (
	"go-novel/db"
)

func main() {
	db.StartApiServer()
}
