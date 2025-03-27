package infra

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type DBConfig struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	UserName string `json:"UserName"`
	Password string `json:"Password"`
	Database string `json:"Database"`

	MaxIdleConn     int `json:"MaxIdleConn"`
	MaxOpenConn     int `json:"MaxOpenConn"`
	ConnMaxLifeTime int `json:"ConnMaxLifeTime"`
}

// Database ...
func Database(mysqlConfig *DBConfig) (*gorm.DB, error) {
	logx.Info("mysql {%#v}", mysqlConfig)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlConfig.UserName,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Database)
	logx.Info("connect to mysql %s", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
		Plugins: map[string]gorm.Plugin{
			metricsName: NewGormMetricsPlugin(WithDataBaseName(mysqlConfig.Database)),
		},
	})
	if err != nil {
		logx.Errorf("conect db err:%v", err)
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		logx.Errorf("get sqlDB err:%v", err)
		return nil, err
	}

	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConn)
	//打开
	sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConn)
	//超时
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(mysqlConfig.ConnMaxLifeTime))

	return db, nil
}
