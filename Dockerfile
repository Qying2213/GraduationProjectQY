# ============================================================
# 智能人才运营平台 - 后端服务 Dockerfile
# 一个容器运行所有后端微服务
# ============================================================

FROM golang:1.23-alpine AS builder

WORKDIR /build

# 安装必要工具
RUN apk add --no-cache git

# 复制整个 backend 目录
COPY backend ./backend

# 构建所有服务
RUN cd backend/gateway && go build -o /out/gateway .
RUN cd backend/user-service && go build -o /out/user-service .
RUN cd backend/job-service && go build -o /out/job-service .
RUN cd backend/talent-service && go build -o /out/talent-service .
RUN cd backend/message-service && go build -o /out/message-service .
RUN cd backend/interview-service && go build -o /out/interview-service .
RUN cd backend/resume-service && go build -o /out/resume-service .
RUN cd backend/recommendation-service && go build -o /out/recommendation-service .
RUN cd backend/log-service && go build -o /out/log-service .
RUN cd backend/evaluator-service && go build -o /out/evaluator-service ./cmd/server

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata python3 py3-pip
RUN pip3 install --no-cache-dir requests requests-toolbelt pycryptodome urllib3
ENV TZ=Asia/Shanghai

WORKDIR /app

# 复制所有编译好的二进制文件
COPY --from=builder /out/* ./
COPY --from=builder /build/backend/evaluator-service/static ./static
COPY --from=builder /build/backend/evaluator-service/internal/script ./internal/script

# 复制启动脚本
COPY start-services.sh ./
RUN chmod +x start-services.sh

# 创建日志和上传目录
RUN mkdir -p /app/logs /app/uploads /app/data

# 暴露所有服务端口
# 8080: Gateway
# 8081: User Service
# 8082: Job Service
# 8083: Interview Service
# 8084: Resume Service
# 8085: Message Service
# 8086: Talent Service
# 8087: Recommendation Service
# 8088: Log Service
# 8090: Evaluator Service
EXPOSE 8080 8081 8082 8083 8084 8085 8086 8087 8088 8090

CMD ["./start-services.sh"]
