# 环境变量配置指南

本文档介绍 Chat Assistant Backend 项目的环境变量配置。

## 📋 目录

- [配置文件说明](#配置文件说明)
- [环境变量列表](#环境变量列表)
- [配置优先级](#配置优先级)
- [使用示例](#使用示例)
- [安全注意事项](#安全注意事项)

## 📁 配置文件说明

### .env.sample
- **用途**：环境变量模板文件
- **位置**：项目根目录
- **Git 状态**：已提交到版本控制
- **内容**：包含所有可用的环境变量及其默认值

### .env
- **用途**：实际的环境变量配置文件
- **位置**：项目根目录
- **Git 状态**：被 .gitignore 忽略（不提交到版本控制）
- **内容**：根据 .env.sample 复制并修改的实际配置

## 🔧 环境变量列表

### 数据库配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `DB_HOST` | `localhost` | 数据库主机地址 |
| `DB_PORT` | `5432` | 数据库端口 |
| `DB_USER` | `postgres` | 数据库用户名 |
| `DB_PASSWORD` | `postgres` | 数据库密码 |
| `DB_NAME` | `chat_assistant` | 数据库名称 |
| `DB_SSLMODE` | `disable` | SSL 连接模式 |

### 服务器配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `SERVER_HOST` | `0.0.0.0` | 服务器监听地址 |
| `SERVER_PORT` | `8080` | 服务器端口 |

### CORS 配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `ALLOWED_ORIGINS` | `http://localhost:3000,http://localhost:3001` | 允许的跨域来源 |

### 日志配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `LOG_LEVEL` | `info` | 日志级别 (debug, info, warn, error) |
| `LOG_FORMAT` | `json` | 日志格式 (json, console) |
| `LOG_OUTPUT` | `stdout` | 日志输出 (stdout, stderr, file) |

### 国际化配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `DEFAULT_LANGUAGE` | `en` | 默认语言 |

### 其他配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `SHUTDOWN_TIMEOUT` | `30s` | 优雅停机超时时间 |
| `GIN_MODE` | `debug` | Gin 框架模式 (debug, release) |

## 📊 配置优先级

配置的加载优先级（从高到低）：

1. **系统环境变量**（最高优先级）
2. **.env 文件**
3. **config/config.yaml 文件**（最低优先级）

### 示例

```bash
# 系统环境变量会覆盖 .env 文件
export DB_PASSWORD=my_secret_password

# .env 文件会覆盖 config.yaml
# .env 文件中的 DB_HOST=production-db 会覆盖 config.yaml 中的 localhost
```

## 🚀 使用示例

### 1. 设置开发环境

```bash
# 复制模板文件
cp .env.sample .env

# 编辑配置文件
nano .env

# 启动应用
make run
```

### 2. 设置生产环境

```bash
# 设置系统环境变量
export DB_HOST=production-db.example.com
export DB_PASSWORD=super_secret_password
export LOG_LEVEL=warn
export GIN_MODE=release

# 启动应用
make run
```

### 3. Docker 环境

```bash
# 使用 .env 文件
docker-compose up -d

# 或直接设置环境变量
docker run -e DB_HOST=db -e DB_PASSWORD=secret chat-assistant-backend
```

### 4. 不同环境配置

```bash
# 开发环境
cp .env.sample .env.dev
# 编辑 .env.dev

# 测试环境
cp .env.sample .env.test
# 编辑 .env.test

# 生产环境
cp .env.sample .env.prod
# 编辑 .env.prod

# 使用特定环境
ENV_FILE=.env.dev make run
```

## 🔒 安全注意事项

### 1. 敏感信息保护

```bash
# ❌ 错误：不要提交包含敏感信息的 .env 文件
git add .env

# ✅ 正确：只提交模板文件
git add .env.sample
```

### 2. 生产环境安全

```bash
# 使用强密码
DB_PASSWORD=very_strong_random_password_here

# 限制 CORS 来源
ALLOWED_ORIGINS=https://yourdomain.com,https://api.yourdomain.com

# 使用 HTTPS
DB_SSLMODE=require

# 生产日志级别
LOG_LEVEL=warn
GIN_MODE=release
```

### 3. 环境变量验证

应用启动时会验证必要的环境变量：

```bash
# 检查配置是否正确加载
make run
# 查看日志中的配置信息
```

## 🛠️ 配置管理最佳实践

### 1. 开发环境

```bash
# 使用默认配置
cp .env.sample .env
# 根据需要修改
```

### 2. 测试环境

```bash
# 使用测试数据库
DB_HOST=test-db
DB_NAME=chat_assistant_test
LOG_LEVEL=debug
```

### 3. 生产环境

```bash
# 使用环境变量而不是文件
export DB_HOST=prod-db.example.com
export DB_PASSWORD=$(cat /secrets/db_password)
export LOG_LEVEL=warn
export GIN_MODE=release
```

### 4. 配置验证

```bash
# 检查配置是否正确
make run
# 查看健康检查
curl http://localhost:8080/health
```

## 🔍 故障排除

### 常见问题

#### 1. 配置未生效

```bash
# 检查环境变量
env | grep DB_

# 检查 .env 文件
cat .env

# 重启应用
make run
```

#### 2. 数据库连接失败

```bash
# 检查数据库配置
echo $DB_HOST
echo $DB_PORT
echo $DB_USER

# 测试数据库连接
docker-compose exec postgres psql -U postgres -d chat_assistant
```

#### 3. 端口冲突

```bash
# 检查端口占用
lsof -i :8080

# 修改端口
export SERVER_PORT=8081
```

## 📚 相关文档

- [开发指南](development.md)
- [Docker Compose 使用指南](docker-compose-guide.md)
- [部署指南](deployment.md)

---

**提示**：在生产环境中，建议使用专门的配置管理工具如 Kubernetes ConfigMaps 或 Docker Secrets 来管理敏感配置。
