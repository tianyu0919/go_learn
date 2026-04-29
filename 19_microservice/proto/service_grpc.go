// Code generated manually to match service.proto. DO NOT EDIT.

package proto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// ========== User Service Messages ==========

type GetUserRequest struct {
	UserId int64 `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
}

type ListUsersRequest struct{}

type UserResponse struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email  string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
}

type ListUsersResponse struct {
	Users []*UserResponse `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
}

// ========== Order Service Messages ==========

type CreateOrderRequest struct {
	UserId  int64   `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Product string  `protobuf:"bytes,2,opt,name=product,proto3" json:"product,omitempty"`
	Amount  float64 `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

type GetUserOrdersRequest struct {
	UserId int64 `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
}

type OrderResponse struct {
	OrderId   int64   `protobuf:"varint,1,opt,name=order_id,proto3" json:"order_id,omitempty"`
	UserId    int64   `protobuf:"varint,2,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Product   string  `protobuf:"bytes,3,opt,name=product,proto3" json:"product,omitempty"`
	Amount    float64 `protobuf:"fixed64,4,opt,name=amount,proto3" json:"amount,omitempty"`
	Status    string  `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	CreatedAt string  `protobuf:"bytes,6,opt,name=created_at,proto3" json:"created_at,omitempty"`
}

type OrderListResponse struct {
	Orders []*OrderResponse `protobuf:"bytes,1,rep,name=orders,proto3" json:"orders,omitempty"`
}

// ========== User Service ==========

type UserServiceClient interface {
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

const userSvcName = "/microservice.UserService/"

func (c *userServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, userSvcName+"GetUser", in, out, opts...)
	return out, err
}

func (c *userServiceClient) ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error) {
	out := new(ListUsersResponse)
	err := c.cc.Invoke(ctx, userSvcName+"ListUsers", in, out, opts...)
	return out, err
}

type UserServiceServer interface {
	GetUser(context.Context, *GetUserRequest) (*UserResponse, error)
	ListUsers(context.Context, *ListUsersRequest) (*ListUsersResponse, error)
}

type UnimplementedUserServiceServer struct{}

func (UnimplementedUserServiceServer) GetUser(context.Context, *GetUserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUserServiceServer) ListUsers(context.Context, *ListUsersRequest) (*ListUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}

func RegisterUserServiceServer(s *grpc.Server, srv UserServiceServer) {
	desc := &grpc.ServiceDesc{
		ServiceName: "microservice.UserService",
		HandlerType: (*UserServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{MethodName: "GetUser", Handler: _UserService_GetUser_Handler},
			{MethodName: "ListUsers", Handler: _UserService_ListUsers_Handler},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "proto/service.proto",
	}
	s.RegisterService(desc, srv)
}

func _UserService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: userSvcName + "GetUser"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUser(ctx, req.(*GetUserRequest))
	})
}

func _UserService_ListUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: userSvcName + "ListUsers"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ListUsers(ctx, req.(*ListUsersRequest))
	})
}

// ========== Order Service ==========

type OrderServiceClient interface {
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
	GetUserOrders(ctx context.Context, in *GetUserOrdersRequest, opts ...grpc.CallOption) (*OrderListResponse, error)
}

type orderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderServiceClient(cc grpc.ClientConnInterface) OrderServiceClient {
	return &orderServiceClient{cc}
}

const orderSvcName = "/microservice.OrderService/"

func (c *orderServiceClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, orderSvcName+"CreateOrder", in, out, opts...)
	return out, err
}

func (c *orderServiceClient) GetUserOrders(ctx context.Context, in *GetUserOrdersRequest, opts ...grpc.CallOption) (*OrderListResponse, error) {
	out := new(OrderListResponse)
	err := c.cc.Invoke(ctx, orderSvcName+"GetUserOrders", in, out, opts...)
	return out, err
}

type OrderServiceServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*OrderResponse, error)
	GetUserOrders(context.Context, *GetUserOrdersRequest) (*OrderListResponse, error)
}

type UnimplementedOrderServiceServer struct{}

func (UnimplementedOrderServiceServer) CreateOrder(context.Context, *CreateOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedOrderServiceServer) GetUserOrders(context.Context, *GetUserOrdersRequest) (*OrderListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserOrders not implemented")
}

func RegisterOrderServiceServer(s *grpc.Server, srv OrderServiceServer) {
	desc := &grpc.ServiceDesc{
		ServiceName: "microservice.OrderService",
		HandlerType: (*OrderServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{MethodName: "CreateOrder", Handler: _OrderService_CreateOrder_Handler},
			{MethodName: "GetUserOrders", Handler: _OrderService_GetUserOrders_Handler},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "proto/service.proto",
	}
	s.RegisterService(desc, srv)
}

func _OrderService_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: orderSvcName + "CreateOrder"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).CreateOrder(ctx, req.(*CreateOrderRequest))
	})
}

func _OrderService_GetUserOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserOrdersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).GetUserOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: orderSvcName + "GetUserOrders"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).GetUserOrders(ctx, req.(*GetUserOrdersRequest))
	})
}
