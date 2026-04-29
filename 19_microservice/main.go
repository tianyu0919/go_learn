// ========================================
// Lesson 19: Microservice Project
// ========================================
// 综合项目：用户服务 + 订单服务（通过 gRPC 通信）
// 这个项目把前面学的知识整合在一起：
// - HTTP API（对外）
// - gRPC（服务间通信）
// - JSON 序列化
// - 并发处理
// - 错误处理

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "go-learning/19_microservice/proto"
)

// ========================================
// User Service（gRPC 服务）
// ========================================

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userService struct {
	pb.UnimplementedUserServiceServer
	mu     sync.RWMutex
	users  map[int64]*User
	nextID int64
}

func newUserService() *userService {
	svc := &userService{
		users:  make(map[int64]*User),
		nextID: 1,
	}
	// 预置数据
	svc.users[1] = &User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	svc.users[2] = &User{ID: 2, Name: "Bob", Email: "bob@example.com"}
	svc.users[3] = &User{ID: 3, Name: "Carol", Email: "carol@example.com"}
	svc.nextID = 4
	return svc
}

func (s *userService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[req.UserId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user %d not found", req.UserId)
	}
	return &pb.UserResponse{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}

func (s *userService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var users []*pb.UserResponse
	for _, u := range s.users {
		users = append(users, &pb.UserResponse{
			UserId: u.ID,
			Name:   u.Name,
			Email:  u.Email,
		})
	}
	return &pb.ListUsersResponse{Users: users}, nil
}

// ========================================
// Order Service（gRPC 服务）
// ========================================

type Order struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	Product   string  `json:"product"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}

type orderService struct {
	pb.UnimplementedOrderServiceServer
	mu     sync.RWMutex
	orders map[int64]*Order
	nextID int64
}

func newOrderService() *orderService {
	svc := &orderService{
		orders: make(map[int64]*Order),
		nextID: 1,
	}
	// 预置数据
	now := time.Now().Format(time.RFC3339)
	svc.orders[1] = &Order{ID: 1, UserID: 1, Product: "Go Book", Amount: 49.99, Status: "completed", CreatedAt: now}
	svc.orders[2] = &Order{ID: 2, UserID: 1, Product: "Keyboard", Amount: 129.99, Status: "shipped", CreatedAt: now}
	svc.orders[3] = &Order{ID: 3, UserID: 2, Product: "Monitor", Amount: 399.99, Status: "pending", CreatedAt: now}
	svc.nextID = 4
	return svc
}

func (s *orderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order := &Order{
		ID:        s.nextID,
		UserID:    req.UserId,
		Product:   req.Product,
		Amount:    req.Amount,
		Status:    "pending",
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	s.orders[s.nextID] = order
	s.nextID++

	return orderToProto(order), nil
}

func (s *orderService) GetUserOrders(ctx context.Context, req *pb.GetUserOrdersRequest) (*pb.OrderListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var orders []*pb.OrderResponse
	for _, o := range s.orders {
		if o.UserID == req.UserId {
			orders = append(orders, orderToProto(o))
		}
	}
	return &pb.OrderListResponse{Orders: orders}, nil
}

func orderToProto(o *Order) *pb.OrderResponse {
	return &pb.OrderResponse{
		OrderId:   o.ID,
		UserId:    o.UserID,
		Product:   o.Product,
		Amount:    o.Amount,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
	}
}

// ========================================
// API Gateway（HTTP -> gRPC）
// ========================================

type APIGateway struct {
	userClient  pb.UserServiceClient
	orderClient pb.OrderServiceClient
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GET /api/users - 获取所有用户
func (g *APIGateway) handleListUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := g.userClient.ListUsers(r.Context(), &pb.ListUsersRequest{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.Users)
}

// GET /api/users/{id} - 获取用户详情（包含订单）
func (g *APIGateway) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// 并发请求用户信息和订单（用 goroutine 加速！）
	type userResult struct {
		user *pb.UserResponse
		err  error
	}
	type ordersResult struct {
		orders *pb.OrderListResponse
		err    error
	}

	userCh := make(chan userResult, 1)
	ordersCh := make(chan ordersResult, 1)

	go func() {
		user, err := g.userClient.GetUser(r.Context(), &pb.GetUserRequest{UserId: id})
		userCh <- userResult{user, err}
	}()

	go func() {
		orders, err := g.orderClient.GetUserOrders(r.Context(), &pb.GetUserOrdersRequest{UserId: id})
		ordersCh <- ordersResult{orders, err}
	}()

	ur := <-userCh
	or := <-ordersCh

	if ur.err != nil {
		st, _ := status.FromError(ur.err)
		if st.Code() == codes.NotFound {
			writeError(w, http.StatusNotFound, "user not found")
		} else {
			writeError(w, http.StatusInternalServerError, ur.err.Error())
		}
		return
	}

	// 组合响应
	response := map[string]any{
		"user": map[string]any{
			"id":    ur.user.UserId,
			"name":  ur.user.Name,
			"email": ur.user.Email,
		},
	}
	if or.err == nil {
		response["orders"] = or.orders.Orders
	} else {
		response["orders"] = []any{}
	}

	writeJSON(w, http.StatusOK, response)
}

// POST /api/orders - 创建订单
func (g *APIGateway) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID  int64   `json:"user_id"`
		Product string  `json:"product"`
		Amount  float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// 先验证用户是否存在
	_, err := g.userClient.GetUser(r.Context(), &pb.GetUserRequest{UserId: input.UserID})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			writeError(w, http.StatusBadRequest, "user not found")
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// 创建订单
	order, err := g.orderClient.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		UserId:  input.UserID,
		Product: input.Product,
		Amount:  input.Amount,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

func main() {
	userAddr := "localhost:50061"
	orderAddr := "localhost:50062"
	httpAddr := ":8081"

	fmt.Println("============================================")
	fmt.Println("  Microservice Demo")
	fmt.Println("============================================")
	fmt.Println("  Architecture:")
	fmt.Println("    Client -> HTTP Gateway -> gRPC Services")
	fmt.Println()
	fmt.Println("  Services:")
	fmt.Printf("    User Service  (gRPC): %s\n", userAddr)
	fmt.Printf("    Order Service (gRPC): %s\n", orderAddr)
	fmt.Printf("    API Gateway   (HTTP): %s\n", httpAddr)
	fmt.Println("============================================")

	// 启动 User Service
	go func() {
		lis, err := net.Listen("tcp", userAddr)
		if err != nil {
			log.Fatal(err)
		}
		s := grpc.NewServer()
		pb.RegisterUserServiceServer(s, newUserService())
		log.Printf("[UserService] Listening on %s", userAddr)
		s.Serve(lis)
	}()

	// 启动 Order Service
	go func() {
		lis, err := net.Listen("tcp", orderAddr)
		if err != nil {
			log.Fatal(err)
		}
		s := grpc.NewServer()
		pb.RegisterOrderServiceServer(s, newOrderService())
		log.Printf("[OrderService] Listening on %s", orderAddr)
		s.Serve(lis)
	}()

	time.Sleep(200 * time.Millisecond)

	// 创建 gRPC 客户端连接
	userConn, err := grpc.NewClient(userAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	orderConn, err := grpc.NewClient(orderAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 启动 API Gateway
	gateway := &APIGateway{
		userClient:  pb.NewUserServiceClient(userConn),
		orderClient: pb.NewOrderServiceClient(orderConn),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/users", gateway.handleListUsers)
	mux.HandleFunc("GET /api/users/{id}", gateway.handleGetUser)
	mux.HandleFunc("POST /api/orders", gateway.handleCreateOrder)

	fmt.Println()
	fmt.Println("  Test commands:")
	fmt.Println(`    curl localhost:8081/api/users`)
	fmt.Println(`    curl localhost:8081/api/users/1`)
	fmt.Println(`    curl -X POST localhost:8081/api/orders -d '{"user_id":1,"product":"Mouse","amount":29.99}'`)
	fmt.Println("============================================")

	log.Fatal(http.ListenAndServe(httpAddr, mux))
}

// ========================================
// 练习:
// 1. 添加认证中间件（检查 API Key）
// 2. 添加请求限流
// 3. 实现订单状态变更通知（gRPC Streaming）
// ========================================
