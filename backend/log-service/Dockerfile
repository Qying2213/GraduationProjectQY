FROM golang:1.23-alpine AS builder

WORKDIR /app

# 复制 common 模块
COPY common/ /app/common/

# 复制 log-service 模块
COPY log-service/ /app/log-service/

WORKDIR /app/log-service

# 下载依赖
RUN go mod download

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/log-service/main .

EXPOSE 8088

CMD ["./main"]
