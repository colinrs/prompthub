package infra

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/metric"
	"gorm.io/gorm"
)

const (
	gormClientNamespace = "gormc"
	metricsName         = "gormc_metrics_plugin"
)

const (
	CmdQuery  = "SELECT"
	CmdUpdate = "UPDATE"
	CmdInsert = "INSERT"
	CmdDelete = "DELETE"
	Unknown   = "unknown"
)

type ctxKey struct{}

var startTimeKey = ctxKey{}

type GormMetricsPlugin struct {
	queryCounter     metric.CounterVec
	errorCounter     metric.CounterVec
	latencyHistogram metric.HistogramVec

	tableNameRewriter func(string) string
	dataBaseName      string
}

type GormMetricsOption func(*GormMetricsPlugin)

func NewGormMetricsPlugin(options ...GormMetricsOption) *GormMetricsPlugin {
	g := &GormMetricsPlugin{
		queryCounter: metric.NewCounterVec(&metric.CounterVecOpts{
			Namespace: gormClientNamespace,
			Subsystem: "requests",
			Name:      "total",
			Help:      "GORM client requests count.",
			Labels:    []string{"db", "table", "operation"},
		}),
		errorCounter: metric.NewCounterVec(&metric.CounterVecOpts{
			Namespace: gormClientNamespace,
			Subsystem: "requests",
			Name:      "err_total",
			Help:      "GORM client requests err code count.",
			Labels:    []string{"db", "table", "operation"},
		}),
		latencyHistogram: metric.NewHistogramVec(&metric.HistogramVecOpts{
			Namespace: gormClientNamespace,
			Subsystem: "requests",
			Name:      "duration_ms",
			Help:      "GORM client requests duration(ms).",
			Labels:    []string{"db", "table", "operation"},
			Buckets:   []float64{0.25, 0.5, 1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000, 10000, 15000},
		}),
		tableNameRewriter: func(s string) string {
			return s
		},
		dataBaseName: "",
	}
	for _, option := range options {
		option(g)
	}
	return g
}

func (p *GormMetricsPlugin) Name() string {
	return metricsName
}

func (p *GormMetricsPlugin) Initialize(db *gorm.DB) error {
	createCallbacksName := "create"
	err := db.Callback().Create().Before(fmt.Sprintf("gorm:%s", createCallbacksName)).
		Register(fmt.Sprintf("%s:before_%s", metricsName, createCallbacksName), p.before)
	if err != nil {
		return err
	}
	err = db.Callback().Create().After(fmt.Sprintf("gorm:%s", createCallbacksName)).
		Register(fmt.Sprintf("%s:after_%s", metricsName, createCallbacksName), p.after)
	if err != nil {
		return err
	}
	queryCallbacksName := "query"
	err = db.Callback().Query().Before(fmt.Sprintf("gorm:%s", queryCallbacksName)).
		Register(fmt.Sprintf("%s:before_%s", metricsName, queryCallbacksName), p.before)
	if err != nil {
		return err
	}
	err = db.Callback().Query().After(fmt.Sprintf("gorm:%s", queryCallbacksName)).
		Register(fmt.Sprintf("%s:after_%s", metricsName, queryCallbacksName), p.after)
	if err != nil {
		return err
	}
	updateCallbacksName := "update"
	err = db.Callback().Update().Before(fmt.Sprintf("gorm:%s", updateCallbacksName)).
		Register(fmt.Sprintf("%s:before_%s", metricsName, updateCallbacksName), p.before)
	if err != nil {
		return err
	}
	err = db.Callback().Update().After(fmt.Sprintf("gorm:%s", updateCallbacksName)).
		Register(fmt.Sprintf("%s:after_%s", metricsName, updateCallbacksName), p.after)
	if err != nil {
		return err
	}
	deleteCallbacksName := "delete"
	err = db.Callback().Delete().Before(fmt.Sprintf("gorm:%s", deleteCallbacksName)).
		Register(fmt.Sprintf("%s:before_%s", metricsName, deleteCallbacksName), p.before)
	if err != nil {
		return err
	}
	err = db.Callback().Delete().After(fmt.Sprintf("gorm:%s", deleteCallbacksName)).
		Register(fmt.Sprintf("%s:after_%s", metricsName, deleteCallbacksName), p.after)
	if err != nil {
		return err
	}
	rowCallbacksName := "row"
	err = db.Callback().Row().Before(fmt.Sprintf("gorm:%s", rowCallbacksName)).
		Register(fmt.Sprintf("%s:before_%s", metricsName, rowCallbacksName), p.before)
	if err != nil {
		return err
	}
	err = db.Callback().Row().After(fmt.Sprintf("gorm:%s", rowCallbacksName)).
		Register(fmt.Sprintf("%s:after_%s", metricsName, rowCallbacksName), p.after)
	if err != nil {
		return err
	}
	rawCallbacksName := "raw"
	err = db.Callback().Raw().Before(fmt.Sprintf("gorm:%s", rawCallbacksName)).
		Register(fmt.Sprintf("%s:before_%s", metricsName, rawCallbacksName), p.before)
	if err != nil {
		return err
	}
	err = db.Callback().Raw().After(fmt.Sprintf("gorm:%s", rawCallbacksName)).
		Register(fmt.Sprintf("%s:after_%s", metricsName, rawCallbacksName), p.after)
	if err != nil {
		return err
	}
	return nil
}

func (p *GormMetricsPlugin) before(db *gorm.DB) {
	ctx := db.Statement.Context
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = context.WithValue(ctx, startTimeKey, time.Now())
	db.Statement.Context = ctx
}

func (p *GormMetricsPlugin) after(db *gorm.DB) {
	ctx := db.Statement.Context
	startTime, ok := ctx.Value(startTimeKey).(time.Time)
	if !ok {
		// 如果没有找到开始时间，就使用当前时间（这种情况不应该发生）
		startTime = time.Now()
	}
	duration := time.Since(startTime)

	operation := matchOperationCommand(db.Statement.SQL.String())
	dbName := db.Name()
	if p.dataBaseName != "" {
		dbName = p.dataBaseName
	}
	tableName := p.tableNameRewriter(db.Statement.Table)
	p.queryCounter.Inc(dbName, tableName, operation)
	p.latencyHistogram.ObserveFloat(float64(duration.Milliseconds()), dbName, tableName, operation)

	if db.Error != nil {
		p.errorCounter.Inc(dbName, tableName, operation)
	}
}

func (p *GormMetricsPlugin) Close() error {
	return nil
}

func matchOperationCommand(sql string) string {
	if len(sql) > 6 {
		cmd := strings.ToUpper(sql[:6])
		switch cmd {
		case CmdQuery, CmdUpdate, CmdInsert, CmdDelete:
			return cmd
		}
	}
	return Unknown
}

func (p *GormMetricsPlugin) TranslateError(err error) error {
	return err
}

func WithTableNameRewriter(rewriter func(string) string) GormMetricsOption {
	return func(p *GormMetricsPlugin) {
		p.tableNameRewriter = rewriter
	}
}

func WithDataBaseName(name string) GormMetricsOption {
	return func(p *GormMetricsPlugin) {
		p.dataBaseName = name
	}
}
