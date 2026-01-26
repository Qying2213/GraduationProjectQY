FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY common ./common
COPY recommendation-service ./recommendation-service

WORKDIR /build/recommendation-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o recommendation-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/recommendation-service/recommendation-service .

EXPOSE 8087
CMD ["./recommendation-service"]
