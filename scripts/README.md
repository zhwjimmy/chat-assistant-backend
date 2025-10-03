# 数据库脚本说明

本目录包含聊天助手后端项目的数据库相关脚本。

## 脚本文件

### 1. `init_database.sql` - 数据库初始化脚本
- 创建所有必要的表
- 创建索引和触发器
- 使用 `IF NOT EXISTS` 避免重复创建
- **推荐用于开发和测试环境**

### 2. `create_tables.sql` - 完整表创建脚本
- 包含详细的注释和说明
- 适合生产环境使用
- 包含表注释和字段注释

### 3. `drop_tables.sql` - 删除表脚本
- 删除所有表、触发器、函数
- **注意：会删除所有数据，请谨慎使用！**
- 用于重置数据库环境

## 使用方法

### 初始化数据库
```bash
# 连接到 PostgreSQL 数据库
psql -h localhost -U postgres -d chat_assistant

# 执行初始化脚本
\i scripts/init_database.sql
```

### 重置数据库
```bash
# 删除所有表
\i scripts/drop_tables.sql

# 重新初始化
\i scripts/init_database.sql
```

## 数据库结构

### 表关系
```
users (1) ──→ (N) conversations (1) ──→ (N) messages
```

### 主要特性
- 使用 UUID 作为主键
- 支持软删除 (deleted_at 字段)
- 自动更新时间戳
- 无外键约束（通过程序保证数据完整性）
- 合理的索引设计

## 注意事项

1. 确保 PostgreSQL 已安装并运行
2. 确保已创建目标数据库
3. 确保有足够的数据库权限
4. 在生产环境使用前请备份数据
