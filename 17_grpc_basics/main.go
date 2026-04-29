// ========================================
// Lesson 17: gRPC Basics
// ========================================
// gRPC 是 Google 开源的高性能 RPC 框架
// 使用 Protocol Buffers 作为序列化格式
//
// 安装步骤:
//   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
//   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
//
// 生成代码:
//   protoc --go_out=. --go-grpc_out=. proto/todo.proto

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "go-learning/17_grpc_basics/proto"
)

// ========================================
// 服务端实现
// ========================================

type todoServer struct {
	pb.UnimplementedTodoServiceServer // 嵌入未实现的服务器（向前兼容）
	mu                                sync.RWMutex
	todos                             map[int64]*pb.Todo
	nextID                            int64
}

func newTodoServer() *todoServer {
	return &todoServer{
		todos:  make(map[int64]*pb.Todo),
		nextID: 1,
	}
}

// CreateTodo 实现 gRPC 方法
func (s *todoServer) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.Todo, error) {
	if req.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	todo := &pb.Todo{
		Id:        s.nextID,
		Title:     req.Title,
		Completed: false,
	}
	s.todos[s.nextID] = todo
	s.nextID++

	log.Printf("[Server] Created todo: %v", todo)
	return todo, nil
}

// GetTodo 获取单个 todo
func (s *todoServer) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.todos[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "todo %d not found", req.Id)
	}
	return todo, nil
}

// ListTodos 列出所有 todo
func (s *todoServer) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]*pb.Todo, 0, len(s.todos))
	for _, t := range s.todos {
		todos = append(todos, t)
	}
	return &pb.ListTodosResponse{Todos: todos}, nil
}

// UpdateTodo 更新 todo
func (s *todoServer) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "todo %d not found", req.Id)
	}

	if req.Title != "" {
		todo.Title = req.Title
	}
	todo.Completed = req.Completed
	return todo, nil
}

// DeleteTodo 删除 todo
func (s *todoServer) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[req.Id]; !ok {
		return nil, status.Errorf(codes.NotFound, "todo %d not found", req.Id)
	}
	delete(s.todos, req.Id)
	return &pb.DeleteTodoResponse{Success: true}, nil
}

// ========================================
// 启动服务器
// ========================================
func startServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(grpcServer, newTodoServer())

	log.Printf("[Server] Listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// ========================================
// 客户端调用
// ========================================
func runClient(addr string) {
	// 建立连接
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)
	ctx := context.Background()

	// ---- Create ----
	fmt.Println("\n--- Create Todos ---")
	todo1, err := client.CreateTodo(ctx, &pb.CreateTodoRequest{Title: "Learn gRPC"})
	if err != nil {
		log.Fatalf("CreateTodo error: %v", err)
	}
	fmt.Printf("Created: %v\n", todo1)

	todo2, _ := client.CreateTodo(ctx, &pb.CreateTodoRequest{Title: "Build microservice"})
	fmt.Printf("Created: %v\n", todo2)

	todo3, _ := client.CreateTodo(ctx, &pb.CreateTodoRequest{Title: "Deploy to production"})
	fmt.Printf("Created: %v\n", todo3)

	// ---- 错误处理 ----
	fmt.Println("\n--- Error Handling ---")
	_, err = client.CreateTodo(ctx, &pb.CreateTodoRequest{Title: ""})
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("Error code: %s, message: %s\n", st.Code(), st.Message())
	}

	// ---- List ----
	fmt.Println("\n--- List Todos ---")
	listResp, _ := client.ListTodos(ctx, &pb.ListTodosRequest{})
	for _, t := range listResp.Todos {
		fmt.Printf("  [%d] %s (completed: %t)\n", t.Id, t.Title, t.Completed)
	}

	// ---- Get ----
	fmt.Println("\n--- Get Todo ---")
	got, _ := client.GetTodo(ctx, &pb.GetTodoRequest{Id: 1})
	fmt.Printf("Got: %v\n", got)

	// ---- Update ----
	fmt.Println("\n--- Update Todo ---")
	updated, _ := client.UpdateTodo(ctx, &pb.UpdateTodoRequest{
		Id: 1, Title: "Learn gRPC ✓", Completed: true,
	})
	fmt.Printf("Updated: %v\n", updated)

	// ---- Delete ----
	fmt.Println("\n--- Delete Todo ---")
	delResp, _ := client.DeleteTodo(ctx, &pb.DeleteTodoRequest{Id: 3})
	fmt.Printf("Deleted: success=%t\n", delResp.Success)

	// ---- Final List ----
	fmt.Println("\n--- Final State ---")
	listResp, _ = client.ListTodos(ctx, &pb.ListTodosRequest{})
	for _, t := range listResp.Todos {
		status := "[ ]"
		if t.Completed {
			status = "[✓]"
		}
		fmt.Printf("  %s %s\n", status, t.Title)
	}
}

func main() {
	addr := "localhost:50051"

	fmt.Println("============================================")
	fmt.Println("  gRPC Todo Service")
	fmt.Println("============================================")

	// 启动服务器（后台）
	go startServer(addr)
	time.Sleep(100 * time.Millisecond) // 等待服务器启动

	// 运行客户端
	runClient(addr)
}

// ========================================
// 练习:
// 1. 添加 deadline/timeout 到客户端调用
// 2. 实现 gRPC 拦截器（interceptor）记录请求日志
// 3. 添加元数据（metadata）传递认证信息
// ========================================
