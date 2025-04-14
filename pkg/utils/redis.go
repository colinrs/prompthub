package utils

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	// 定义 Lua 脚本
	// 1. 使用 INCR 增加值
	// 2. 如果是第一次设置（值等于 1），则设置过期时间
	luaScript = `
		local val = redis.call('INCR', KEYS[1])
		if val == 1 then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end
		return val
	`
)

var (
	script = redis.NewScript(luaScript)
)

// IncrementAndSetTTLWithLua 使用 Lua 脚本执行 INCR 和 EXPIRE
func IncrementAndSetTTLWithLua(ctx context.Context, rdb *redis.Redis, key string, ttlSeconds int64) (int64, error) {
	// 执行 Lua 脚本
	// KEYS[1] 是 key，ARGV[1] 是 ttl（秒）
	result, err := rdb.ScriptRunCtx(ctx, script, []string{key}, ttlSeconds)
	if err != nil {
		return 0, err
	}
	if int64Result, ok := result.(int64); ok {
		return int64Result, nil
	}
	return 0, errors.New("redis incr failed")
}
