# Build stage - 后端
FROM golang:1.23-alpine AS backend-builder

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

# 构建后端
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backend ./cmd/server

# Build stage - 前端
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# 设置npm镜像
RUN npm config set registry https://registry.npmmirror.com

# 复制前端package文件
COPY frontend/package.json frontend/package-lock.json ./

# 安装依赖
RUN npm install

# 复制前端源代码
COPY frontend/ .

# 构建H5版本
RUN npm run build:h5

# Runtime stage
FROM alpine:latest

WORKDIR /app

# 安装必要工具
RUN apk add --no-cache nginx ca-certificates tzdata && \
    mkdir -p /run/nginx /app/cert /var/log/nginx /usr/share/nginx/html && \
    touch /app/cert/certificate.crt

# 设置时区
ENV TZ=Asia/Shanghai

# 从后端构建阶段复制二进制文件
COPY --from=backend-builder /app/backend /app/backend
COPY --from=backend-builder /app/config.yaml /app/config.yaml
COPY --from=backend-builder /app/api/swagger /app/api/swagger

# 从前端构建阶段复制静态文件
COPY --from=frontend-builder /app/frontend/dist/build/h5 /usr/share/nginx/html

# 复制nginx配置
COPY frontend/nginx.conf /etc/nginx/http.d/default.conf

# 复制启动脚本
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh

# 暴露端口
EXPOSE 80 8080

# 启动
CMD ["/app/start.sh"]
