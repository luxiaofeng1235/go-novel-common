package gredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go-novel/global"
	"strconv"
	"time"
)

var ctx = context.Background()

// Set a key/value
func Set(key string, data interface{}, expiration time.Duration) error {
	err := global.Redis.Set(ctx, key, data, expiration).Err()
	if err != nil {
		global.Errlog.Errorf("gredis set failed %v", err)
		return err
	}
	return nil
}

// Set a key
func Get(key string) string {
	val, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		//global.Errlog.Errorf("gredis get failed %v", err)
		return ""
	}
	return val
}
func Del(key string) (res int64, err error) {
	val, err := global.Redis.Del(ctx, key).Result()
	if err != nil {
		global.Errlog.Errorf("gredis Del failed %v", err)
		return 0, err
	}
	return val, nil
}

func Zadd(key string, members redis.Z) (int64, error) {
	val, err := global.Redis.ZAdd(ctx, key, &members).Result()
	if err != nil {
		global.Errlog.Errorf("gredis zadd failed %v", err)
		return 0, err
	}
	return val, nil
}
func Zscore(key string, uid int) (int64, int64) {
	val, err := global.Redis.ZScore(ctx, key, strconv.Itoa(uid)).Result()
	if err != nil {
		global.Errlog.Errorf("gredis ZScore failed %v", err)
		return 0, 0
	}
	return int64(val), 0

}

func Zrem(key string, members ...interface{}) (int64, error) {
	val, err := global.Redis.ZRem(ctx, key, members).Result()
	if err != nil {
		global.Errlog.Errorf("gredis zrem failed %v", err)
		return 0, err
	}
	return val, nil
}

func Incr(key string) (int64, error) {
	val, err := global.Redis.Incr(ctx, key).Result()
	if err != nil {
		global.Errlog.Errorf("gredis incr failed %v", err)
		return 0, err
	}
	return val, nil
}

func Decr(key string) (int64, error) {
	val, err := global.Redis.Decr(ctx, key).Result()
	if err != nil {
		global.Errlog.Errorf("gredis decr failed %v", err)
		return 0, err
	}
	return val, nil
}

// 无序集合相关
func SAdd(key string, members ...interface{}) (res int64, err error) {
	val, err := global.Redis.SAdd(ctx, key, members).Result()
	if err != nil {
		global.Errlog.Errorf("gredis sadd failed %v", err)
		return 0, err
	}
	return val, nil
}

func SRem(key string, members ...interface{}) (res int64, err error) {
	val, err := global.Redis.SRem(ctx, key, members).Result()
	if err != nil {
		global.Errlog.Errorf("gredis srem failed %v", err)
		return 0, err
	}
	return val, nil
}

func SIsMember(key string, members interface{}) (res bool) {
	val, err := global.Redis.SIsMember(ctx, key, members).Result()
	if err != nil {
		global.Errlog.Errorf("gredis SIsMember failed %v", err)
		return false
	}
	return val
}

func SMembers(key string) (res []string) {

	val, err := global.Redis.SMembers(ctx, key).Result()
	if err != nil {
		global.Errlog.Errorf("gredis SMembers failed %v", err)
		var res []string
		return res
	}
	return val
}

func Zcard(key string) (res int64, err error) {

	val, err := global.Redis.ZCard(ctx, key).Result()
	if err != nil {
		global.Errlog.Errorf("gredis ZCard failed %v", err)
		return 0, err
	}
	return val, nil
}

func ZunionStore(key string, z *redis.ZStore) (res int64, err error) {

	val, err := global.Redis.ZUnionStore(ctx, key, z).Result()
	if err != nil {
		global.Errlog.Errorf("gredis ZUnionStore failed %v", err)
		return 2212312, err
	}
	return val, nil
}

func Zremrangebyrank(key string, start int64, end int64) (res int64, err error) {

	val, err := global.Redis.ZRemRangeByRank(ctx, key, start, end).Result()
	if err != nil {
		global.Errlog.Errorf("gredis ZRemRangeByRank failed %v", err)
		return 0, err
	}
	return val, nil
}
func Zrangebyscore(key string, opt *redis.ZRangeBy) (res []string, err error) {

	val, err := global.Redis.ZRangeByScore(ctx, key, opt).Result()
	if err != nil {
		global.Errlog.Errorf("gredis ZRangeByScore failed %v", err)
		return nil, err
	}
	return val, nil
}
func Zrange(key string, start int64, end int64) (res []string, err error) {

	val, err := global.Redis.ZRange(ctx, key, start, end).Result()
	if err != nil {
		global.Errlog.Errorf("gredis ZRange failed %v", err)
		return nil, err
	}
	return val, nil

}
func Zrevrange(key string, start int64, end int64) (res []string, err error) {

	val, err := global.Redis.ZRevRange(ctx, key, start, end).Result()
	if err != nil {
		global.Errlog.Errorf("gredis ZRevRange failed %v", err)
		return nil, err
	}
	return val, nil
}

// 队列  入队列
func LPush(key string, values ...interface{}) (res int64, err error) {

	result, err := global.Redis.LPush(ctx, key, values).Result()
	if err != nil {
		//global.Errlog.Errorf("gredis LPush failed %v", err)
		return 0, err
	}
	return result, nil

}

// 队列  出队列
func RPop(key string) (res string, err error) {

	result, err := global.Redis.RPop(ctx, key).Result()
	if err != nil {
		//global.Errlog.Errorf("gredis LPush failed %v", err)
		return "", err
	}
	return result, nil

}
