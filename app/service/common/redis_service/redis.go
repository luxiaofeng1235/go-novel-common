package redis_service

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go-novel/global"
	"time"
)

func Get(key string) string {
	ctx := context.Background()

	value, err := global.Redis.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return ""
	}
	return value
}

// key name zhangsan 2*time.Hour
func Set(key string, value interface{}, expiration time.Duration) error {
	valueData, err := json.Marshal(value)
	ctx := context.Background()

	if err != nil {
		return err
	}
	return global.Redis.Set(ctx, key, valueData, expiration).Err()
}

func Exist(key string) bool {
	ctx := context.Background()

	num, err := global.Redis.Exists(ctx, key).Result()
	if err != nil || num < 1 {
		return false
	}
	return true
}

// Del 删除key keys := redis_service.Keys() redis_service.Del(keys...)
func Del(key ...string) error {
	ctx := context.Background()

	return global.Redis.Del(ctx, key...).Err()
}

// Keys 返回与pattern匹配的key
func Keys() []string {
	ctx := context.Background()

	return global.Redis.Keys(ctx, "*").Val()
}

func GetListValueAll(key string) []string {
	ctx := context.Background()

	info, _ := global.Redis.LRange(ctx, key, 0, -1).Result()
	return info
}

func GetListValueRange(key string, start int, end int) []string {
	ctx := context.Background()

	info, _ := global.Redis.LRange(ctx, key, int64(start), int64(end)).Result()
	return info
}

func AddSetValue(key string, value string) int64 {
	ctx := context.Background()

	info, _ := global.Redis.SAdd(ctx, key, value).Result()
	return info
}

func DeleteSetValue(key string, setKey string) int64 {
	ctx := context.Background()

	info, _ := global.Redis.SRem(ctx, key, setKey).Result()
	return info
}

func Rename(key string, newName string) string {
	ctx := context.Background()

	t, _ := global.Redis.Rename(ctx, key, newName).Result()
	return t
}

func GetType(key string) string {
	ctx := context.Background()

	t, _ := global.Redis.Type(ctx, key).Result()
	return t
}

// set类型指定返回固定条数
func SRandMemberN(key string, count int64) (vals []string) {
	ctx := context.Background()

	vals = global.Redis.SRandMemberN(ctx, key, count).Val()
	return
}

// set类型删除key值
func SRem(key string, v ...interface{}) (err error) {
	ctx := context.Background()

	err = global.Redis.SRem(ctx, key, v...).Err()
	return
}

// set类型批量设置
func HMSet(key string, fields map[string]interface{}, expiration time.Duration) (err error) {
	ctx := context.Background()
	err = global.Redis.HMSet(ctx, key, fields, expiration).Err()
	return
}

// get类型设置
func HGet(key, field string) (err error) {
	ctx := context.Background()
	global.Redis.HGet(ctx, key, field).Val()
	return
}

func HExists(key, field string) (isExists bool, err error) {
	ctx := context.Background()
	isExists = global.Redis.HExists(ctx, key, field).Val()
	return
}

func HGetAll(key string) (val map[string]string, err error) {
	ctx := context.Background()
	val = global.Redis.HGetAll(ctx, key).Val()
	return
}

// set类型设置
func HSet(key, field string, value interface{}) (err error) {
	ctx := context.Background()
	err = global.Redis.HSet(ctx, key, field, value).Err()
	return
}

// hash类型删除
func HDel(key string, fields ...string) (err error) {
	ctx := context.Background()

	err = global.Redis.HDel(ctx, key, fields...).Err()
	return
}

// set类型指定返回固定条数
func ZMemberN(key string, count int64) (vals []string) {
	ctx := context.Background()

	vals = global.Redis.ZRange(ctx, key, 0, count).Val()
	return
}
