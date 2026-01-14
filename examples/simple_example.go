package main

import (
	"context"
	"fmt"
	"log"

	"github.com/darwinOrg/go-cdp-sdk"
)

func main() {
	// 创建 WebSocket 客户端
	client := cdpsdk.NewWebSocketClient("ws://localhost:3001", "")

	// 连接到服务器
	fmt.Println("📌 连接到 WebSocket 服务器...")
	if err := client.Connect(context.Background()); err != nil {
		log.Fatalf("❌ 连接失败: %v", err)
	}
	defer client.Close()
	fmt.Println("✅ 连接成功")

	// 启动浏览器
	fmt.Println("\n📌 启动浏览器...")
	resp, err := client.StartBrowser(false)
	if err != nil {
		log.Fatalf("❌ 启动浏览器失败: %v", err)
	}
	fmt.Printf("✅ 浏览器已启动: port=%v\n", resp.Data["port"])

	// 创建新页面
	fmt.Println("\n📌 创建新页面...")
	pageID := "page-1"
	resp, err = client.NewPage(pageID)
	if err != nil {
		log.Fatalf("❌ 创建页面失败: %v", err)
	}
	fmt.Printf("✅ 页面已创建: %s\n", pageID)

	// 导航到简单的页面
	fmt.Println("\n📌 导航到 example.com...")
	resp, err = client.Navigate(pageID, "https://example.com")
	if err != nil {
		log.Printf("❌ 导航失败: %v\n", err)
	} else if resp.Success {
		fmt.Printf("✅ 导航成功\n")
	}

	// 获取页面标题
	fmt.Println("\n📌 获取页面标题...")
	resp, err = client.GetTitle(pageID)
	if err != nil {
		log.Printf("❌ 获取标题失败: %v\n", err)
	} else if resp.Success {
		fmt.Printf("✅ 页面标题: %v\n", resp.Data["title"])
	}

	// 停止浏览器
	fmt.Println("\n📌 停止浏览器...")
	resp, err = client.StopBrowser()
	if err != nil {
		log.Printf("❌ 停止浏览器失败: %v\n", err)
	} else if resp.Success {
		fmt.Println("✅ 浏览器已停止")
	}

	fmt.Println("\n✅ 测试完成！")
}
