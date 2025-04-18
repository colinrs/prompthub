#!/bin/bash

# 检查是否提供了tag参数
if [ -z "$1" ]; then
    echo "错误: 请提供镜像标签作为参数"
    echo "用法: $0 <tag>"
    exit 1
fi

# 设置变量
TAG=$1
REGISTRY="crpi-uuv5fw7uy2j5kw76.cn-heyuan.personal.cr.aliyuncs.com"
NAMESPACE="feitiansh"
REPO="prompthub"
IMAGE_NAME="${REGISTRY}/${NAMESPACE}/${REPO}:${TAG}"

echo "aaaa $ALIYUNCS_REGISTRY_PASSSWD"
# 登录到阿里云容器镜像服务
echo "登录到阿里云容器镜像服务..."
docker login --username=947834020@qq.com --password=$ALIYUNCS_REGISTRY_PASSSWD ${REGISTRY} || {
    echo "错误: 登录失败"
    exit 1
}

# 构建Docker镜像
echo "构建Docker镜像..."
docker build -t ${IMAGE_NAME} . || {
    echo "错误: 镜像构建失败"
    exit 1
}

# 推送镜像到阿里云容器镜像服务
echo "推送镜像到阿里云容器镜像服务..."
docker push ${IMAGE_NAME} || {
    echo "错误: 镜像推送失败"
    exit 1
}

echo "成功: 镜像已构建并推送到阿里云容器镜像服务"
echo "镜像地址: ${IMAGE_NAME}"
