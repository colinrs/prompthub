# 第一阶段：构建阶段
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置Go环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN go build -o prompthub .

# 第二阶段：运行阶段
FROM alpine:latest

# 安装基础工具和SSL证书
RUN apk --no-cache add ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/prompthub /app/
# 复制配置文件
COPY --from=builder /app/etc/prompthub-api.yaml /app/etc/
# 复制模板文件
COPY --from=builder /app/template /app/template

# 暴露端口（根据配置文件中的端口设置）
EXPOSE 8080

# 启动应用
CMD ["/app/prompthub", "-f", "/app/etc/prompthub-api.yaml"]