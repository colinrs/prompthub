package user

import (
	"context"
	"errors"

	"time"

	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/model"
	"github.com/colinrs/prompthub/pkg/code"
	"github.com/colinrs/prompthub/pkg/utils"

	"github.com/rs/xid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RegisterUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterUserLogic {
	return &RegisterUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterUserLogic) RegisterUser(req *types.RegisterUserRequest) (resp types.LoginResponse, err error) {
	userTable := gen.Use(l.svcCtx.DB).UsersTable
	userInfo, err := userTable.Where(userTable.Email.Eq(req.Email)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp, err
	}
	if userInfo != nil {
		return resp, code.ErrUserExist
	}
	value := &model.UsersTable{
		UserName:    req.Name,
		Password:    utils.HashPassword(req.Password, l.svcCtx.Config.PasswdSecret),
		Email:       req.Email,
		UserStatus:  config.UserStatusNormal,
		Description: "",
	}
	err = userTable.Save(value)
	if err != nil {
		return resp, err
	}
	resp.UserId = uint(value.ID)
	resp.Name = value.UserName
	resp.Email = value.Email
	expiredAt := time.Now().Add(time.Duration(l.svcCtx.Config.JwtExpired) * time.Second).Unix()
	claimsInfo := map[string]interface{}{
		"user_id": value.ID,
		"name":    value.UserName,
		"email":   value.Email,
		"tid":     xid.New().String(),
	}
	resp.ExpiredAt = expiredAt
	resp.Token, err = utils.GenerateJWT(claimsInfo, []byte(l.svcCtx.Config.JwtSecret), expiredAt)

	if err != nil {
		return resp, err
	}
	return resp, nil
}
