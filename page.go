package cdpsdk

import (
	"fmt"
)

// Page 页面结构体，封装页面相关操作
type Page struct {
	client *HTTPClient
	pageId string
}

// NewPage 创建页面实例
func NewPage(client *HTTPClient, pageId string) *Page {
	return &Page{
		client: client,
		pageId: pageId,
	}
}

// GetPageID 获取页面ID
func (p *Page) GetPageID() string {
	return p.pageId
}

// GetClient 获取 HTTP 客户端
func (p *Page) GetClient() *HTTPClient {
	return p.client
}

// ========== 导航操作 ==========

// Navigate 导航到 URL
func (p *Page) Navigate(url string) error {
	return p.client.Navigate(p, url)
}

// NavigateWithLoadedState 导航并等待加载完成
func (p *Page) NavigateWithLoadedState(url string) error {
	return p.client.NavigateWithLoadedState(p, url)
}

// Reload 刷新页面
func (p *Page) Reload() error {
	return p.client.Reload(p)
}

// ReloadWithLoadedState 刷新并等待加载完成
func (p *Page) ReloadWithLoadedState() error {
	return p.client.ReloadWithLoadedState(p)
}

// ========== 页面信息 ==========

// GetTitle 获取页面标题
func (p *Page) GetTitle() (string, error) {
	return p.client.GetTitle(p)
}

// GetURL 获取页面 URL
func (p *Page) GetURL() (string, error) {
	return p.client.GetURL(p)
}

// GetHTML 获取页面 HTML
func (p *Page) GetHTML() (string, error) {
	return p.client.GetHTML(p)
}

// ========== 脚本执行 ==========

// ExecuteScript 执行 JavaScript 并返回结果
func (p *Page) ExecuteScript(script string) (any, error) {
	return p.client.ExecuteScript(p, script)
}

// ========== 等待操作 ==========

// WaitForLoadStateLoad 等待页面加载完成
func (p *Page) WaitForLoadStateLoad() error {
	return p.client.WaitForLoadStateLoad(p)
}

// WaitForDomContentLoaded 等待 DOM 加载完成
func (p *Page) WaitForDomContentLoaded() error {
	return p.client.WaitForDomContentLoaded(p)
}

// WaitForSelectorVisible 等待元素可见
func (p *Page) WaitForSelectorVisible(selector string) error {
	return p.client.WaitForSelectorVisible(p, selector)
}

// Wait 等待元素
func (p *Page) Wait(selector string) error {
	return p.client.ElementWait(p, selector, 10000)
}

// ========== 高级功能 ==========

// ExpectResponseText 等待响应文本
func (p *Page) ExpectResponseText(urlOrPredicate, callback string) (string, error) {
	return p.client.ExpectResponseText(p, urlOrPredicate, callback)
}

// MustInnerText 强制获取内部文本
func (p *Page) MustInnerText(selector string) (string, error) {
	return p.client.MustInnerText(p, selector)
}

// MustTextContent 强制获取文本内容
func (p *Page) MustTextContent(selector string) (string, error) {
	return p.client.MustTextContent(p, selector)
}

// ExpectExtPage 等待新页面
func (p *Page) ExpectExtPage(callback string) (string, error) {
	return p.client.ExpectExtPage(p, callback)
}

// Release 释放页面
func (p *Page) Release() error {
	return p.client.Release(p)
}

// CloseAll 关闭所有页面
func (p *Page) CloseAll() error {
	return p.client.CloseAll(p)
}

// ========== 截图 ==========

// Screenshot 截图
func (p *Page) Screenshot(format string) ([]byte, error) {
	return p.client.Screenshot(p, format)
}

// ========== 元素操作快捷方式 ==========

// Locator 创建定位器
func (p *Page) Locator(selector string) *Locator {
	return p.client.Locator(p, selector)
}

// Exists 检查元素是否存在
func (p *Page) Exists(selector string) (bool, error) {
	return p.client.ElementExists(p, selector)
}

// Text 获取元素文本
func (p *Page) Text(selector string) (string, error) {
	return p.client.ElementText(p, selector)
}

// Click 点击元素
func (p *Page) Click(selector string) error {
	return p.client.ElementClick(p, selector)
}

// SetValue 设置元素值
func (p *Page) SetValue(selector, value string) error {
	return p.client.ElementSetValue(p, selector, value)
}

// Attribute 获取元素属性
func (p *Page) Attribute(selector, attr string) (string, error) {
	return p.client.ElementAttribute(p, selector, attr)
}

// AllTexts 获取所有匹配元素的文本
func (p *Page) AllTexts(selector string) ([]string, error) {
	return p.client.ElementAllTexts(p, selector)
}

// AllAttributes 获取所有匹配元素的属性
func (p *Page) AllAttributes(selector, attr string) ([]string, error) {
	return p.client.ElementAllAttributes(p, selector, attr)
}

// Count 获取元素数量
func (p *Page) Count(selector string) (int, error) {
	return p.client.ElementCount(p, selector)
}

// ========== 链式操作 ==========

// NavigateThen 导航后执行操作
func (p *Page) NavigateThen(url string, callback func(*Page) error) error {
	if err := p.Navigate(url); err != nil {
		return err
	}
	return callback(p)
}

// NavigateAndWait 导航并等待
func (p *Page) NavigateAndWait(url string, waitFunc func(*Page) error) error {
	if err := p.Navigate(url); err != nil {
		return err
	}
	return waitFunc(p)
}

// ClickThen 点击后执行操作
func (p *Page) ClickThen(selector string, callback func(*Page) error) error {
	if err := p.Click(selector); err != nil {
		return err
	}
	return callback(p)
}

// SetValueThen 设置值后执行操作
func (p *Page) SetValueThen(selector, value string, callback func(*Page) error) error {
	if err := p.SetValue(selector, value); err != nil {
		return err
	}
	return callback(p)
}

// ========== 实用方法 ==========

// PrintTitle 打印页面标题
func (p *Page) PrintTitle() error {
	title, err := p.GetTitle()
	if err != nil {
		return err
	}
	fmt.Printf("页面标题: %s\n", title)
	return nil
}

// PrintURL 打印页面 URL
func (p *Page) PrintURL() error {
	url, err := p.GetURL()
	if err != nil {
		return err
	}
	fmt.Printf("页面 URL: %s\n", url)
	return nil
}
