# 定义项目名称
PROJECT_NAME := prompthub.api

GOCTL := goctl

# 设置go模块环境变量
export GO111MODULE := on


# 格式化API定义文件
format:
	$(GOCTL) api format --dir ./
	$(GOCTL) api format --dir ./desc

# 生成Go API代码
gen-go-api:
	$(GOCTL) api go --api ./desc/$(PROJECT_NAME) --dir ./ --style=go_zero
	$(GOCTL) api plugin -plugin goctl-swagger="swagger -filename prompthub.json" -api ./desc/$(PROJECT_NAME) -dir swagger

# 编译项目
build:
	mkdir -p bin
	rm -rf bin/*
	go build -o bin/prompthub prompthub.go

lint:
	golangci-lint run --timeout=10m

sql:
	go run ./scripts/gen.go

# 默认目标
all: format gen-go-api

api: format gen-go-api

.PHONY: format gen-go-api lint