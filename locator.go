package cdpsdk

import (
	"fmt"
)

// Locator 元素定位器，支持链式调用
type Locator struct {
	client    *HTTPClient
	pageID    string
	selector  string
	selectors []string // 选择器链
}

// Locator 创建定位器
func (hc *HTTPClient) Locator(pageID, selector string) *Locator {
	return &Locator{
		client:    hc,
		pageID:    pageID,
		selector:  selector,
		selectors: []string{selector},
	}
}

// Locator 嵌套定位器，支持多级定位
func (l *Locator) Locator(selector string) *Locator {
	newSelector := fmt.Sprintf("%s %s", l.selector, selector)
	return &Locator{
		client:    l.client,
		pageID:    l.pageID,
		selector:  newSelector,
		selectors: append(l.selectors, selector),
	}
}

// GetSelectors 获取选择器链
func (l *Locator) GetSelectors() []string {
	return l.selectors
}

// GetSelector 获取最终的选择器
func (l *Locator) GetSelector() string {
	return l.selector
}

// Exists 检查元素是否存在
func (l *Locator) Exists() (bool, error) {
	return l.client.ElementExists(l.pageID, l.selector)
}

// Text 获取元素文本
func (l *Locator) Text() (string, error) {
	return l.client.ElementText(l.pageID, l.selector)
}

// Click 点击元素
func (l *Locator) Click() error {
	return l.client.ElementClick(l.pageID, l.selector)
}

// Hover 鼠标悬停
func (l *Locator) Hover() error {
	return l.client.ElementHover(l.pageID, l.selector)
}

// SetValue 设置元素值
func (l *Locator) SetValue(value string) error {
	return l.client.ElementSetValue(l.pageID, l.selector, value)
}

// Attribute 获取元素属性
func (l *Locator) Attribute(attr string) (string, error) {
	return l.client.ElementAttribute(l.pageID, l.selector, attr)
}

// AllTexts 获取所有匹配元素的文本
func (l *Locator) AllTexts() ([]string, error) {
	return l.client.ElementAllTexts(l.pageID, l.selector)
}

// AllAttributes 获取所有匹配元素的属性
func (l *Locator) AllAttributes(attr string) ([]string, error) {
	return l.client.ElementAllAttributes(l.pageID, l.selector, attr)
}

// Count 获取元素数量
func (l *Locator) Count() (int, error) {
	return l.client.ElementCount(l.pageID, l.selector)
}
