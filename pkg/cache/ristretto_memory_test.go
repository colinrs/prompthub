package cache

import (
	"context"
	"testing"
	"time"

	"github.com/colinrs/prompthub/pkg/codec"
)

func TestRistrettoCache(t *testing.T) {

	memCache, err := NewRistrettoCache(RistrettoCacheConfig{
		NumCounters: 100,
		Capacity:    100,
		CostFunc:    func(value interface{}) int64 { return 1 },
	}, codec.NewSonicCodec())

	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	err = memCache.Set(ctx, "key", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	var value int
	err = memCache.Get(ctx, "key", &value)
	if err != nil {
		t.Fatal(err)
	}
	if value != 1 {
		t.Fatal("value not equal")
	}
	_ = memCache.Close(ctx)
}
