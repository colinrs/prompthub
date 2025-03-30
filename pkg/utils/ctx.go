package utils

import (
	"context"
	"github.com/colinrs/prompthub/pkg/code"
	"github.com/colinrs/prompthub/pkg/constant"
)

func GetUserIDFromCtx(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value(constant.UserId).(float64)
	if !ok {
		return 0, code.ErrUserNoLogin
	}
	if userID <= 0 {
		return 0, code.ErrUserNoLogin
	}
	return uint(userID), nil
}

func GetUserEmailFromCtx(ctx context.Context) string {
	email, _ := ctx.Value(constant.Email).(string)
	return email
}

func GetUserNameFromCtx(ctx context.Context) string {
	name, _ := ctx.Value(constant.UserName).(string)
	return name
}
