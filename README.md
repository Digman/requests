# Requests

一个强大的 Go HTTP 客户端库，支持 TLS 指纹模拟，基于 [tls-client](https://github.com/Digman/tls-client) 构建。

## 特性

- **TLS 指纹模拟** - 支持 Chrome (103-133)、Firefox (102-135)、Safari (15-18) 浏览器指纹
- **代理支持** - 支持 HTTP、HTTPS 和 SOCKS5 代理
- **Cookie 管理** - 完整的 Cookie jar 支持，支持基于域名的操作
- **自定义 Headers** - 支持设置自定义请求头和请求头顺序
- **证书固定** - 可选的证书固定功能，增强安全性
- **链式 API** - 流畅的链式调用接口构建请求
- **多种数据格式** - 支持 URL 参数、JSON 和 multipart 文件上传
- **灵活的响应处理** - 支持获取字符串、字节、JSON 格式响应或保存到文件
- **自动重定向控制** - 可启用或禁用自动重定向
- **调试模式** - 内置请求调试功能

## 安装

```bash
go get github.com/Digman/requests
```

## 快速开始

### 基本 GET 请求

```go
package main

import (
    "fmt"
    "github.com/Digman/requests"
)

func main() {
    client := requests.DefaultClient()
    _, body, err := client.NewRequest().
        Get("https://httpbin.org/get").
        Send().
        End()

    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(body)
}
```

### POST JSON 数据

```go
client := requests.DefaultClient()

data := map[string]interface{}{
    "username": "testuser",
    "password": "123456",
}

_, resp, err := client.NewRequest().
    Post("https://httpbin.org/post").
    SetJson(data).
    Send().
    EndJson()

if err != nil {
    fmt.Println("Error:", err)
    return
}

fmt.Println(resp.Get("json.username").String())
```

### 使用代理

```go
client := requests.DefaultClient()
err := client.SetProxy("http://127.0.0.1:8080")
if err != nil {
    fmt.Println("Proxy error:", err)
    return
}

_, body, _ := client.NewRequest().
    Get("https://httpbin.org/ip").
    Send().
    End()

fmt.Println(body)
```

### 自定义请求头

```go
client := requests.DefaultClient()

_, body, err := client.NewRequest().
    Get("https://httpbin.org/headers").
    SetHeader("Referer", "https://google.com").
    SetHeader("X-Custom-Header", "custom-value").
    Send().
    End()
```

## API 文档

### Client API

#### 客户端创建

| 方法 | 说明 | 参数 |
|------|------|------|
| `NewClient(userAgent string, cp ...*CertPinning) *Client` | 使用自定义用户代理创建新客户端 | `userAgent`: 自定义用户代理字符串<br>`cp`: 可选的证书固定配置 |
| `DefaultClient() *Client` | 使用默认设置创建客户端 | 无 |
| `TimeoutClient(timeout int) *Client` | 创建带自定义超时的客户端 | `timeout`: 超时时间（毫秒） |

**示例:**
```go
// 默认客户端
client := requests.DefaultClient()

// 自定义用户代理
client := requests.NewClient("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

// 自定义超时（10秒）
client := requests.TimeoutClient(10000)
```

#### 客户端配置

| 方法 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `SetProxy(proxyUrl string) error` | 设置 HTTP/HTTPS/SOCKS5 代理 | `proxyUrl`: 代理 URL（如 "http://127.0.0.1:8080"） | error |
| `SetKeepAlive(b bool)` | 启用或禁用 keep-alive | `b`: true 启用，false 禁用 | - |
| `SetAutoRedirect(b bool)` | 启用或禁用自动重定向 | `b`: true 启用，false 禁用 | - |

**示例:**
```go
client := requests.DefaultClient()
client.SetProxy("socks5://127.0.0.1:1080")
client.SetKeepAlive(true)
client.SetAutoRedirect(false)
```

#### Cookie 管理

| 方法 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `NewCookies()` | 创建新的 cookie jar | 无 | - |
| `SetCookies(domain, cookies)` | 为域名设置 cookies | `domain`: 域名<br>`cookies`: Cookie 切片 | - |
| `SetUrlCookies(cookieUrl, cookies)` | 为 URL 设置 cookies | `cookieUrl`: 完整 URL<br>`cookies`: Cookie 切片 | - |
| `GetCookies(domain) []*http.Cookie` | 获取域名的 cookies | `domain`: 域名 | Cookie 切片 |
| `GetUrlCookies(cookieUrl) []*http.Cookie` | 获取 URL 的 cookies | `cookieUrl`: 完整 URL | Cookie 切片 |

**示例:**
```go
client := requests.DefaultClient()

// 设置 cookies
cookies := []*http.Cookie{
    {Name: "session", Value: "abc123"},
    {Name: "user_id", Value: "12345"},
}
client.SetCookies("example.com", cookies)

// 获取 cookies
retrievedCookies := client.GetCookies("example.com")
```

#### 工具方法

| 方法 | 说明 | 返回值 |
|------|------|--------|
| `GetRequestInfo() (bool, string)` | 测试请求到 httpbin.org/get | 成功标志，响应内容 |
| `GetFingerPrint() (bool, string)` | 从 tls.peet.ws 获取 TLS 指纹 | 成功标志，指纹数据 |
| `GetIPLocation() (bool, string)` | 从 ip-api.com 获取 IP 位置 | 成功标志，位置数据 |

#### 请求创建

| 方法 | 说明 | 返回值 |
|------|------|--------|
| `NewRequest() *Request` | 创建带默认请求头的新请求 | `*Request` |

### Request API

#### HTTP 方法

| 方法 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `Get(url string) *Request` | 创建 GET 请求 | `url`: 目标 URL | `*Request` |
| `Post(url string) *Request` | 创建 POST 请求 | `url`: 目标 URL | `*Request` |
| `Put(url string) *Request` | 创建 PUT 请求 | `url`: 目标 URL | `*Request` |
| `Head(url string) *Request` | 创建 HEAD 请求 | `url`: 目标 URL | `*Request` |
| `Options(url string) *Request` | 创建 OPTIONS 请求 | `url`: 目标 URL | `*Request` |
| `SetMethod(name string) *Request` | 设置自定义 HTTP 方法 | `name`: HTTP 方法名 | `*Request` |

**示例:**
```go
client := requests.DefaultClient()

// GET 请求
client.NewRequest().Get("https://api.example.com/users")

// POST 请求
client.NewRequest().Post("https://api.example.com/users")

// 自定义方法
client.NewRequest().SetMethod("PATCH").SetUrl("https://api.example.com/users/1")
```

#### 请求头

| 方法 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `SetHeader(name, value string) *Request` | 设置单个请求头 | `name`: 请求头名称<br>`value`: 请求头值 | `*Request` |
| `SetHeaders(values map[string]string) *Request` | 批量设置请求头 | `values`: 请求头映射 | `*Request` |
| `SetHeaderOrder(order []string) *Request` | 设置请求头顺序 | `order`: 请求头名称切片 | `*Request` |

**示例:**
```go
client.NewRequest().
    Get("https://api.example.com/data").
    SetHeader("Authorization", "Bearer token123").
    SetHeader("Content-Type", "application/json").
    Send()

// 批量设置请求头
headers := map[string]string{
    "Authorization": "Bearer token123",
    "X-API-Key": "key123",
}
client.NewRequest().
    Get("https://api.example.com/data").
    SetHeaders(headers).
    Send()

// 自定义请求头顺序
order := []string{"authorization", "content-type", "accept"}
client.NewRequest().
    Post("https://api.example.com/data").
    SetHeaderOrder(order).
    Send()
```

#### 身份认证

| 方法 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `SetBasicAuth(userName, password string) *Request` | 设置 HTTP 基本认证 | `userName`: 用户名<br>`password`: 密码 | `*Request` |

**示例:**
```go
client.NewRequest().
    Get("https://api.example.com/private").
    SetBasicAuth("user", "pass").
    Send()
```

#### 请求数据

##### URL 参数 / 表单数据

| 方法 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `SetData(name, value string) *Request` | 设置单个数据字段 | `name`: 字段名<br>`value`: 字段值 | `*Request` |
| `SetAllData(data url.Values) *Request` | 一次性设置所有数据 | `data`: url.Values | `*Request` |
| `GetAllData() url.Values` | 获取所有当前数据 | 无 | `url.Values` |

**示例:**
```go
// GET 请求带查询参数
client.NewRequest().
    Get("https://api.example.com/search").
    SetData("q", "golang").
    SetData("page", "1").
    Send()
// 结果为: https://api.example.com/search?q=golang&page=1

// POST 表单数据
client.NewRequest().
    Post("https://api.example.com/login").
    SetData("username", "user").
    SetData("password", "pass").
    Send()
```

##### JSON 数据

| 方法 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `SetJsonData(s string) *Request` | 设置原始 JSON 字符串 | `s`: JSON 字符串 | `*Request` |
| `SetJson(data any) *Request` | 从结构体/映射设置 JSON | `data`: 任何可序列化类型 | `*Request` |

**示例:**
```go
// 使用 JSON 字符串
jsonStr := `{"username":"test","password":"123"}`
client.NewRequest().
    Post("https://api.example.com/login").
    SetJsonData(jsonStr).
    Send()

// 使用结构体/映射
data := map[string]interface{}{
    "username": "test",
    "password": "123",
}
client.NewRequest().
    Post("https://api.example.com/login").
    SetJson(data).
    Send()
```

##### 文件上传

| 方法 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `SetFileData(name, value string, isFile bool) *Request` | 添加文件或表单字段 | `name`: 字段名<br>`value`: 文件路径或字段值<br>`isFile`: true 表示文件，false 表示文本字段 | `*Request` |

**示例:**
```go
client.NewRequest().
    Post("https://api.example.com/upload").
    SetFileData("file", "/path/to/image.png", true).
    SetFileData("description", "My image", false).
    SetFileData("category", "photos", false).
    Send("file")
```

#### Cookies

| 方法 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `SetCookies(cookies *[]*http.Cookie) *Request` | 为此请求设置 cookies | `cookies`: Cookie 切片指针 | `*Request` |

#### 调试

| 方法 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `SetDebug(d bool) *Request` | 启用调试输出 | `d`: true 启用 | `*Request` |

#### 执行

| 方法 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `Send(a ...interface{}) *Request` | 发送请求 | `a`: 可选的数据类型（"url"、"json"、"file"） | `*Request` |

**注意:** 数据类型会根据您使用的 Set 方法自动检测。只有在文件上传时需要手动指定:
```go
// 自动检测为 "json"
req.SetJson(data).Send()

// 文件上传必须指定 "file"
req.SetFileData("file", "/path/file.txt", true).Send("file")
```

#### 响应处理

| 方法 | 说明 | 返回值 |
|------|------|--------|
| `End() (*http.Response, string, error)` | 获取字符串格式响应 | response, 响应内容字符串, error |
| `EndJson() (*http.Response, gjson.Result, error)` | 获取 JSON 格式响应 | response, gjson.Result, error |
| `EndByte() (*http.Response, []byte, error)` | 获取字节格式响应 | response, 响应内容字节, error |
| `EndResponse() (*http.Response, error)` | 仅获取响应（自动读取 body 以支持 keep-alive） | response, error |
| `EndFile(savePath, saveFileName string) (*http.Response, error)` | 保存响应到文件 | response, error |

**示例:**
```go
// 字符串响应
_, body, err := client.NewRequest().
    Get("https://api.example.com/data").
    Send().
    End()

// JSON 响应，使用 gjson
_, result, err := client.NewRequest().
    Get("https://api.example.com/users/1").
    Send().
    EndJson()
username := result.Get("name").String()
age := result.Get("age").Int()

// 字节响应（用于二进制数据）
_, data, err := client.NewRequest().
    Get("https://example.com/image.png").
    Send().
    EndByte()

// 下载文件
_, err := client.NewRequest().
    Get("https://example.com/file.zip").
    Send().
    EndFile("/downloads/", "myfile.zip")
```

## 高级示例

### 完整的 API 请求与认证

```go
client := requests.DefaultClient()

// 登录并获取 token
loginData := map[string]interface{}{
    "username": "user@example.com",
    "password": "secret123",
}

_, loginResp, err := client.NewRequest().
    Post("https://api.example.com/auth/login").
    SetJson(loginData).
    Send().
    EndJson()

if err != nil {
    panic(err)
}

token := loginResp.Get("token").String()

// 使用 token 进行认证请求
_, userData, err := client.NewRequest().
    Get("https://api.example.com/user/profile").
    SetHeader("Authorization", "Bearer "+token).
    Send().
    EndJson()

fmt.Println(userData.Get("email").String())
```

### 文件上传

```go
client := requests.DefaultClient()

_, resp, err := client.NewRequest().
    Post("https://api.example.com/upload").
    SetFileData("document", "/path/to/document.pdf", true).
    SetFileData("title", "My Document", false).
    SetFileData("visibility", "private", false).
    SetHeader("Authorization", "Bearer token123").
    Send("file").
    EndJson()

if err != nil {
    fmt.Println("上传失败:", err)
    return
}

fileId := resp.Get("file_id").String()
fmt.Println("文件上传成功，ID:", fileId)
```

### 使用自定义浏览器指纹

```go
// Chrome 133
userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
client := requests.NewClient(userAgent)

// Firefox 135
userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:135.0) Gecko/20100101 Firefox/135.0"
client := requests.NewClient(userAgent)

// Safari iOS 18
userAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Mobile/15E148 Safari/604.1"
client := requests.NewClient(userAgent)
```

### 证书固定

```go
certPinning := &requests.CertPinning{
    CertificatePins: map[string][]string{
        "example.com": {
            "sha256/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
            "sha256/BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB=",
        },
    },
    BadPinHandler: func(host string, actualPins []string) error {
        fmt.Printf("证书固定不匹配: %s\n", host)
        return fmt.Errorf("证书验证失败")
    },
}

client := requests.NewClient("Mozilla/5.0 ...", certPinning)
```

### 会话管理

```go
client := requests.DefaultClient()

// 登录以建立会话
client.NewRequest().
    Post("https://example.com/login").
    SetData("username", "user").
    SetData("password", "pass").
    Send().
    End()

// Cookies 自动维护
// 后续请求将自动包含会话 cookies
_, body, _ := client.NewRequest().
    Get("https://example.com/dashboard").
    Send().
    End()

// 获取当前 cookies
cookies := client.GetCookies("example.com")
for _, cookie := range cookies {
    fmt.Printf("%s = %s\n", cookie.Name, cookie.Value)
}
```

### 调试请求

```go
client := requests.DefaultClient()

_, body, err := client.NewRequest().
    Post("https://httpbin.org/post").
    SetData("key", "value").
    SetDebug(true).  // 启用调试输出
    Send().
    End()

// 输出显示:
// [Request Debug]
// -------------------------------------------------------------------
// Request: POST https://httpbin.org/post
// Header: map[...]
// Cookies: [...]
// Body: key=value
// -------------------------------------------------------------------
```

## 支持的浏览器指纹

### Chrome
103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 133

### Firefox
102, 104, 105, 106, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 132, 133, 135

### Safari
- macOS: Version 15, 16
- iOS: 15.5, 15.6, 16.x, 17.x, 18.0, 18.x
- iPad: 18.x

## 依赖

- [tls-client](https://github.com/Digman/tls-client) - TLS 指纹 HTTP 客户端
- [fhttp](https://github.com/bogdanfinn/fhttp) - 支持请求头排序的 HTTP/2
- [gjson](https://github.com/tidwall/gjson) - JSON 解析

## 许可证

本项目继承其依赖项的许可证。
