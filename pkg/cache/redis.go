package cache

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/prompthub/pkg/codec"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mathx"
)

var (
	// make the unstable expiry to be [0.95, 1.05] * seconds
	expiryDeviation = 0.05
	defaultExpire   = 5 * time.Minute
)

type RedisConfig struct {
	Addr        string `yaml:"addr" json:"addr"`
	DB          int    `yaml:"db" json:"db"`
	PoolSize    int    `yaml:"pool_size" json:"pool_size"`
	IdleTimeout int    `yaml:"idle_timeout" json:"idle_timeout"`
	Prefix      string `yaml:"prefix" json:"prefix"`
	Username    string `yaml:"username" json:"username"`
	Password    string `yaml:"password" json:"password"`
}

type redisCache struct {
	client         *redis.Client
	prefix         string
	plugins        []Plugin
	unstableExpiry mathx.Unstable
	DefaultExpire  time.Duration
	codec          codec.Codec
}

var _ Cache = (*redisCache)(nil)

func getFullKey(prefix, key string) string {
	return prefix + "_" + key
}

func NewRedisCache(conf *RedisConfig, codec codec.Codec) Cache {
	redisCacheInstance := &redisCache{
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		DefaultExpire:  defaultExpire,
		codec:          codec,
	}
	redisCacheInstance.client = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		DB:       conf.DB,
		PoolSize: conf.PoolSize,
		Username: conf.Username,
		Password: conf.Password,
	})
	if conf.Prefix != "" {
		redisCacheInstance.prefix = conf.Prefix
	}
	return redisCacheInstance
}

func (r *redisCache) Get(ctx context.Context, key string, receiver interface{}, opts ...OperationOption) error {
	var byteValue []byte
	fullKey := getFullKey(r.prefix, key)
	startTime := time.Now()
	byteValue, err := r.client.Get(ctx, fullKey).Bytes()
	elapsed := time.Since(startTime).Milliseconds()
	for _, p := range r.plugins {
		p.OnGetRequestEnd(ctx, "cmdGet", elapsed, fullKey, err)
	}
	if err != nil {
		logx.WithContext(ctx).Errorf("get redis key: %v, error: %v", fullKey, err)
		return err
	}
	err = r.codec.Unmarshal(byteValue, receiver)
	// something err get key from redis
	if err != nil {
		logx.WithContext(ctx).Errorf("Unmarshal: %v, error: %v", fullKey, err)
		return err
	}
	return nil
}

func (r *redisCache) GetMany(ctx context.Context, receiverMap map[string]interface{}, opts ...OperationOption) error {
	for key, receiver := range receiverMap {
		err := r.Get(ctx, key, receiver, opts...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration, opts ...OperationOption) error {
	var byteValue []byte
	startTime := time.Now()
	fullKey := getFullKey(r.prefix, key)
	byteValue, err := r.codec.Marshal(value)
	if err != nil {
		logx.WithContext(ctx).Errorf("json.Marshal redis value: %v, error: %v", value, err)
		return err
	}
	expiration = r.unstableExpiry.AroundDuration(expiration)
	err = r.client.Set(ctx, fullKey, byteValue, expiration).Err()
	elapsed := time.Since(startTime).Milliseconds()
	for _, p := range r.plugins {
		p.OnSetRequestEnd(ctx, "cmdSet", elapsed, fullKey, err)
	}
	if err != nil {
		logx.WithContext(ctx).Errorf("set redis key: %v, error: %v", fullKey, err)
		return err
	}
	return nil
}

func (r *redisCache) SetMany(ctx context.Context, valueMap map[string]interface{}, expire time.Duration, opts ...OperationOption) error {
	for key, value := range valueMap {
		err := r.Set(ctx, key, value, expire, opts...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *redisCache) Delete(ctx context.Context, key string, opts ...OperationOption) error {
	fullKey := getFullKey(r.prefix, key)
	startTime := time.Now()
	_, err := r.client.Del(ctx, key).Result()
	elapsed := time.Since(startTime).Milliseconds()
	for _, p := range r.plugins {
		p.OnGetRequestEnd(ctx, "cmdDel", elapsed, fullKey, err)
	}
	// not found key
	if errors.Is(err, redis.Nil) {
		return nil
	}
	// something err get key from redis
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) DeleteMany(ctx context.Context, keys []string, opts ...OperationOption) error {
	for _, key := range keys {
		err := r.Delete(ctx, key, opts...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *redisCache) Load(ctx context.Context, loader SourceLoaders, key string, receiver interface{},
	expire time.Duration, opts ...OperationOption) error {
	err := r.Get(ctx, key, receiver, opts...)
	if err == nil {
		return nil
	}
	sourceData, err := loader(ctx, []string{key})
	if err != nil {
		return err
	}
	b, err := r.codec.Marshal(sourceData)
	if err != nil {
		return err
	}
	err = r.codec.Unmarshal(b, receiver)
	if err != nil {
		return err
	}
	go func() {
		_ = r.Set(context.WithoutCancel(ctx), key, sourceData, expire, opts...)
	}()
	return nil
}
func (r *redisCache) LoadMany(ctx context.Context, loader SourceLoaders,
	receiverMap map[string]interface{}, expire time.Duration, opts ...OperationOption) error {
	for key, receiver := range receiverMap {
		err := r.Load(ctx, loader, key, receiver, expire, opts...)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *redisCache) Flush(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}
func (r *redisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
func (r *redisCache) Close(ctx context.Context) error {
	return r.client.Close()
}

func (r *redisCache) AddPlugin(p Plugin) {
	r.plugins = append(r.plugins, p)
}
