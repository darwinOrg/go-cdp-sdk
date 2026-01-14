package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darwinOrg/go-cdp-sdk"
)

func main() {
	// 创建 WebSocket 客户端
	client := cdp.NewWebSocketClient("ws://localhost:3001", "")

	// 连接到服务器
	fmt.Println("📌 连接到 WebSocket 服务器...")
	if err := client.Connect(context.Background()); err != nil {
		log.Fatalf("❌ 连接失败: %v", err)
	}
	defer client.Close()
	fmt.Println("✅ 连接成功")

	// 注册事件处理器
	client.RegisterEventHandler("load", func(event *cdp.Response) {
		fmt.Printf("📄 页面加载事件: %s\n", event.PageID)
	})

	client.RegisterEventHandler("console", func(event *cdp.Response) {
		fmt.Printf("🖥️  控制台事件: %v\n", event.EventData)
	})

	// 启动浏览器
	fmt.Println("\n📌 启动浏览器...")
	resp, err := client.StartBrowser(false) // 不使用无头模式
	if err != nil {
		log.Fatalf("❌ 启动浏览器失败: %v", err)
	}
	fmt.Printf("✅ 浏览器已启动\n")

	// 创建新页面
	fmt.Println("\n📌 创建新页面...")
	pageID := "page-1"
	resp, err = client.NewPage(pageID)
	if err != nil {
		log.Fatalf("❌ 创建页面失败: %v", err)
	}
	fmt.Printf("✅ 页面已创建: %s\n", pageID)

	// 导航到 URL
	fmt.Println("\n📌 导航到 BOSS直聘...")
	url := "https://www.zhipin.com/gongsi/job/5d627415a46b4a750nJ9.html?ka=company-jobs"
	resp, err = client.Navigate(pageID, url)
	if err != nil {
		log.Fatalf("❌ 导航失败: %v", err)
	}
	fmt.Printf("✅ 导航成功\n")

	// 等待页面加载
	fmt.Println("\n⏳ 等待页面加载...")
	time.Sleep(5 * time.Second)

	// 获取页面标题
	fmt.Println("\n📌 获取页面标题...")
	resp, err = client.GetTitle(pageID)
	if err != nil {
		log.Printf("❌ 获取标题失败: %v\n", err)
	} else if resp.Success {
		fmt.Printf("✅ 页面标题: %v\n", resp.Data["title"])
	}

	// 获取页面 URL
	fmt.Println("\n📌 获取页面 URL...")
	resp, err = client.GetURL(pageID)
	if err != nil {
		log.Printf("❌ 获取 URL 失败: %v\n", err)
	} else if resp.Success {
		fmt.Printf("✅ 页面 URL: %v\n", resp.Data["url"])
	}

	// 执行 JavaScript
	fmt.Println("\n📌 执行 JavaScript...")
	resp, err = client.ExecuteScript(pageID, "document.title")
	if err != nil {
		log.Printf("❌ 执行脚本失败: %v\n", err)
	} else if resp.Success {
		fmt.Printf("✅ 执行结果: %v\n", resp.Data["result"])
	}

	// 检查元素是否存在
	fmt.Println("\n📌 检查元素是否存在...")
	resp, err = client.ElementExists(pageID, "h1")
	if err != nil {
		log.Printf("❌ 检查元素失败: %v\n", err)
	} else if resp.Success {
		fmt.Printf("✅ 元素存在: %v\n", resp.Data["exists"])
	}

	// 截图
	fmt.Println("\n📌 截图...")
	resp, err = client.Screenshot(pageID, "png")
	if err != nil {
		log.Printf("❌ 截图失败: %v\n", err)
	} else if resp.Success {
		fmt.Printf("✅ 截图成功\n")
	}

	// 等待用户中断
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║  按下 Ctrl+C 停止程序                                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 停止浏览器
	fmt.Println("\n📌 停止浏览器...")
	resp, err = client.StopBrowser()
	if err != nil {
		log.Printf("❌ 停止浏览器失败: %v\n", err)
	} else if resp.Success {
		fmt.Println("✅ 浏览器已停止")
	}

	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    程序结束 ✅                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
}
