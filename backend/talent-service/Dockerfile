FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY common ./common
COPY talent-service ./talent-service

WORKDIR /build/talent-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o talent-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/talent-service/talent-service .

EXPOSE 8086
CMD ["./talent-service"]
