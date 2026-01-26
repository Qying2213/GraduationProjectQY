#!/bin/sh

echo "启动所有服务..."

# 后台启动所有服务
./user-service &
./job-service &
./talent-service &
./message-service &
./interview-service &
./resume-service &
./recommendation-service &

# 前台启动 gateway（保持容器运行）
./gateway
