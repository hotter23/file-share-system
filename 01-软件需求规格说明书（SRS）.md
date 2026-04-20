# 软件需求规格说明书（SRS）

## 项目名称

**基于 Docker + 腾讯云 COS 的文件共享系统**

---

## 1. 概述

### 1.1 项目背景

随着云计算、大数据及移动互联网技术的飞速发展，文件数据量呈爆炸式增长，用户对文件存储、传输、共享的便捷性、安全性及可扩展性需求日益提升。本系统旨在构建一个基于 Docker 容器化技术和腾讯云 COS 云存储服务的私有文件共享平台，实现文件的安全存储、便捷共享与高效传输。

### 1.2 项目目标

1. 实现支持多用户认证的文件上传、下载、共享功能
2. 采用 Docker 容器化部署，实现快速搭建与跨环境迁移
3. 集成腾讯云 COS，提供海量、安全、可靠的云端存储能力
4. 实现精细化的文件权限控制和多用户并发访问支持
5. 保障数据传输和存储的安全性

### 1.3 项目范围

**纳入范围：**
- 用户管理模块（注册、登录、认证）
- 文件管理模块（上传、下载、删除、预览）
- 文件共享模块（链接分享、权限控制）
- 存储管理模块（对接腾讯云 COS）
- 容器化部署（Docker Compose 编排）

**不纳入范围：**
- 移动端原生 APP
- 实时音视频通信
- 商业 CDN 加速配置

---

## 2. 用户角色

| 角色 | 说明 |
|------|------|
| 普通用户 | 注册登录后，可上传、下载、管理自己的文件，进行文件分享 |
| 管理员 | 管理系统用户，监控文件存储情况，管理系统配置 |

---

## 3. 功能需求

### 3.1 用户管理模块

#### 3.1.1 用户注册
- 用户名（唯一）、邮箱、密码
- 密码加密存储（BCrypt）
- 注册后发送激活邮件（可选）

#### 3.1.2 用户登录
- 用户名/邮箱 + 密码登录
- 登录成功后返回 JWT Token
- 支持记住登录状态

#### 3.1.3 用户认证
- JWT Token 认证
- Token 有效期 24 小时
- Token 刷新机制

#### 3.1.4 管理员功能
- 用户列表查询
- 用户启用/禁用
- 用户删除

### 3.2 文件管理模块

#### 3.2.1 文件上传
- 支持拖拽上传和点击上传
- 支持大文件分片上传（每个分片 5MB）
- 支持断点续传
- 上传进度实时显示
- 文件类型和大小校验（单文件最大 5GB）
- 上传完成后计算 MD5 校验

#### 3.2.2 文件下载
- 普通下载
- 分片下载（断点续传）
- 下载时显示进度

#### 3.2.3 文件管理
- 文件列表展示（分页）
- 文件搜索（按文件名）
- 文件夹管理（创建文件夹）
- 文件删除（支持批量删除）
- 文件重命名
- 文件详情查看（大小、上传时间、下载次数）

#### 3.2.4 文件预览
- 图片预览（JPEG、PNG、GIF、WebP）
- 文本文件预览（TXT、JSON、XML、MD）

### 3.3 文件共享模块

#### 3.3.1 创建分享链接
- 生成分享链接（基于 UUID）
- 设置链接有效期（1小时/24小时/7天/永久）
- 设置访问密码（可选）
- 记录分享次数

#### 3.3.2 访问分享链接
- 通过链接访问共享文件
- 密码验证（如设置了密码）
- 下载文件

#### 3.3.3 分享管理
- 查看我创建的分享
- 取消分享
- 复制分享链接

### 3.4 存储模块（腾讯云 COS 集成）

#### 3.4.1 COS 文件操作
- 文件上传至 COS
- 文件下载自 COS
- 文件删除
- 文件复制

#### 3.4.2 传输优化
- 分片上传（最大支持 20GB 单文件）
- 并行上传分片
- 传输限速（可配置）

#### 3.4.3 CDN 加速
- 热门文件开启 CDN 加速
- 边缘节点缓存

### 3.5 权限控制模块

#### 3.5.1 文件访问权限
- 私有文件：仅文件所有者可访问
- 共享文件：通过分享链接访问

#### 3.5.2 目录权限
- 所有者：读、写、删、管理权限
- 无权限用户：无法访问

---

## 4. 非功能需求

### 4.1 性能需求

| 指标 | 要求 |
|------|------|
| 系统并发用户数 | ≥ 100 |
| 文件上传速度 | 取决于网络带宽，瓶颈在 COS 上传 |
| 文件下载速度 | 取决于网络带宽，瓶颈在 COS 下载 |
| 页面响应时间 | ≤ 2 秒 |
| 大文件分片上传 | 支持单文件最大 5GB |

### 4.2 安全性需求

| 需求 | 说明 |
|------|------|
| 身份认证 | JWT Token |
| 密码加密 | BCrypt |
| 传输加密 | HTTPS |
| 文件存储加密 | AES-256（腾讯云 COS 服务端加密） |
| 文件完整性 | MD5 校验 |
| SQL 注入防护 | 参数化查询 |
| XSS 防护 | 输入过滤与转义 |

### 4.3 可用性需求

| 需求 | 说明 |
|------|------|
| 系统可用性 | 99.9% |
| 部署方式 | Docker Compose 一键部署 |
| 数据备份 | COS 多地域备份 |

### 4.4 兼容性需求

| 类型 | 要求 |
|------|------|
| 浏览器 | Chrome、Firefox、Safari、Edge 最新两个版本 |
| 操作系统 | Windows、macOS、Linux |
| 移动端 | 支持响应式布局 |

---

## 5. 系统架构

### 5.1 整体架构

```
┌─────────────────────────────────────────────────────┐
│                    前端 (Vue.js)                     │
│              Web 浏览器 / 移动端 H5                  │
└─────────────────────┬───────────────────────────────┘
                      │ HTTPS
┌─────────────────────▼───────────────────────────────┐
│                  Nginx 反向代理                       │
│              静态资源 + API 路由                      │
└─────────────────────┬───────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────┐
│              Spring Cloud 微服务架构                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────┐│
│  │用户服务   │  │文件服务   │  │分享服务   │  │存储服务││
│  │:8081     │  │:8082     │  │:8083     │  │:8084 ││
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └──┬───┘│
└───────┼────────────┼────────────┼─────────────┼────┘
        │            │            │             │
┌───────▼────────────▼────────────▼─────────────▼────┐
│                    MySQL 数据库                      │
│              (用户、文件、权限、分享数据)             │
├─────────────────────────────────────────────────────┤
│                    Redis 缓存                        │
│              (会话、Token、热点数据)                  │
├─────────────────────────────────────────────────────┤
│                 Docker Compose 编排                   │
│              各微服务容器化部署                       │
└─────────────────────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────┐
│                  腾讯云 COS 存储                      │
│              (文件本体、CDN 加速)                     │
└─────────────────────────────────────────────────────┘
```

### 5.2 微服务拆分

| 服务 | 端口 | 职责 |
|------|------|------|
| gateway | 8000 | API 网关，统一路由，认证鉴权 |
| user-service | 8081 | 用户管理，登录认证 |
| file-service | 8082 | 文件上传下载管理 |
| share-service | 8083 | 分享链接管理 |
| cos-service | 8084 | 腾讯云 COS 操作封装 |

### 5.3 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue.js 3 + Element Plus + Axios |
| 后端 | Spring Boot 3 + Spring Cloud |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis 7.0 |
| 容器 | Docker + Docker Compose |
| 云存储 | 腾讯云 COS SDK |
| 认证 | JWT + Spring Security |
| 构建 | Maven |

---

## 6. 数据库设计

### 6.1 ER 图（实体关系）

```
用户 (user)
  ├── id (PK, BIGINT)
  ├── username (UK, VARCHAR)
  ├── password (VARCHAR)
  ├── email (VARCHAR)
  ├── role (ENUM: USER, ADMIN)
  ├── status (ENUM: ACTIVE, DISABLED)
  └── created_at, updated_at

文件 (file)
  ├── id (PK, BIGINT)
  ├── file_uuid (UK, VARCHAR) — 业务主键
  ├── file_name (VARCHAR)
  ├── file_size (BIGINT)
  ├── file_type (VARCHAR)
  ├── md5 (VARCHAR)
  ├── cos_key (VARCHAR) — COS 存储路径
  ├── bucket_name (VARCHAR)
  ├── storage_type (ENUM: COS, LOCAL)
  ├── folder_id (FK → folder.id, NULLABLE)
  ├── user_id (FK → user.id)
  ├── download_count (INT)
  └── created_at, updated_at

文件夹 (folder)
  ├── id (PK, BIGINT)
  ├── folder_uuid (UK, VARCHAR)
  ├── folder_name (VARCHAR)
  ├── parent_id (FK → folder.id, NULLABLE) — 上级目录
  ├── user_id (FK → user.id)
  └── created_at, updated_at

分享 (share)
  ├── id (PK, BIGINT)
  ├── share_uuid (UK, VARCHAR)
  ├── file_id (FK → file.id)
  ├── share_code (VARCHAR) — 短码
  ├── password (VARCHAR, NULLABLE) — 分享密码
  ├── expire_time (DATETIME, NULLABLE)
  ├── view_count (INT)
  ├── created_by (FK → user.id)
  └── created_at

文件分片记录 (file_chunk)
  ├── id (PK, BIGINT)
  ├── file_uuid (VARCHAR)
  ├── chunk_index (INT)
  ├── chunk_size (BIGINT)
  ├── upload_id (VARCHAR) — COS 分片上传 ID
  └── created_at
```

### 6.2 表结构

```sql
-- 用户表
CREATE TABLE `user` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `username` VARCHAR(64) NOT NULL UNIQUE,
    `password` VARCHAR(255) NOT NULL,
    `email` VARCHAR(128),
    `role` ENUM('USER', 'ADMIN') DEFAULT 'USER',
    `status` ENUM('ACTIVE', 'DISABLED') DEFAULT 'ACTIVE',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 文件夹表
CREATE TABLE `folder` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `folder_uuid` VARCHAR(64) NOT NULL UNIQUE,
    `folder_name` VARCHAR(255) NOT NULL,
    `parent_id` BIGINT,
    `user_id` BIGINT NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`parent_id`) REFERENCES `folder`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE
);

-- 文件表
CREATE TABLE `file` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `file_uuid` VARCHAR(64) NOT NULL UNIQUE,
    `file_name` VARCHAR(255) NOT NULL,
    `file_size` BIGINT NOT NULL DEFAULT 0,
    `file_type` VARCHAR(64),
    `md5` VARCHAR(64),
    `cos_key` VARCHAR(512),
    `bucket_name` VARCHAR(128),
    `storage_type` ENUM('COS', 'LOCAL') DEFAULT 'COS',
    `folder_id` BIGINT,
    `user_id` BIGINT NOT NULL,
    `download_count` INT DEFAULT 0,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`folder_id`) REFERENCES `folder`(`id`) ON DELETE SET NULL,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE
);

-- 分享表
CREATE TABLE `share` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `share_uuid` VARCHAR(64) NOT NULL UNIQUE,
    `file_id` BIGINT NOT NULL,
    `share_code` VARCHAR(16),
    `password` VARCHAR(255),
    `expire_time` DATETIME,
    `view_count` INT DEFAULT 0,
    `created_by` BIGINT NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`file_id`) REFERENCES `file`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`created_by`) REFERENCES `user`(`id`) ON DELETE CASCADE
);

-- 分片上传记录表
CREATE TABLE `file_chunk` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `file_uuid` VARCHAR(64) NOT NULL,
    `chunk_index` INT NOT NULL,
    `chunk_size` BIGINT NOT NULL,
    `upload_id` VARCHAR(255),
    `status` TINYINT DEFAULT 0,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY `uk_file_chunk` (`file_uuid`, `chunk_index`)
);
```

---

## 7. API 接口设计

### 7.1 用户模块

#### 7.1.1 用户注册
- **POST** `/api/user/register`
- Request:
```json
{
  "username": "string",
  "password": "string",
  "email": "string"
}
```
- Response (200):
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "userId": "string",
    "username": "string"
  }
}
```

#### 7.1.2 用户登录
- **POST** `/api/user/login`
- Request:
```json
{
  "username": "string",
  "password": "string"
}
```
- Response (200):
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "string",
    "userId": "string",
    "username": "string"
  }
}
```

#### 7.1.3 获取用户信息
- **GET** `/api/user/info`
- Headers: `Authorization: Bearer <token>`
- Response (200):
```json
{
  "code": 200,
  "data": {
    "userId": "string",
    "username": "string",
    "email": "string",
    "role": "USER"
  }
}
```

### 7.2 文件模块

#### 7.2.1 上传文件（简单上传）
- **POST** `/api/file/upload`
- Headers: `Authorization: Bearer <token>`
- Form-Data: `file` (binary), `folderId` (optional)
- Response (200):
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "fileUuid": "string",
    "fileName": "string",
    "fileSize": 123456,
    "md5": "string"
  }
}
```

#### 7.2.2 分片上传初始化
- **POST** `/api/file/upload/init`
- Headers: `Authorization: Bearer <token>`
- Request:
```json
{
  "fileName": "string",
  "fileSize": 5242880000,
  "chunkSize": 5242880,
  "folderId": "string"
}
```
- Response (200):
```json
{
  "code": 200,
  "data": {
    "fileUuid": "string",
    "uploadId": "string",
    "chunkCount": 1000
  }
}
```

#### 7.2.3 分片上传
- **POST** `/api/file/upload/chunk`
- Headers: `Authorization: Bearer <token>`
- Form-Data: `fileUuid`, `uploadId`, `chunkIndex`, `file` (binary)
- Response (200):
```json
{
  "code": 200,
  "message": "分片上传成功",
  "data": {
    "chunkIndex": 0,
    "uploadId": "string"
  }
}
```

#### 7.2.4 分片上传完成
- **POST** `/api/file/upload/complete`
- Headers: `Authorization: Bearer <token>`
- Request:
```json
{
  "fileUuid": "string",
  "uploadId": "string"
}
```
- Response (200):
```json
{
  "code": 200,
  "message": "上传完成",
  "data": {
    "fileUuid": "string",
    "fileName": "string",
    "md5": "string"
  }
}
```

#### 7.2.5 下载文件
- **GET** `/api/file/download/{fileUuid}`
- Headers: `Authorization: Bearer <token>`
- Response: 文件流

#### 7.2.6 获取文件列表
- **GET** `/api/file/list`
- Headers: `Authorization: Bearer <token>`
- Query: `folderId` (optional), `page` (default 1), `size` (default 20), `keyword` (optional)
- Response (200):
```json
{
  "code": 200,
  "data": {
    "records": [
      {
        "fileUuid": "string",
        "fileName": "string",
        "fileSize": 123456,
        "fileType": "string",
        "createdAt": "2026-04-17T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 20
  }
}
```

#### 7.2.7 删除文件
- **DELETE** `/api/file/{fileUuid}`
- Headers: `Authorization: Bearer <token>`
- Response (200):
```json
{
  "code": 200,
  "message": "删除成功"
}
```

### 7.3 文件夹模块

#### 7.3.1 创建文件夹
- **POST** `/api/folder`
- Headers: `Authorization: Bearer <token>`
- Request:
```json
{
  "folderName": "string",
  "parentId": "string"
}
```
- Response (200):
```json
{
  "code": 200,
  "data": {
    "folderUuid": "string",
    "folderName": "string"
  }
}
```

#### 7.3.2 获取文件夹列表
- **GET** `/api/folder/list`
- Headers: `Authorization: Bearer <token>`
- Query: `parentId` (optional)
- Response (200):
```json
{
  "code": 200,
  "data": [
    {
      "folderUuid": "string",
      "folderName": "string",
      "parentId": "string",
      "createdAt": "2026-04-17T10:00:00Z"
    }
  ]
}
```

### 7.4 分享模块

#### 7.4.1 创建分享链接
- **POST** `/api/share`
- Headers: `Authorization: Bearer <token>`
- Request:
```json
{
  "fileUuid": "string",
  "password": "string",
  "expireType": "1H|24H|7D|PERMANENT"
}
```
- Response (200):
```json
{
  "code": 200,
  "data": {
    "shareUuid": "string",
    "shareUrl": "/share/{shareUuid}",
    "shareCode": "string",
    "expireTime": "2026-04-18T10:00:00Z"
  }
}
```

#### 7.4.2 访问分享
- **GET** `/api/share/{shareUuid}`
- Query: `password` (optional)
- Response (200):
```json
{
  "code": 200,
  "data": {
    "fileName": "string",
    "fileSize": 123456,
    "fileType": "string",
    "requirePassword": true
  }
}
```

#### 7.4.3 下载分享文件
- **GET** `/api/share/{shareUuid}/download`
- Query: `password` (optional)
- Response: 文件流

#### 7.4.4 获取我的分享列表
- **GET** `/api/share/my`
- Headers: `Authorization: Bearer <token>`
- Response (200):
```json
{
  "code": 200,
  "data": [
    {
      "shareUuid": "string",
      "fileName": "string",
      "viewCount": 10,
      "expireTime": "2026-04-18T10:00:00Z",
      "createdAt": "2026-04-17T10:00:00Z"
    }
  ]
}
```

#### 7.4.5 删除分享
- **DELETE** `/api/share/{shareUuid}`
- Headers: `Authorization: Bearer <token>`
- Response (200):
```json
{
  "code": 200,
  "message": "删除成功"
}
```

### 7.5 状态码规范

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或 Token 失效 |
| 403 | 无权限访问 |
| 404 | 资源不存在 |
| 413 | 文件大小超出限制 |
| 415 | 不支持的文件类型 |
| 500 | 服务器内部错误 |

---

## 8. 验收标准

### 8.1 功能验收

- [ ] 用户可以完成注册、登录、登出
- [ ] 登录后获取 JWT Token，Token 在请求头中有效
- [ ] 用户可以上传文件（支持拖拽）
- [ ] 用户可以上传大文件（>100MB）并显示进度
- [ ] 支持断点续传
- [ ] 用户可以下载自己的文件
- [ ] 用户可以创建文件夹
- [ ] 用户可以删除文件和文件夹
- [ ] 用户可以创建分享链接
- [ ] 分享链接可以设置密码和有效期
- [ ] 任何人可通过分享链接下载文件
- [ ] 管理员可以查看所有用户和文件列表

### 8.2 非功能验收

- [ ] 系统可通过 Docker Compose 一键部署
- [ ] 前后端分离架构，API 响应正常
- [ ] 文件上传下载使用 HTTPS 加密
- [ ] 密码使用 BCrypt 加密存储
- [ ] 支持主流浏览器访问

---

## 9. 项目目录结构

```
file-share-system/
├── docs/                          # 文档目录
│   ├── 01-软件需求规格说明书（SRS）.md
│   ├── 02-系统设计文档.md
│   ├── 03-数据库设计文档.md
│   └── 04-API接口文档.md
├── backend/                       # 后端项目
│   ├── file-share-gateway/         # API 网关
│   ├── file-share-user/            # 用户服务
│   ├── file-share-file/            # 文件服务
│   ├── file-share-share/            # 分享服务
│   └── pom.xml                     # 父 POM
├── frontend/                      # 前端项目
│   ├── src/
│   ├── public/
│   └── package.json
├── docker/                        # Docker 配置
│   ├── docker-compose.yml
│   ├── mysql/
│   └── redis/
└── README.md
```
