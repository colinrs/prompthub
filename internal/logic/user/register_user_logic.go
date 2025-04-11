package user

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rs/xid"
	"html/template"
	"path/filepath"

	"github.com/colinrs/prompthub/gen"
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/svc"
	"github.com/colinrs/prompthub/internal/types"
	"github.com/colinrs/prompthub/model"
	"github.com/colinrs/prompthub/pkg/code"
	"github.com/colinrs/prompthub/pkg/utils"

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
		UserStatus:  config.UserStatusPending,
		Description: "",
	}
	err = userTable.Save(value)
	if err != nil {
		return resp, err
	}
	go func() {
		emailCode := xid.New().String()
		key := fmt.Sprintf("%s:%s", config.VerificationCode, emailCode)
		err = l.svcCtx.RedisClient.SetexCtx(l.ctx, key, req.Email,
			l.svcCtx.Config.CodeTime.VerificationCodeExpire)
		if err != nil {
			return
		}
		emailData := EmailData{
			EmailVerificationLink: fmt.Sprintf("%s?code=%s&email=%s", l.svcCtx.Config.WebsiteUrl, emailCode, req.Email),
			EffectiveTime:         fmt.Sprintf("%d minutes", l.svcCtx.Config.CodeTime.VerificationCodeExpire/60),
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
	return resp, nil
}

// EmailData 定义模板中需要的数据结构
type EmailData struct {
	EmailVerificationLink string `json:"EmailVerificationLink"`
	EffectiveTime         string `json:"EffectiveTime"`
}

// RenderEmailTemplate 渲染 HTML 模板
func RenderEmailTemplate(templatePath string, data EmailData) (string, error) {
	absPath, err := filepath.Abs(templatePath)
	if err != nil {
		return "", fmt.Errorf("resolve absolute path: %w", err)
	}

	// 解析模板文件
	tmpl, err := template.ParseFiles(absPath)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	// 创建缓冲区来存储渲染后的 HTML
	var buf bytes.Buffer

	// 执行模板渲染
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	// 返回渲染后的 HTML 字符串
	return buf.String(), nil
}
