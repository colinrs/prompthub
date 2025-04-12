package user

import (
	"context"

	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/model"
	"github.com/colinrs/prompthub/pkg/code"
	"github.com/colinrs/prompthub/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) error {
	userTable := gen.Use(l.svcCtx.DB).UsersTable
	userId, err := utils.GetUserIDFromCtx(l.ctx)
	if err != nil {
		return err
	}
	if userId == 0 {
		return code.ErrUserNoLogin
	}
	userInfo, err := userTable.Where(userTable.ID.Eq(int32(userId))).First()
	if err != nil {
		return err
	}
	if userInfo == nil {
		return code.ErrUserNotExist
	}
	newUser := &model.UsersTable{
		UserName: req.Name,
		Avatar:   req.Avatar,
	}
	_, err = userTable.WithContext(l.ctx).Where(userTable.ID.Eq(int32(userId))).Updates(newUser)
	if err != nil {
		return err
	}
	return nil
}
