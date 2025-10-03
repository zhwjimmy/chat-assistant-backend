# 数据导入功能

这个目录包含了数据导入功能的相关文件和示例。

## 目录结构

```
scripts/import/
├── sample_data/              # 示例数据文件
│   ├── chatgpt_sample.json  # ChatGPT导出格式示例
│   ├── claude_sample.json   # Claude导出格式示例
│   └── gemini_sample.json   # Gemini导出格式示例
└── README.md                # 本文件
```

## 使用方法

### 1. 构建导入工具

```bash
make build-importer
```

### 2. 运行导入命令

```bash
# 导入ChatGPT数据
go run cmd/importer/main.go --platform=chatgpt --file=./scripts/import/sample_data/chatgpt_sample.json --user-id=123e4567-e89b-12d3-a456-426614174000

# 导入Claude数据
go run cmd/importer/main.go --platform=claude --file=./scripts/import/sample_data/claude_sample.json --user-id=123e4567-e89b-12d3-a456-426614174000

# 导入Gemini数据
go run cmd/importer/main.go --platform=gemini --file=./scripts/import/sample_data/gemini_sample.json --user-id=123e4567-e89b-12d3-a456-426614174000
```

### 3. 干运行（不写入数据库）

```bash
go run cmd/importer/main.go --platform=chatgpt --file=./scripts/import/sample_data/chatgpt_sample.json --user-id=123e4567-e89b-12d3-a456-426614174000 --dry-run
```

### 4. 详细日志

```bash
go run cmd/importer/main.go --platform=chatgpt --file=./scripts/import/sample_data/chatgpt_sample.json --user-id=123e4567-e89b-12d3-a456-426614174000 --verbose
```

## 支持的平台

- **chatgpt**: ChatGPT导出格式
- **claude**: Claude导出格式  
- **gemini**: Gemini导出格式

## 数据格式

### 标准化格式

所有平台的数据都会被转换为以下标准化格式：

```json
{
  "conversations": [
    {
      "id": "conversation_id",
      "title": "对话标题",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "provider": "平台名称",
      "model": "模型名称",
      "messages": [
        {
          "id": "message_id",
          "role": "user|assistant|system",
          "content": "消息内容",
          "created_at": "2024-01-01T00:00:00Z"
        }
      ]
    }
  ]
}
```

## 注意事项

1. **用户ID**: 必须提供有效的UUID格式的用户ID
2. **文件格式**: 支持JSON格式的导出文件
3. **数据验证**: 导入前会进行数据格式验证
4. **事务处理**: 使用数据库事务确保数据一致性
5. **错误处理**: 详细的错误信息和日志记录

## 扩展新平台

要添加新的平台支持，需要：

1. 在 `internal/importer/parsers/` 下创建新的平台目录
2. 实现 `Parser` 接口
3. 在 `registry.go` 中注册新解析器
4. 添加对应的类型定义

## 故障排除

### 常见错误

1. **文件不存在**: 检查文件路径是否正确
2. **无效的用户ID**: 确保用户ID是有效的UUID格式
3. **不支持的平台**: 检查平台名称是否正确
4. **数据格式错误**: 检查JSON文件格式是否正确

### 日志查看

使用 `--verbose` 参数可以查看详细的导入日志，帮助诊断问题。