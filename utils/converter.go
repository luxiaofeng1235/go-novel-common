/*
 * @Descripttion: Model 转换工具（读取 config.yml 的 MySQL 配置）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 10:25:00
 */
package utils

import (
	"fmt"
	"go-novel/config"
	"strings"

	"github.com/gohouse/converter"
)

func AutoModel() {
	// 初始化
	t2t := converter.NewTable2Struct()
	// 个性化配置
	t2t.Config(&converter.T2tConfig{
		// 如果字段首字母本来就是大写, 就不添加tag, 默认false添加, true不添加
		RmTagIfUcFirsted: false,
		// tag的字段名字是否转换为小写, 如果本身有大写字母的话, 默认false不转
		TagToLower: false,
		// 字段首字母大写的同时, 是否要把其他字母转换为小写,默认false不转换
		UcFirstOnly: false,
		//// 每个struct放入单独的文件,默认false,放入同一个文件(暂未提供)
		//SeperatFile: false,
	})
	// 开始迁移转换
	mysqlAddress := strings.TrimSpace(config.GetString("mysql.address"))
	if mysqlAddress == "" {
		host := strings.TrimSpace(config.GetString("mysql.host"))
		port := config.GetInt("mysql.port")
		if port == 0 {
			port = 3306
		}
		if host == "" {
			host = "127.0.0.1"
		}
		if strings.Contains(host, ":") {
			mysqlAddress = host
		} else {
			mysqlAddress = fmt.Sprintf("%s:%d", host, port)
		}
	}

	dbName := strings.TrimSpace(config.GetString("mysql.database"))
	if dbName == "" {
		dbName = strings.TrimSpace(config.GetString("mysql.dbname"))
	}
	if dbName == "" {
		dbName = strings.TrimSpace(config.GetString("mysql.name"))
	}
	if dbName == "" {
		dbName = DBNAME
	}

	mysqlUser := strings.TrimSpace(config.GetString("mysql.user"))
	if mysqlUser == "" {
		mysqlUser = MYSQLUSER
	}
	mysqlPassword := config.GetString("mysql.password")
	params := strings.TrimSpace(config.GetString("mysql.params"))
	params = strings.TrimPrefix(params, "?")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", mysqlUser, mysqlPassword, mysqlAddress, dbName, params)

	err := t2t.
		// 指定某个表,如果不指定,则默认全部表都迁移
		Table("mc_book_feedback").
		// 表前缀
		Prefix("").
		// 是否添加json tag
		EnableJsonTag(true).
		// 生成struct的包名(默认为空的话, 则取名为: package model)
		PackageName("model").
		// tag字段的key值,默认是orm
		TagKey("gorm").
		// 是否添加结构体方法获取表名
		RealNameMethod("TableName").
		// 生成的结构体保存路径
		SavePath("./model.go").
		// 数据库dsn,这里可以使用 t2t.DB() 代替,参数为 *sql.DB 对象
		Dsn(dsn).
		// 执行
		Run()

	fmt.Println(err)
}
