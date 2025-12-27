FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM python:3.11-alpine
WORKDIR /app
RUN pip install --no-cache-dir requests requests-toolbelt pycryptodome urllib3
COPY --from=builder /app/server .
COPY --from=builder /app/internal/script/*.py ./internal/script/
COPY --from=builder /app/static ./static
RUN mkdir -p data
EXPOSE 8090
CMD ["./server"]
