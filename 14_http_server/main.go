// ========================================
// Lesson 14: HTTP Server
// ========================================
// 用 Go 标准库 net/http 构建 Web 服务
// Go 的 HTTP 标准库非常强大，很多生产项目直接使用

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// ---- 数据模型 ----
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// ---- 内存存储 ----
type TodoStore struct {
	mu     sync.RWMutex
	todos  map[int]Todo
	nextID int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos:  make(map[int]Todo),
		nextID: 1,
	}
}

func (s *TodoStore) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo := Todo{
		ID:        s.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	s.todos[s.nextID] = todo
	s.nextID++
	return todo
}

func (s *TodoStore) GetAll() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]Todo, 0, len(s.todos))
	for _, t := range s.todos {
		todos = append(todos, t)
	}
	return todos
}

func (s *TodoStore) GetByID(id int) (Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.todos[id]
	return todo, ok
}

func (s *TodoStore) Update(id int, title string, completed bool) (Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return Todo{}, false
	}
	todo.Title = title
	todo.Completed = completed
	s.todos[id] = todo
	return todo, true
}

func (s *TodoStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.todos[id]
	if ok {
		delete(s.todos, id)
	}
	return ok
}

// ---- HTTP 处理器 ----
type TodoHandler struct {
	store *TodoStore
}

// JSON 响应帮助函数
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// GET /api/todos
func (h *TodoHandler) handleList(w http.ResponseWriter, r *http.Request) {
	todos := h.store.GetAll()
	writeJSON(w, http.StatusOK, todos)
}

// POST /api/todos
func (h *TodoHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if input.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	todo := h.store.Create(input.Title)
	writeJSON(w, http.StatusCreated, todo)
}

// GET /api/todos/{id}
func (h *TodoHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	todo, ok := h.store.GetByID(id)
	if !ok {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

// PUT /api/todos/{id}
func (h *TodoHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	todo, ok := h.store.Update(id, input.Title, input.Completed)
	if !ok {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

// DELETE /api/todos/{id}
func (h *TodoHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if !h.store.Delete(id) {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---- 中间件 ----
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("-> %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("<- %s %s (%v)", r.Method, r.URL.Path, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	store := NewTodoStore()

	// 预置一些数据
	store.Create("Learn Go basics")
	store.Create("Build HTTP server")
	store.Create("Learn gRPC")

	handler := &TodoHandler{store: store}

	// ---- 路由（Go 1.22+ 支持方法和路径参数）----
	mux := http.NewServeMux()

	// RESTful API 路由
	mux.HandleFunc("GET /api/todos", handler.handleList)
	mux.HandleFunc("POST /api/todos", handler.handleCreate)
	mux.HandleFunc("GET /api/todos/{id}", handler.handleGet)
	mux.HandleFunc("PUT /api/todos/{id}", handler.handleUpdate)
	mux.HandleFunc("DELETE /api/todos/{id}", handler.handleDelete)

	// 健康检查
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// 应用中间件
	var finalHandler http.Handler = mux
	finalHandler = loggingMiddleware(finalHandler)
	finalHandler = corsMiddleware(finalHandler)

	// 启动服务器
	addr := ":8080"
	fmt.Println("============================================")
	fmt.Println("  Todo API Server")
	fmt.Println("============================================")
	fmt.Printf("  Listening on http://localhost%s\n", addr)
	fmt.Println()
	fmt.Println("  Endpoints:")
	fmt.Println("    GET    /api/todos      - List all todos")
	fmt.Println("    POST   /api/todos      - Create a todo")
	fmt.Println("    GET    /api/todos/{id}  - Get a todo")
	fmt.Println("    PUT    /api/todos/{id}  - Update a todo")
	fmt.Println("    DELETE /api/todos/{id}  - Delete a todo")
	fmt.Println("    GET    /health          - Health check")
	fmt.Println()
	fmt.Println("  Test with curl:")
	fmt.Println(`    curl localhost:8080/api/todos`)
	fmt.Println(`    curl -X POST localhost:8080/api/todos -d '{"title":"New task"}'`)
	fmt.Println("============================================")

	log.Fatal(http.ListenAndServe(addr, finalHandler))
}

// ========================================
// 练习:
// 1. 添加分页功能：GET /api/todos?page=1&per_page=10
// 2. 添加搜索功能：GET /api/todos?q=keyword
// 3. 添加 rate limiting 中间件
// ========================================
