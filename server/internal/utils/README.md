# Utils 工具函数目录
这个目录用于存放服务端的工具函数。

## 用途

- 字符串处理工具
- 时间格式化工具
- 数据验证工具
- 其他通用工具函数

## 示例

```go
package utils

import "time"

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
    return t.Format("2006-01-02 15:04:05")
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
    // 实现邮箱验证逻辑
    return true
}
```

## 注意

- 工具函数应该是无状态的
- 避免在工具函数中包含业务逻辑
- 保持函数简单和可测试
