#!/bin/bash

# 设定 Goctl 的版本
GOCTL_VERSION="v1.6.3"

# 安装 Goctl 工具
echo "Installing goctl version $GOCTL_VERSION..."
go install github.com/zeromicro/go-zero/tools/goctl@$GOCTL_VERSION

# 检查环境配置
echo "Checking goctl environment..."
goctl env check --install --verbose --force

# 检查命令执行是否成功
if [ $? -eq 0 ]; then
    echo "Goctl installation and environment check completed successfully!"
else
    echo "There was an error during the goctl installation or environment check."
    exit 1
fi

# 安装 golangci-lint 工具
echo "Installing golangci-lint..."
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.1
echo "Golangci-lint installation completed successfully!"