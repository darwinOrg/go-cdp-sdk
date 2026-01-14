package main

import (
	"fmt"
	"log"
	"time"

	"github.com/darwinOrg/go-cdp-sdk"
)

func main() {
	// 创建 HTTP 客户端
	client := cdpsdk.NewHTTPClient("http://localhost:3000", "test-zhipin-session")

	// 目标 URL
	targetURL := "https://www.zhipin.com/gongsi/job/5d627415a46b4a750nJ9.html?ka=company-jobs"

	fmt.Println("🚀 开始测试 BOSS 直聘 URL...")

	// 1. 连接到浏览器（9222 端口）
	fmt.Println("\n📌 步骤 1: 连接到浏览器（端口 9222）...")
	resp, err := client.ConnectBrowser(9222)
	if err != nil {
		log.Fatalf("❌ 连接浏览器失败: %v", err)
	}
	fmt.Printf("✅ 已连接到浏览器: sessionId=%s, port=%v\n", client.GetSessionID(), resp.Data["port"])

	// 使用默认页面
	pageID := "default"

	// 2. 导航到目标 URL
	fmt.Printf("\n📌 步骤 2: 导航到 %s...\n", targetURL)
	resp, err = client.Navigate(pageID, targetURL)
	if err != nil {
		log.Printf("❌ 导航失败: %v\n", err)
		return
	}
	fmt.Println("✅ 导航成功")

	// 3. 等待页面加载
	fmt.Println("\n📌 步骤 3: 等待页面加载...")
	time.Sleep(5 * time.Second) // 等待 5 秒让页面完全加载
	fmt.Println("✅ 等待完成")

	// 4. 获取页面标题
	fmt.Println("\n📌 步骤 4: 获取页面标题...")
	resp, err = client.GetTitle(pageID)
	if err != nil {
		log.Printf("❌ 获取标题失败: %v\n", err)
	} else if title, ok := resp.Data["title"].(string); ok {
		fmt.Printf("✅ 页面标题: %s\n", title)
	}

	// 5. 获取页面 URL
	fmt.Println("\n📌 步骤 5: 获取页面 URL...")
	resp, err = client.GetURL(pageID)
	if err != nil {
		log.Printf("❌ 获取 URL 失败: %v\n", err)
	} else if url, ok := resp.Data["url"].(string); ok {
		fmt.Printf("✅ 页面 URL: %s\n", url)
	}

	// 6. 检查页面标题元素
	fmt.Println("\n📌 步骤 6: 检查页面标题元素...")
	resp, err = client.ElementExists(pageID, "h1")
	if err != nil {
		log.Printf("❌ 检查元素失败: %v\n", err)
	} else if exists, ok := resp.Data["exists"].(bool); ok {
		fmt.Printf("✅ h1 元素存在: %v\n", exists)
	}

	// 7. 检查职位标题元素（BOSS 直聘的职位标题）
	fmt.Println("\n📌 步骤 7: 检查职位标题元素...")
	jobTitleSelectors := []string{
		".job-primary .job-name",
		".job-name",
		"div.job-name",
		"[class*='job-name']",
	}

	for _, selector := range jobTitleSelectors {
		resp, err = client.ElementExists(pageID, selector)
		if err == nil && resp.Success {
			if exists, ok := resp.Data["exists"].(bool); ok && exists {
				fmt.Printf("✅ 找到职位标题元素: %s\n", selector)
				// 尝试获取文本
				resp, err = client.ElementText(pageID, selector)
				if err == nil && resp.Success {
					if text, ok := resp.Data["text"].(string); ok {
						fmt.Printf("   职位标题: %s\n", text)
					}
				}
				break
			}
		}
	}

	// 8. 检查公司名称元素
	fmt.Println("\n📌 步骤 8: 检查公司名称元素...")
	companySelectors := []string{
		".job-primary .company-name",
		".company-name",
		"div.company-name",
		"[class*='company-name']",
	}

	for _, selector := range companySelectors {
		resp, err = client.ElementExists(pageID, selector)
		if err == nil && resp.Success {
			if exists, ok := resp.Data["exists"].(bool); ok && exists {
				fmt.Printf("✅ 找到公司名称元素: %s\n", selector)
				// 尝试获取文本
				resp, err = client.ElementText(pageID, selector)
				if err == nil && resp.Success {
					if text, ok := resp.Data["text"].(string); ok {
						fmt.Printf("   公司名称: %s\n", text)
					}
				}
				break
			}
		}
	}

	// 9. 检查薪资元素
	fmt.Println("\n📌 步骤 9: 检查薪资元素...")
	salarySelectors := []string{
		".job-primary .salary",
		".salary",
		"span.salary",
		"[class*='salary']",
	}

	for _, selector := range salarySelectors {
		resp, err = client.ElementExists(pageID, selector)
		if err == nil && resp.Success {
			if exists, ok := resp.Data["exists"].(bool); ok && exists {
				fmt.Printf("✅ 找到薪资元素: %s\n", selector)
				// 尝试获取文本
				resp, err = client.ElementText(pageID, selector)
				if err == nil && resp.Success {
					if text, ok := resp.Data["text"].(string); ok {
						fmt.Printf("   薪资: %s\n", text)
					}
				}
				break
			}
		}
	}

	// 10. 获取页面 HTML（前 500 字符）
	fmt.Println("\n📌 步骤 10: 获取页面 HTML（前 500 字符）...")
	resp, err = client.GetHTML(pageID)
	if err != nil {
		log.Printf("❌ 获取 HTML 失败: %v\n", err)
	} else if html, ok := resp.Data["html"].(string); ok {
		preview := html
		if len(preview) > 500 {
			preview = preview[:500]
		}
		fmt.Printf("✅ HTML 预览:\n%s...\n", preview)
	}

	// 11. 截图
	fmt.Println("\n📌 步骤 11: 截图...")
	screenshotData, err := client.Screenshot(pageID, "png")
	if err != nil {
		log.Printf("❌ 截图失败: %v\n", err)
	} else {
		fmt.Printf("✅ 截图成功（数据大小: %d 字节）\n", len(screenshotData))
	}

	// 12. 停止浏览器
	fmt.Println("\n📌 步骤 12: 停止浏览器...")
	resp, err = client.StopBrowser()
	if err != nil {
		log.Printf("❌ 停止浏览器失败: %v\n", err)
	} else if resp.Success {
		fmt.Println("✅ 浏览器已停止")
	}

	fmt.Println("\n🎉 测试完成！")
}
