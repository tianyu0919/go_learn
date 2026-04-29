package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

const fileName = "names.json"
const tmpFileName = "names.json.tmp"

// 为了应对可能的并发请求，使用互斥锁保护文件读写
var fileMutex sync.Mutex

// RequestPayload 用于解析前端传来的 JSON 数据
type RequestPayload struct {
	Names string `json:"names"`
}

// ResponsePayload 用于给前端返回处理结果
type ResponsePayload struct {
	Added int    `json:"added"`
	Total int    `json:"total"`
	Error string `json:"error,omitempty"`
}

func main() {
	// 静态文件服务：直接响应 index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// API 接口：处理提交的人名
	http.HandleFunc("/api/names", handleNames)

	fmt.Println("服务已启动，请在浏览器访问 http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("启动服务失败:", err)
	}
}

func handleNames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, ResponsePayload{Error: "只支持 POST 请求"}, http.StatusMethodNotAllowed)
		return
	}

	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		sendJSON(w, ResponsePayload{Error: "无效的 JSON 数据"}, http.StatusBadRequest)
		return
	}

	// 保证读写文件和处理的原子性
	fileMutex.Lock()
	defer fileMutex.Unlock()

	var names []string
	seen := make(map[string]struct{})

	// 读取历史数据并建立去重 map
	data, err := os.ReadFile(fileName)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &names)
		if err == nil {
			for _, name := range names {
				seen[name] = struct{}{}
			}
		}
	}

	// 使用换行符分割字符串
	lines := strings.Split(payload.Names, "\n")
	addedCount := 0

	for _, line := range lines {
		// 去除首尾的空白字符（比如 \r 或者多余空格），但保留中间的空格
		line = strings.TrimSpace(line)
		if line != "" {
			// 去重逻辑
			if _, exists := seen[line]; !exists {
				seen[line] = struct{}{}
				names = append(names, line)
				addedCount++
			}
		}
	}

	// 如果有新数据，则写入 JSON 文件
	if addedCount > 0 {
		fileData, err := json.MarshalIndent(names, "", "  ")
		if err != nil {
			sendJSON(w, ResponsePayload{Error: "生成 JSON 失败"}, http.StatusInternalServerError)
			return
		}

		// 原子写入策略
		err = os.WriteFile(tmpFileName, fileData, 0644)
		if err != nil {
			sendJSON(w, ResponsePayload{Error: "写入临时文件失败"}, http.StatusInternalServerError)
			return
		}
		err = os.Rename(tmpFileName, fileName)
		if err != nil {
			sendJSON(w, ResponsePayload{Error: "保存文件失败"}, http.StatusInternalServerError)
			return
		}
	}

	// 成功响应，返回当前新增数和总数
	sendJSON(w, ResponsePayload{
		Added: addedCount,
		Total: len(names),
	}, http.StatusOK)
}

func sendJSON(w http.ResponseWriter, resp ResponsePayload, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
