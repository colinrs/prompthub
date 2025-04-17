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

type ForgotPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForgotPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForgotPasswordLogic {
	return &ForgotPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForgotPasswordLogic) ForgotPassword(req *types.ForgotPasswordnRequest) (resp *types.ForgotPasswordnResponse, err error) {
	userTable := gen.Use(l.svcCtx.DB).UsersTable
	userInfo, err := userTable.Where(userTable.Email.Eq(req.Email), userTable.UserStatus.Eq(
		config.UserStatusNormal)).First()
	if err != nil {
		return nil, err
	}
	if userInfo == nil {
		return nil, code.ErrUserNotExist
	}
	err = l.ForgotPasswordEvent(req)
	if err != nil {
		return nil, err
	}
	return
}

func (l *ForgotPasswordLogic) ForgotPasswordEvent(req *types.ForgotPasswordnRequest) (err error) {
	ctx := context.WithoutCancel(l.ctx)
	emailCode := xid.New().String()
	key := fmt.Sprintf("%s:%s", config.VerificationCode, emailCode)
	err = l.svcCtx.RedisClient.SetexCtx(ctx, key, req.Email,
		l.svcCtx.Config.CodeTime.VerificationCodeExpire)
	if err != nil {
		l.Errorf("Failed to set email:%s verification code: %v", req.Email, err)
		return
	}
	key = fmt.Sprintf("%s:%s", config.SingleVerificationEmailLimit, req.Email)
	val, err := utils.IncrementAndSetTTLWithLua(
		context.WithoutCancel(l.ctx), l.svcCtx.RedisClient, key, 86400)
	if err != nil {
		l.Errorf("Failed to set single verification email:%s limit: %v", req.Email, err)
		return
	}
	if val > int64(l.svcCtx.Config.SingleVerificationEmailLimit) {
		l.Errorf("Single verification email:%s limit exceeded", req.Email)
		return code.ErrVerificationLimitExceed
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
		return code.ErrVerificationLimitExceed
	}
	emailData := ForgotPasswordData{
		VerificationCode: emailCode,
		EffectiveTime:    fmt.Sprintf("%d mins", l.svcCtx.Config.CodeTime.PasswordResetCodeExpire/60),
	}
	htmlBody, err := RenderEmailTemplate("template/email_code_zh.html", emailData)
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
	l.Debugf("Email:%s verification code: %s", req.Email, emailCode)
	return
}
