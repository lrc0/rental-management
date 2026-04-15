# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 设置Go代理
ENV GOPROXY=https://goproxy.cn,direct

# 安装依赖
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 安装swag工具
RUN go install github.com/swaggo/swag/cmd/swag@latest

# 复制源代码
COPY . .

# 生成swagger文档
RUN /go/bin/swag init -g cmd/server/main.go -o api/swagger

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# 安装ca证书
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 从builder复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

# 暴露端口
EXPOSE 8080

# 运行
CMD ["./main"]
