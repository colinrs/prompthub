package gosafe

import (
	"github.com/zeromicro/go-zero/core/metric"
)

// Metrics instance
var (
	goSafePanicTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "go_pkg",
		Subsystem: "go_panic",
		Name:      "total",
		Help:      "go safe panic custom error count.",
	})
)
