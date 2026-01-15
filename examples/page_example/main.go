package main

import (
	"fmt"
	"log"

	"github.com/darwinOrg/go-cdp-sdk"
)

func main() {
	// 创建 HTTP 客户端
	client := cdpsdk.NewHTTPClient("http://localhost:3000", "page-test")

	// 连接到浏览器
	fmt.Println("🚀 测试 Page 结构体功能...")
	if err := client.ConnectBrowser(9222); err != nil {
		log.Fatalf("❌ 连接浏览器失败: %v", err)
	}
	defer client.StopBrowser()

	// 创建页面实例
	page := client.NewPage("default")

	// 测试页面操作
	fmt.Println("\n📌 测试页面操作...")

	// 1. 导航
	fmt.Println("1️⃣ 导航到 example.com...")
	if err := page.Navigate("https://example.com"); err != nil {
		log.Printf("❌ 导航失败: %v\n", err)
	} else {
		fmt.Println("✅ 导航成功")
	}

	// 2. 等待加载
	fmt.Println("\n2️⃣ 等待页面加载...")
	if err := page.WaitForLoadStateLoad(); err != nil {
		log.Printf("❌ 等待加载失败: %v\n", err)
	} else {
		fmt.Println("✅ 页面加载完成")
	}

	// 3. 获取页面信息
	fmt.Println("\n3️⃣ 获取页面信息...")
	if err := page.PrintTitle(); err != nil {
		log.Printf("❌ 打印标题失败: %v\n", err)
	}
	if err := page.PrintURL(); err != nil {
		log.Printf("❌ 打印 URL 失败: %v\n", err)
	}

	// 4. 使用 Locator
	fmt.Println("\n4️⃣ 使用 Locator 操作元素...")
	h1Locator := page.Locator("h1")
	h1Text, err := h1Locator.Text()
	if err != nil {
		log.Printf("❌ 获取文本失败: %v\n", err)
	} else {
		fmt.Printf("✅ h1 文本: %s\n", h1Text)
	}

	// 5. 多级 Locator + 链式操作
	fmt.Println("\n5️⃣ 多级 Locator + 链式操作...")
	linkLocator := page.Locator("div").Locator("p").Locator("a")
	exists, err := linkLocator.Exists()
	if err != nil {
		log.Printf("❌ 检查存在失败: %v\n", err)
	} else if exists {
		fmt.Println("✅ 找到链接元素")
	}

	// 6. 元素操作快捷方式
	fmt.Println("\n6️⃣ 元素操作快捷方式...")
	pText, err := page.Text("p")
	if err != nil {
		log.Printf("❌ 获取 p 文本失败: %v\n", err)
	} else {
		fmt.Printf("✅ p 文本: %s\n", pText)
	}

	// 7. 链式操作
	fmt.Println("\n7️⃣ 链式操作...")
	if err := page.NavigateThen("https://www.baidu.com", func(p *cdpsdk.Page) error {
		fmt.Println("导航到百度完成")
		return p.WaitForLoadStateLoad()
	}); err != nil {
		log.Printf("❌ 链式操作失败: %v\n", err)
	} else {
		fmt.Println("✅ 链式操作成功")
	}

	// 9. 执行脚本
	fmt.Println("\n9️⃣ 执行 JavaScript...")
	result, err := page.ExecuteScript("document.title")
	if err != nil {
		log.Printf("❌ 执行脚本失败: %v\n", err)
	} else {
		fmt.Printf("✅ 脚本结果: %v\n", result)
	}

	// 10. 截图
	fmt.Println("\n🔟 截图...")
	screenshotData, err := page.Screenshot("png")
	if err != nil {
		log.Printf("❌ 截图失败: %v\n", err)
	} else {
		fmt.Printf("✅ 截图成功（数据大小: %d 字节）\n", len(screenshotData))
	}

	fmt.Println("\n✅ Page 结构体测试完成！")
}
