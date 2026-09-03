# ZZHAN 博客系统

一个基于 Go + Vue 3 的前后端分离博客系统，支持文章管理、分类标签、评论点赞、用户认证等功能。

## 🚀 技术栈

### 后端 (ZZHAN)

| 组件 | 技术 |
|------|------|
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | MySQL |
| 缓存 | Redis |
| 认证 | JWT |
| 日志 | Zap + Lumberjack |
| 配置 | Viper |
| 验证码 | base64Captcha |

### 前端 (ZZHAN_web)

- **前台** (user): Vue 3 + TypeScript + Vite
- **后台** (admin): Vue 3 + TypeScript + Vite

## 📁 项目结构

```
博客/
├── ZZHAN/                    # 后端项目
│   ├── cmd/server/           # 程序入口
│   ├── internal/
│   │   ├── api/              # 路由和控制器
│   │   │   ├── admin/        # 后台管理接口
│   │   │   └── web/          # 前台接口
│   │   ├── app/              # 应用初始化
│   │   ├── middleware/       # 中间件
│   │   ├── model/            # 数据模型
│   │   │   ├── dto/          # 数据传输对象
│   │   │   └── entity/       # 实体定义
│   │   ├── repository/       # 数据访问层
│   │   └── service/          # 业务逻辑层
│   ├── pkg/                  # 公共工具包
│   │   ├── b64c/             # 验证码
│   │   ├── config/           # 配置加载
│   │   ├── database/         # 数据库连接
│   │   ├── errors/           # 错误处理
│   │   ├── jwt/              # JWT 认证
│   │   ├── logger/           # 日志工具
│   │   ├── response/         # 响应封装
│   │   ├── storage/          # 存储服务
│   │   └── utils/            # 工具函数
│   ├── config.yaml           # 配置文件
│   ├── Makefile              # 构建脚本
│   └── go.mod                # Go 依赖
│
├── ZZHAN_web/                # 前端项目
│   ├── user/                 # 前台用户端
│   └── admin/                # 后台管理端
│
└── docker-compose.yml        # Docker 编排
```

## ✨ 功能特性

### 前台功能
- 📝 文章浏览（Markdown 渲染）
- 🔍 文章搜索
- 📂 分类筛选
- 🏷️ 标签筛选
- 👍 点赞/取消点赞
- 💬 评论/回复
- 📊 站点统计
- 📅 文章归档
- 🔐 用户注册/登录
- 👤 个人中心

### 后台功能
- 📊 数据仪表盘
- 📝 文章管理（发布/编辑/下架）
- 📂 分类管理
- 🏷️ 标签管理
- 💬 评论管理
- 👥 用户管理
- ⚙️ 站点设置
- 📋 操作日志

## 🛠️ 快速开始

### 环境要求

- Go 1.24+
- Node.js 18+
- MySQL 8.0+
- Redis 7.0+

### 1. 克隆项目

```bash
git clone <repository-url>
cd 博客
```

### 2. 配置后端

```bash
cd ZZHAN

# 复制配置文件
cp .env.example .env

# 编辑配置文件，填入数据库和 Redis 连接信息
vim config.yaml
```

### 3. 初始化数据库

```bash
# 导入数据库结构
mysql -u root -p < zzhan_dump.sql
```

### 4. 启动后端

```bash
# 使用 Makefile
make run

# 或直接运行
go run cmd/server/main.go

# 指定配置文件
go run cmd/server/main.go -c config.yaml
```

### 5. 启动前端

```bash
# 前台
cd ZZHAN_web/user
npm install
npm run dev

# 后台
cd ZZHAN_web/admin
npm install
npm run dev
```

## 📦 Docker 部署

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

## 🔧 Makefile 命令

| 命令 | 说明 |
|------|------|
| `make build` | 编译项目 |
| `make run` | 运行项目 |
| `make test` | 运行测试 |
| `make test-coverage` | 生成测试覆盖率报告 |
| `make clean` | 清理构建文件 |
| `make deps` | 下载依赖 |
| `make fmt` | 格式化代码 |
| `make vet` | 静态检查 |
| `make lint` | 代码检查（需安装 golangci-lint） |
| `make help` | 显示帮助信息 |

## 📡 API 接口

### 前台接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/articles` | 文章列表 |
| GET | `/api/v1/articles/:slug` | 文章详情 |
| GET | `/api/v1/categories` | 分类列表 |
| GET | `/api/v1/tags` | 标签列表 |
| GET | `/api/v1/archives` | 归档列表 |
| GET | `/api/v1/stats` | 站点统计 |
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/login` | 用户登录 |
| POST | `/api/v1/articles/:id/like` | 点赞 |
| POST | `/api/v1/articles/:id/comments` | 发表评论 |

### 后台接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/admin/auth/login` | 管理员登录 |
| GET | `/api/admin/dashboard` | 仪表盘数据 |
| GET/POST/PUT/DELETE | `/api/admin/articles` | 文章管理 |
| GET/POST/PUT/DELETE | `/api/admin/categories` | 分类管理 |
| GET/POST/PUT/DELETE | `/api/admin/tags` | 标签管理 |
| GET/PUT/DELETE | `/api/admin/comments` | 评论管理 |
| GET/PUT | `/api/admin/users` | 用户管理 |
| GET/PUT | `/api/admin/settings` | 站点设置 |

## 🔐 认证机制

- 使用 JWT 进行身份验证
- 支持 Access Token + Refresh Token 双 Token 机制
- 单设备登录（同一用户只允许一个活跃会话）
- Token 黑名单（登出时失效）

## 📊 浏览量统计

采用 **IP 去重** 机制：
- 同一 IP 在 10 分钟内多次访问只计 1 次
- 使用 Redis 存储访问记录
- Redis 不可用时降级为直接计数

## 📝 许可证

MIT License

## 👥 作者

- 23107

## 🙏 致谢

感谢所有为这个项目做出贡献的人。
