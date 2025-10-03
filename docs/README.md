# 文档目录

本目录包含 Chat Assistant Backend 项目的详细文档。

## 📚 文档列表

- [开发指南](development.md) - 完整的开发流程和环境配置
- [Docker Compose 使用指南](docker-compose-guide.md) - Docker Compose 的详细使用说明
- [Makefile 命令参考](makefile-commands.md) - 所有 Makefile 命令的详细说明
- [环境变量配置指南](environment-configuration.md) - 环境变量配置和最佳实践

## 🚀 快速开始

### 本地开发（推荐）

```bash
# 1. 启动数据库
docker-compose up postgres -d

# 2. 运行应用
make run-dev

# 3. 测试
curl http://localhost:8080/health
```

### 完整环境

```bash
# 启动所有服务
make docker-compose-up
```

## 📖 详细文档

请查看各个文档文件获取详细信息：

- **开发指南**：环境准备、开发流程、故障排除
- **Docker Compose 指南**：使用场景、配置说明、最佳实践
- **Makefile 命令**：所有可用命令的详细说明

## 🤝 贡献

如需更新文档，请：

1. 修改对应的 `.md` 文件
2. 提交 Pull Request
3. 确保文档格式正确

---

**Happy Coding! 🚀**
