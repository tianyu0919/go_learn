// Code generated manually to match todo.proto. DO NOT EDIT.
// In production, run: protoc --go_out=. --go-grpc_out=. proto/todo.proto

package proto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Todo message
type Todo struct {
	Id        int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title     string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Completed bool   `protobuf:"varint,3,opt,name=completed,proto3" json:"completed,omitempty"`
}

type CreateTodoRequest struct {
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
}

type GetTodoRequest struct {
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

type ListTodosRequest struct{}

type ListTodosResponse struct {
	Todos []*Todo `protobuf:"bytes,1,rep,name=todos,proto3" json:"todos,omitempty"`
}

type UpdateTodoRequest struct {
	Id        int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title     string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Completed bool   `protobuf:"varint,3,opt,name=completed,proto3" json:"completed,omitempty"`
}

type DeleteTodoRequest struct {
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

type DeleteTodoResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

// TodoServiceClient is the client API for TodoService.
type TodoServiceClient interface {
	CreateTodo(ctx context.Context, in *CreateTodoRequest, opts ...grpc.CallOption) (*Todo, error)
	GetTodo(ctx context.Context, in *GetTodoRequest, opts ...grpc.CallOption) (*Todo, error)
	ListTodos(ctx context.Context, in *ListTodosRequest, opts ...grpc.CallOption) (*ListTodosResponse, error)
	UpdateTodo(ctx context.Context, in *UpdateTodoRequest, opts ...grpc.CallOption) (*Todo, error)
	DeleteTodo(ctx context.Context, in *DeleteTodoRequest, opts ...grpc.CallOption) (*DeleteTodoResponse, error)
}

type todoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTodoServiceClient(cc grpc.ClientConnInterface) TodoServiceClient {
	return &todoServiceClient{cc}
}

const todoServiceName = "/todo.TodoService/"

func (c *todoServiceClient) CreateTodo(ctx context.Context, in *CreateTodoRequest, opts ...grpc.CallOption) (*Todo, error) {
	out := new(Todo)
	err := c.cc.Invoke(ctx, todoServiceName+"CreateTodo", in, out, opts...)
	return out, err
}

func (c *todoServiceClient) GetTodo(ctx context.Context, in *GetTodoRequest, opts ...grpc.CallOption) (*Todo, error) {
	out := new(Todo)
	err := c.cc.Invoke(ctx, todoServiceName+"GetTodo", in, out, opts...)
	return out, err
}

func (c *todoServiceClient) ListTodos(ctx context.Context, in *ListTodosRequest, opts ...grpc.CallOption) (*ListTodosResponse, error) {
	out := new(ListTodosResponse)
	err := c.cc.Invoke(ctx, todoServiceName+"ListTodos", in, out, opts...)
	return out, err
}

func (c *todoServiceClient) UpdateTodo(ctx context.Context, in *UpdateTodoRequest, opts ...grpc.CallOption) (*Todo, error) {
	out := new(Todo)
	err := c.cc.Invoke(ctx, todoServiceName+"UpdateTodo", in, out, opts...)
	return out, err
}

func (c *todoServiceClient) DeleteTodo(ctx context.Context, in *DeleteTodoRequest, opts ...grpc.CallOption) (*DeleteTodoResponse, error) {
	out := new(DeleteTodoResponse)
	err := c.cc.Invoke(ctx, todoServiceName+"DeleteTodo", in, out, opts...)
	return out, err
}

// TodoServiceServer is the server API for TodoService.
type TodoServiceServer interface {
	CreateTodo(context.Context, *CreateTodoRequest) (*Todo, error)
	GetTodo(context.Context, *GetTodoRequest) (*Todo, error)
	ListTodos(context.Context, *ListTodosRequest) (*ListTodosResponse, error)
	UpdateTodo(context.Context, *UpdateTodoRequest) (*Todo, error)
	DeleteTodo(context.Context, *DeleteTodoRequest) (*DeleteTodoResponse, error)
}

// UnimplementedTodoServiceServer should be embedded for forward compatibility.
type UnimplementedTodoServiceServer struct{}

func (UnimplementedTodoServiceServer) CreateTodo(context.Context, *CreateTodoRequest) (*Todo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTodo not implemented")
}
func (UnimplementedTodoServiceServer) GetTodo(context.Context, *GetTodoRequest) (*Todo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTodo not implemented")
}
func (UnimplementedTodoServiceServer) ListTodos(context.Context, *ListTodosRequest) (*ListTodosResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTodos not implemented")
}
func (UnimplementedTodoServiceServer) UpdateTodo(context.Context, *UpdateTodoRequest) (*Todo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTodo not implemented")
}
func (UnimplementedTodoServiceServer) DeleteTodo(context.Context, *DeleteTodoRequest) (*DeleteTodoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTodo not implemented")
}

func RegisterTodoServiceServer(s *grpc.Server, srv TodoServiceServer) {
	desc := &grpc.ServiceDesc{
		ServiceName: "todo.TodoService",
		HandlerType: (*TodoServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{MethodName: "CreateTodo", Handler: _TodoService_CreateTodo_Handler},
			{MethodName: "GetTodo", Handler: _TodoService_GetTodo_Handler},
			{MethodName: "ListTodos", Handler: _TodoService_ListTodos_Handler},
			{MethodName: "UpdateTodo", Handler: _TodoService_UpdateTodo_Handler},
			{MethodName: "DeleteTodo", Handler: _TodoService_DeleteTodo_Handler},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "proto/todo.proto",
	}
	s.RegisterService(desc, srv)
}

func _TodoService_CreateTodo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTodoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).CreateTodo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: todoServiceName + "CreateTodo"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).CreateTodo(ctx, req.(*CreateTodoRequest))
	})
}

func _TodoService_GetTodo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTodoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).GetTodo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: todoServiceName + "GetTodo"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).GetTodo(ctx, req.(*GetTodoRequest))
	})
}

func _TodoService_ListTodos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTodosRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).ListTodos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: todoServiceName + "ListTodos"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).ListTodos(ctx, req.(*ListTodosRequest))
	})
}

func _TodoService_UpdateTodo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTodoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).UpdateTodo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: todoServiceName + "UpdateTodo"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).UpdateTodo(ctx, req.(*UpdateTodoRequest))
	})
}

func _TodoService_DeleteTodo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTodoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).DeleteTodo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: todoServiceName + "DeleteTodo"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).DeleteTodo(ctx, req.(*DeleteTodoRequest))
	})
}
