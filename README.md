# 租房管理系统

一个完整的租房管理解决方案，前后端集成在一个容器中，支持 H5 和微信小程序。

## 项目结构

```
rental-management/
├── cmd/server/main.go          # 后端入口
├── internal/                   # 后端内部代码
│   ├── config/                 # 配置管理
│   ├── handler/                # HTTP处理器
│   ├── middleware/             # 中间件
│   ├── model/                  # 数据模型
│   ├── repository/             # 数据访问层
│   ├── service/                # 业务逻辑层
│   └── pkg/                    # 内部工具包
├── frontend/                   # 前端项目
│   ├── src/                    # 源代码
│   │   └── utils/
│   │       ├── config.js       # 微信云托管配置
│   │       └── request.js      # API请求封装
│   ├── nginx.conf              # Nginx配置
│   └── project.config.json     # 小程序配置
├── scripts/                    # 脚本文件
├── docker-compose.yml          # Docker编排文件
├── Dockerfile                  # 前后端一体化镜像
├── start.sh                    # 容器启动脚本
├── config.yaml                 # 配置文件
└── go.mod                      # Go模块定义
```

## 技术栈

### 后端
- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 7.0
- **认证**: JWT
- **文档**: Swagger
- **日志**: Zap

### 前端
- **框架**: uni-app (Vue 3)
- **UI**: 自适应移动端
- **支持**: H5 / 微信小程序

## 功能模块

- 用户认证（注册/登录）
- 房源管理（整栋/单套/商铺）
- 房间管理（支持月租/季租/年租）
- 租客管理
- 合同管理
- 抄表记录（水/电/气）
- 账单管理
- 收款记录
- 收入统计

---

## 部署方式

### 方式一：Docker Compose（推荐本地开发）

**一键启动所有服务：**

```bash
# 进入项目目录
cd rental-management

# 启动服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

**访问地址：**

| 服务 | 地址 |
|------|------|
| 前端 H5 | http://localhost |
| 后端 API | http://localhost:8080 |
| Swagger 文档 | http://localhost/swagger/index.html |

> 注意：前后端在同一容器中，Nginx 代理后端 API，所有请求可通过 80 端口访问。

**停止服务：**

```bash
# 停止服务
docker-compose down

# 停止并删除数据
docker-compose down -v
```

---

### 方式二：微信云托管部署（推荐生产环境）

#### 1. 微信云托管配置信息

```
服务ID：prod-6gomfrgz543fc0f3
服务名称：rental-managent
```

#### 2. 构建并推送镜像

```bash
# 构建镜像
docker build -t rental-management:latest .

# 标记镜像（替换为你的镜像仓库地址）
docker tag rental-management:latest ccr.ccs.tencentyun.com/你的命名空间/rental-management:latest

# 推送镜像
docker push ccr.ccs.tencentyun.com/你的命名空间/rental-management:latest
```

#### 3. 微信云托管控制台配置

| 配置项 | 推荐值 |
|--------|--------|
| 内存 | ≥ 512MB |
| CPU | 0.5核 |
| 最小实例数 | 1 |
| 最大实例数 | 10 |
| 健康检查路径 | `/health` |
| 健康检查端口 | 80 |

#### 4. 环境变量配置

在微信云托管控制台设置以下环境变量：

| 环境变量 | 说明 | 示例 |
|----------|------|------|
| DB_HOST | 数据库地址 | 你的MySQL地址 |
| DB_PORT | 数据库端口 | 3306 |
| DB_USERNAME | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | your_password |
| DB_DATABASE | 数据库名 | rental_management |
| REDIS_HOST | Redis地址 | 你的Redis地址 |
| REDIS_PORT | Redis端口 | 6379 |
| REDIS_PASSWORD | Redis密码 | your_redis_password |

#### 5. 前端配置说明

前端已配置微信云托管调用，配置文件位于 `frontend/src/utils/config.js`：

```javascript
// 微信云托管配置
export const cloudConfig = {
  serviceId: 'prod-6gomfrgz543fc0f3',
  serviceName: 'rental-managent',
  // ...
}
```

**调用方式：**
- **H5 部署**：使用相对路径 `/api/v1`，通过 Nginx 代理到后端
- **微信小程序**：使用服务名调用 `http://rental-managent/api/v1`

---

### 方式三：微信小程序单独部署

```bash
cd frontend

# 安装依赖
npm install

# 构建微信小程序
npm run build:mp-weixin
```

构建产物在 `dist/build/mp-weixin/`，使用微信开发者工具导入并发布。

**小程序调用云托管服务：**

小程序会自动通过服务名调用同一环境下的后端服务：
```javascript
// 微信小程序环境自动使用服务名调用
const url = 'http://rental-managent/api/v1/...'
```

---

## 本地开发

### 后端开发

```bash
# 安装依赖
go mod download

# 启动服务（需要先启动 MySQL 和 Redis）
go run cmd/server/main.go

# 或者使用 Makefile
make dev
```

### 前端开发

```bash
cd frontend

# 安装依赖
npm install

# H5 开发模式
npm run dev:h5

# 微信小程序开发模式
npm run dev:mp-weixin
```

---

## 配置说明

### 后端配置 (config.yaml)

```yaml
server:
  port: 8080
  mode: release    # debug/release

database:
  host: localhost
  port: 3306
  database: rental_management
  username: root
  password: root123456
  charset: utf8mb4

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-secret-key
  issuer: rental-management
  expire_hours: 24

log:
  level: info
  format: json
  output: stdout
```

### 环境变量覆盖

支持以下环境变量覆盖配置：

| 环境变量 | 说明 |
|----------|------|
| DB_HOST | 数据库地址 |
| DB_PORT | 数据库端口 |
| DB_USERNAME | 数据库用户名 |
| DB_PASSWORD | 数据库密码 |
| REDIS_HOST | Redis地址 |
| REDIS_PORT | Redis端口 |
| REDIS_PASSWORD | Redis密码 |

---

## API 文档

启动服务后访问 Swagger 文档：
- http://localhost/swagger/index.html（通过 Nginx 代理）
- http://localhost:8080/swagger/index.html（直接访问后端）

### 主要接口

| 模块 | 接口 | 说明 |
|------|------|------|
| 认证 | POST /api/v1/auth/register | 用户注册 |
| 认证 | POST /api/v1/auth/login | 用户登录 |
| 房源 | GET /api/v1/properties | 房源列表 |
| 房源 | POST /api/v1/properties | 创建房源 |
| 房间 | GET /api/v1/rooms | 房间列表 |
| 房间 | POST /api/v1/rooms | 创建房间 |
| 租客 | GET /api/v1/tenants | 租客列表 |
| 合同 | POST /api/v1/contracts | 签订合同 |
| 抄表 | POST /api/v1/meter-readings | 录入抄表 |
| 账单 | GET /api/v1/bills | 账单列表 |
| 账单 | GET /api/v1/bills/statistics | 收入统计 |

---

## Make 命令

```bash
make build      # 编译
make dev        # 开发运行
make test       # 运行测试
make clean      # 清理
make docker-up  # 启动Docker
make docker-down# 停止Docker
make swag       # 生成Swagger文档
```

---

## License

MIT
