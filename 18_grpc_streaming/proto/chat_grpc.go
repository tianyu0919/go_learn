// Code generated manually to match chat.proto. DO NOT EDIT.

package proto

import (
	context "context"
	"io"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type ChatMessage struct {
	User      string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Content   string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	Timestamp int64  `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

type SubscribeRequest struct {
	Topic string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
}

type MessageSummary struct {
	MessageCount int32  `protobuf:"varint,1,opt,name=message_count,proto3" json:"message_count,omitempty"`
	LastUser     string `protobuf:"bytes,2,opt,name=last_user,proto3" json:"last_user,omitempty"`
}

// ChatServiceClient
type ChatServiceClient interface {
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (ChatService_SubscribeClient, error)
	SendMessages(ctx context.Context, opts ...grpc.CallOption) (ChatService_SendMessagesClient, error)
	Chat(ctx context.Context, opts ...grpc.CallOption) (ChatService_ChatClient, error)
}

type chatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatServiceClient(cc grpc.ClientConnInterface) ChatServiceClient {
	return &chatServiceClient{cc}
}

const chatServiceName = "/chat.ChatService/"

// Subscribe - Server streaming
type ChatService_SubscribeClient interface {
	Recv() (*ChatMessage, error)
	grpc.ClientStream
}

type chatServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *chatServiceSubscribeClient) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chatServiceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (ChatService_SubscribeClient, error) {
	desc := &grpc.StreamDesc{ServerStreams: true}
	stream, err := c.cc.NewStream(ctx, desc, chatServiceName+"Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// SendMessages - Client streaming
type ChatService_SendMessagesClient interface {
	Send(*ChatMessage) error
	CloseAndRecv() (*MessageSummary, error)
	grpc.ClientStream
}

type chatServiceSendMessagesClient struct {
	grpc.ClientStream
}

func (x *chatServiceSendMessagesClient) Send(m *ChatMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatServiceSendMessagesClient) CloseAndRecv() (*MessageSummary, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(MessageSummary)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chatServiceClient) SendMessages(ctx context.Context, opts ...grpc.CallOption) (ChatService_SendMessagesClient, error) {
	desc := &grpc.StreamDesc{ClientStreams: true}
	stream, err := c.cc.NewStream(ctx, desc, chatServiceName+"SendMessages", opts...)
	if err != nil {
		return nil, err
	}
	return &chatServiceSendMessagesClient{stream}, nil
}

// Chat - Bidirectional streaming
type ChatService_ChatClient interface {
	Send(*ChatMessage) error
	Recv() (*ChatMessage, error)
	grpc.ClientStream
}

type chatServiceChatClient struct {
	grpc.ClientStream
}

func (x *chatServiceChatClient) Send(m *ChatMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatServiceChatClient) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chatServiceClient) Chat(ctx context.Context, opts ...grpc.CallOption) (ChatService_ChatClient, error) {
	desc := &grpc.StreamDesc{ServerStreams: true, ClientStreams: true}
	stream, err := c.cc.NewStream(ctx, desc, chatServiceName+"Chat", opts...)
	if err != nil {
		return nil, err
	}
	return &chatServiceChatClient{stream}, nil
}

// ChatServiceServer
type ChatServiceServer interface {
	Subscribe(*SubscribeRequest, ChatService_SubscribeServer) error
	SendMessages(ChatService_SendMessagesServer) error
	Chat(ChatService_ChatServer) error
}

type UnimplementedChatServiceServer struct{}

func (UnimplementedChatServiceServer) Subscribe(*SubscribeRequest, ChatService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedChatServiceServer) SendMessages(ChatService_SendMessagesServer) error {
	return status.Errorf(codes.Unimplemented, "method SendMessages not implemented")
}
func (UnimplementedChatServiceServer) Chat(ChatService_ChatServer) error {
	return status.Errorf(codes.Unimplemented, "method Chat not implemented")
}

// Subscribe server stream
type ChatService_SubscribeServer interface {
	Send(*ChatMessage) error
	grpc.ServerStream
}

type chatServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *chatServiceSubscribeServer) Send(m *ChatMessage) error {
	return x.ServerStream.SendMsg(m)
}

// SendMessages server stream
type ChatService_SendMessagesServer interface {
	SendAndClose(*MessageSummary) error
	Recv() (*ChatMessage, error)
	grpc.ServerStream
}

type chatServiceSendMessagesServer struct {
	grpc.ServerStream
}

func (x *chatServiceSendMessagesServer) SendAndClose(m *MessageSummary) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatServiceSendMessagesServer) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Chat server stream
type ChatService_ChatServer interface {
	Send(*ChatMessage) error
	Recv() (*ChatMessage, error)
	grpc.ServerStream
}

type chatServiceChatServer struct {
	grpc.ServerStream
}

func (x *chatServiceChatServer) Send(m *ChatMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatServiceChatServer) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func RegisterChatServiceServer(s *grpc.Server, srv ChatServiceServer) {
	desc := &grpc.ServiceDesc{
		ServiceName: "chat.ChatService",
		HandlerType: (*ChatServiceServer)(nil),
		Methods:     []grpc.MethodDesc{},
		Streams: []grpc.StreamDesc{
			{
				StreamName:    "Subscribe",
				Handler:       _ChatService_Subscribe_Handler,
				ServerStreams:  true,
			},
			{
				StreamName:    "SendMessages",
				Handler:       _ChatService_SendMessages_Handler,
				ClientStreams:  true,
			},
			{
				StreamName:    "Chat",
				Handler:       _ChatService_Chat_Handler,
				ServerStreams:  true,
				ClientStreams:  true,
			},
		},
		Metadata: "proto/chat.proto",
	}
	s.RegisterService(desc, srv)
}

func _ChatService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChatServiceServer).Subscribe(m, &chatServiceSubscribeServer{stream})
}

func _ChatService_SendMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServiceServer).SendMessages(&chatServiceSendMessagesServer{stream})
}

func _ChatService_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServiceServer).Chat(&chatServiceChatServer{stream})
}

// Ensure io.EOF is accessible
var _ = io.EOF
