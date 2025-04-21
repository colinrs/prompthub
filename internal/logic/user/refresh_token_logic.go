package user

import (
	"context"
	"time"

	"github.com/colinrs/prompthub/pkg/constant"
	"github.com/colinrs/prompthub/pkg/utils"
	"github.com/rs/xid"

	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(_ *types.RefreshTokenRequeste) (resp *types.RefreshTokenResponse, err error) {
	resp = &types.RefreshTokenResponse{}
	expiredAt := time.Now().Add(time.Duration(l.svcCtx.Config.JwtExpired) * time.Second).Unix()
	userId, err := utils.GetUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	userEmail, err := utils.GetUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	claimsInfo := map[string]interface{}{
		constant.UserId:   userId,
		constant.UserName: utils.GetUserNameFromCtx(l.ctx),
		constant.Email:    userEmail,
		"tid":             xid.New().String(),
	}
	resp.ExpiredAt = expiredAt
	resp.Token, err = utils.GenerateJWT(claimsInfo, []byte(l.svcCtx.Config.JwtSecret), expiredAt)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
