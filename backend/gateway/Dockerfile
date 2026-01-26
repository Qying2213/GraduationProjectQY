FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

# 复制 common 模块和当前服务
COPY common ./common
COPY gateway ./gateway

WORKDIR /build/gateway

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/gateway/gateway .

EXPOSE 8080
CMD ["./gateway"]
