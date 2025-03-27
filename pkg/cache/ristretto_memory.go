package cache

import (
	"context"
	"time"

	"github.com/colinrs/prompthub/pkg/codec"
	"github.com/dgraph-io/ristretto"
	"github.com/zeromicro/go-zero/core/logx"
)

// RistrettoCacheConfig defines config used to build InMemoryCache backed by Ristretto cache
type RistrettoCacheConfig struct {
	Capacity        int64                       `yaml:"capacity" json:"capacity"`
	NumCounters     int64                       `yaml:"num_counters" json:"num_counters"`
	UseInternalCost bool                        `yaml:"use_internal_cost" json:"use_internal_cost"`
	CostFunc        func(val interface{}) int64 `yaml:"-" json:"-"`
}

// RistrettoCache is a in-memory cache, internally it uses Ristretto
type RistrettoCache struct {
	inner       *ristretto.Cache
	NumCounters int64
	CostFunc    func(val interface{}) int64
	plugins     []Plugin
	codec       codec.Codec
}

// NewRistrettoCache creates a new RistrettoCache
func NewRistrettoCache(ristrettoConfig RistrettoCacheConfig, codec codec.Codec) (Cache, error) {
	inner, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: ristrettoConfig.NumCounters,
		MaxCost:     ristrettoConfig.Capacity,
		BufferItems: 64,
		Cost:        ristrettoConfig.CostFunc,
	})
	if err != nil {
		return nil, err
	}
	return &RistrettoCache{
		inner:       inner,
		NumCounters: ristrettoConfig.NumCounters,
		CostFunc:    ristrettoConfig.CostFunc,
		codec:       codec,
	}, nil
}

func (c *RistrettoCache) Get(ctx context.Context, key string, receiver interface{}, opts ...OperationOption) error {
	value, ok := c.inner.Get(key)
	if !ok {
		return ErrNotFound
	}
	data, err := c.codec.Marshal(value)
	if err != nil {
		return err
	}
	err = c.codec.Unmarshal(data, receiver)
	if err != nil {
		return err
	}
	return nil
}

func (c *RistrettoCache) GetMany(ctx context.Context, receiverMap map[string]interface{}, opts ...OperationOption) error {
	for key, receiver := range receiverMap {
		err := c.Get(ctx, key, receiver, opts...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RistrettoCache) Set(ctx context.Context, key string, value interface{}, expire time.Duration, opts ...OperationOption) error {
	return c.setInner(key, value, expire)
}

func (c *RistrettoCache) SetMany(ctx context.Context, valueMap map[string]interface{}, expire time.Duration, opts ...OperationOption) error {
	for key, value := range valueMap {
		_ = c.setInner(key, value, expire)
	}
	return nil
}

func (c *RistrettoCache) setInner(key string, value interface{}, expire time.Duration) error {
	c.inner.SetWithTTL(key, value, c.CostFunc(value), expire)
	return nil
}

func (c *RistrettoCache) Delete(ctx context.Context, key string, opts ...OperationOption) error {
	c.inner.Del(key)
	return nil
}

func (c *RistrettoCache) DeleteMany(ctx context.Context, keys []string, opts ...OperationOption) error {
	for _, key := range keys {
		c.inner.Del(key)
	}
	return nil
}

func (c *RistrettoCache) Load(ctx context.Context, loader SourceLoaders, key string, receiver interface{}, expire time.Duration, opts ...OperationOption) error {
	err := c.Get(ctx, key, receiver, opts...)
	if err == nil {
		return nil
	}
	logx.WithContext(ctx).Debugf("load from local cache err: %s, now from source key:%+v", err.Error(), key)
	values, err := loader(ctx, []string{key})
	if err == nil && len(values) > 0 {
		value := values[0]
		b, err := c.codec.Marshal(value)
		if err != nil {
			return err
		}
		err = c.codec.Unmarshal(b, receiver)
		if err != nil {
			return err
		}
		go func() {
			c.inner.SetWithTTL(key, value, c.CostFunc(value), expire)
		}()
		return nil
	}
	return ErrFromSource
}

func (c *RistrettoCache) LoadMany(ctx context.Context, loader SourceLoaders, receiverMap map[string]interface{}, expire time.Duration, opts ...OperationOption) error {
	for key, receiver := range receiverMap {
		err := c.Load(ctx, loader, key, receiver, expire)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RistrettoCache) Ping(_ context.Context) error {
	return nil
}

func (c *RistrettoCache) Flush(_ context.Context) error {
	c.inner.Clear()
	return nil
}

func (c *RistrettoCache) Close(_ context.Context) error {
	c.inner.Close()
	return nil
}

func (c *RistrettoCache) AddPlugin(p Plugin) {
	c.plugins = append(c.plugins, p)
}
