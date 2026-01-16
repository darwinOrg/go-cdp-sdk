package cdpsdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient HTTP 客户端
type HTTPClient struct {
	baseURL    string
	sessionID  string
	httpClient *http.Client
	pages      []string // 页面ID列表
}

// HTTPResponse HTTP 响应
type HTTPResponse struct {
	Success bool           `json:"success"`
	Data    map[string]any `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}

// NewHTTPClient 创建新的 HTTP 客户端
// sessionID 可以为空，会在调用 StartBrowser 或 ConnectBrowser 时自动生成
func NewHTTPClient(baseURL, sessionID string) *HTTPClient {
	return &HTTPClient{
		baseURL:   baseURL,
		sessionID: sessionID,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute, // 增加超时时间到 5 分钟
		},
		pages: []string{}, // 初始化页面列表
	}
}

// doRequest 执行 HTTP 请求
func (hc *HTTPClient) doRequest(method, endpoint string, body any) (*HTTPResponse, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := hc.baseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := hc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var httpResp HTTPResponse
	if err := json.Unmarshal(respBody, &httpResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !httpResp.Success {
		return nil, fmt.Errorf("server error: %s", httpResp.Error)
	}

	return &httpResp, nil
}

// doRequestBinary 执行 HTTP 请求并返回原始数据
func (hc *HTTPClient) doRequestBinary(method, endpoint string, body any) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := hc.baseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := hc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, nil
}

// StartBrowser 启动浏览器
func (hc *HTTPClient) StartBrowser(headless bool) error {
	body := map[string]any{}
	if headless {
		body["headless"] = "new"
	}

	resp, err := hc.doRequest("POST", "/api/browser/start", body)
	if err != nil {
		return err
	}

	// 从响应中获取 sessionId
	if sessionId, ok := resp.Data["sessionId"].(string); ok {
		hc.sessionID = sessionId
	} else {
		return fmt.Errorf("sessionId not found in response")
	}

	// 从响应中获取页面列表
	if pages, ok := resp.Data["pages"].([]any); ok {
		hc.pages = make([]string, 0, len(pages))
		for _, p := range pages {
			if pageID, ok := p.(string); ok {
				hc.pages = append(hc.pages, pageID)
			}
		}
	}

	return nil
}

// ConnectBrowser 连接到现有浏览器
func (hc *HTTPClient) ConnectBrowser(port int) error {
	body := map[string]any{
		"port": port,
	}

	resp, err := hc.doRequest("POST", "/api/browser/connect", body)
	if err != nil {
		return err
	}

	// 从响应中获取 sessionId
	if sessionId, ok := resp.Data["sessionId"].(string); ok {
		hc.sessionID = sessionId
	} else {
		return fmt.Errorf("sessionId not found in response")
	}

	// 从响应中获取页面列表
	if pages, ok := resp.Data["pages"].([]any); ok {
		hc.pages = make([]string, 0, len(pages))
		for _, p := range pages {
			if pageID, ok := p.(string); ok {
				hc.pages = append(hc.pages, pageID)
			}
		}
	}

	return nil
}

// StopBrowser 停止浏览器
func (hc *HTTPClient) StopBrowser() error {
	body := map[string]any{
		"sessionId": hc.sessionID,
	}

	_, err := hc.doRequest("POST", "/api/browser/stop", body)
	return err
}

// ClosePage 关闭页面
func (hc *HTTPClient) ClosePage(page *Page) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/close", body)
	return err
}

// Navigate 导航到 URL
func (hc *HTTPClient) Navigate(page *Page, url string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"url":       url,
	}

	_, err := hc.doRequest("POST", "/api/page/navigate", body)
	return err
}

// NavigateWithLoadedState 导航并等待加载完成
func (hc *HTTPClient) NavigateWithLoadedState(page *Page, url string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"url":       url,
	}

	_, err := hc.doRequest("POST", "/api/page/navigate-with-loaded-state", body)
	return err
}

// Reload 刷新页面
func (hc *HTTPClient) Reload(page *Page) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/reload", body)
	return err
}

// ReloadWithLoadedState 刷新并等待加载完成
func (hc *HTTPClient) ReloadWithLoadedState(page *Page) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/reload-with-loaded-state", body)
	return err
}

// ExecuteScript 执行 JavaScript
func (hc *HTTPClient) ExecuteScript(page *Page, script string) (any, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"script":    script,
	}

	resp, err := hc.doRequest("POST", "/api/page/execute", body)
	if err != nil {
		return nil, err
	}

	return resp.Data["result"], nil
}

// GetTitle 获取页面标题
func (hc *HTTPClient) GetTitle(page *Page) (string, error) {
	endpoint := fmt.Sprintf("/api/page/title?sessionId=%s", hc.sessionID)
	if page.pageID != "" {
		endpoint += fmt.Sprintf("&pageId=%s", page.pageID)
	}

	resp, err := hc.doRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	if title, ok := resp.Data["title"].(string); ok {
		return title, nil
	}

	return "", fmt.Errorf("title not found in response")
}

// GetURL 获取页面 URL
func (hc *HTTPClient) GetURL(page *Page) (string, error) {
	endpoint := fmt.Sprintf("/api/page/url?sessionId=%s", hc.sessionID)
	if page.pageID != "" {
		endpoint += fmt.Sprintf("&pageId=%s", page.pageID)
	}

	resp, err := hc.doRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	if url, ok := resp.Data["url"].(string); ok {
		return url, nil
	}

	return "", fmt.Errorf("url not found in response")
}

// GetHTML 获取页面 HTML
func (hc *HTTPClient) GetHTML(page *Page) (string, error) {
	endpoint := fmt.Sprintf("/api/page/html?sessionId=%s", hc.sessionID)
	if page.pageID != "" {
		endpoint += fmt.Sprintf("&pageId=%s", page.pageID)
	}

	resp, err := hc.doRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	if html, ok := resp.Data["html"].(string); ok {
		return html, nil
	}

	return "", fmt.Errorf("html not found in response")
}

// Screenshot 截图
func (hc *HTTPClient) Screenshot(page *Page, format string) ([]byte, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"format":    format,
	}

	return hc.doRequestBinary("POST", "/api/page/screenshot", body)
}

// WaitForLoadStateLoad 等待页面加载完成
func (hc *HTTPClient) WaitForLoadStateLoad(page *Page) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/wait-for-load-state-load", body)
	return err
}

// WaitForDomContentLoaded 等待 DOM 加载完成
func (hc *HTTPClient) WaitForDomContentLoaded(page *Page) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/wait-for-dom-content-loaded", body)
	return err
}

// WaitForSelectorVisible 等待选择器可见
func (hc *HTTPClient) WaitForSelectorVisible(page *Page, selector string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
	}

	_, err := hc.doRequest("POST", "/api/page/wait-for-selector-visible", body)
	return err
}

// ExpectResponseText 等待响应文本
func (hc *HTTPClient) ExpectResponseText(page *Page, urlOrPredicate, callback string) (string, error) {
	body := map[string]any{
		"sessionId":      hc.sessionID,
		"pageId":         page.pageID,
		"urlOrPredicate": urlOrPredicate,
		"callback":       callback,
	}

	resp, err := hc.doRequest("POST", "/api/page/expect-response-text", body)
	if err != nil {
		return "", err
	}

	if text, ok := resp.Data["text"].(string); ok {
		return text, nil
	}

	return "", fmt.Errorf("text not found in response")
}

// MustInnerText 必须获取内部文本
func (hc *HTTPClient) MustInnerText(page *Page, selector string) (string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
	}

	resp, err := hc.doRequest("POST", "/api/page/must-inner-text", body)
	if err != nil {
		return "", err
	}

	if text, ok := resp.Data["text"].(string); ok {
		return text, nil
	}

	return "", fmt.Errorf("text not found in response")
}

// MustTextContent 必须获取文本内容
func (hc *HTTPClient) MustTextContent(page *Page, selector string) (string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
	}

	resp, err := hc.doRequest("POST", "/api/page/must-text-content", body)
	if err != nil {
		return "", err
	}

	if text, ok := resp.Data["text"].(string); ok {
		return text, nil
	}

	return "", fmt.Errorf("text not found in response")
}

// Release 释放页面锁
func (hc *HTTPClient) Release(page *Page) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/release", body)
	return err
}

// CloseAll 关闭所有页面
func (hc *HTTPClient) CloseAll(page *Page) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/close-all", body)
	return err
}

// ExpectExtPage 等待新页面
func (hc *HTTPClient) ExpectExtPage(page *Page, callback string) (string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"callback":  callback,
	}

	resp, err := hc.doRequest("POST", "/api/page/expect-ext-page", body)
	if err != nil {
		return "", err
	}

	if pageID, ok := resp.Data["pageId"].(string); ok {
		return pageID, nil
	}

	return "", fmt.Errorf("pageId not found in response")
}

// ElementExists 检查元素是否存在
func (hc *HTTPClient) ElementExists(page *Page, selector string) (bool, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
	}

	resp, err := hc.doRequest("POST", "/api/element/exists", body)
	if err != nil {
		return false, err
	}

	if exists, ok := resp.Data["exists"].(bool); ok {
		return exists, nil
	}

	return false, fmt.Errorf("exists not found in response")
}

// ElementText 获取元素文本
func (hc *HTTPClient) ElementText(page *Page, selector string) (string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
	}

	resp, err := hc.doRequest("POST", "/api/element/text", body)
	if err != nil {
		return "", err
	}

	if text, ok := resp.Data["text"].(string); ok {
		return text, nil
	}

	return "", fmt.Errorf("text not found in response")
}

// ElementClick 点击元素
func (hc *HTTPClient) ElementClick(page *Page, selector string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
	}

	_, err := hc.doRequest("POST", "/api/element/click", body)
	return err
}

// ElementHover 鼠标悬停
func (hc *HTTPClient) ElementHover(page *Page, selector string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
	}

	_, err := hc.doRequest("POST", "/api/element/hover", body)
	return err
}

// ElementSetValue 设置元素值
func (hc *HTTPClient) ElementSetValue(page *Page, selector, value string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
		"value":     value,
	}

	_, err := hc.doRequest("POST", "/api/element/setValue", body)
	return err
}

// ElementWait 等待元素
func (hc *HTTPClient) ElementWait(page *Page, selector string, timeout int) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
		"timeout":   timeout,
	}

	_, err := hc.doRequest("POST", "/api/element/wait", body)
	return err
}

// ElementAttribute 获取元素属性
func (hc *HTTPClient) ElementAttribute(page *Page, selector, attribute string) (string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
		"attribute": attribute,
	}

	resp, err := hc.doRequest("POST", "/api/element/attribute", body)
	if err != nil {
		return "", err
	}

	if value, ok := resp.Data["value"].(string); ok {
		return value, nil
	}

	return "", fmt.Errorf("value not found in response")
}

// ElementAllTexts 获取所有匹配元素的文本
func (hc *HTTPClient) ElementAllTexts(page *Page, selector string) ([]string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
	}

	resp, err := hc.doRequest("POST", "/api/element/all-texts", body)
	if err != nil {
		return nil, err
	}

	if texts, ok := resp.Data["texts"].([]any); ok {
		result := make([]string, len(texts))
		for i, t := range texts {
			if s, ok := t.(string); ok {
				result[i] = s
			}
		}
		return result, nil
	}

	return nil, fmt.Errorf("texts not found in response")
}

// ElementAllAttributes 获取所有匹配元素的属性
func (hc *HTTPClient) ElementAllAttributes(page *Page, selector, attribute string) ([]string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
		"attribute": attribute,
	}

	resp, err := hc.doRequest("POST", "/api/element/all-attributes", body)
	if err != nil {
		return nil, err
	}

	if attributes, ok := resp.Data["attributes"].([]any); ok {
		result := make([]string, len(attributes))
		for i, a := range attributes {
			if s, ok := a.(string); ok {
				result[i] = s
			}
		}
		return result, nil
	}

	return nil, fmt.Errorf("attributes not found in response")
}

// ElementCount 获取元素数量
func (hc *HTTPClient) ElementCount(page *Page, selector string) (int, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    page.pageID,
		"selector":  selector,
	}

	resp, err := hc.doRequest("POST", "/api/element/count", body)
	if err != nil {
		return 0, err
	}

	if count, ok := resp.Data["count"].(float64); ok {
		return int(count), nil
	}

	return 0, fmt.Errorf("count not found in response")
}

// GetSessionID 获取会话 ID
func (hc *HTTPClient) GetSessionID() string {
	return hc.sessionID
}

// NewPage 创建新页面
func (hc *HTTPClient) NewPage() (*Page, error) {
	resp, err := hc.doRequest("POST", "/api/page/new", nil)
	if err != nil {
		return nil, err
	}

	// 从响应中获取 pageId
	if pageID, ok := resp.Data["pageId"].(string); ok {
		hc.pages = append(hc.pages, pageID)
		return NewPage(hc, pageID), nil
	}

	return nil, fmt.Errorf("pageId not found in response")
}

// GetDefaultPage 获取默认页面实例（第一个页面）
func (hc *HTTPClient) GetDefaultPage() (*Page, error) {
	if len(hc.pages) == 0 {
		return nil, fmt.Errorf("no pages available")
	}
	return NewPage(hc, hc.pages[0]), nil
}

// GetPage 根据页面ID获取页面实例
func (hc *HTTPClient) GetPage(pageID string) (*Page, error) {
	for _, pid := range hc.pages {
		if pid == pageID {
			return NewPage(hc, pageID), nil
		}
	}
	return nil, fmt.Errorf("page not found: %s", pageID)
}

// GetPages 获取所有页面ID
func (hc *HTTPClient) GetPages() []string {
	return hc.pages
}

// SetTimeout 设置请求超时时间
func (hc *HTTPClient) SetTimeout(timeout time.Duration) {
	hc.httpClient.Timeout = timeout
}
