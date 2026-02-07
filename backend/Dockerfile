FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装 CA 证书
RUN apk --no-cache add ca-certificates

# 复制编译后的二进制文件
COPY --from=builder /app/server .

# 创建静态文件目录
RUN mkdir -p /app/static

EXPOSE 8080

CMD ["./server"]
