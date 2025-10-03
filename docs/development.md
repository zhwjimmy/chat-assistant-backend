# 开发指南

本文档介绍 Chat Assistant Backend 项目的开发流程、环境配置和最佳实践。

## 📋 目录

- [开发环境准备](#开发环境准备)
- [本地开发流程](#本地开发流程)
- [Docker Compose 使用场景](#docker-compose-使用场景)
- [数据库管理](#数据库管理)
- [代码质量](#代码质量)
- [测试](#测试)
- [部署](#部署)
- [故障排除](#故障排除)

## 🛠️ 开发环境准备

### 系统要求

- Go 1.23.1+ (推荐使用 1.23.4)
- Docker & Docker Compose
- Git

### Go 版本管理

如果 Go 1.23.1 不可用，可以使用版本管理器：

```bash
# 使用 gvm
gvm install go1.23.1
gvm use go1.23.1

# 使用 asdf
asdf install golang 1.23.1
asdf global golang 1.23.1

# 使用 goenv
goenv install 1.23.1
goenv global 1.23.1
```

### 安装开发工具

```bash
# 一键安装所有开发工具
make install-tools

# 或手动安装
go install github.com/cosmtrek/air@latest                    # 热重载
go install github.com/pressly/goose/v3/cmd/goose@latest      # 数据库迁移
go install github.com/swaggo/swag/cmd/swag@latest            # API 文档生成
go install github.com/google/wire/cmd/wire@latest            # 依赖注入
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # 代码检查
```

### 环境配置

```bash
# 1. 克隆项目
git clone <repository-url>
cd chat-assistant-backend

# 2. 设置 Go 代理（国内用户）
export GOPROXY=https://goproxy.cn,direct

# 3. 下载依赖
make deps

# 4. 配置环境变量
cp .env.example .env
# 编辑 .env 文件
```

## 🚀 本地开发流程

### 推荐开发方式

**只启动数据库容器，本地运行应用** - 这是最高效的开发方式。

#### 1. 启动开发环境

```bash
# 只启动 PostgreSQL 数据库
docker-compose up postgres -d

# 验证数据库启动
docker-compose ps
```

#### 2. 运行应用

```bash
# 标准运行模式
make run

# 或热重载模式（推荐）
make run-dev
```

#### 3. 测试应用

```bash
# 健康检查
curl http://localhost:8080/health

# 预期响应
{
  "service": "chat-assistant-backend",
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 4. 开发结束

```bash
# 停止数据库
docker-compose down
```

## 🐳 Docker Compose 使用场景

### 场景 1: 本地开发（推荐）

```bash
# 只启动数据库，本地运行应用
docker-compose up postgres -d
make run-dev
```

**优势：**
- 快速启动
- 资源节省
- 调试方便
- 支持热重载

### 场景 2: 完整环境测试

```bash
# 启动所有服务（数据库 + 应用）
make docker-compose-up
```

**使用场景：**
- 集成测试
- 端到端测试
- 演示环境

### 场景 3: 数据库迁移

```bash
# 启动数据库
docker-compose up postgres -d

# 执行迁移
make migrate-up

# 或使用迁移服务
docker-compose --profile migrate up migrate
```

### 场景 4: 生产部署

```bash
# 生产环境部署
docker-compose -f docker-compose.prod.yml up -d
```

## 🗄️ 数据库管理

### 启动数据库

```bash
# 启动 PostgreSQL
docker-compose up postgres -d

# 查看数据库状态
docker-compose ps postgres

# 查看数据库日志
docker-compose logs -f postgres
```

### 数据库连接

```bash
# 连接信息
Host: localhost
Port: 5432
Database: chat_assistant
Username: postgres
Password: postgres

# 连接字符串
postgres://postgres:postgres@localhost:5432/chat_assistant?sslmode=disable
```

### 迁移管理

```bash
# 创建迁移文件
# 在 internal/migrations/ 目录下创建 .sql 文件

# 执行迁移
make migrate-up

# 回滚迁移
make migrate-down

# 查看迁移状态
goose -dir internal/migrations postgres "postgres://postgres:postgres@localhost:5432/chat_assistant?sslmode=disable" status
```

## 🔍 代码质量

### 代码检查

```bash
# 运行 linter
make lint

# 格式化代码
make fmt

# 代码检查
make vet

# 检查 Go 版本
make check-go-version
```

## 🧪 测试

### 单元测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 运行特定包的测试
go test ./internal/config/...

# 运行特定测试
go test -run TestConfigLoad ./internal/config/
```

## 🔧 常用命令

### 开发命令

```bash
# 环境管理
make setup              # 设置开发环境
make deps               # 下载依赖
make run                # 运行应用
make run-dev            # 热重载模式
make build              # 构建应用

# 代码质量
make lint               # 代码检查
make fmt                # 格式化
make vet                # 代码检查

# 测试
make test               # 运行测试
make test-coverage      # 测试覆盖率

# 数据库
make migrate-up         # 执行迁移
make migrate-down       # 回滚迁移

# 文档生成
make gen-swagger        # 生成 API 文档
make gen-wire           # 生成依赖注入代码
```

### Docker 命令

```bash
# 服务管理
docker-compose up postgres -d          # 启动数据库
docker-compose up -d                   # 启动所有服务
docker-compose down                    # 停止所有服务
docker-compose ps                      # 查看服务状态
docker-compose logs -f postgres        # 查看数据库日志

# 数据库操作
docker-compose exec postgres psql -U postgres -d chat_assistant  # 连接数据库
```

## 🚨 故障排除

### 常见问题

#### 1. 依赖下载失败

```bash
# 设置 Go 代理
export GOPROXY=https://goproxy.cn,direct

# 重新下载依赖
go mod tidy
```

#### 2. 数据库连接失败

```bash
# 检查数据库状态
docker-compose ps postgres

# 查看数据库日志
docker-compose logs postgres

# 重启数据库
docker-compose restart postgres
```

#### 3. 端口被占用

```bash
# 查看端口占用
lsof -i :8080
lsof -i :5432

# 杀死占用进程
kill -9 <PID>

# 或修改配置文件中的端口
```

#### 4. 编译错误

```bash
# 检查 Go 版本
go version

# 清理构建缓存
go clean -cache

# 重新构建
make clean
make build
```

---

**Happy Coding! 🚀**
