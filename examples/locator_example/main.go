package main

import (
	"fmt"
	"log"

	"github.com/darwinOrg/go-cdp-sdk"
)

func main() {
	// 创建 HTTP 客户端
	sessionId := "locator-example-session"
	client := cdpsdk.NewHTTPClient("http://localhost:3000", sessionId)

	// 启动浏览器
	fmt.Println("🚀 测试 Locator 功能...")
	if err := client.StartBrowser(false); err != nil {
		log.Fatalf("❌ 启动浏览器失败: %v", err)
	}
	defer client.StopBrowser()

	page := cdpsdk.NewPage(client)

	// 导航到测试页面
	fmt.Println("\n📌 导航到示例页面...")
	if err := page.Navigate("https://example.com"); err != nil {
		log.Fatalf("❌ 导航失败: %v", err)
	}

	// 测试 Locator 链式调用
	fmt.Println("\n📌 测试 Locator 链式调用...")

	// 1. 单级 Locator
	fmt.Println("1️⃣ 单级 Locator:")
	h1Locator := page.Locator("h1")
	fmt.Printf("   选择器: %v\n", h1Locator.GetSelectors())
	h1Text, err := h1Locator.Text()
	if err != nil {
		log.Printf("❌ 获取文本失败: %v\n", err)
	} else {
		fmt.Printf("   文本: %s", h1Text)
	}

	// 2. 二级 Locator
	fmt.Println("\n2️⃣ 二级 Locator:")
	bodyLocator := page.Locator("body")
	pLocator := bodyLocator.Locator("p")
	fmt.Printf("   选择器链: %v\n", pLocator.GetSelectors())
	fmt.Printf("   最终选择器: %s\n", pLocator.GetSelector())
	pText, err := pLocator.Text()
	if err != nil {
		log.Printf("❌ 获取文本失败: %v\n", err)
	} else {
		fmt.Printf("   文本: %s\n", pText)
	}

	// 3. 三级 Locator
	fmt.Println("\n3️⃣ 三级 Locator:")
	divLocator := page.Locator("div")
	pLocator2 := divLocator.Locator("p")
	aLocator := pLocator2.Locator("a")
	fmt.Printf("   选择器链: %v\n", aLocator.GetSelectors())
	fmt.Printf("   最终选择器: %s\n", aLocator.GetSelector())
	exists, err := aLocator.Exists()
	if err != nil {
		log.Printf("❌ 检查存在失败: %v\n", err)
	} else {
		fmt.Printf("   存在: %v\n", exists)
	}

	// 4. 使用链式调用点击元素
	fmt.Println("\n4️⃣ 链式调用 + 点击:")
	linkLocator := page.Locator("div").Locator("p").Locator("a")
	exists, err = linkLocator.Exists()
	if err != nil {
		log.Printf("❌ 检查存在失败: %v\n", err)
	} else if exists {
		fmt.Printf("   找到链接，准备点击...\n")
		// 注意：在 example.com 上点击可能会离开页面，这里只演示
		// if err := linkLocator.Click(); err != nil {
		//     log.Printf("❌ 点击失败: %v\n", err)
		// } else {
		//     fmt.Println("   ✅ 点击成功")
		// }
	}

	// 5. 获取所有匹配元素
	fmt.Println("\n5️⃣ 获取所有 div 元素:")
	divCount, err := divLocator.Count()
	if err != nil {
		log.Printf("❌ 获取数量失败: %v\n", err)
	} else {
		fmt.Printf("   div 数量: %d\n", divCount)
	}

	fmt.Println("\n✅ Locator 测试完成！")
}
