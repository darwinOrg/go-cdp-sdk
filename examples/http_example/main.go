package main

import (
	"fmt"
	"log"

	"github.com/darwinOrg/go-cdp-sdk"
)

func main() {
	// 创建 HTTP 客户端
	sessionId := "http-example-session"
	client := cdpsdk.NewHTTPClient("http://localhost:3000", sessionId)

	// 启动浏览器
	fmt.Println("📌 启动浏览器...")
	if err := client.StartBrowser(false); err != nil {
		log.Fatalf("❌ 启动浏览器失败: %v", err)
	}
	fmt.Println("✅ 浏览器已启动")

	// 创建页面（使用默认页面）
	page := cdpsdk.NewPage(client)

	// 导航到 example.com
	fmt.Println("\n📌 导航到 example.com...")
	if err := page.Navigate("https://example.com"); err != nil {
		log.Printf("❌ 导航失败: %v\n", err)
	} else {
		fmt.Println("✅ 导航成功")
	}

	// 等待页面加载完成
	fmt.Println("\n📌 等待页面加载完成...")
	if err := page.WaitForLoadStateLoad(); err != nil {
		log.Printf("❌ 等待加载失败: %v\n", err)
	} else {
		fmt.Println("✅ 页面加载完成")
	}

	// 获取页面标题
	fmt.Println("\n📌 获取页面标题...")
	title, err := page.GetTitle()
	if err != nil {
		log.Printf("❌ 获取标题失败: %v\n", err)
	} else {
		fmt.Printf("✅ 页面标题: %s\n", title)
	}

	// 获取页面 URL
	fmt.Println("\n📌 获取页面 URL...")
	url, err := page.GetURL()
	if err != nil {
		log.Printf("❌ 获取 URL 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 页面 URL: %s\n", url)
	}

	// 检查元素是否存在
	fmt.Println("\n📌 检查 h1 元素是否存在...")
	locator := page.Locator("h1")
	exists, err := locator.Exists()
	if err != nil {
		log.Printf("❌ 检查元素失败: %v\n", err)
	} else {
		fmt.Printf("✅ 元素存在: %v\n", exists)
	}

	// 获取元素文本
	fmt.Println("\n📌 获取 h1 元素的文本...")
	text, err := locator.Text()
	if err != nil {
		log.Printf("❌ 获取元素文本失败: %v\n", err)
	} else {
		fmt.Printf("✅ 元素文本: %s\n", text)
	}

	// 截图
	fmt.Println("\n📌 截图...")
	screenshotData, err := page.Screenshot("png")
	if err != nil {
		log.Printf("❌ 截图失败: %v\n", err)
	} else {
		fmt.Printf("✅ 截图成功（数据大小: %d 字节）\n", len(screenshotData))
	}

	// 停止浏览器
	fmt.Println("\n📌 停止浏览器...")
	if err := client.StopBrowser(); err != nil {
		log.Printf("❌ 停止浏览器失败: %v\n", err)
	} else {
		fmt.Println("✅ 浏览器已停止")
	}

	fmt.Println("\n✅ 测试完成！")
}
