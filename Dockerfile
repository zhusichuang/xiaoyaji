# 基于当前项目依赖版本，使用 Go 1.23 构建。
FROM golang:1.23-alpine3.20 AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo Asia/Shanghai > /etc/timezone

WORKDIR /app

COPY --from=builder /app/main /app/index.html /app/

CMD ["/app/main"]
