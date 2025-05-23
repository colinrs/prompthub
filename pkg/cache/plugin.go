package cache

import (
	"context"
	"errors"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
)

type Plugin interface {
	OnSetRequestEnd(ctx context.Context, cmd string, elapsed int64, fullKey string, err error)
	OnGetRequestEnd(ctx context.Context, cmd string, elapsed int64, fullKey string, err error)
}

var _ Plugin = (*defaultCachePlugin)(nil)

type defaultCachePlugin struct {
}

const (
	OK      = "OK"
	Success = "Success"
	Failed  = "Failed"
)

var DefaultCachePlugin = &defaultCachePlugin{}

func (d *defaultCachePlugin) OnSetRequestEnd(ctx context.Context, cmd string, elapsed int64, fullKey string, err error) {
	var errString, cmdStatus string
	errString = OK
	cmdStatus = Success
	if err != nil {
		errString = err.Error()
		cmdStatus = Failed
	}
	logx.WithContext(ctx).Infof("redis_cmd:%s|fullKey:%s|err:%s|elapsed:%dms,cmd status:%s", cmd, fullKey, errString, elapsed, cmdStatus)
}

func (d *defaultCachePlugin) OnGetRequestEnd(ctx context.Context, cmd string, elapsed int64, fullKey string, err error) {
	var errString, isHit, cmdStatus string
	errString = OK
	isHit = Success
	cmdStatus = Success
	if err != nil {
		if errors.Is(err, redis.Nil) {
			errString = "not found key"
			isHit = Failed
		} else {
			errString = err.Error()
			cmdStatus = Failed
			isHit = Failed // not count hit
		}
	}
	logx.WithContext(ctx).Infof("redis_cmd:%s|is_hit:%s|fullKey:%s|err:%s|elapsed:%dms, cmd status:%s",
		cmd, isHit, fullKey, errString, elapsed, cmdStatus)
}
