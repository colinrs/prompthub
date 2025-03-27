package svc

import (
	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/infra"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	RedisClient *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		DB:          initDB(c),
		RedisClient: initRedis(c),
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
