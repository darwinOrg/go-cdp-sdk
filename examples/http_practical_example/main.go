package main

import (
	"fmt"
	"log"
	"time"

	"github.com/darwinOrg/go-cdp-sdk"
)

func main() {
	// 创建 HTTP 客户端
	client := cdpsdk.NewHTTPClient("http://localhost:3000", "my-automation-session")

	// 连接到现有浏览器（9222 端口）
	fmt.Println("🚀 开始自动化流程...")
	fmt.Println("📌 步骤 1: 连接到浏览器（端口 9222）...")
	if err := client.ConnectBrowser(9222); err != nil {
		log.Fatalf("❌ 连接浏览器失败: %v", err)
	}
	fmt.Println("✅ 已连接到浏览器")

	// 使用默认页面
	pageID := "default"

	// 导航到百度
	fmt.Println("\n📌 步骤 2: 导航到百度首页...")
	if err := client.Navigate(pageID, "https://www.baidu.com"); err != nil {
		log.Printf("❌ 导航失败: %v\n", err)
	} else {
		fmt.Println("✅ 导航成功")
	}

	// 等待页面加载
	fmt.Println("\n📌 步骤 3: 等待页面加载完成...")
	time.Sleep(2 * time.Second) // 简单等待

	// 获取页面标题
	fmt.Println("\n📌 步骤 4: 获取页面标题...")
	title, err := client.GetTitle(pageID)
	if err != nil {
		log.Printf("❌ 获取标题失败: %v\n", err)
	} else {
		fmt.Printf("✅ 页面标题: %s\n", title)
	}

	// 检查搜索框是否存在
	fmt.Println("\n📌 步骤 5: 检查搜索框是否存在...")
	exists, err := client.ElementExists(pageID, "#kw")
	if err != nil {
		log.Printf("❌ 检查元素失败: %v\n", err)
	} else if exists {
		fmt.Println("✅ 搜索框存在")
	} else {
		fmt.Println("⚠️  搜索框不存在")
	}

	// 在搜索框中输入文本
	fmt.Println("\n📌 步骤 6: 在搜索框中输入文本...")
	if err := client.ElementSetValue(pageID, "#kw", "TypeScript CDP 自动化"); err != nil {
		log.Printf("❌ 输入文本失败: %v\n", err)
	} else {
		fmt.Println("✅ 输入成功")
	}

	// 随机等待（模拟人类行为）
	fmt.Println("\n📌 步骤 7: 随机等待...")
	if err := client.RandomWait(pageID, "middle"); err != nil {
		log.Printf("❌ 随机等待失败: %v\n", err)
	} else {
		fmt.Println("✅ 等待完成")
	}

	// 点击搜索按钮
	fmt.Println("\n📌 步骤 8: 点击搜索按钮...")
	if err := client.ElementClick(pageID, "#su"); err != nil {
		log.Printf("❌ 点击失败: %v\n", err)
	} else {
		fmt.Println("✅ 点击成功")
	}

	// 等待搜索结果加载
	fmt.Println("\n📌 步骤 9: 等待搜索结果加载...")
	time.Sleep(3 * time.Second)

	// 获取搜索结果数量
	fmt.Println("\n📌 步骤 10: 获取搜索结果数量...")
	count, err := client.ElementCount(pageID, ".result")
	if err != nil {
		log.Printf("❌ 获取结果数量失败: %v\n", err)
	} else {
		fmt.Printf("✅ 搜索结果数量: %d\n", count)
	}

	// 获取所有搜索结果的标题
	fmt.Println("\n📌 步骤 11: 获取搜索结果标题...")
	texts, err := client.ElementAllTexts(pageID, ".result h3 a")
	if err != nil {
		log.Printf("❌ 获取标题失败: %v\n", err)
	} else {
		fmt.Printf("✅ 找到 %d 个结果:\n", len(texts))
		for i, text := range texts {
			if i < 5 { // 只显示前5个
				fmt.Printf("   %d. %s\n", i+1, text)
			}
		}
	}

	// 截图保存当前状态
	fmt.Println("\n📌 步骤 12: 截图...")
	screenshotData, err := client.Screenshot(pageID, "png")
	if err != nil {
		log.Printf("❌ 截图失败: %v\n", err)
	} else {
		fmt.Printf("✅ 截图成功（数据大小: %d 字节）\n", len(screenshotData))
		// 可以将 screenshotData 保存到文件
		// err := os.WriteFile("screenshot.png", screenshotData, 0644)
	}

	// 获取页面 HTML（可选）
	fmt.Println("\n📌 步骤 13: 获取页面 HTML...")
	html, err := client.GetHTML(pageID)
	if err != nil {
		log.Printf("❌ 获取 HTML 失败: %v\n", err)
	} else {
		fmt.Printf("✅ HTML 长度: %d 字符\n", len(html))
	}

	// 停止浏览器
	fmt.Println("\n📌 步骤 14: 停止浏览器...")
	if err := client.StopBrowser(); err != nil {
		log.Printf("❌ 停止浏览器失败: %v\n", err)
	} else {
		fmt.Println("✅ 浏览器已停止")
	}

	fmt.Println("\n🎉 自动化流程完成！")
}
