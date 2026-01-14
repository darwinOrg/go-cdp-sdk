package cdp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketClient WebSocket 客户端
type WebSocketClient struct {
	conn          *websocket.Conn
	url           string
	sessionID     string
	mu            sync.Mutex
	requestID     int
	pendingReqs   map[int]chan *Response
	eventHandlers map[string][]EventHandler
	done          chan struct{}
}

// Request WebSocket 请求
type Request struct {
	Type      string                 `json:"type"`
	SessionID string                 `json:"sessionId,omitempty"`
	PageID    string                 `json:"pageId,omitempty"`
	RequestID string                 `json:"requestId,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// Response WebSocket 响应
type Response struct {
	Type      string                 `json:"type"`
	SessionID string                 `json:"sessionId,omitempty"`
	PageID    string                 `json:"pageId,omitempty"`
	RequestID string                 `json:"requestId,omitempty"`
	Success   bool                   `json:"success"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Event     string                 `json:"event,omitempty"`
	EventData map[string]interface{} `json:"eventData,omitempty"`
	Timestamp string                 `json:"timestamp,omitempty"`
}

// EventHandler 事件处理器函数类型
type EventHandler func(event *Response)

// NewWebSocketClient 创建新的 WebSocket 客户端
func NewWebSocketClient(url, sessionID string) *WebSocketClient {
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().UnixNano())
	}

	return &WebSocketClient{
		url:           url,
		sessionID:     sessionID,
		pendingReqs:   make(map[int]chan *Response),
		eventHandlers: make(map[string][]EventHandler),
		done:          make(chan struct{}),
	}
}

// Connect 连接到 WebSocket 服务器
func (wsc *WebSocketClient) Connect(ctx context.Context) error {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, wsc.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket server: %w", err)
	}

	wsc.conn = conn

	// 启动消息接收协程
	go wsc.receiveMessages()

	return nil
}

// Close 关闭连接
func (wsc *WebSocketClient) Close() error {
	close(wsc.done)
	
	wsc.mu.Lock()
	defer wsc.mu.Unlock()

	if wsc.conn != nil {
		return wsc.conn.Close()
	}

	return nil
}

// receiveMessages 接收消息的协程
func (wsc *WebSocketClient) receiveMessages() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("WebSocket receive panic: %v\n", r)
		}
	}()

	for {
		select {
		case <-wsc.done:
			return
		default:
			_, message, err := wsc.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Printf("WebSocket read error: %v\n", err)
				}
				return
			}

			var resp Response
			if err := json.Unmarshal(message, &resp); err != nil {
				log.Printf("Failed to unmarshal message: %v, message: %s\n", err, string(message))
				continue
			}

			wsc.handleResponse(&resp)
		}
	}
}

// handleResponse 处理响应
func (wsc *WebSocketClient) handleResponse(resp *Response) {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()

	// 如果是事件消息（有 event 字段）
	if resp.Event != "" {
		handlers := wsc.eventHandlers[resp.Event]
		for _, handler := range handlers {
			go handler(resp)
		}
		return
	}

	// 如果是响应消息（有 requestId）
	if resp.RequestID != "" {
		// 查找对应的等待通道
		for id, ch := range wsc.pendingReqs {
			if fmt.Sprintf("%d", id) == resp.RequestID {
				select {
				case ch <- resp:
					delete(wsc.pendingReqs, id)
				default:
				}
				return
			}
		}
	}
}

// sendRequest 发送请求并等待响应
func (wsc *WebSocketClient) sendRequest(ctx context.Context, req *Request) (*Response, error) {
	wsc.mu.Lock()

	if wsc.conn == nil {
		wsc.mu.Unlock()
		return nil, fmt.Errorf("not connected to WebSocket server")
	}

	// 生成请求 ID
	wsc.requestID++
	req.RequestID = fmt.Sprintf("%d", wsc.requestID)
	req.SessionID = wsc.sessionID

	// 创建响应通道
	respCh := make(chan *Response, 1)
	wsc.pendingReqs[wsc.requestID] = respCh

	// 发送请求
	message, err := json.Marshal(req)
	if err != nil {
		delete(wsc.pendingReqs, wsc.requestID)
		wsc.mu.Unlock()
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if err := wsc.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		delete(wsc.pendingReqs, wsc.requestID)
		wsc.mu.Unlock()
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// 释放锁，等待响应
	wsc.mu.Unlock()

	// 等待响应（不持有锁）
	select {
	case resp := <-respCh:
		return resp, nil
	case <-ctx.Done():
		wsc.mu.Lock()
		delete(wsc.pendingReqs, wsc.requestID)
		wsc.mu.Unlock()
		return nil, fmt.Errorf("request canceled")
	case <-time.After(30 * time.Second):
		wsc.mu.Lock()
		delete(wsc.pendingReqs, wsc.requestID)
		wsc.mu.Unlock()
		return nil, fmt.Errorf("request timeout after 30 seconds")
	}
}

// RegisterEventHandler 注册事件处理器
func (wsc *WebSocketClient) RegisterEventHandler(eventType string, handler EventHandler) {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()

	wsc.eventHandlers[eventType] = append(wsc.eventHandlers[eventType], handler)
}

// StartBrowser 启动浏览器
func (wsc *WebSocketClient) StartBrowser(headless bool) (*Response, error) {
	req := &Request{
		Type: "start_browser",
		Data: map[string]interface{}{
			"headless": headless,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// StopBrowser 停止浏览器
func (wsc *WebSocketClient) StopBrowser() (*Response, error) {
	req := &Request{
		Type: "stop_browser",
	}

	return wsc.sendRequest(context.Background(), req)
}

// NewPage 创建新页面
func (wsc *WebSocketClient) NewPage(pageID string) (*Response, error) {
	req := &Request{
		Type:   "new_page",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// ClosePage 关闭页面
func (wsc *WebSocketClient) ClosePage(pageID string) (*Response, error) {
	req := &Request{
		Type:   "close_page",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// Navigate 导航到 URL
func (wsc *WebSocketClient) Navigate(pageID, url string) (*Response, error) {
	req := &Request{
		Type:   "navigate",
		PageID: pageID,
		Data: map[string]interface{}{
			"url": url,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// Reload 刷新页面
func (wsc *WebSocketClient) Reload(pageID string) (*Response, error) {
	req := &Request{
		Type:   "reload",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// ExecuteScript 执行 JavaScript
func (wsc *WebSocketClient) ExecuteScript(pageID, script string) (*Response, error) {
	req := &Request{
		Type:   "execute_script",
		PageID: pageID,
		Data: map[string]interface{}{
			"script": script,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// GetTitle 获取页面标题
func (wsc *WebSocketClient) GetTitle(pageID string) (*Response, error) {
	req := &Request{
		Type:   "get_title",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// GetURL 获取页面 URL
func (wsc *WebSocketClient) GetURL(pageID string) (*Response, error) {
	req := &Request{
		Type:   "get_url",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// Screenshot 截图
func (wsc *WebSocketClient) Screenshot(pageID string, format string) (*Response, error) {
	req := &Request{
		Type:   "screenshot",
		PageID: pageID,
		Data: map[string]interface{}{
			"format": format,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// ElementExists 检查元素是否存在
func (wsc *WebSocketClient) ElementExists(pageID, selector string) (*Response, error) {
	req := &Request{
		Type:   "element_exists",
		PageID: pageID,
		Data: map[string]interface{}{
			"selector": selector,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// ElementText 获取元素文本
func (wsc *WebSocketClient) ElementText(pageID, selector string) (*Response, error) {
	req := &Request{
		Type:   "element_text",
		PageID: pageID,
		Data: map[string]interface{}{
			"selector": selector,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// ElementClick 点击元素
func (wsc *WebSocketClient) ElementClick(pageID, selector string) (*Response, error) {
	req := &Request{
		Type:   "element_click",
		PageID: pageID,
		Data: map[string]interface{}{
			"selector": selector,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// ElementSetValue 设置元素值
func (wsc *WebSocketClient) ElementSetValue(pageID, selector, value string) (*Response, error) {
	req := &Request{
		Type:   "element_set_value",
		PageID: pageID,
		Data: map[string]interface{}{
			"selector": selector,
			"value":    value,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// GetSessionID 获取会话 ID
func (wsc *WebSocketClient) GetSessionID() string {
	return wsc.sessionID
}

// NavigateWithLoadedState 导航并等待加载完成
func (wsc *WebSocketClient) NavigateWithLoadedState(pageID, url string) (*Response, error) {
	req := &Request{
		Type:   "navigate_with_loaded_state",
		PageID: pageID,
		Data: map[string]interface{}{
			"url": url,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// ReloadWithLoadedState 刷新并等待加载完成
func (wsc *WebSocketClient) ReloadWithLoadedState(pageID string) (*Response, error) {
	req := &Request{
		Type:   "reload_with_loaded_state",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// WaitForLoadStateLoad 等待页面加载完成
func (wsc *WebSocketClient) WaitForLoadStateLoad(pageID string) (*Response, error) {
	req := &Request{
		Type:   "wait_for_load_state_load",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// WaitForDomContentLoaded 等待 DOM 加载完成
func (wsc *WebSocketClient) WaitForDomContentLoaded(pageID string) (*Response, error) {
	req := &Request{
		Type:   "wait_for_dom_content_loaded",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// WaitForSelectorStateVisible 等待元素可见
func (wsc *WebSocketClient) WaitForSelectorStateVisible(pageID, selector string) (*Response, error) {
	req := &Request{
		Type:   "wait_for_selector_state_visible",
		PageID: pageID,
		Data: map[string]interface{}{
			"selector": selector,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// ExpectResponseText 等待响应文本
func (wsc *WebSocketClient) ExpectResponseText(pageID, urlOrPredicate string, callback func() error) (string, error) {
	req := &Request{
		Type:   "expect_response_text",
		PageID: pageID,
		Data: map[string]interface{}{
			"urlOrPredicate": urlOrPredicate,
			"callback":       callback,
		},
	}

	resp, err := wsc.sendRequest(context.Background(), req)
	if err != nil {
		return "", err
	}

	if resp.Data != nil {
		if text, ok := resp.Data["text"].(string); ok {
			return text, nil
		}
	}

	return "", fmt.Errorf("response text not found")
}

// MustInnerText 必须获取内部文本
func (wsc *WebSocketClient) MustInnerText(pageID, selector string) (string, error) {
	req := &Request{
		Type:   "must_inner_text",
		PageID: pageID,
		Data: map[string]interface{}{
			"selector": selector,
		},
	}

	resp, err := wsc.sendRequest(context.Background(), req)
	if err != nil {
		return "", err
	}

	if resp.Data != nil {
		if text, ok := resp.Data["text"].(string); ok {
			return text, nil
		}
	}

	return "", fmt.Errorf("inner text not found")
}

// MustTextContent 必须获取文本内容
func (wsc *WebSocketClient) MustTextContent(pageID, selector string) (string, error) {
	req := &Request{
		Type:   "must_text_content",
		PageID: pageID,
		Data: map[string]interface{}{
			"selector": selector,
		},
	}

	resp, err := wsc.sendRequest(context.Background(), req)
	if err != nil {
		return "", err
	}

	if resp.Data != nil {
		if text, ok := resp.Data["text"].(string); ok {
			return text, nil
		}
	}

	return "", fmt.Errorf("text content not found")
}

// Suspend 暂停页面
func (wsc *WebSocketClient) Suspend(pageID string) (*Response, error) {
	req := &Request{
		Type:   "suspend",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// Continue 继续页面
func (wsc *WebSocketClient) Continue(pageID string) (*Response, error) {
	req := &Request{
		Type:   "continue",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// Release 释放页面锁
func (wsc *WebSocketClient) Release(pageID string) (*Response, error) {
	req := &Request{
		Type:   "release",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// CloseAll 关闭所有页面
func (wsc *WebSocketClient) CloseAll(pageID string) (*Response, error) {
	req := &Request{
		Type:   "close_all",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// ExpectExtPage 等待新页面
func (wsc *WebSocketClient) ExpectExtPage(pageID string, callback func() error) (*Response, error) {
	req := &Request{
		Type:   "expect_ext_page",
		PageID: pageID,
		Data: map[string]interface{}{
			"callback": callback,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// ElementWait 等待元素
func (wsc *WebSocketClient) ElementWait(pageID, selector string) (*Response, error) {
	req := &Request{
		Type:   "element_wait",
		PageID: pageID,
		Data: map[string]interface{}{
			"selector": selector,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// ElementAttribute 获取元素属性
func (wsc *WebSocketClient) ElementAttribute(pageID, selector, attribute string) (*Response, error) {
	req := &Request{
		Type:   "element_attribute",
		PageID: pageID,
		Data: map[string]interface{}{
			"selector":  selector,
			"attribute": attribute,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// SubscribeEvents 订阅事件
func (wsc *WebSocketClient) SubscribeEvents(pageID string, events []string) (*Response, error) {
	req := &Request{
		Type:   "subscribe_events",
		PageID: pageID,
		Data: map[string]interface{}{
			"events": events,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// RandomWait 随机等待
func (wsc *WebSocketClient) RandomWait(pageID string, min, max int) (*Response, error) {
	req := &Request{
		Type:   "random_wait",
		PageID: pageID,
		Data: map[string]interface{}{
			"min": min,
			"max": max,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// GetHTML 获取 HTML
func (wsc *WebSocketClient) GetHTML(pageID string) (*Response, error) {
	req := &Request{
		Type:   "get_html",
		PageID: pageID,
	}

	return wsc.sendRequest(context.Background(), req)
}

// ElementAllTexts 获取所有元素文本
func (wsc *WebSocketClient) ElementAllTexts(pageID, selector string) (*Response, error) {
	req := &Request{
		Type:   "element_all_texts",
		PageID: pageID,
		Data: map[string]interface{}{
			"selector": selector,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// ElementAllAttributes 获取所有元素属性
func (wsc *WebSocketClient) ElementAllAttributes(pageID, selector string) (*Response, error) {
	req := &Request{
		Type:   "element_all_attributes",
		PageID: pageID,
		Data: map[string]interface{}{
			"selector": selector,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// ElementCount 获取元素数量
func (wsc *WebSocketClient) ElementCount(pageID, selector string) (*Response, error) {
	req := &Request{
		Type:   "element_count",
		PageID: pageID,
		Data: map[string]interface{}{
			"selector": selector,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}

// ConnectBrowser 连接到现有浏览器
func (wsc *WebSocketClient) ConnectBrowser(port int) (*Response, error) {
	req := &Request{
		Type: "connect_browser",
		Data: map[string]interface{}{
			"port": port,
		},
	}

	return wsc.sendRequest(context.Background(), req)
}
