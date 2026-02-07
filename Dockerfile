# 构建阶段
FROM golang:1.24-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 复制源代码
COPY backend/ ./

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags '-extldflags "-static" -s -w' \
    -o novelforge cmd/server/main.go

# 运行阶段
FROM alpine:latest

# 安装 CA 证书和时区数据
RUN apk --no-cache add ca-certificates tzdata

# 创建非 root 用户
RUN addgroup -g 1000 novelforge && \
    adduser -D -u 1000 -G novelforge novelforge

# 设置工作目录
WORKDIR /app

# 复制二进制文件
COPY --from=builder /app/novelforge .

# 复制前端文件
COPY frontend/ ./frontend/

# 设置权限
RUN chown -R novelforge:novelforge /app

# 切换到非 root 用户
USER novelforge

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# 启动应用
CMD ["./novelforge"]
