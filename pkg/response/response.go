package response

import (
	"context"
	"errors"

	"github.com/colinrs/prompthub/pkg/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type Response[T any] struct {
	Code   int           `json:"code"`
	Msg    string        `json:"msg"`
	Data   T             `json:"data"`
	Errors []*code.Error `json:"errors,omitempty"`
}

func ErrHandle(ctx context.Context, err error) (int, any) {
	logx.WithContext(ctx).Errorf("err: %v", err)
	var v *code.Err
	switch {
	case errors.As(err, &v):
		return v.GetHTTPCode(), Response[any]{
			Code:   v.GetCode(),
			Msg:    v.Error(),
			Data:   nil,
			Errors: v.GetErrors(),
		}
	default:
		return code.UnknownErr.GetHTTPCode(), Response[any]{
			Code:   code.UnknownErr.GetCode(),
			Msg:    code.UnknownErr.Error(),
			Data:   nil,
			Errors: code.UnknownErr.GetErrors(),
		}
	}

}

func OKHandle(ctx context.Context, data any) any {
	return Response[any]{
		Code: code.OKErr.GetCode(),
		Msg:  code.OKErr.Error(),
		Data: data,
	}
}
