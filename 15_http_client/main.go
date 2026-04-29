// ========================================
// Lesson 15: HTTP Client
// ========================================
// 学习如何用 Go 调用外部 HTTP API
// 包含超时控制、重试、并发请求等实用模式

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// ---- 基础 GET 请求 ----
func basicGet() {
	fmt.Println("--- Basic GET ---")

	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close() // 重要！必须关闭 Body

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Body (first 200 chars): %.200s...\n", string(body))
}

// ---- POST JSON ----
func postJSON() {
	fmt.Println("\n--- POST JSON ---")

	payload := map[string]any{
		"title": "Learn Go HTTP Client",
		"done":  false,
	}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(
		"https://httpbin.org/post",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %s\n", resp.Status)

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	if data, ok := result["json"]; ok {
		fmt.Printf("Server received: %v\n", data)
	}
}

// ---- 自定义请求（带 Header）----
func customRequest() {
	fmt.Println("\n--- Custom Request ---")

	req, err := http.NewRequest("GET", "https://httpbin.org/headers", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 添加自定义 Header
	req.Header.Set("User-Agent", "GoLearner/1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Custom-Header", "hello-from-go")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	prettyJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(prettyJSON))
}

// ---- 带超时和取消的请求 ----
func requestWithTimeout() {
	fmt.Println("\n--- Request with Timeout ---")

	// 使用 context 控制超时
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/delay/1", nil)

	client := &http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error after %v: %v\n", time.Since(start), err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %s (took %v)\n", resp.Status, time.Since(start))
}

// ---- API 客户端封装 ----
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewAPIClient(baseURL string, timeout time.Duration) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *APIClient) Get(path string, result any) error {
	resp, err := c.HTTPClient.Get(c.BaseURL + path)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *APIClient) Post(path string, payload any, result any) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+path,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// ---- 并发请求 ----
func concurrentRequests() {
	fmt.Println("\n--- Concurrent Requests ---")

	urls := []string{
		"https://httpbin.org/get?id=1",
		"https://httpbin.org/get?id=2",
		"https://httpbin.org/get?id=3",
	}

	type Result struct {
		URL    string
		Status string
		Err    error
	}

	results := make(chan Result, len(urls))
	var wg sync.WaitGroup

	start := time.Now()
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Get(url)
			if err != nil {
				results <- Result{URL: url, Err: err}
				return
			}
			defer resp.Body.Close()
			results <- Result{URL: url, Status: resp.Status}
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if r.Err != nil {
			fmt.Printf("  %s -> Error: %v\n", r.URL, r.Err)
		} else {
			fmt.Printf("  %s -> %s\n", r.URL, r.Status)
		}
	}
	fmt.Printf("  Total time: %v (concurrent!)\n", time.Since(start))
}

// ---- 重试机制 ----
func requestWithRetry(url string, maxRetries int) (*http.Response, error) {
	var lastErr error
	client := &http.Client{Timeout: 5 * time.Second}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(attempt) * time.Second
			fmt.Printf("  Retry %d after %v...\n", attempt, delay)
			time.Sleep(delay)
		}

		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode >= 500 {
			resp.Body.Close()
			lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
			continue
		}
		return resp, nil
	}

	return nil, fmt.Errorf("all %d retries failed: %w", maxRetries, lastErr)
}

func main() {
	// 注意：这些示例需要网络连接
	fmt.Println("============================================")
	fmt.Println("  HTTP Client Examples")
	fmt.Println("============================================")
	fmt.Println("  (requires internet connection)")
	fmt.Println()

	basicGet()
	postJSON()
	customRequest()
	requestWithTimeout()

	// API 客户端
	fmt.Println("\n--- API Client ---")
	client := NewAPIClient("https://httpbin.org", 10*time.Second)
	var getResult map[string]any
	if err := client.Get("/get", &getResult); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Origin: %v\n", getResult["origin"])
	}

	concurrentRequests()
}

// ========================================
// 练习:
// 1. 封装一个完整的 REST API 客户端（支持 GET/POST/PUT/DELETE）
// 2. 实现请求/响应拦截器（类似 Axios interceptors）
// 3. 用 HTTP Client 配合 Lesson 14 的服务器做一个完整的增删改查
// ========================================
