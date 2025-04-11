package user

import (
	"context"
	"fmt"
	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/pkg/code"
	"github.com/colinrs/prompthub/pkg/utils"
	"github.com/rs/xid"

	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendVerificationCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendVerificationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendVerificationCodeLogic {
	return &SendVerificationCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendVerificationCodeLogic) SendVerificationCode(req *types.SendVerificationCodeRequest) (resp *types.SendVerificationCodeResponse, err error) {
	userTable := gen.Use(l.svcCtx.DB).UsersTable
	count, err := userTable.Where(userTable.Email.Eq(req.Email), userTable.UserStatus.In(
		config.UserStatusNormal, config.UserStatusPending)).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, code.ErrUserNotExist
	}
	key := fmt.Sprintf("%s:%s", config.VerificationCode, xid.New().String())
	err = l.svcCtx.RedisClient.SetexCtx(l.ctx, key, req.Email,
		l.svcCtx.Config.CodeTime.VerificationCodeExpire)
	if err != nil {
		return nil, err
	}
	l.Debugf("send verification code, email: %s, code: %s", req.Email, key)
	go func() {
		if req.Event != "emailVerification" {
			return
		}
		ctx := context.WithoutCancel(l.ctx)
		emailCode := xid.New().String()
		key := fmt.Sprintf("%s:%s", config.VerificationCode, emailCode)
		err = l.svcCtx.RedisClient.SetexCtx(ctx, key, req.Email,
			l.svcCtx.Config.CodeTime.VerificationCodeExpire)
		if err != nil {
			l.Errorf("Failed to set email:%s verification code: %v", req.Email, err)
			return
		}
		emailData := EmailData{
			EmailVerificationLink: fmt.Sprintf("%s?code=%s&email=%s", l.svcCtx.Config.WebsiteUrl, emailCode, req.Email),
			EffectiveTime:         fmt.Sprintf("%d Hours", l.svcCtx.Config.CodeTime.VerificationCodeExpire/3600),
		}
		htmlBody, err := RenderEmailTemplate("template/email_verification_zh.html", emailData)
		if err != nil {
			l.Errorf("Failed to render email:%s template: %v", req.Email, err)
			return
		}
		err = utils.SendEmail(l.svcCtx.Config.EmailAccessKeyId, l.svcCtx.Config.EmailAccessSecret, &utils.SendMailRequest{
			AccountName: l.svcCtx.Config.EmailAccountName,
			ToAddress:   req.Email,
			Subject:     l.svcCtx.Config.EmailSubject,
			HtmlBody:    htmlBody,
		})
		if err != nil {
			l.Errorf("Failed to send email:%s: %v", req.Email, err)
			return
		}
		l.Debugf("Email:%s verification link: %s", req.Email, emailData.EmailVerificationLink)
	}()
	resp = &types.SendVerificationCodeResponse{}
	return
}
