# Data Sync Tool

数据同步工具，用于将数据库中的数据同步到 Elasticsearch。

## 功能

- 全量同步 conversations 和 messages 到 Elasticsearch
- 支持试运行模式，查看同步统计信息
- 简单易用的命令行界面

## 使用方法

### 构建工具

```bash
make build-data-sync
```

### 同步数据

```bash
# 同步所有数据到 Elasticsearch
make sync-data

# 或者直接运行
./bin/chat-assistant-data-sync
```

### 试运行

```bash
# 试运行，查看将要同步的数据统计
make sync-data-dry

# 或者直接运行
./bin/chat-assistant-data-sync -dry-run
```

### 查看帮助

```bash
./bin/chat-assistant-data-sync -help
```

## 工作流程

1. **连接数据库**: 从 PostgreSQL 读取所有 conversations 和 messages
2. **数据转换**: 将数据库模型转换为 Elasticsearch 文档格式
3. **批量索引**: 使用 Elasticsearch 的批量 API 进行索引
4. **结果报告**: 显示同步结果和统计信息

## 注意事项

- 确保 Elasticsearch 服务正在运行
- 确保数据库连接正常
- 同步前建议先运行试运行模式查看数据统计
- 全量同步会覆盖 Elasticsearch 中的现有数据

## 错误处理

- 如果数据库连接失败，工具会退出并显示错误信息
- 如果 Elasticsearch 连接失败，工具会退出并显示错误信息
- 如果同步过程中出现错误，工具会显示详细的错误信息

## 性能

- 适用于几千条数据的小规模同步
- 全量同步通常在几秒内完成
- 内存使用量取决于数据量大小
