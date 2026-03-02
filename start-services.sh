#!/bin/sh
# ============================================================
# Docker 容器内服务启动脚本
# ============================================================

echo "=========================================="
echo "  智能人才运营平台 - 启动所有服务"
echo "=========================================="

# 等待数据库就绪
echo "[INFO] 等待数据库连接..."
sleep 5

# 后台启动所有微服务
echo "[INFO] 启动 User Service (8081)..."
./user-service &

echo "[INFO] 启动 Job Service (8082)..."
./job-service &

echo "[INFO] 启动 Interview Service (8083)..."
./interview-service &

echo "[INFO] 启动 Resume Service (8084)..."
./resume-service &

echo "[INFO] 启动 Message Service (8085)..."
./message-service &

echo "[INFO] 启动 Talent Service (8086)..."
./talent-service &

echo "[INFO] 启动 Recommendation Service (8087)..."
./recommendation-service &

echo "[INFO] 启动 Log Service (8088)..."
./log-service &

echo "[INFO] 启动 Evaluator Service (8090)..."
./evaluator-service &

# 等待微服务启动
sleep 3

# 前台启动 Gateway（保持容器运行）
echo "[INFO] 启动 Gateway (8080)..."
echo "=========================================="
echo "  所有服务已启动！"
echo "=========================================="
./gateway
