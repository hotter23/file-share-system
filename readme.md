# 基于 Docker + 腾讯云 COS 的文件共享系统

## 项目简介

这是一个使用 Go + Gin 框架开发的文件共享系统，采用微服务架构，支持 Docker 容器化部署和腾讯云 COS 云存储。

## 技术栈

### 后端 (Go)
- **框架**: Gin (Web 框架)
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis
- **云存储**: 腾讯云 COS SDK
- **认证**: JWT

### 前端
- **框架**: Vue.js 3
- **UI**: Element Plus
- **HTTP**: Axios
- **构建**: Vite

### 基础设施
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx

## 系统架构

```
┌─────────────────────────────────────────────────────┐
│                    前端 (Vue.js)                     │
└─────────────────────┬───────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────┐
│                    Nginx (:80)                      │
└─────────────────────┬───────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────┐
│                  Go Gateway (:8000)                  │
│                  路由转发 + JWT 认证                   │
└────┬─────────────────┬─────────────────┬────────────┘
     │                 │                 │
┌────▼────┐      ┌─────▼────┐     ┌─────▼────┐
│用户服务  │      │文件服务   │     │分享服务   │
│:8081    │      │:8082    │     │:8083    │
│Go+Gin   │      │Go+Gin   │     │Go+Gin   │
└─────────┘      └─────────┘     └─────────┘
```

## 快速开始

### 环境要求

- Go 1.21+
- Docker & Docker Compose
- MySQL 8.0 (Docker 方式)

### 1. 克隆项目

```bash
git clone <repository-url>
cd file-share-system
```

### 2. 启动基础设施

```bash
cd docker
docker-compose up -d mysql redis
```

### 3. 初始化数据库

MySQL 容器启动时会自动执行 `mysql/init.sql` 初始化脚本。

### 4. 配置腾讯云 COS（可选）

创建 `.env` 文件：

```bash
COS_SECRET_ID=your-secret-id
COS_SECRET_KEY=your-secret-key
COS_BUCKET=your-bucket
COS_REGION=ap-guangzhou
```

### 5. 编译后端

```bash
cd backend
go mod download
go build -o bin/gateway ./cmd/gateway
go build -o bin/user-service ./cmd/user-service
go build -o bin/file-service ./cmd/file-service
go build -o bin/share-service ./cmd/share-service
```

### 6. 启动服务

```bash
# 启动网关
./bin/gateway

# 启动用户服务
./bin/user-service

# 启动文件服务
./bin/file-service

# 启动分享服务
./bin/share-service
```

### 7. 启动前端

```bash
cd frontend
npm install
npm run dev
```

### 8. 访问系统

- 前端地址: http://localhost:3000
- API 网关: http://localhost:8000

### 9. 初始账号

- 用户名: `admin`
- 密码: `admin123`

## Docker Compose 部署

```bash
cd docker
docker-compose up -d
```

## 项目结构

```
file-share-system/
├── docs/                          # 设计文档
├── backend/                      # Go 后端
│   ├── go.mod
│   ├── cmd/
│   │   ├── gateway/              # API 网关
│   │   ├── user-service/         # 用户服务
│   │   ├── file-service/         # 文件服务
│   │   └── share-service/        # 分享服务
│   ├── internal/
│   │   ├── model/                # 数据模型
│   │   └── dto/                  # 数据传输对象
│   └── pkg/
│       ├── jwt/                  # JWT 工具
│       └── response/             # 统一响应
├── frontend/                     # Vue.js 前端
├── docker/                       # Docker 配置
│   ├── docker-compose.yml
│   ├── mysql/
│   └── nginx/
└── README.md
```

## API 接口

### 用户模块

| 接口 | 方法 | 说明 |
|------|------|------|
| `POST /api/user/register` | 注册 |
| `POST /api/user/login` | 登录 |
| `GET /api/user/info` | 获取用户信息 |

### 文件模块

| 接口 | 方法 | 说明 |
|------|------|------|
| `POST /api/file/upload` | 上传文件 |
| `POST /api/file/upload/init` | 分片上传初始化 |
| `POST /api/file/upload/chunk` | 上传分片 |
| `POST /api/file/upload/complete` | 完成分片上传 |
| `GET /api/file/list` | 文件列表 |
| `GET /api/file/:fileUuid` | 文件详情 |
| `GET /api/file/download/:fileUuid` | 下载文件 |
| `DELETE /api/file/:fileUuid` | 删除文件 |

### 文件夹模块

| 接口 | 方法 | 说明 |
|------|------|------|
| `POST /api/file/folder` | 创建文件夹 |
| `GET /api/file/folder/list` | 文件夹列表 |
| `DELETE /api/file/folder/:folderUuid` | 删除文件夹 |

### 分享模块

| 接口 | 方法 | 说明 |
|------|------|------|
| `POST /api/share` | 创建分享 |
| `GET /share/:shareUuid` | 访问分享 |
| `GET /api/share/my` | 我的分享 |
| `DELETE /api/share/:shareUuid` | 删除分享 |

## License

MIT License
