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
}

// HTTPResponse HTTP 响应
type HTTPResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// NewHTTPClient 创建新的 HTTP 客户端
func NewHTTPClient(baseURL, sessionID string) *HTTPClient {
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().UnixNano())
	}

	return &HTTPClient{
		baseURL:   baseURL,
		sessionID: sessionID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest 执行 HTTP 请求
func (hc *HTTPClient) doRequest(method, endpoint string, body interface{}) (*HTTPResponse, error) {
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
func (hc *HTTPClient) doRequestBinary(method, endpoint string, body interface{}) ([]byte, error) {
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
func (hc *HTTPClient) StartBrowser(headless bool) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"headless":  headless,
	}

	return hc.doRequest("POST", "/api/browser/start", body)
}

// ConnectBrowser 连接到现有浏览器
func (hc *HTTPClient) ConnectBrowser(port int) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"port":      port,
	}

	return hc.doRequest("POST", "/api/browser/connect", body)
}

// StopBrowser 停止浏览器
func (hc *HTTPClient) StopBrowser() (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
	}

	return hc.doRequest("POST", "/api/browser/stop", body)
}

// NewPage 创建新页面
func (hc *HTTPClient) NewPage(pageID string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	return hc.doRequest("POST", "/api/page/new", body)
}

// ClosePage 关闭页面
func (hc *HTTPClient) ClosePage(pageID string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	return hc.doRequest("POST", "/api/page/close", body)
}

// Navigate 导航到 URL
func (hc *HTTPClient) Navigate(pageID, url string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"url":       url,
	}

	return hc.doRequest("POST", "/api/page/navigate", body)
}

// NavigateWithLoadedState 导航并等待加载完成
func (hc *HTTPClient) NavigateWithLoadedState(pageID, url string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"url":       url,
	}

	return hc.doRequest("POST", "/api/page/navigate-with-loaded-state", body)
}

// Reload 刷新页面
func (hc *HTTPClient) Reload(pageID string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	return hc.doRequest("POST", "/api/page/reload", body)
}

// ReloadWithLoadedState 刷新并等待加载完成
func (hc *HTTPClient) ReloadWithLoadedState(pageID string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	return hc.doRequest("POST", "/api/page/reload-with-loaded-state", body)
}

// ExecuteScript 执行 JavaScript
func (hc *HTTPClient) ExecuteScript(pageID, script string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"script":    script,
	}

	return hc.doRequest("POST", "/api/page/execute", body)
}

// GetTitle 获取页面标题
func (hc *HTTPClient) GetTitle(pageID string) (*HTTPResponse, error) {
	endpoint := fmt.Sprintf("/api/page/title?sessionId=%s", hc.sessionID)
	if pageID != "" {
		endpoint += fmt.Sprintf("&pageId=%s", pageID)
	}

	return hc.doRequest("GET", endpoint, nil)
}

// GetURL 获取页面 URL
func (hc *HTTPClient) GetURL(pageID string) (*HTTPResponse, error) {
	endpoint := fmt.Sprintf("/api/page/url?sessionId=%s", hc.sessionID)
	if pageID != "" {
		endpoint += fmt.Sprintf("&pageId=%s", pageID)
	}

	return hc.doRequest("GET", endpoint, nil)
}

// GetHTML 获取页面 HTML
func (hc *HTTPClient) GetHTML(pageID string) (*HTTPResponse, error) {
	endpoint := fmt.Sprintf("/api/page/html?sessionId=%s", hc.sessionID)
	if pageID != "" {
		endpoint += fmt.Sprintf("&pageId=%s", pageID)
	}

	return hc.doRequest("GET", endpoint, nil)
}

// Screenshot 截图
func (hc *HTTPClient) Screenshot(pageID, format string) ([]byte, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"format":    format,
	}

	return hc.doRequestBinary("POST", "/api/page/screenshot", body)
}

// RandomWait 随机等待
func (hc *HTTPClient) RandomWait(pageID string, duration interface{}) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"duration":  duration,
	}

	return hc.doRequest("POST", "/api/page/random-wait", body)
}

// WaitForLoadStateLoad 等待页面加载完成
func (hc *HTTPClient) WaitForLoadStateLoad(pageID string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	return hc.doRequest("POST", "/api/page/wait-for-load-state-load", body)
}

// WaitForDomContentLoaded 等待 DOM 加载完成
func (hc *HTTPClient) WaitForDomContentLoaded(pageID string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	return hc.doRequest("POST", "/api/page/wait-for-dom-content-loaded", body)
}

// WaitForSelectorVisible 等待选择器可见
func (hc *HTTPClient) WaitForSelectorVisible(pageID, selector string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
	}

	return hc.doRequest("POST", "/api/page/wait-for-selector-visible", body)
}

// ExpectResponseText 等待响应文本
func (hc *HTTPClient) ExpectResponseText(pageID, urlOrPredicate, callback string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId":      hc.sessionID,
		"pageId":         pageID,
		"urlOrPredicate": urlOrPredicate,
		"callback":       callback,
	}

	return hc.doRequest("POST", "/api/page/expect-response-text", body)
}

// MustInnerText 必须获取内部文本
func (hc *HTTPClient) MustInnerText(pageID, selector string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
	}

	return hc.doRequest("POST", "/api/page/must-inner-text", body)
}

// MustTextContent 必须获取文本内容
func (hc *HTTPClient) MustTextContent(pageID, selector string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
	}

	return hc.doRequest("POST", "/api/page/must-text-content", body)
}

// Suspend 暂停页面
func (hc *HTTPClient) Suspend(pageID string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	return hc.doRequest("POST", "/api/page/suspend", body)
}

// Continue 继续页面
func (hc *HTTPClient) Continue(pageID string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	return hc.doRequest("POST", "/api/page/continue", body)
}

// Release 释放页面锁
func (hc *HTTPClient) Release(pageID string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	return hc.doRequest("POST", "/api/page/release", body)
}

// CloseAll 关闭所有页面
func (hc *HTTPClient) CloseAll(pageID string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	return hc.doRequest("POST", "/api/page/close-all", body)
}

// ExpectExtPage 等待新页面
func (hc *HTTPClient) ExpectExtPage(pageID, callback string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"callback":  callback,
	}

	return hc.doRequest("POST", "/api/page/expect-ext-page", body)
}

// ElementExists 检查元素是否存在
func (hc *HTTPClient) ElementExists(pageID, selector string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
	}

	return hc.doRequest("POST", "/api/element/exists", body)
}

// ElementText 获取元素文本
func (hc *HTTPClient) ElementText(pageID, selector string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
	}

	return hc.doRequest("POST", "/api/element/text", body)
}

// ElementClick 点击元素
func (hc *HTTPClient) ElementClick(pageID, selector string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
	}

	return hc.doRequest("POST", "/api/element/click", body)
}

// ElementSetValue 设置元素值
func (hc *HTTPClient) ElementSetValue(pageID, selector, value string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
		"value":     value,
	}

	return hc.doRequest("POST", "/api/element/setValue", body)
}

// ElementWait 等待元素
func (hc *HTTPClient) ElementWait(pageID, selector string, timeout int) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
		"timeout":   timeout,
	}

	return hc.doRequest("POST", "/api/element/wait", body)
}

// ElementAttribute 获取元素属性
func (hc *HTTPClient) ElementAttribute(pageID, selector, attribute string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
		"attribute": attribute,
	}

	return hc.doRequest("POST", "/api/element/attribute", body)
}

// ElementAllTexts 获取所有匹配元素的文本
func (hc *HTTPClient) ElementAllTexts(pageID, selector string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
	}

	return hc.doRequest("POST", "/api/element/all-texts", body)
}

// ElementAllAttributes 获取所有匹配元素的属性
func (hc *HTTPClient) ElementAllAttributes(pageID, selector, attribute string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
		"attribute": attribute,
	}

	return hc.doRequest("POST", "/api/element/all-attributes", body)
}

// ElementCount 获取元素数量
func (hc *HTTPClient) ElementCount(pageID, selector string) (*HTTPResponse, error) {
	body := map[string]interface{}{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
	}

	return hc.doRequest("POST", "/api/element/count", body)
}

// GetSessionID 获取会话 ID
func (hc *HTTPClient) GetSessionID() string {
	return hc.sessionID
}

// SetTimeout 设置请求超时时间
func (hc *HTTPClient) SetTimeout(timeout time.Duration) {
	hc.httpClient.Timeout = timeout
}
