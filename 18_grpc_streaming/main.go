// ========================================
// Lesson 18: gRPC Streaming
// ========================================
// gRPC 支持 4 种通信模式:
// 1. Unary（一元）: 请求-响应（Lesson 17 已学）
// 2. Server Streaming: 服务端流
// 3. Client Streaming: 客户端流
// 4. Bidirectional Streaming: 双向流

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "go-learning/18_grpc_streaming/proto"
)

// ========================================
// 服务端实现
// ========================================

type chatServer struct {
	pb.UnimplementedChatServiceServer
}

// ServerStreaming: 服务端返回多条消息（如实时价格推送）
func (s *chatServer) Subscribe(req *pb.SubscribeRequest, stream pb.ChatService_SubscribeServer) error {
	log.Printf("[Server] Client subscribed to topic: %s", req.Topic)

	for i := 0; i < 5; i++ {
		msg := &pb.ChatMessage{
			User:      "System",
			Content:   fmt.Sprintf("[%s] Update #%d - Price: %.2f", req.Topic, i+1, 100+rand.Float64()*50),
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(msg); err != nil {
			return err
		}
		log.Printf("[Server] Sent: %s", msg.Content)
		time.Sleep(500 * time.Millisecond) // 模拟实时推送
	}
	return nil
}

// ClientStreaming: 客户端发送多条消息，服务端返回汇总
func (s *chatServer) SendMessages(stream pb.ChatService_SendMessagesServer) error {
	var count int
	var lastUser string

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// 客户端结束发送，返回汇总
			return stream.SendAndClose(&pb.MessageSummary{
				MessageCount: int32(count),
				LastUser:     lastUser,
			})
		}
		if err != nil {
			return err
		}

		count++
		lastUser = msg.User
		log.Printf("[Server] Received from %s: %s", msg.User, msg.Content)
	}
}

// BidirectionalStreaming: 双向流（聊天室）
func (s *chatServer) Chat(stream pb.ChatService_ChatServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("[Server] %s: %s", msg.User, msg.Content)

		// 服务端回复
		reply := &pb.ChatMessage{
			User:      "Bot",
			Content:   fmt.Sprintf("Echo: %s", msg.Content),
			Timestamp: time.Now().Unix(),
		}
		if err := stream.Send(reply); err != nil {
			return err
		}
	}
}

// ========================================
// 客户端调用
// ========================================

func demoServerStreaming(client pb.ChatServiceClient) {
	fmt.Println("--- Server Streaming ---")
	fmt.Println("(Client subscribes, server pushes updates)")

	stream, err := client.Subscribe(context.Background(), &pb.SubscribeRequest{
		Topic: "BTC/USD",
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  <- %s: %s\n", msg.User, msg.Content)
	}
}

func demoClientStreaming(client pb.ChatServiceClient) {
	fmt.Println("\n--- Client Streaming ---")
	fmt.Println("(Client sends multiple messages, server replies once)")

	stream, err := client.SendMessages(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	messages := []struct{ user, content string }{
		{"Alice", "Hello!"},
		{"Alice", "How's the server?"},
		{"Bob", "I'm here too!"},
		{"Alice", "Great, let's chat"},
	}

	for _, m := range messages {
		fmt.Printf("  -> %s: %s\n", m.user, m.content)
		stream.Send(&pb.ChatMessage{
			User:      m.user,
			Content:   m.content,
			Timestamp: time.Now().Unix(),
		})
		time.Sleep(200 * time.Millisecond)
	}

	summary, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Summary: %d messages, last from: %s\n",
		summary.MessageCount, summary.LastUser)
}

func demoBidirectionalStreaming(client pb.ChatServiceClient) {
	fmt.Println("\n--- Bidirectional Streaming ---")
	fmt.Println("(Client and server exchange messages freely)")

	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 在 goroutine 中接收消息
	done := make(chan bool)
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Printf("Recv error: %v", err)
				done <- true
				return
			}
			fmt.Printf("  <- %s: %s\n", msg.User, msg.Content)
		}
	}()

	// 发送消息
	messages := []string{"Hi there!", "What's the weather?", "Thanks, bye!"}
	for _, content := range messages {
		fmt.Printf("  -> Alice: %s\n", content)
		stream.Send(&pb.ChatMessage{
			User:      "Alice",
			Content:   content,
			Timestamp: time.Now().Unix(),
		})
		time.Sleep(300 * time.Millisecond)
	}

	stream.CloseSend()
	<-done
}

func main() {
	addr := "localhost:50052"

	fmt.Println("============================================")
	fmt.Println("  gRPC Streaming Examples")
	fmt.Println("============================================")

	// 启动服务器
	go func() {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterChatServiceServer(grpcServer, &chatServer{})
		log.Printf("[Server] Listening on %s", addr)
		grpcServer.Serve(lis)
	}()
	time.Sleep(100 * time.Millisecond)

	// 创建客户端
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	demoServerStreaming(client)
	demoClientStreaming(client)
	demoBidirectionalStreaming(client)

	fmt.Println("\n============================================")
	fmt.Println("  All streaming demos complete!")
	fmt.Println("============================================")
}

// ========================================
// 练习:
// 1. 实现一个文件上传服务（Client Streaming）
// 2. 实现一个日志监控服务（Server Streaming）
// 3. 实现一个多人聊天室（Bidirectional + 广播）
// ========================================
