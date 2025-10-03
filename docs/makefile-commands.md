# Makefile 命令参考

本文档详细介绍 Chat Assistant Backend 项目中 Makefile 提供的所有命令及其用法。

## 📋 目录

- [构建和运行](#构建和运行)
- [开发工具](#开发工具)
- [数据库操作](#数据库操作)
- [代码质量](#代码质量)
- [测试](#测试)
- [Docker 操作](#docker-操作)
- [文档生成](#文档生成)
- [环境管理](#环境管理)

## 🏗️ 构建和运行

### 基础构建

```bash
# 构建应用
make build

# 清理构建产物
make clean

# 下载依赖
make deps
```

### 运行应用

```bash
# 标准运行模式
make run

# 热重载模式（推荐开发使用）
make run-dev

# 检查 Go 版本
make check-go-version
```

**说明：**
- `make run`：直接运行编译后的二进制文件
- `make run-dev`：使用 air 工具实现热重载，代码变更时自动重启
- 需要先安装 air：`go install github.com/cosmtrek/air@latest`

## 🛠️ 开发工具

### 工具安装

```bash
# 安装所有开发工具
make install-tools
```

安装的工具包括：
- `air`：热重载工具
- `goose`：数据库迁移工具
- `swag`：API 文档生成工具
- `wire`：依赖注入工具
- `golangci-lint`：代码检查工具

### 环境设置

```bash
# 设置开发环境
make setup
```

等同于执行：
1. `make check-go-version`
2. `make install-tools`
3. `make deps`

## 🗄️ 数据库操作

### 迁移管理

```bash
# 执行数据库迁移
make migrate-up

# 回滚数据库迁移
make migrate-down
```

**前提条件：**
- 需要安装 goose：`go install github.com/pressly/goose/v3/cmd/goose@latest`
- 需要启动 PostgreSQL 数据库
- 迁移文件位于 `internal/migrations/` 目录

### 数据库连接

默认连接配置：
- Host: localhost
- Port: 5432
- User: postgres
- Password: postgres
- Database: chat_assistant

## 🔍 代码质量

### 代码检查

```bash
# 运行 linter
make lint

# 格式化代码
make fmt

# 代码检查
make vet
```

**说明：**
- `make lint`：使用 golangci-lint 进行代码检查
- `make fmt`：使用 gofmt 格式化代码
- `make vet`：使用 go vet 进行代码检查

### 配置文件

- `.golangci.yml`：golangci-lint 配置文件
- 包含的检查项：gofmt, goimports, govet, errcheck, staticcheck, unused, gosimple, ineffassign, typecheck, gocyclo, goconst, misspell, lll

## 🧪 测试

### 单元测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage
```

**说明：**
- `make test`：运行所有包的测试
- `make test-coverage`：生成覆盖率报告并保存为 `coverage.html`

### 测试文件

测试文件应放在对应的包目录下，以 `_test.go` 结尾。

## 🐳 Docker 操作

### 镜像构建

```bash
# 构建 Docker 镜像
make docker-build

# 运行 Docker 容器
make docker-run
```

### 容器编排

```bash
# 启动 docker-compose 服务
make docker-compose-up

# 停止 docker-compose 服务
make docker-compose-down
```

**说明：**
- `make docker-build`：构建应用 Docker 镜像
- `make docker-run`：运行单个容器
- `make docker-compose-up`：启动完整的服务栈（数据库 + 应用）
- `make docker-compose-down`：停止所有服务

## 📚 文档生成

### API 文档

```bash
# 生成 Swagger 文档
make gen-swagger
```

**前提条件：**
- 需要安装 swag：`go install github.com/swaggo/swag/cmd/swag@latest`
- 需要在代码中添加 Swagger 注释

### 依赖注入

```bash
# 生成 Wire 依赖注入代码
make gen-wire
```

**前提条件：**
- 需要安装 wire：`go install github.com/google/wire/cmd/wire@latest`
- 需要在 `internal/wire.go` 中定义依赖关系

## 🌍 环境管理

### 环境变量

项目支持通过以下方式配置环境变量：

1. **配置文件**：`config/config.yaml`
2. **环境变量**：`.env` 文件
3. **系统环境变量**：直接设置

### 配置优先级

1. 系统环境变量（最高优先级）
2. `.env` 文件
3. `config/config.yaml`（最低优先级）

## 📝 命令示例

### 完整开发流程

```bash
# 1. 设置开发环境
make setup

# 2. 启动数据库
docker-compose up postgres -d

# 3. 运行应用（热重载模式）
make run-dev

# 4. 在另一个终端运行测试
make test

# 5. 代码检查
make lint

# 6. 生成文档
make gen-swagger
```

### 日常开发

```bash
# 启动开发环境
docker-compose up postgres -d
make run-dev

# 修改代码后自动重启
# 运行测试
make test

# 提交前检查
make lint
make fmt
make vet
```

### 部署准备

```bash
# 构建生产镜像
make docker-build

# 运行完整测试
make test-coverage

# 生成文档
make gen-swagger
make gen-wire
```

## 🚨 故障排除

### 常见问题

#### 1. 命令不存在

```bash
# 检查工具是否安装
which air
which goose
which swag
which wire
which golangci-lint

# 安装缺失的工具
make install-tools
```

#### 2. 权限问题

```bash
# 检查文件权限
ls -la bin/

# 修复权限
chmod +x bin/chat-assistant-backend
```

#### 3. 依赖问题

```bash
# 清理并重新下载依赖
make clean
make deps
```

#### 4. 数据库连接问题

```bash
# 检查数据库状态
docker-compose ps postgres

# 重启数据库
docker-compose restart postgres
```

## 📊 命令参考表

| 命令 | 功能 | 前提条件 |
|------|------|----------|
| `make build` | 构建应用 | Go 环境 |
| `make run` | 运行应用 | 已构建 |
| `make run-dev` | 热重载运行 | 安装 air |
| `make test` | 运行测试 | Go 环境 |
| `make lint` | 代码检查 | 安装 golangci-lint |
| `make migrate-up` | 数据库迁移 | 安装 goose + 数据库 |
| `make docker-build` | 构建镜像 | Docker |
| `make gen-swagger` | 生成 API 文档 | 安装 swag |
| `make gen-wire` | 生成依赖注入 | 安装 wire |

## 💡 最佳实践

1. **开发时使用**：`make run-dev` 进行热重载开发
2. **提交前检查**：`make lint && make test`
3. **部署前准备**：`make build && make test-coverage`
4. **定期更新依赖**：`make deps`
5. **保持工具更新**：定期运行 `make install-tools`

---

**提示**：使用 `make help` 查看所有可用命令的简要说明。
