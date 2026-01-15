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
	Success bool           `json:"success"`
	Data    map[string]any `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
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
	body := map[string]any{
		"sessionId": hc.sessionID,
	}
	if headless {
		body["headless"] = "new"
	}

	_, err := hc.doRequest("POST", "/api/browser/start", body)
	return err
}

// ConnectBrowser 连接到现有浏览器
func (hc *HTTPClient) ConnectBrowser(port int) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"port":      port,
	}

	_, err := hc.doRequest("POST", "/api/browser/connect", body)
	return err
}

// StopBrowser 停止浏览器
func (hc *HTTPClient) StopBrowser() error {
	body := map[string]any{
		"sessionId": hc.sessionID,
	}

	_, err := hc.doRequest("POST", "/api/browser/stop", body)
	return err
}

// NewPage 创建新页面
func (hc *HTTPClient) NewPage(pageID string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/new", body)
	return err
}

// ClosePage 关闭页面
func (hc *HTTPClient) ClosePage(pageID string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/close", body)
	return err
}

// Navigate 导航到 URL
func (hc *HTTPClient) Navigate(pageID, url string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"url":       url,
	}

	_, err := hc.doRequest("POST", "/api/page/navigate", body)
	return err
}

// NavigateWithLoadedState 导航并等待加载完成
func (hc *HTTPClient) NavigateWithLoadedState(pageID, url string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"url":       url,
	}

	_, err := hc.doRequest("POST", "/api/page/navigate-with-loaded-state", body)
	return err
}

// Reload 刷新页面
func (hc *HTTPClient) Reload(pageID string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/reload", body)
	return err
}

// ReloadWithLoadedState 刷新并等待加载完成
func (hc *HTTPClient) ReloadWithLoadedState(pageID string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/reload-with-loaded-state", body)
	return err
}

// ExecuteScript 执行 JavaScript
func (hc *HTTPClient) ExecuteScript(pageID, script string) (any, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"script":    script,
	}

	resp, err := hc.doRequest("POST", "/api/page/execute", body)
	if err != nil {
		return nil, err
	}

	return resp.Data["result"], nil
}

// GetTitle 获取页面标题
func (hc *HTTPClient) GetTitle(pageID string) (string, error) {
	endpoint := fmt.Sprintf("/api/page/title?sessionId=%s", hc.sessionID)
	if pageID != "" {
		endpoint += fmt.Sprintf("&pageId=%s", pageID)
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
func (hc *HTTPClient) GetURL(pageID string) (string, error) {
	endpoint := fmt.Sprintf("/api/page/url?sessionId=%s", hc.sessionID)
	if pageID != "" {
		endpoint += fmt.Sprintf("&pageId=%s", pageID)
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
func (hc *HTTPClient) GetHTML(pageID string) (string, error) {
	endpoint := fmt.Sprintf("/api/page/html?sessionId=%s", hc.sessionID)
	if pageID != "" {
		endpoint += fmt.Sprintf("&pageId=%s", pageID)
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
func (hc *HTTPClient) Screenshot(pageID, format string) ([]byte, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"format":    format,
	}

	return hc.doRequestBinary("POST", "/api/page/screenshot", body)
}

// RandomWait 随机等待
func (hc *HTTPClient) RandomWait(pageID string, duration any) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"duration":  duration,
	}

	_, err := hc.doRequest("POST", "/api/page/random-wait", body)
	return err
}

// WaitForLoadStateLoad 等待页面加载完成
func (hc *HTTPClient) WaitForLoadStateLoad(pageID string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/wait-for-load-state-load", body)
	return err
}

// WaitForDomContentLoaded 等待 DOM 加载完成
func (hc *HTTPClient) WaitForDomContentLoaded(pageID string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/wait-for-dom-content-loaded", body)
	return err
}

// WaitForSelectorVisible 等待选择器可见
func (hc *HTTPClient) WaitForSelectorVisible(pageID, selector string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
	}

	_, err := hc.doRequest("POST", "/api/page/wait-for-selector-visible", body)
	return err
}

// ExpectResponseText 等待响应文本
func (hc *HTTPClient) ExpectResponseText(pageID, urlOrPredicate, callback string) (string, error) {
	body := map[string]any{
		"sessionId":      hc.sessionID,
		"pageId":         pageID,
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
func (hc *HTTPClient) MustInnerText(pageID, selector string) (string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
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
func (hc *HTTPClient) MustTextContent(pageID, selector string) (string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
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

// Suspend 暂停页面
func (hc *HTTPClient) Suspend(pageID string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/suspend", body)
	return err
}

// Continue 继续页面
func (hc *HTTPClient) Continue(pageID string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/continue", body)
	return err
}

// Release 释放页面锁
func (hc *HTTPClient) Release(pageID string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/release", body)
	return err
}

// CloseAll 关闭所有页面
func (hc *HTTPClient) CloseAll(pageID string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
	}

	_, err := hc.doRequest("POST", "/api/page/close-all", body)
	return err
}

// ExpectExtPage 等待新页面
func (hc *HTTPClient) ExpectExtPage(pageID, callback string) (string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
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
func (hc *HTTPClient) ElementExists(pageID, selector string) (bool, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
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
func (hc *HTTPClient) ElementText(pageID, selector string) (string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
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
func (hc *HTTPClient) ElementClick(pageID, selector string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
	}

	_, err := hc.doRequest("POST", "/api/element/click", body)
	return err
}

// ElementSetValue 设置元素值
func (hc *HTTPClient) ElementSetValue(pageID, selector, value string) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
		"value":     value,
	}

	_, err := hc.doRequest("POST", "/api/element/setValue", body)
	return err
}

// ElementWait 等待元素
func (hc *HTTPClient) ElementWait(pageID, selector string, timeout int) error {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
		"selector":  selector,
		"timeout":   timeout,
	}

	_, err := hc.doRequest("POST", "/api/element/wait", body)
	return err
}

// ElementAttribute 获取元素属性
func (hc *HTTPClient) ElementAttribute(pageID, selector, attribute string) (string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
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
func (hc *HTTPClient) ElementAllTexts(pageID, selector string) ([]string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
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
func (hc *HTTPClient) ElementAllAttributes(pageID, selector, attribute string) ([]string, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
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
func (hc *HTTPClient) ElementCount(pageID, selector string) (int, error) {
	body := map[string]any{
		"sessionId": hc.sessionID,
		"pageId":    pageID,
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

// SetTimeout 设置请求超时时间
func (hc *HTTPClient) SetTimeout(timeout time.Duration) {
	hc.httpClient.Timeout = timeout
}
