package main

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"reflect"
)

type LevelDBStore struct {
	db *leveldb.DB
}

type CustomLogger struct {
	logger *log.Logger
}

func NewCustomLogger() *CustomLogger {
	return &CustomLogger{
		logger: log.New(os.Stdout, "CUSTOM LOG: ", log.LstdFlags),
	}
}

func (l *CustomLogger) Logf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func NewLevelDBStore() (*LevelDBStore, error) {

	path := "./dbcache/leveldb" //定义缓存的路径
	// 检查数据库文件是否存在
	// 使用 os.MkdirAll 递归创建目录
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// 打开 LevelDB 数据库
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBStore{db: db}, nil
}

func (s *LevelDBStore) Close() error {
	return s.db.Close()
}

func (s *LevelDBStore) Put(key, value []byte) error {
	return s.db.Put(key, value, nil)
}

func (s *LevelDBStore) Get(key []byte) ([]byte, error) {
	return s.db.Get(key, nil)
}

func (s *LevelDBStore) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

func (s *LevelDBStore) PutString(key, value string) error {
	return s.Put([]byte(key), []byte(value))
}

func (s *LevelDBStore) GetString(key string) (string, error) {
	value, err := s.Get([]byte(key))
	return string(value), err
}

func main() {

	store, err := NewLevelDBStore()
	if err != nil {
		fmt.Println("Error creating LevelDBStore:", err)
		return
	}
	defer store.Close()
	//// 写入数据
	if err := store.PutString("hh", "334"); err != nil {
		fmt.Println("Error putting data:", err)
		return
	}

	// 读取数据
	value, _ := store.GetString("hh")
	//if err != nil {
	//	fmt.Println("Error getting data:", err)
	//	return
	//}
	if value == "" {
		fmt.Println("empty value")
	} else {
		fmt.Println("value:", value)
		fmt.Println(reflect.TypeOf(value))
	}
}
