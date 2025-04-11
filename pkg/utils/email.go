package utils

import (
	"bytes"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm20151123 "github.com/alibabacloud-go/dm-20151123/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/zeromicro/go-zero/core/logx"
	"html/template"
	"path/filepath"
)

func createEmailClient(accessKeyId, accessKeySecret string) (client *dm20151123.Client, err error) {
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dm
	config.Endpoint = tea.String("dm.aliyuncs.com")
	client = &dm20151123.Client{}
	client, err = dm20151123.NewClient(config)
	return client, err
}

type SendMailRequest struct {
	AccountName string `json:"accountName"`
	ToAddress   string `json:"toAddress"`
	Subject     string `json:"subject"`
	HtmlBody    string `json:"htmlBody"`
}

func SendEmail(accessKeyId, accessKeySecret string, sendMailRequest *SendMailRequest) (err error) {
	client, err := createEmailClient(accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	singleSendMailRequest := &dm20151123.SingleSendMailRequest{
		AccountName:            tea.String(sendMailRequest.AccountName),
		AddressType:            tea.Int32(1),
		TagName:                tea.String(""),
		ReplyToAddress:         tea.Bool(false),
		ToAddress:              tea.String(sendMailRequest.ToAddress),
		Subject:                tea.String(sendMailRequest.Subject),
		HtmlBody:               tea.String(sendMailRequest.HtmlBody),
		TextBody:               tea.String(""),
		FromAlias:              tea.String(""),
		ReplyAddress:           tea.String(""),
		ReplyAddressAlias:      tea.String(""),
		ClickTrace:             tea.String(""),
		UnSubscribeLinkType:    tea.String(""),
		UnSubscribeFilterLevel: tea.String(""),
		Headers:                tea.String(""),
	}
	runtime := &util.RuntimeOptions{}
	_, err = client.SingleSendMailWithOptions(singleSendMailRequest, runtime)
	if err != nil {
		logx.Errorf("Failed to send email:%s: %v", sendMailRequest.ToAddress, err)
		return err
	}
	return nil
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
