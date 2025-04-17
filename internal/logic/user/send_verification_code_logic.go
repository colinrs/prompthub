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
	"github.com/rs/xid"
	"github.com/spf13/cast"
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
	ok, err := l.checkEmailLimit(req.Email)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, code.ErrVerificationLimitExceed
	}
	userTable := gen.Use(l.svcCtx.DB).UsersTable
	userInfo, err := userTable.Where(userTable.Email.Eq(req.Email), userTable.UserStatus.In(
		config.UserStatusNormal, config.UserStatusPending)).First()
	if err != nil {
		return nil, err
	}
	if userInfo == nil {
		return nil, code.ErrUserNotExist
	}
	if userInfo.UserStatus == config.UserStatusNormal {
		return nil, code.ErrUserExist
	}
	key := fmt.Sprintf("%s:%s", config.VerificationCode, xid.New().String())
	err = l.svcCtx.RedisClient.SetexCtx(l.ctx, key, req.Email,
		l.svcCtx.Config.CodeTime.VerificationCodeExpire)
	if err != nil {
		return nil, err
	}
	l.Debugf("send verification code, email: %s, code: %s", req.Email, key)
	go func() {
		key = fmt.Sprintf("%s:%s", config.SingleVerificationEmailLimit, req.Email)
		val, err := utils.IncrementAndSetTTLWithLua(
			context.WithoutCancel(l.ctx), l.svcCtx.RedisClient, key, 86400)
		if err != nil {
			l.Errorf("Failed to set single verification email:%s limit: %v", req.Email, err)
			return
		}
		if val > int64(l.svcCtx.Config.SingleVerificationEmailLimit) {
			l.Errorf("Single verification email:%s limit exceeded", req.Email)
		}
		key = fmt.Sprintf("%s:%s", config.VerificationEmailLimitKey, req.Email)
		val, err = utils.IncrementAndSetTTLWithLua(
			context.WithoutCancel(l.ctx), l.svcCtx.RedisClient, key, 86400)
		if err != nil {
			l.Errorf("Failed to set all verification email:%s limit: %v", req.Email, err)
			return
		}
		if val > int64(l.svcCtx.Config.SingleVerificationEmailLimit) {
			l.Errorf("All verification email:%s limit exceeded", req.Email)
			return
		}
		switch req.Event {
		case config.EmailVerificationEvent:
			err = l.EmailVerificationEvent(req)
			if err != nil {
				l.Errorf("Failed to set verification email:%s verification code: %v", req.Email, err)
			}
		}
	}()
	resp = &types.SendVerificationCodeResponse{}
	return
}

func (l *SendVerificationCodeLogic) checkEmailLimit(email string) (bool, error) {
	key := fmt.Sprintf("%s:%s", config.SingleVerificationEmailLimit, email)
	limitString, err := l.svcCtx.RedisClient.GetCtx(l.ctx, key)
	if err != nil {
		return false, err
	}
	limit, err := cast.ToIntE(limitString)
	if err != nil {
		return false, err
	}

	if limit > l.svcCtx.Config.SingleVerificationEmailLimit {
		return false, code.ErrVerificationLimitExceed
	}
	limitString, err = l.svcCtx.RedisClient.GetCtx(l.ctx, config.VerificationEmailLimitKey)
	if err != nil {
		return false, err
	}
	limit, err = cast.ToIntE(limitString)
	if err != nil {
		return false, err
	}
	if limit > l.svcCtx.Config.VerificationEmailLimit {
		return false, code.ErrVerificationLimitExceed
	}
	return true, nil

}

func (l *SendVerificationCodeLogic) EmailVerificationEvent(req *types.SendVerificationCodeRequest) (err error) {
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
	return
}
