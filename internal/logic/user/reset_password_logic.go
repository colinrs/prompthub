package user

import (
	"context"
	"fmt"
	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/pkg/code"
	"github.com/colinrs/prompthub/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordRequest) (resp *types.ResetPasswordResponse, err error) {
	userTable := gen.Use(l.svcCtx.DB).UsersTable
	userInfo, err := userTable.Where(userTable.Email.Eq(req.Email), userTable.UserStatus.Eq(
		config.UserStatusNormal)).First()
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
	if !(value == req.Email && value == userInfo.Email) {
		return nil, code.ErrVerificationCodeInvalid
	}
	newPasswd := utils.HashPassword(req.NewPassword, l.svcCtx.Config.PasswdSecret)
	if newPasswd == req.NewPassword {
		return nil, code.ErrPasswordSameInvalid
	}
	userInfo.Password = utils.HashPassword(req.NewPassword, l.svcCtx.Config.PasswdSecret)
	_, err = userTable.Where(userTable.ID.Eq(userInfo.ID)).Updates(userInfo)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.RedisClient.DelCtx(l.ctx, key)
	if err != nil {
		l.Errorf("delete verification code failed, err: %v", err)
	}
	return
}
