#!/bin/sh

# 启动脚本 - 同时启动前端(Nginx)和后端(Go)服务

echo "Starting rental management system..."

# 初始化证书（微信云托管需要）
if [ -f /app/cert/initenv.sh ]; then
    /bin/sh /app/cert/initenv.sh
fi

# 启动后端服务
echo "Starting backend on port 8080..."
/app/backend &
BACKEND_PID=$!

# 启动Nginx
echo "Starting frontend on port 80..."
nginx
NGINX_PID=$!

# 等待任意进程退出
wait -n $BACKEND_PID $NGINX_PID

# 如果有一个进程退出，杀死另一个
kill $BACKEND_PID $NGINX_PID 2>/dev/null

exit 1
