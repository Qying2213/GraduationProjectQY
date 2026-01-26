# 一个 Dockerfile 构建整个后端项目
FROM golang:1.23-alpine AS builder

WORKDIR /build

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

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

# 复制所有编译好的二进制文件
COPY --from=builder /out/* ./

# 复制启动脚本
COPY start-services.sh ./
RUN chmod +x start-services.sh

EXPOSE 8080 8081 8082 8083 8084 8085 8086 8087

CMD ["./start-services.sh"]
