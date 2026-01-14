# Go CDP SDK - HTTP 客户端使用指南

## 概述

`go-cdp-sdk` 提供了 HTTP 和 WebSocket 两种客户端来调用 TypeScript CDP 服务。HTTP 客户端更适合大多数自动化场景，因为它简单易用、调试方便。

## HTTP vs WebSocket 对比

| 特性 | HTTP 客户端 | WebSocket 客户端 |
|------|------------|-----------------|
| 简单性 | ✅ 非常简单 | ⚠️ 需要维护连接 |
| 调试 | ✅ 可用 curl 测试 | ❌ 需要专用工具 |
| 状态管理 | ✅ 无状态 | ⚠️ 需要维护连接状态 |
| 实时事件 | ❌ 不支持 | ✅ 支持事件推送 |
| 适用场景 | 线性自动化 | 需要实时监听 |

## 快速开始

### 1. 安装 SDK

```bash
go get github.com/darwinOrg/go-cdp-sdk
```

### 2. 启动 TypeScript CDP HTTP 服务

```bash
cd /path/to/ts-cdp
npm run server
```

服务将在 `http://localhost:3000` 启动。

### 3. 创建 HTTP 客户端

```go
package main

import (
    "log"
    cdpsdk "github.com/darwinOrg/go-cdp-sdk"
)

func main() {
    // 创建客户端
    client := cdpsdk.NewHTTPClient("http://localhost:3000", "my-session")

    // 连接到现有浏览器（9222 端口）
    resp, err := client.ConnectBrowser(9222)
    if err != nil {
        log.Fatal(err)
    }

    // 使用客户端...
}
```

## 基本用法

### 连接到浏览器

```go
// 连接到现有浏览器
resp, err := client.ConnectBrowser(9222)
if err != nil {
    log.Fatal(err)
}

// 或者启动新浏览器
resp, err := client.StartBrowser(false) // false = 显示浏览器窗口
```

### 页面操作

```go
pageID := "default"

// 导航到 URL
resp, err := client.Navigate(pageID, "https://example.com")

// 获取页面标题
resp, err := client.GetTitle(pageID)
title := resp.Data["title"]

// 获取页面 URL
resp, err := client.GetURL(pageID)
url := resp.Data["url"]

// 刷新页面
resp, err := client.Reload(pageID)

// 执行 JavaScript
resp, err := client.ExecuteScript(pageID, "document.title")

// 截图
screenshotData, err := client.Screenshot(pageID, "png")
// screenshotData 是 []byte 类型的图片数据
// 可以保存到文件:
// err := os.WriteFile("screenshot.png", screenshotData, 0644)
```

### 元素操作

```go
// 检查元素是否存在
resp, err := client.ElementExists(pageID, "#element-id")
exists := resp.Data["exists"]

// 获取元素文本
resp, err := client.ElementText(pageID, "h1")
text := resp.Data["text"]

// 点击元素
resp, err := client.ElementClick(pageID, "#button")

// 设置输入框值
resp, err := client.ElementSetValue(pageID, "#input", "hello world")

// 等待元素出现
resp, err := client.ElementWait(pageID, "#loading", 10000)

// 获取元素属性
resp, err := client.ElementAttribute(pageID, "#link", "href")
href := resp.Data["value"]
```

### 高级功能

```go
// 导航并等待加载完成
resp, err := client.NavigateWithLoadedState(pageID, "https://example.com")

// 等待页面加载完成
resp, err := client.WaitForLoadStateLoad(pageID)

// 等待 DOM 加载完成
resp, err := client.WaitForDomContentLoaded(pageID)

// 等待元素可见
resp, err := client.WaitForSelectorVisible(pageID, "#element")

// 随机等待（模拟人类行为）
resp, err := client.RandomWait(pageID, "middle") // short/middle/long

// 获取所有匹配元素的文本
resp, err := client.ElementAllTexts(pageID, ".item")
texts := resp.Data["texts"]

// 获取元素数量
resp, err := client.ElementCount(pageID, ".item")
count := resp.Data["count"]
```

### 多页面管理

```go
// 创建新页面
resp, err := client.NewPage("page-1")

// 在不同页面操作
client.Navigate("page-1", "https://example.com")
client.Navigate("page-2", "https://google.com")

// 关闭页面
resp, err := client.ClosePage("page-1")
```

## 完整示例

### 示例 1: 基础自动化

```go
package main

import (
    "log"
    cdpsdk "github.com/darwinOrg/go-cdp-sdk"
)

func main() {
    client := cdpsdk.NewHTTPClient("http://localhost:3000", "")

    // 连接到浏览器
    if _, err := client.ConnectBrowser(9222); err != nil {
        log.Fatal(err)
    }

    pageID := "default"

    // 导航
    client.Navigate(pageID, "https://example.com")

    // 获取标题
    resp, _ := client.GetTitle(pageID)
    log.Printf("Title: %v", resp.Data["title"])

    // 停止
    client.StopBrowser()
}
```

### 示例 2: 表单提交

```go
// 导航到表单页面
client.Navigate(pageID, "https://example.com/form")

// 填写表单
client.ElementSetValue(pageID, "#username", "john")
client.ElementSetValue(pageID, "#password", "secret")

// 提交表单
client.ElementClick(pageID, "#submit")

// 等待结果
client.WaitForSelectorVisible(pageID, "#success-message")
```

### 示例 3: 数据抓取

```go
// 导航到列表页面
client.Navigate(pageID, "https://example.com/items")

// 获取所有项目标题
resp, _ := client.ElementAllTexts(pageID, ".item-title")
titles := resp.Data["texts"].([]interface{})

// 获取所有项目链接
resp, _ = client.ElementAllAttributes(pageID, ".item-title", "href")
links := resp.Data["attributes"].([]interface{})

// 处理数据
for i, title := range titles {
    log.Printf("%d: %s -> %s", i+1, title, links[i])
}
```

## API 参考

### 浏览器管理

- `ConnectBrowser(port int)` - 连接到现有浏览器
- `StartBrowser(headless bool)` - 启动新浏览器
- `StopBrowser()` - 停止浏览器

### 页面管理

- `NewPage(pageID string)` - 创建新页面
- `ClosePage(pageID string)` - 关闭页面
- `Navigate(pageID, url string)` - 导航到 URL
- `NavigateWithLoadedState(pageID, url string)` - 导航并等待加载
- `Reload(pageID string)` - 刷新页面
- `GetTitle(pageID string)` - 获取页面标题
- `GetURL(pageID string)` - 获取页面 URL
- `GetHTML(pageID string)` - 获取页面 HTML
- `ExecuteScript(pageID, script string)` - 执行 JavaScript
- `Screenshot(pageID, format string) ([]byte, error)` - 截图，返回图片的二进制数据

### 元素操作

- `ElementExists(pageID, selector string)` - 检查元素是否存在
- `ElementText(pageID, selector string)` - 获取元素文本
- `ElementClick(pageID, selector string)` - 点击元素
- `ElementSetValue(pageID, selector, value string)` - 设置元素值
- `ElementWait(pageID, selector string, timeout int)` - 等待元素
- `ElementAttribute(pageID, selector, attribute string)` - 获取元素属性
- `ElementAllTexts(pageID, selector string)` - 获取所有匹配元素的文本
- `ElementAllAttributes(pageID, selector, attribute string)` - 获取所有匹配元素的属性
- `ElementCount(pageID, selector string)` - 获取元素数量

### 等待功能

- `WaitForLoadStateLoad(pageID string)` - 等待页面加载
- `WaitForDomContentLoaded(pageID string)` - 等待 DOM 加载
- `WaitForSelectorVisible(pageID, selector string)` - 等待元素可见
- `RandomWait(pageID string, duration interface{})` - 随机等待

### 高级功能

- `MustInnerText(pageID, selector string)` - 强制获取内部文本
- `MustTextContent(pageID, selector string)` - 强制获取文本内容
- `ExpectResponseText(pageID, urlOrPredicate, callback string)` - 等待响应文本
- `ExpectExtPage(pageID, callback string)` - 等待新页面
- `Suspend(pageID string)` - 暂停页面
- `Continue(pageID string)` - 继续页面
- `Release(pageID string)` - 释放页面
- `CloseAll(pageID string)` - 关闭所有页面

## 错误处理

```go
resp, err := client.Navigate(pageID, "https://example.com")
if err != nil {
    log.Printf("请求失败: %v", err)
    return
}

if !resp.Success {
    log.Printf("操作失败: %s", resp.Error)
    return
}

// 处理成功响应
log.Printf("成功: %v", resp.Data)
```

## 配置选项

### 设置超时时间

```go
client := cdpsdk.NewHTTPClient("http://localhost:3000", "")
client.SetTimeout(60 * time.Second) // 设置为 60 秒
```

### 使用自定义会话 ID

```go
client := cdpsdk.NewHTTPClient("http://localhost:3000", "my-custom-session-id")
```

## 运行示例

```bash
# 基础示例
go run examples/http_example.go

# 实用示例
go run examples/http_practical_example.go
```

## 注意事项

1. **确保 TypeScript CDP 服务正在运行**：HTTP 客户端需要连接到运行中的服务
2. **浏览器端口**：默认使用 9222 端口连接 Chrome
3. **会话管理**：每个客户端实例有一个唯一的 sessionId
4. **页面 ID**：默认使用 "default" 页面，也可以创建新页面
5. **错误处理**：始终检查返回的错误和响应的 Success 字段

## 最佳实践

1. **使用 RandomWait**：在操作之间添加随机等待，模拟人类行为
2. **等待元素**：使用 WaitForSelectorVisible 等待元素出现后再操作
3. **错误处理**：始终检查错误并适当处理
4. **资源清理**：使用完毕后调用 StopBrowser 释放资源
5. **使用 pageId**：在多页面场景中明确指定 pageId

## 与 WebSocket 客户端对比

如果你的需求是：
- ✅ 简单的线性自动化流程
- ✅ 需要快速开发和调试
- ✅ 不需要实时事件监听

**使用 HTTP 客户端**

如果你的需求是：
- ⚠️ 需要监听浏览器事件（console、error、dialog）
- ⚠️ 需要实时响应页面变化
- ⚠️ 高频操作（每秒多次请求）

**考虑使用 WebSocket 客户端**

## 许可证

MIT License