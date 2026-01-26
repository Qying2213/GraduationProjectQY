FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY common ./common
COPY user-service ./user-service

WORKDIR /build/user-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/user-service/user-service .

EXPOSE 8081
CMD ["./user-service"]
