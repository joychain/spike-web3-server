package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

func SetFromRedis(key string, value interface{}, duration time.Duration, client *redis.Client) error {
	key = strings.ToLower(key)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	if err := client.Set(ctx, key, value, duration).Err(); err != nil {
		log.Errorf("redis set key : %s value %v : err : %v", key, value, err)
		return err
	}
	return nil
}

func IncrFromRedis(key string, client *redis.Client) error {
	key = strings.ToLower(key)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	if err := client.Incr(ctx, key).Err(); err != nil {
		log.Errorf("redis set key : %s : err : %v", key, err)
		return err
	}
	return nil
}

func GetStringFromRedis(key string, client *redis.Client) (string, bool, error) {
	key = strings.ToLower(key)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", true, nil
		}
		log.Error(err.Error())
		return "", false, err
	}
	return val, false, nil
}

func GetIntFromRedis(key string, client *redis.Client) (int64, bool, error) {
	key = strings.ToLower(key)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	val, err := client.Get(ctx, key).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, true, nil
		}
		log.Error(err.Error())
		return 0, false, err
	}
	return val, false, nil
}

func RmKeyByPrefix(prefix string, client *redis.Client) {
	prefix = strings.ToLower(prefix)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	keys, err := client.Keys(ctx, fmt.Sprintf("%s*", prefix)).Result()
	if err != nil {
		log.Errorf("query redis keys err : %v", err)
	}
	if len(keys) == 0 {
		return
	}
	if err := client.Del(ctx, keys...).Err(); err != nil {
		log.Errorf("delete redis key err : %v", err)
	}
}

func Lock(key string, value interface{}, duration time.Duration, client *redis.Client) (bool, error) {
	key = strings.ToLower(key)

	result, err := client.SetNX(context.Background(), key, value, duration).Result()

	if err != nil {
		log.Errorf("redis lock key : %s value %v : err : %v", key, value, err)
		return false, err
	}

	return result, nil
}

func UnLock(key string, client *redis.Client) error {
	key = strings.ToLower(key)

	_, err := client.Del(context.Background(), key).Result()
	if err != nil {
		log.Errorf("redis unLock key : %s err : %v", key, err)
		return err
	}
	return nil
}
