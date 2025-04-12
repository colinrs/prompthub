package user

import (
	"context"
	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/pkg/code"
	"github.com/colinrs/prompthub/pkg/utils"

	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordRequest) error {
	if req.NewPassword == req.OldPassword {
		return code.ErrPasswordSameInvalid
	}
	userTable := gen.Use(l.svcCtx.DB).UsersTable
	userId, err := utils.GetUserIDFromCtx(l.ctx)
	if err != nil {
		return code.ErrUserNoLogin
	}
	userInfo, err := userTable.Where(userTable.ID.Eq(int32(userId)), userTable.UserStatus.Eq(
		config.UserStatusNormal)).First()
	if err != nil {
		return err
	}
	if userInfo == nil {
		return code.ErrUserNotExist
	}
	if userInfo.Password != utils.HashPassword(req.OldPassword, l.svcCtx.Config.PasswdSecret) {
		return code.ErrPasswordIncorrect
	}
	userInfo.Password = utils.HashPassword(req.NewPassword, l.svcCtx.Config.PasswdSecret)
	_, err = userTable.Where(userTable.ID.Eq(userInfo.ID)).Updates(userInfo)
	if err != nil {
		return err
	}
	return nil
}
