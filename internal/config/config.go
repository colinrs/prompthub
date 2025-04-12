package config

import (
	"github.com/colinrs/prompthub/internal/infra"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DataBase                     *infra.DBConfig `json:"Database" yaml:"Database"`
	Redis                        redis.RedisConf `json:"Redis" yaml:"Redis"`
	PasswdSecret                 string          `json:"PasswdSecret" yaml:"PasswdSecret"`
	JwtSecret                    string          `json:"JwtSecret" yaml:"JwtSecret"`
	JwtExpired                   int64           `json:"JwtExpired" yaml:"JwtExpired"`
	CodeTime                     CodeTime        `json:"CodeTime" yaml:"CodeTime"`
	WebsiteUrl                   string          `json:"WebsiteUrl" yaml:"WebsiteUrl"`
	EmailAccessKeyId             string          `json:",env=EMAIL_ACCESS_KEY_ID"`
	EmailAccessSecret            string          `json:",env=EMAIL_ACCESS_SECRET"`
	EmailAccountName             string          `json:",env=EMAIL_ACCOUNT_NAME"`
	EmailSubject                 string          `json:"EmailSubject" yaml:"EmailSubject"`
	VerificationEmailLimit       int             `json:"VerificationEmailLimit" yaml:"VerificationEmailLimit"`
	SingleVerificationEmailLimit int             `json:"SingleVerificationEmailLimit" yaml:"SingleVerificationEmailLimit"`
}

type CodeTime struct {
	VerificationCodeExpire  int
	PasswordResetCodeExpire int
}
