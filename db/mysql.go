/*
 * @Descripttion: MySQL/Redis 初始化（从 config.yml 读取）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 10:25:00
 */
package db

import (
	"database/sql"
	"fmt"
	"go-novel/config"
	"go-novel/global"
	"go-novel/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"time"
)

func GetDB() (mysqlAddress, dbName, mysqlUser, mysqlPasswd string) {
	mysqlAddress = strings.TrimSpace(config.GetString("mysql.address"))
	if mysqlAddress == "" {
		host := strings.TrimSpace(config.GetString("mysql.host"))
		port := config.GetInt("mysql.port")
		if port == 0 {
			port = 3306
		}
		if host == "" {
			log.Fatal("缺少 MySQL 配置：请在 config.yml 中设置 mysql.host（以及 mysql.port）或 mysql.address")
		}
		if strings.Contains(host, ":") {
			mysqlAddress = host
		} else {
			mysqlAddress = fmt.Sprintf("%s:%d", host, port)
		}
	}

	dbName = strings.TrimSpace(config.GetString("mysql.database"))
	if dbName == "" {
		dbName = strings.TrimSpace(config.GetString("mysql.dbname"))
	}
	if dbName == "" {
		dbName = strings.TrimSpace(config.GetString("mysql.name"))
	}
	if dbName == "" {
		dbName = utils.DBNAME
	}

	mysqlUser = strings.TrimSpace(config.GetString("mysql.user"))
	if mysqlUser == "" {
		mysqlUser = "root"
	}
	mysqlPasswd = config.GetString("mysql.password")
	return
}

func InitDB(args ...string) {
	//连接mysql
	var host, name, user, passwd string
	switch len(args) {
	case 0:
		host, name, user, passwd = GetDB()
	case 4:
		host, name, user, passwd = args[0], args[1], args[2], args[3]
	default:
		log.Fatalf("InitDB 参数数量错误：期望 0 或 4 个参数，实际 %d 个", len(args))
	}
	fmt.Println("mysql addr:", host, "mysql db:", name, "mysql user:", user)
	InitMysql(host, name, user, passwd)

	//连接redis
	addr, passwd, defaultdb := GetRedis()
	fmt.Println("redis addr:", addr, "redis db:", defaultdb)
	InitRedis(addr, passwd, defaultdb)

}

func InitMysql(address, dbname, user, passwd string) *gorm.DB {
	params := strings.TrimSpace(config.GetString("mysql.params"))
	if params == "" {
		params = "charset=utf8mb4&parseTime=True&loc=Local"
	}
	params = strings.TrimPrefix(params, "?")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", user, passwd, address, dbname, params)

	//fmt.Println("dsn:",dsn)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags|log.Llongfile), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // 开启彩色打印
		},
	)

	var sqlDB *sql.DB
	var err error

	if sqlDB != nil {
		if err := sqlDB.Ping(); err != nil {
			log.Println("数据库断开重连")
		}
	}

	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "mc_",
			SingularTable: true,
		},
	})

	//Db.LogMode(true)
	// Error
	if err != nil {
		log.Fatal("连接数据库失败，请检查参数：", err)
	}
	fmt.Println("连接数据库成功")

	sqlDB, _ = global.DB.DB()

	//设置连接池
	//空闲 SetMaxIdleCons 设置连接池中的最大闲置连接数。
	maxIdle := config.GetInt("mysql.pool.maxIdleConns")
	if maxIdle <= 0 {
		maxIdle = 25
	}
	sqlDB.SetMaxIdleConns(maxIdle)
	//打开 SetMaxOpenCons 设置数据库的最大连接数量。
	maxOpen := config.GetInt("mysql.pool.maxOpenConns")
	if maxOpen <= 0 {
		maxOpen = 100
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	//超时 SetConnMaxLifetiment 设置连接的最大可复用时间。
	//sqlDB.SetConnMaxLifetime(-1)
	connMaxLifetimeSeconds := config.GetInt("mysql.pool.connMaxLifetimeSeconds")
	if connMaxLifetimeSeconds <= 0 {
		connMaxLifetimeSeconds = 600
	}
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetimeSeconds) * time.Second)

	connMaxIdleTimeSeconds := config.GetInt("mysql.pool.connMaxIdleTimeSeconds")
	if connMaxIdleTimeSeconds > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(connMaxIdleTimeSeconds) * time.Second)
	}
	return global.DB
}
