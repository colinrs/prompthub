package svc

import (
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/infra"
	"github.com/colinrs/prompthub/internal/middleware"

	swd "github.com/kirklin/go-swd"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config                 config.Config
	DB                     *gorm.DB
	RedisClient            *redis.Redis
	UserNonLoginMiddleware rest.Middleware
	UserLoginMiddleware    rest.Middleware
	DetectorSWD            *swd.SWD
}

func NewServiceContext(c config.Config) *ServiceContext {
	detector, err := swd.New()
	if err != nil {
		logx.Must(err)
	}
	return &ServiceContext{
		Config:                 c,
		DB:                     initDB(c),
		RedisClient:            initRedis(c),
		UserNonLoginMiddleware: middleware.NewUserNonLoginMiddleware(c).Handle,
		UserLoginMiddleware:    middleware.NewUserLoginMiddleware(c).Handle,
		DetectorSWD:            detector,
	}

}

func initDB(c config.Config) *gorm.DB {
	db, err := infra.Database(c.DataBase)
	logx.Must(err)
	return db
}

func initRedis(c config.Config) *redis.Redis {
	return redis.MustNewRedis(c.Redis)
}
