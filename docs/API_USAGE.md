# 用户详情查询API使用说明

## API端点

### GET /api/v1/users/{id}

获取指定用户的详细信息。

#### 请求参数

- **路径参数**:
  - `id` (string, required): 用户ID，必须是有效的UUID格式

#### 响应格式

**成功响应 (200 OK)**:
```json
{
  "success": true,
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "john_doe",
    "avatar": "https://example.com/avatar.jpg",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**错误响应**:

1. **无效UUID格式 (400 Bad Request)**:
```json
{
  "success": false,
  "error": {
    "code": "INVALID_UUID",
    "message": "Invalid user ID format",
    "details": "User ID must be a valid UUID"
  }
}
```

2. **用户不存在 (404 Not Found)**:
```json
{
  "success": false,
  "error": {
    "code": "USER_NOT_FOUND",
    "message": "User not found",
    "details": "No user found with the specified ID"
  }
}
```

3. **服务器内部错误 (500 Internal Server Error)**:
```json
{
  "success": false,
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "Internal server error",
    "details": "Failed to retrieve user"
  }
}
```

## 使用示例

### cURL示例

```bash
# 获取用户详情
curl -X GET "http://localhost:8080/api/v1/users/123e4567-e89b-12d3-a456-426614174000" \
  -H "Accept: application/json"
```

### JavaScript示例

```javascript
const userId = '123e4567-e89b-12d3-a456-426614174000';

fetch(`http://localhost:8080/api/v1/users/${userId}`)
  .then(response => response.json())
  .then(data => {
    if (data.success) {
      console.log('用户信息:', data.data);
    } else {
      console.error('错误:', data.error);
    }
  })
  .catch(error => {
    console.error('请求失败:', error);
  });
```

### Python示例

```python
import requests
import json

user_id = '123e4567-e89b-12d3-a456-426614174000'
url = f'http://localhost:8080/api/v1/users/{user_id}'

try:
    response = requests.get(url, headers={'Accept': 'application/json'})
    data = response.json()
    
    if data['success']:
        print('用户信息:', data['data'])
    else:
        print('错误:', data['error'])
        
except requests.exceptions.RequestException as e:
    print('请求失败:', e)
```

## 启动服务

1. 确保PostgreSQL数据库正在运行
2. 配置数据库连接信息（在`config/config.yaml`中）
3. 运行服务：

```bash
# 开发模式
make run

# 或者直接运行
go run cmd/server/main.go
```

## 健康检查

服务启动后，可以通过以下端点检查服务状态：

```bash
curl http://localhost:8080/health
```

响应：
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z",
  "service": "chat-assistant-backend"
}
```
