#!/bin/bash
# 下载所有后端服务依赖

# 设置 Go 代理（国内加速）
export GOPROXY=https://goproxy.cn,direct

echo "========================================="
echo "下载后端服务依赖"
echo "========================================="

cd "$(dirname "$0")"

echo ""
echo "[1/7] common 模块..."
cd common && go mod tidy && cd ..

echo ""
echo "[2/7] user-service..."
cd user-service && go mod tidy && cd ..

echo ""
echo "[3/7] job-service..."
cd job-service && go mod tidy && cd ..

echo ""
echo "[4/7] talent-service..."
cd talent-service && go mod tidy && cd ..

echo ""
echo "[5/7] message-service..."
cd message-service && go mod tidy && cd ..

echo ""
echo "[6/7] interview-service..."
cd interview-service && go mod tidy && cd ..

echo ""
echo "[7/7] resume-service..."
cd resume-service && go mod tidy && cd ..

echo ""
echo "========================================="
echo "依赖下载完成！"
echo "========================================="
