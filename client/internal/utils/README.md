# Utils 工具函数目录
这个目录用于存放客户端后端的工具函数。

## 用途

- 字符串处理工具
- 数据转换工具
- 验证工具
- 其他通用工具函数

## 示例

```go
package utils

// Base64Encode Base64编码
func Base64Encode(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode Base64解码
func Base64Decode(str string) ([]byte, error) {
    return base64.StdEncoding.DecodeString(str)
}
```

## 注意

- 工具函数应该是无状态的
- 避免在工具函数中包含业务逻辑
- 保持函数简单和可测试
