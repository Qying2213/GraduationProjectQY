FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY common ./common
COPY message-service ./message-service

WORKDIR /build/message-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o message-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/message-service/message-service .

EXPOSE 8085
CMD ["./message-service"]
