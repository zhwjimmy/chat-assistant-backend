# Elasticsearch Infrastructure

这个包提供了 Elasticsearch 的基础设施组件，包括客户端封装、健康检查和配置管理。

## 组件

### Client
- 封装了 go-elasticsearch/v8 客户端
- 提供连接管理、健康检查、索引操作等功能
- 支持认证和超时配置

### HealthChecker
- 提供 Elasticsearch 集群健康检查
- 支持等待集群变为健康状态
- 返回详细的健康状态信息

### Config
- 统一的配置结构
- 支持多主机、认证、索引配置
- 提供默认配置

## 使用方法

### 1. 基本使用

```go
import "chat-assistant-backend/internal/infra/elasticsearch"

// 创建配置
cfg := &elasticsearch.Config{
    Hosts:    []string{"http://localhost:9200"},
    Username: "",
    Password: "",
    Timeout:  30 * time.Second,
    Index: elasticsearch.IndexConfig{
        Conversations: "conversations",
        Messages:      "messages",
    },
}

// 创建客户端
client, err := elasticsearch.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}

// 测试连接
ctx := context.Background()
if err := client.Ping(ctx); err != nil {
    log.Fatal(err)
}
```

### 2. 健康检查

```go
// 创建健康检查器
healthChecker := elasticsearch.NewHealthChecker(client)

// 检查健康状态
status := healthChecker.Check(ctx)
fmt.Printf("Status: %s\n", status.Status)

// 等待集群健康
if err := healthChecker.WaitForHealthy(ctx, 60*time.Second); err != nil {
    log.Fatal(err)
}
```

### 3. 索引操作

```go
// 检查索引是否存在
exists, err := client.IndexExists(ctx, "conversations")
if err != nil {
    log.Fatal(err)
}

// 创建索引
mapping := `{
    "mappings": {
        "properties": {
            "id": {"type": "keyword"},
            "title": {"type": "text"}
        }
    }
}`
if err := client.CreateIndex(ctx, "conversations", mapping); err != nil {
    log.Fatal(err)
}
```

## 配置

### 环境变量
- `ELASTICSEARCH_HOSTS`: 逗号分隔的主机列表
- `ELASTICSEARCH_USERNAME`: 用户名
- `ELASTICSEARCH_PASSWORD`: 密码
- `ELASTICSEARCH_TIMEOUT`: 超时时间

### 配置文件
```yaml
elasticsearch:
  hosts:
    - "http://localhost:9200"
  username: ""
  password: ""
  timeout: 30s
  index:
    conversations: "conversations"
    messages: "messages"
```

## 依赖注入

通过 Wire 进行依赖注入：

```go
// wire.go
func NewElasticsearchClient(cfg *config.Config) (*elasticsearch.Client, error) {
    esConfig := &elasticsearch.Config{
        Hosts:    cfg.Elasticsearch.Hosts,
        Username: cfg.Elasticsearch.Username,
        Password: cfg.Elasticsearch.Password,
        Timeout:  cfg.Elasticsearch.Timeout,
        Index: elasticsearch.IndexConfig{
            Conversations: cfg.Elasticsearch.Index.Conversations,
            Messages:      cfg.Elasticsearch.Index.Messages,
        },
    }
    return elasticsearch.NewClient(esConfig)
}
```

## 测试

运行测试需要 Elasticsearch 服务运行：

```bash
# 启动 Elasticsearch
docker-compose up -d elasticsearch

# 运行测试
go test ./internal/infra/elasticsearch/...
```

如果没有 Elasticsearch 服务，测试会跳过并记录日志。
