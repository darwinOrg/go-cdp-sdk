package main

import (
	"fmt"
	"log"

	"github.com/darwinOrg/go-cdp-sdk"
)

func main() {
	// 创建 HTTP 客户端
	client := cdpsdk.NewHTTPClient("http://localhost:3000", "")

	// 连接到现有浏览器（9222 端口）
	fmt.Println("📌 连接到现有浏览器（端口 9222）...")
	resp, err := client.ConnectBrowser(9222)
	if err != nil {
		log.Fatalf("❌ 连接浏览器失败: %v", err)
	}
	fmt.Printf("✅ 已连接到浏览器: sessionId=%s, port=%v\n", client.GetSessionID(), resp.Data["port"])

	// 创建新页面（可选，也可以使用默认页面）
	fmt.Println("\n📌 创建新页面...")
	pageID := "page-1"
	resp, err = client.NewPage(pageID)
	if err != nil {
		log.Printf("❌ 创建页面失败: %v\n", err)
	} else {
		fmt.Printf("✅ 页面已创建: %s\n", pageID)
	}

	// 导航到 example.com
	fmt.Println("\n📌 导航到 example.com...")
	resp, err = client.Navigate(pageID, "https://example.com")
	if err != nil {
		log.Printf("❌ 导航失败: %v\n", err)
	} else {
		fmt.Println("✅ 导航成功")
	}

	// 等待页面加载完成
	fmt.Println("\n📌 等待页面加载完成...")
	resp, err = client.WaitForLoadStateLoad(pageID)
	if err != nil {
		log.Printf("❌ 等待加载失败: %v\n", err)
	} else {
		fmt.Println("✅ 页面加载完成")
	}

	// 获取页面标题
	fmt.Println("\n📌 获取页面标题...")
	resp, err = client.GetTitle(pageID)
	if err != nil {
		log.Printf("❌ 获取标题失败: %v\n", err)
	} else {
		fmt.Printf("✅ 页面标题: %v\n", resp.Data["title"])
	}

	// 获取页面 URL
	fmt.Println("\n📌 获取页面 URL...")
	resp, err = client.GetURL(pageID)
	if err != nil {
		log.Printf("❌ 获取 URL 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 页面 URL: %v\n", resp.Data["url"])
	}

	// 检查元素是否存在
	fmt.Println("\n📌 检查 h1 元素是否存在...")
	resp, err = client.ElementExists(pageID, "h1")
	if err != nil {
		log.Printf("❌ 检查元素失败: %v\n", err)
	} else {
		fmt.Printf("✅ 元素存在: %v\n", resp.Data["exists"])
	}

	// 获取元素文本
	fmt.Println("\n📌 获取 h1 元素的文本...")
	resp, err = client.ElementText(pageID, "h1")
	if err != nil {
		log.Printf("❌ 获取元素文本失败: %v\n", err)
	} else {
		fmt.Printf("✅ 元素文本: %v\n", resp.Data["text"])
	}

	// 截图
	fmt.Println("\n📌 截图...")
	screenshotData, err := client.Screenshot(pageID, "png")
	if err != nil {
		log.Printf("❌ 截图失败: %v\n", err)
	} else {
		fmt.Printf("✅ 截图成功（数据大小: %d 字节）\n", len(screenshotData))
	}

	// 停止浏览器
	fmt.Println("\n📌 停止浏览器...")
	resp, err = client.StopBrowser()
	if err != nil {
		log.Printf("❌ 停止浏览器失败: %v\n", err)
	} else {
		fmt.Println("✅ 浏览器已停止")
	}

	fmt.Println("\n✅ 测试完成！")
}
