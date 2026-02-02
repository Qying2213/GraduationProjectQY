package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"common/database"

	"github.com/redis/go-redis/v9"
)

// DefaultTTL 默认缓存时间
const DefaultTTL = 5 * time.Minute

// Get 从缓存获取值
func Get(ctx context.Context, key string) (string, error) {
	client := database.GetRedis()
	if client == nil {
		return "", fmt.Errorf("redis client not initialized")
	}
	return client.Get(ctx, key).Result()
}

// Set 设置缓存值
func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	client := database.GetRedis()
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	var data string
	switch v := value.(type) {
	case string:
		data = v
	case []byte:
		data = string(v)
	default:
		bytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		data = string(bytes)
	}

	return client.Set(ctx, key, data, ttl).Err()
}

// GetJSON 从缓存获取JSON并反序列化
func GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// Delete 删除缓存
func Delete(ctx context.Context, key string) error {
	client := database.GetRedis()
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return client.Del(ctx, key).Err()
}

// GetOrSet 获取缓存，如果不存在则调用 loader 获取并缓存
func GetOrSet(ctx context.Context, key string, loader func() (interface{}, error), ttl time.Duration) (string, error) {
	// 尝试从缓存获取
	data, err := Get(ctx, key)
	if err == nil {
		return data, nil
	}

	// 缓存未命中，调用 loader
	if err != redis.Nil {
		// Redis 错误，降级直接调用 loader
		result, loaderErr := loader()
		if loaderErr != nil {
			return "", loaderErr
		}
		bytes, _ := json.Marshal(result)
		return string(bytes), nil
	}

	// 调用 loader 获取数据
	result, err := loader()
	if err != nil {
		return "", err
	}

	// 缓存结果
	_ = Set(ctx, key, result, ttl)

	bytes, _ := json.Marshal(result)
	return string(bytes), nil
}

// BuildKey 构建缓存键
func BuildKey(prefix string, parts ...interface{}) string {
	key := prefix
	for _, p := range parts {
		key = fmt.Sprintf("%s:%v", key, p)
	}
	return key
}

// CacheStats 缓存统计
type CacheStats struct {
	Hits   int64 `json:"hits"`
	Misses int64 `json:"misses"`
}

// IsAvailable 检查 Redis 是否可用
func IsAvailable() bool {
	client := database.GetRedis()
	if client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return client.Ping(ctx).Err() == nil
}
