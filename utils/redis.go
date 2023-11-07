package utils

import (
	"encoding/json"
	"gopan/global"
	"time"
)

const (
	placeHolder = "*"
)

/**
redis工具
*/

func GetFromRedis[T any](key string) (*T, error) {
	var target T
	bytes, err := global.RedisClient.Get(global.Ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &target)
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func SetToRedis[T any](key string, data *T, expiration time.Duration) error {
	// 占位符防止缓存穿透
	if data == nil {
		global.RedisClient.Set(global.Ctx, key, []byte(placeHolder), time.Minute)
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	global.RedisClient.Set(global.Ctx, key, bytes, expiration)
	return nil
}
