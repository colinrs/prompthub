#!/bin/bash

# 安装 golangci-lint 工具
echo "Installing golangci-lint..."
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.1
echo "Golangci-lint installation completed successfully!"