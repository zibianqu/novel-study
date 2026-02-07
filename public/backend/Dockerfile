FROM golang:1.22-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git ca-certificates

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o /novelforge ./cmd/server

# 运行镜像
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /novelforge .
COPY config.yaml .

# 前端静态文件将通过volume挂载到 /app/static

EXPOSE 8080

CMD ["./novelforge"]
