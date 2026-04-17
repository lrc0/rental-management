# 租房管理系统

一个完整的租房管理解决方案，支持 H5 和微信小程序。

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
├── frontend/                   # 前端项目 (uni-app)
│   ├── src/                    # 源代码
│   │   └── utils/
│   │       ├── config.js       # 微信云托管配置
│   │       └── request.js      # API请求封装
│   ├── dist/build/mp-weixin/   # 小程序构建产物
│   └── project.config.json     # 小程序配置
├── scripts/init_db.sql         # 数据库初始化脚本
├── docker-compose.yml          # Docker编排文件
├── Dockerfile                  # Docker镜像构建
├── config.yaml                 # 后端配置文件
└── go.mod                      # Go模块定义
```

## 技术栈

### 后端
- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 7.0（可选）
- **认证**: JWT
- **文档**: Swagger

### 前端
- **框架**: uni-app (Vue 3)
- **支持**: H5 / 微信小程序

---

## 本地调试

### 前置条件

1. 安装 Go 1.21+
2. 安装 Node.js 18+
3. 安装 MySQL 8.0
4. 安装微信开发者工具（调试小程序）

### 1. 数据库准备

```bash
# 创建数据库
mysql -u root -p

# 在 MySQL 中执行
CREATE DATABASE rental_management DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 或者使用初始化脚本
mysql -u root -p < scripts/init_db.sql
```

### 2. 后端配置

修改 `config.yaml`：

```yaml
server:
  port: 8080
  mode: debug    # debug 模式便于调试

database:
  host: localhost
  port: 3306
  database: rental_management
  username: root
  password: 你的密码
  charset: utf8mb4

jwt:
  secret: your-secret-key
  issuer: rental-management
  expire_hours: 24
```

### 3. 启动后端

```bash
# 进入项目目录
cd /Users/ruicai.li/go/src/AIAgent/Dev-AI/rental-management

# 下载依赖
go mod download

# 启动后端服务
go run ./cmd/server/main.go

# 看到以下输出表示启动成功：
# Database connected successfully
# Database migrated successfully
# Server starting addr=:8080
```

### 4. 测试后端接口

```bash
# 健康检查
curl http://localhost:8080/health

# 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"123456","name":"测试房东"}'

# 登录获取 Token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"123456"}'
```

### 5. 启动前端开发

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# H5 开发模式（浏览器访问）
npm run dev:h5

# 微信小程序开发模式
npm run dev:mp-weixin
```

### 6. 微信开发者工具调试

1. 打开微信开发者工具
2. 导入项目：选择 `frontend/dist/dev/mp-weixin` 目录
3. 在 `frontend/src/utils/config.js` 中设置**本地调试模式**：

```javascript
// 微信云托管配置
export const cloudConfig = {
  env: 'prod-6gomfrgz543fc0f3',
  serviceName: 'rental-managent',
  apiPrefix: '/api/v1'
}

export default cloudConfig
```

4. 修改 `frontend/src/utils/request.js`，在本地开发时直接请求本地后端：

```javascript
// 在文件末尾的请求函数中，本地调试时修改：
// #ifndef MP-WEIXIN || H5
// 其他环境（开发环境）- 请求本地后端
uni.request({
  url: 'http://localhost:8080' + cloudConfig.apiPrefix + options.url,
  // ...
})
// #endif
```

5. 在微信开发者工具中，点击右上角「详情」→「本地设置」→ 勾选「不校验合法域名」

---

## 微信云托管部署

### 1. 微信云托管配置信息

```
环境ID：prod-6gomfrgz543fc0f3
服务名称：rental-managent
```

### 2. 前端配置

`frontend/src/utils/config.js`：

```javascript
export const cloudConfig = {
  env: 'prod-6gomfrgz543fc0f3',      // 微信云托管环境ID
  serviceName: 'rental-managent',     // 服务名称
  apiPrefix: '/api/v1'                // API 前缀
}

export default cloudConfig
```

### 3. 构建前端

```bash
cd frontend
npm install
npm run build:mp-weixin
```

构建产物在 `frontend/dist/build/mp-weixin/`

### 4. 构建并推送后端镜像

```bash
# 进入项目目录
cd rental-management

# 构建镜像
docker build -t rental-management:latest .

# 推送到腾讯云容器镜像服务
docker tag rental-management:latest ccr.ccs.tencentyun.com/你的命名空间/rental-management:latest
docker push ccr.ccs.tencentyun.com/你的命名空间/rental-management:latest
```

### 5. 微信云托管控制台配置

| 配置项 | 推荐值 |
|--------|--------|
| 内存 | ≥ 512MB |
| CPU | 0.5核 |
| 最小实例数 | 1 |
| 最大实例数 | 10 |
| 健康检查路径 | `/health` |
| 健康检查端口 | 8080 |

### 6. 环境变量配置

在微信云托管控制台设置：

| 环境变量 | 说明 | 示例 |
|----------|------|------|
| DB_HOST | 数据库地址 | 你的MySQL地址 |
| DB_PORT | 数据库端口 | 3306 |
| DB_USERNAME | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | your_password |
| DB_DATABASE | 数据库名 | rental_management |

### 7. 小程序调用云托管服务

小程序通过 `wx.cloud.callContainer` 调用云托管服务：

```javascript
// frontend/src/utils/request.js 已封装
wx.cloud.callContainer({
  config: { env: 'prod-6gomfrgz543fc0f3' },
  path: '/api/v1/...',           // API 路径
  method: 'GET',                 // 请求方法
  header: {
    'X-WX-SERVICE': 'rental-managent'  // 服务名称
  },
  data: {},
  success: (res) => { /* 处理响应 */ }
})
```

### 8. 常见问题排查

**错误 102002: 系统错误**
- 检查服务是否正常启动
- 查看微信云托管控制台的「服务日志」
- 确认数据库连接配置正确

**错误 404: 接口不存在**
- 确认后端镜像已更新并重新部署
- 检查服务日志是否有启动错误

**CORS 跨域错误**
- 后端已配置 CORS 中间件，支持跨域请求

---

## Docker Compose 一键部署

```bash
# 启动所有服务（MySQL + Redis + 后端 + 前端）
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

访问地址：
- 前端 H5: http://localhost
- 后端 API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html

---

## API 接口文档

启动后端后访问 Swagger 文档：
- http://localhost:8080/swagger/index.html

### 主要接口

| 模块 | 方法 | 接口 | 说明 |
|------|------|------|------|
| 认证 | POST | /api/v1/auth/register | 用户注册 |
| 认证 | POST | /api/v1/auth/login | 用户登录 |
| 认证 | GET | /api/v1/auth/profile | 获取用户信息 |
| 认证 | GET | /api/v1/statistics | 首页统计数据 |
| 房源 | GET | /api/v1/properties | 房源列表 |
| 房源 | POST | /api/v1/properties | 创建房源 |
| 房源 | PUT | /api/v1/properties/:id | 更新房源 |
| 房源 | DELETE | /api/v1/properties/:id | 删除房源 |
| 房间 | GET | /api/v1/rooms | 房间列表 |
| 房间 | POST | /api/v1/rooms | 创建房间 |
| 房间 | PUT | /api/v1/rooms/:id | 更新房间 |
| 房间 | DELETE | /api/v1/rooms/:id | 删除房间 |
| 租客 | GET | /api/v1/tenants | 租客列表 |
| 合同 | GET | /api/v1/contracts | 合同列表 |
| 合同 | POST | /api/v1/contracts | 签订合同 |
| 合同 | PUT | /api/v1/contracts/:id | 更新合同 |
| 合同 | DELETE | /api/v1/contracts/:id | 删除合同 |
| 抄表 | GET | /api/v1/meter-readings | 抄表记录 |
| 抄表 | POST | /api/v1/meter-readings | 录入抄表 |
| 抄表 | DELETE | /api/v1/meter-readings/:id | 删除抄表 |
| 账单 | GET | /api/v1/bills | 账单列表 |
| 账单 | POST | /api/v1/bills | 创建账单 |
| 账单 | PUT | /api/v1/bills/:id/pay | 确认收款 |
| 账单 | DELETE | /api/v1/bills/:id | 删除账单 |
| 账单 | GET | /api/v1/bills/statistics | 账单统计 |
| 账单 | GET | /api/v1/bills/monthly-statistics | 月度统计 |

---

## 功能清单

| 功能 | 状态 | 说明 |
|------|------|------|
| 用户注册/登录 | ✅ | JWT 认证 |
| 房源管理 | ✅ | 增删改查 |
| 房间管理 | ✅ | 支持月租/季租/年租 |
| 租客管理 | ✅ | 增删改查 |
| 合同管理 | ✅ | 签订/解约/编辑/删除 |
| 抄表记录 | ✅ | 水/电/气三表 |
| 账单管理 | ✅ | 自动计算费用 |
| 收款记录 | ✅ | 支持多种支付方式 |
| 收入统计 | ✅ | 年度/月度统计 |

---

## License

MIT
