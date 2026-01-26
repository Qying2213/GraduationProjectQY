FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY common ./common
COPY job-service ./job-service

WORKDIR /build/job-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o job-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/job-service/job-service .

EXPOSE 8082
CMD ["./job-service"]
