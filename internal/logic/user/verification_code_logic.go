package user

import (
	"context"
	"fmt"

	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerificationCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationCodeLogic {
	return &VerificationCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationCodeLogic) VerificationCode(req *types.VerificationCodeRequest) (resp *types.VerificationCodeResponse, err error) {
	userTable := gen.Use(l.svcCtx.DB).UsersTable
	userInfo, err := userTable.Where(userTable.Email.Eq(req.Email), userTable.UserStatus.In(
		config.UserStatusNormal, config.UserStatusPending)).First()
	if err != nil {
		return nil, err
	}
	if userInfo == nil {
		return nil, code.ErrUserNotExist
	}
	key := fmt.Sprintf("%s:%s", config.VerificationCode, req.Code)
	value, err := l.svcCtx.RedisClient.GetCtx(l.ctx, key)
	if err != nil {
		return nil, err
	}
	if value != req.Email {
		return nil, code.ErrVerificationCodeInvalid
	}
	_, err = l.svcCtx.RedisClient.DelCtx(l.ctx, key)
	l.Debugf("check verification code, email: %s, code: %s", req.Email, key)
	if userInfo.UserStatus == config.UserStatusPending {
		userInfo.UserStatus = config.UserStatusNormal
		err = userTable.Save(userInfo)
		if err != nil {
			return nil, err
		}
	}
	resp = &types.VerificationCodeResponse{}
	return
}
