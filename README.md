# 租房管理系统 (Rental Management System)

一个面向个人房东的租房管理后端系统，支持房源管理、租客管理、水电气收费、租金账单等功能。

## 功能特性

- **多租户架构**: 支持多个房东独立管理数据
- **房源管理**: 房源发布、编辑、上下架
- **房间管理**: 房间信息、状态管理
- **租客管理**: 租客档案、联系方式
- **合同管理**: 签订、续租、解约
- **抄表收费**: 水电气抄表、自动计算用量
- **账单管理**: 租金账单、费用统计
- **安全认证**: JWT认证、密码加密

## 技术栈

- Go 1.21+
- Gin Web Framework
- GORM ORM
- MySQL 8.0
- Redis 7.0
- JWT认证
- Swagger API文档

## 快速开始

### 环境要求 

- Go 1.21+
- MySQL 8.0+
- Redis 7.0+
- Docker & Docker Compose (可选)

### 使用Docker Compose启动

```bash
# 启动所有服务
make docker-up

# 查看日志
make docker-logs

# 停止服务
make docker-down
```

### 本地开发

1. 安装依赖
```bash
make deps
```

2. 创建数据库
```bash
mysql -u root -p < scripts/init_db.sql
```

3. 修改配置文件 `config.yaml`

4. 运行服务
```bash
make dev
```

### API测试

服务启动后访问:
- 健康检查: http://localhost:8080/health
- API文档: http://localhost:8080/swagger/index.html (需生成)

## API接口

### 认证模块
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/auth/profile` - 获取用户信息
- `PUT /api/v1/auth/profile` - 更新用户信息

### 房源管理
- `POST /api/v1/properties` - 创建房源
- `GET /api/v1/properties` - 房源列表
- `GET /api/v1/properties/:id` - 房源详情
- `PUT /api/v1/properties/:id` - 更新房源
- `DELETE /api/v1/properties/:id` - 删除房源

### 房间管理
- `POST /api/v1/rooms` - 创建房间
- `GET /api/v1/rooms` - 房间列表
- `GET /api/v1/rooms/:id` - 房间详情
- `PUT /api/v1/rooms/:id` - 更新房间
- `DELETE /api/v1/rooms/:id` - 删除房间

### 租客管理
- `POST /api/v1/tenants` - 添加租客
- `GET /api/v1/tenants` - 租客列表
- `PUT /api/v1/tenants/:id` - 更新租客

### 合同管理
- `POST /api/v1/contracts` - 签订合同
- `GET /api/v1/contracts` - 合同列表
- `PUT /api/v1/contracts/:id/terminate` - 解约

### 抄表与账单
- `POST /api/v1/meter-readings` - 录入抄表
- `GET /api/v1/meter-readings` - 抄表记录
- `POST /api/v1/bills` - 生成账单
- `GET /api/v1/bills` - 账单列表
- `PUT /api/v1/bills/:id/pay` - 标记已付
- `GET /api/v1/bills/statistics` - 收入统计

## 项目结构

```
rental-management/
├── cmd/server/main.go          # 入口文件
├── internal/
│   ├── config/                 # 配置管理
│   ├── middleware/             # 中间件
│   ├── model/                  # 数据模型
│   ├── handler/                # HTTP处理器
│   ├── service/                # 业务逻辑
│   ├── repository/             # 数据访问
│   └── pkg/                    # 内部工具
├── pkg/                        # 可导出包
├── migrations/                 # 数据库迁移
├── scripts/                    # 脚本
├── config.yaml                 # 配置文件
├── docker-compose.yml          # Docker编排
├── Dockerfile                  # Docker构建
└── Makefile                    # 构建脚本
```

## Make命令

```bash
make build      # 编译
make dev        # 开发运行
make test       # 运行测试
make clean      # 清理
make docker-up  # 启动Docker
make docker-down# 停止Docker
make swag       # 生成Swagger文档
make lint       # 代码检查
```

## 部署

### Docker部署

```bash
# 构建镜像
docker build -t rental-management:latest .

# 运行容器
docker run -d -p 8080:8080 rental-management:latest
```

### 生产环境配置

1. 修改 `config.yaml` 中的敏感配置
2. 使用环境变量覆盖配置
3. 配置HTTPS
4. 配置日志收集
5. 配置监控告警

## License

MIT
