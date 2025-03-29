package user

import (
	"context"
	"errors"
	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/pkg/code"
	"github.com/colinrs/prompthub/pkg/utils"
	"github.com/rs/xid"
	"gorm.io/gorm"
	"time"

	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	userTable := gen.Use(l.svcCtx.DB).UsersTable
	userInfo, err := userTable.Where(userTable.Email.Eq(req.Email)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp, err
	}
	if userInfo == nil {
		return resp, code.ErrUserNotExist
	}
	if !utils.CheckPassword(req.Password, userInfo.Password, l.svcCtx.Config.PasswdSecret) {
		return resp, code.ErrPasswordIncorrect
	}
	resp = &types.LoginResponse{}
	resp.UserId = uint(userInfo.ID)
	resp.Name = userInfo.UserName
	resp.Email = userInfo.Email
	expiredAt := time.Now().Add(time.Duration(l.svcCtx.Config.JwtExpired) * time.Second).Unix()
	claimsInfo := map[string]interface{}{
		"user_id": userInfo.ID,
		"name":    userInfo.UserName,
		"email":   userInfo.Email,
		"tid":     xid.New().String(),
	}
	resp.ExpiredAt = expiredAt
	resp.Token, err = utils.GenerateJWT(claimsInfo, []byte(l.svcCtx.Config.JwtSecret), expiredAt)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
