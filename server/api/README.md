# API 文档目录
这个目录用于存放API文档和接口定义。

## 用途

- OpenAPI/Swagger 规范文件
- API 文档
- 接口示例
- Postman 集合

## 推荐工具

### Swagger/OpenAPI
```yaml
openapi: 3.0.0
info:
  title: IM System API
  version: 1.0.0
paths:
  /api/auth/register:
    post:
      summary: 用户注册
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                username:
                  type: string
                password:
                  type: string
```

### 文档生成

可以使用以下工具生成API文档：
- `swag` - Go Swagger 文档生成器
- `redoc` - API 文档渲染工具
- `postman` - API 测试和文档工具

## 使用方法

1. 安装 swag
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. 在代码中添加注释
```go
// @Summary 用户注册
// @Description 注册新用户
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 200 {object} AuthResponse
// @Router /api/auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
    // ...
}
```

3. 生成文档
```bash
swag init -g cmd/server/main.go -o api/docs
