package cache

import (
	"context"
	"time"
)

type SourceLoaders func(ctx context.Context, keys []string) ([]interface{}, error)

// Cache is the interface of a cache store
type Cache interface {
	Get(ctx context.Context, key string, receiver interface{}, opts ...OperationOption) error
	GetMany(ctx context.Context, receiverMap map[string]interface{}, opts ...OperationOption) error
	Set(ctx context.Context, key string, value interface{}, expire time.Duration, opts ...OperationOption) error
	SetMany(ctx context.Context, valueMap map[string]interface{}, expire time.Duration, opts ...OperationOption) error
	Delete(ctx context.Context, key string, opts ...OperationOption) error
	DeleteMany(ctx context.Context, keys []string, opts ...OperationOption) error
	Load(ctx context.Context, loader SourceLoaders, key string, receiver interface{}, expire time.Duration, opts ...OperationOption) error
	LoadMany(ctx context.Context, loader SourceLoaders, receiverMap map[string]interface{}, expire time.Duration, opts ...OperationOption) error
	Flush(ctx context.Context) error
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}
