# Docker Compose 使用指南

本文档详细介绍 Chat Assistant Backend 项目中 Docker Compose 的使用场景和最佳实践。

## 📋 目录

- [概述](#概述)
- [使用场景](#使用场景)
- [配置文件说明](#配置文件说明)
- [常用命令](#常用命令)
- [最佳实践](#最佳实践)
- [故障排除](#故障排除)

## 🎯 概述

Docker Compose 在项目中的主要作用是：

1. **本地开发环境**：提供 PostgreSQL 数据库服务
2. **集成测试**：完整的应用 + 数据库环境
3. **生产部署**：容器化部署方案
4. **数据库迁移**：自动化数据库迁移

## 🚀 使用场景

### 场景 1: 本地开发（推荐）

**目标**：只启动数据库，本地运行 Go 应用

```bash
# 启动数据库
docker-compose up postgres -d

# 本地运行应用
make run-dev

# 停止数据库
docker-compose down
```

**优势：**
- ✅ 快速启动（只需启动数据库）
- ✅ 资源节省（不需要构建应用镜像）
- ✅ 调试方便（可以直接在 IDE 中调试）
- ✅ 热重载支持（使用 air 工具）

### 场景 2: 完整环境测试

**目标**：启动完整的应用环境进行集成测试

```bash
# 启动所有服务
make docker-compose-up

# 测试 API
curl http://localhost:8080/health

# 停止所有服务
make docker-compose-down
```

**使用场景：**
- 集成测试
- 端到端测试
- 演示环境
- CI/CD 流水线

### 场景 3: 数据库迁移

**目标**：执行数据库迁移操作

```bash
# 方法 1: 使用 Makefile
docker-compose up postgres -d
make migrate-up

# 方法 2: 使用迁移服务
docker-compose --profile migrate up migrate

# 回滚迁移
make migrate-down
```

### 场景 4: 生产部署

**目标**：生产环境部署

```bash
# 构建生产镜像
make docker-build

# 生产环境部署
docker-compose -f docker-compose.prod.yml up -d
```

## 📁 配置文件说明

### docker-compose.yaml

```yaml
version: '3.8'

services:
  # PostgreSQL 数据库服务
  postgres:
    image: postgres:15-alpine
    container_name: chat-assistant-postgres
    environment:
      POSTGRES_DB: chat_assistant
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  # 应用服务（可选）
  chat-assistant-backend:
    build: .
    container_name: chat-assistant-backend
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=chat_assistant
    depends_on:
      postgres:
        condition: service_healthy

  # 数据库迁移服务（可选）
  migrate:
    image: migrate/migrate
    container_name: chat-assistant-migrate
    command: [
      "-path", "/migrations",
      "-database", "postgres://postgres:postgres@postgres:5432/chat_assistant?sslmode=disable",
      "up"
    ]
    volumes:
      - ./internal/migrations:/migrations
    depends_on:
      postgres:
        condition: service_healthy
    profiles:
      - migrate

volumes:
  postgres_data:
    driver: local
```

### 关键配置说明

1. **健康检查**：确保数据库完全启动后再启动应用
2. **数据持久化**：使用 volume 保存数据库数据
3. **环境变量**：通过环境变量配置数据库连接
4. **服务依赖**：应用等待数据库健康检查通过
5. **Profiles**：迁移服务使用 profile 控制

## 🛠️ 常用命令

### 基础操作

```bash
# 启动服务
docker-compose up postgres -d          # 只启动数据库
docker-compose up -d                   # 启动所有服务
docker-compose up                      # 启动并查看日志

# 停止服务
docker-compose down                    # 停止所有服务
docker-compose stop                    # 暂停服务
docker-compose restart postgres        # 重启数据库

# 查看状态
docker-compose ps                      # 查看服务状态
docker-compose logs -f postgres        # 查看数据库日志
docker-compose logs -f chat-assistant-backend  # 查看应用日志
```

### 数据库操作

```bash
# 连接数据库
docker-compose exec postgres psql -U postgres -d chat_assistant

# 备份数据库
docker-compose exec postgres pg_dump -U postgres chat_assistant > backup.sql

# 恢复数据库
docker-compose exec -T postgres psql -U postgres -d chat_assistant < backup.sql

# 查看数据库大小
docker-compose exec postgres psql -U postgres -d chat_assistant -c "SELECT pg_size_pretty(pg_database_size('chat_assistant'));"
```

### 数据管理

```bash
# 查看数据卷
docker volume ls

# 删除数据卷（会丢失所有数据）
docker volume rm chat-assistant-backend_postgres_data

# 备份数据卷
docker run --rm -v chat-assistant-backend_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_data.tar.gz -C /data .

# 恢复数据卷
docker run --rm -v chat-assistant-backend_postgres_data:/data -v $(pwd):/backup alpine tar xzf /backup/postgres_data.tar.gz -C /data
```

### 迁移操作

```bash
# 执行迁移
docker-compose --profile migrate up migrate

# 查看迁移状态
docker-compose exec postgres psql -U postgres -d chat_assistant -c "SELECT * FROM goose_db_version;"

# 回滚迁移
docker-compose exec postgres psql -U postgres -d chat_assistant -c "DELETE FROM goose_db_version WHERE version_id = (SELECT MAX(version_id) FROM goose_db_version);"
```

## 💡 最佳实践

### 1. 开发环境

```bash
# 推荐：只启动数据库，本地运行应用
docker-compose up postgres -d
make run-dev

# 不推荐：启动完整环境进行开发
docker-compose up -d  # 会构建镜像，较慢
```

### 2. 环境隔离

```bash
# 使用不同的 compose 文件
docker-compose -f docker-compose.dev.yml up -d      # 开发环境
docker-compose -f docker-compose.test.yml up -d     # 测试环境
docker-compose -f docker-compose.prod.yml up -d     # 生产环境
```

### 3. 数据管理

```bash
# 定期备份数据
docker-compose exec postgres pg_dump -U postgres chat_assistant > backup_$(date +%Y%m%d).sql

# 清理旧数据
docker system prune -f
docker volume prune -f
```

### 4. 性能优化

```bash
# 限制资源使用
docker-compose up -d --scale postgres=1

# 监控资源使用
docker stats
```

## 🚨 故障排除

### 常见问题

#### 1. 端口冲突

```bash
# 检查端口占用
lsof -i :5432
lsof -i :8080

# 修改端口
# 在 docker-compose.yaml 中修改 ports 配置
ports:
  - "5433:5432"  # 使用 5433 端口
```

#### 2. 数据库连接失败

```bash
# 检查数据库状态
docker-compose ps postgres

# 查看数据库日志
docker-compose logs postgres

# 重启数据库
docker-compose restart postgres

# 检查网络连接
docker-compose exec postgres ping postgres
```

#### 3. 数据丢失

```bash
# 检查数据卷
docker volume ls | grep postgres

# 恢复数据
docker volume rm chat-assistant-backend_postgres_data
docker-compose up postgres -d
make migrate-up
```

#### 4. 镜像构建失败

```bash
# 清理构建缓存
docker builder prune -f

# 重新构建
docker-compose build --no-cache

# 查看构建日志
docker-compose build --progress=plain
```

#### 5. 服务启动失败

```bash
# 查看详细日志
docker-compose logs --tail=100 postgres

# 检查配置文件
docker-compose config

# 验证服务定义
docker-compose config --services
```

### 调试技巧

```bash
# 进入容器调试
docker-compose exec postgres bash
docker-compose exec chat-assistant-backend sh

# 查看容器资源使用
docker stats

# 查看网络配置
docker network ls
docker network inspect chat-assistant-backend_default
```

## 📊 监控和维护

### 健康检查

```bash
# 检查服务健康状态
docker-compose ps

# 手动健康检查
curl http://localhost:8080/health

# 数据库连接测试
docker-compose exec postgres pg_isready -U postgres
```

### 日志管理

```bash
# 查看实时日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f postgres

# 限制日志行数
docker-compose logs --tail=50 postgres

# 保存日志到文件
docker-compose logs postgres > postgres.log
```

### 性能监控

```bash
# 查看资源使用
docker stats

# 查看容器详细信息
docker inspect chat-assistant-postgres

# 查看数据卷使用情况
docker system df -v
```

---

**提示**：在生产环境中，建议使用专门的监控工具如 Prometheus + Grafana 来监控容器和应用的性能指标。
