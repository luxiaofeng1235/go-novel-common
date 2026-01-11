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
	"time"
)

func GetDB() (mysqlAddress, dbName, mysqlUser, mysqlPasswd string) {
	//初始化数据库
	env := config.GetString("server.env")
	mysqlAddress = "127.0.0.1:3306"
	mysqlUser = "novel"
	mysqlPasswd = "dkke4DJjELz3Tccs"
	if env == utils.Local {
		mysqlAddress = "127.0.0.1"
		mysqlUser = "root"
		mysqlPasswd = "root"
		dbName = utils.DbName
	} else if env == utils.Dev {
		mysqlAddress = "127.0.0.1"
		dbName = utils.DbName
		mysqlUser = "novel"
		mysqlPasswd = "dkke4DJjELz3Tccs"
	} else if env == utils.Prod {
		//ip, _ := utils.GetLocalIP()
		////判断是否为36那台机器，就加载36的方便测试
		//if ip == utils.AdminHostIp {
		//	mysqlAddress = "192.168.10.16:3306" //web的主机配置
		//} else {
		//线上API数据库连接
		mysqlAddress = "192.168.10.15:3306" //线上的api地址
		//}
		dbName = utils.DbName
		mysqlUser = "root"
		mysqlPasswd = "HM9GO3JH3XrLoouh"
	}
	return
}

func InitDB() {
	//连接mysql
	host, name, user, passwd := GetDB()
	fmt.Println("mysql host:", host, "mysql name:", name, "mysql user:", user, "mysql passwd:", passwd)
	InitMysql(host, name, user, passwd)

	//连接redis
	addr, passwd, defaultdb := GetRedis()
	fmt.Println("gredis addr:", addr, "gredis passwd:", passwd, "gredis defaultdb:", defaultdb)
	InitRedis(addr, passwd, defaultdb)

}

func InitMysql(address, dbname, user, passwd string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, passwd, address, dbname)

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
	sqlDB.SetMaxIdleConns(10000)
	//打开 SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(10000)
	//超时 SetConnMaxLifetiment 设置连接的最大可复用时间。
	//sqlDB.SetConnMaxLifetime(-1)
	sqlDB.SetConnMaxLifetime(60 * time.Second)
	return global.DB
}
