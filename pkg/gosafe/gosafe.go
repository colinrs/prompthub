package gosafe

import (
	"context"
	"fmt"

	"github.com/colinrs/prompthub/pkg/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

// Option 定义配置选项的接口
type Option func(*safeOptions)

// safeOptions 存储配置选项
type safeOptions struct {
	recoverHandler func(ctx context.Context, err error, str string)
}

// WithRecoverHandler 自定义 recover 处理函数的选项
func WithRecoverHandler(handler func(ctx context.Context, err error, stack string)) Option {
	return func(opts *safeOptions) {
		opts.recoverHandler = handler
	}
}

// 默认的 recover 处理函数
func defaultRecoverHandler(ctx context.Context, err error, stack string) {
	logx.WithContext(ctx).Errorf("Recovered from panic: %v\n%s", err, stack)
}

// GoSafe 运行一个可能会 panic 的函数，并提供安全的恢复机制
func GoSafe(ctx context.Context, fn func(), options ...Option) {
	// 初始化默认选项
	opts := &safeOptions{
		recoverHandler: defaultRecoverHandler,
	}

	// 应用用户提供的选项
	for _, opt := range options {
		opt(opts)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				// 将 interface{} 转换为 error
				var err error
				switch x := r.(type) {
				case error:
					err = x
				case string:
					err = fmt.Errorf("%s", x)
				default:
					err = fmt.Errorf("%v", x)
				}
				// 调用恢复处理函数
				opts.recoverHandler(ctx, err, utils.Stack())
				goSafePanicTotal.Inc()
			}
		}()
		// 执行传入的函数
		fn()
	}()
}
