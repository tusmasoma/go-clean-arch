// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: task.proto

package gateway

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	TaskService_GetTask_FullMethodName    = "/task.TaskService/GetTask"
	TaskService_ListTasks_FullMethodName  = "/task.TaskService/ListTasks"
	TaskService_CreateTask_FullMethodName = "/task.TaskService/CreateTask"
	TaskService_UpdateTask_FullMethodName = "/task.TaskService/UpdateTask"
	TaskService_DeleteTask_FullMethodName = "/task.TaskService/DeleteTask"
)

// TaskServiceClient is the client API for TaskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TaskServiceClient interface {
	GetTask(ctx context.Context, in *GetTaskRequest, opts ...grpc.CallOption) (*GetTaskResponse, error)
	ListTasks(ctx context.Context, in *ListTasksRequest, opts ...grpc.CallOption) (*ListTasksResponse, error)
	CreateTask(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*CreateTaskResponse, error)
	UpdateTask(ctx context.Context, in *UpdateTaskRequest, opts ...grpc.CallOption) (*UpdateTaskResponse, error)
	DeleteTask(ctx context.Context, in *DeleteTaskRequest, opts ...grpc.CallOption) (*DeleteTaskResponse, error)
}

type taskServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskServiceClient(cc grpc.ClientConnInterface) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) GetTask(ctx context.Context, in *GetTaskRequest, opts ...grpc.CallOption) (*GetTaskResponse, error) {
	out := new(GetTaskResponse)
	err := c.cc.Invoke(ctx, TaskService_GetTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) ListTasks(ctx context.Context, in *ListTasksRequest, opts ...grpc.CallOption) (*ListTasksResponse, error) {
	out := new(ListTasksResponse)
	err := c.cc.Invoke(ctx, TaskService_ListTasks_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) CreateTask(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*CreateTaskResponse, error) {
	out := new(CreateTaskResponse)
	err := c.cc.Invoke(ctx, TaskService_CreateTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) UpdateTask(ctx context.Context, in *UpdateTaskRequest, opts ...grpc.CallOption) (*UpdateTaskResponse, error) {
	out := new(UpdateTaskResponse)
	err := c.cc.Invoke(ctx, TaskService_UpdateTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) DeleteTask(ctx context.Context, in *DeleteTaskRequest, opts ...grpc.CallOption) (*DeleteTaskResponse, error) {
	out := new(DeleteTaskResponse)
	err := c.cc.Invoke(ctx, TaskService_DeleteTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskServiceServer is the server API for TaskService service.
// All implementations must embed UnimplementedTaskServiceServer
// for forward compatibility
type TaskServiceServer interface {
	GetTask(context.Context, *GetTaskRequest) (*GetTaskResponse, error)
	ListTasks(context.Context, *ListTasksRequest) (*ListTasksResponse, error)
	CreateTask(context.Context, *CreateTaskRequest) (*CreateTaskResponse, error)
	UpdateTask(context.Context, *UpdateTaskRequest) (*UpdateTaskResponse, error)
	DeleteTask(context.Context, *DeleteTaskRequest) (*DeleteTaskResponse, error)
	mustEmbedUnimplementedTaskServiceServer()
}

// UnimplementedTaskServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTaskServiceServer struct{}

func (UnimplementedTaskServiceServer) GetTask(context.Context, *GetTaskRequest) (*GetTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTask not implemented")
}

func (UnimplementedTaskServiceServer) ListTasks(context.Context, *ListTasksRequest) (*ListTasksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTasks not implemented")
}

func (UnimplementedTaskServiceServer) CreateTask(context.Context, *CreateTaskRequest) (*CreateTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTask not implemented")
}

func (UnimplementedTaskServiceServer) UpdateTask(context.Context, *UpdateTaskRequest) (*UpdateTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTask not implemented")
}

func (UnimplementedTaskServiceServer) DeleteTask(context.Context, *DeleteTaskRequest) (*DeleteTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTask not implemented")
}
func (UnimplementedTaskServiceServer) mustEmbedUnimplementedTaskServiceServer() {}

// UnsafeTaskServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TaskServiceServer will
// result in compilation errors.
type UnsafeTaskServiceServer interface {
	mustEmbedUnimplementedTaskServiceServer()
}

func RegisterTaskServiceServer(s grpc.ServiceRegistrar, srv TaskServiceServer) {
	s.RegisterService(&TaskService_ServiceDesc, srv)
}

func _TaskService_GetTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).GetTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_GetTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).GetTask(ctx, req.(*GetTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_ListTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTasksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).ListTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_ListTasks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).ListTasks(ctx, req.(*ListTasksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_CreateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).CreateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_CreateTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).CreateTask(ctx, req.(*CreateTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_UpdateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).UpdateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_UpdateTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).UpdateTask(ctx, req.(*UpdateTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_DeleteTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).DeleteTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_DeleteTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).DeleteTask(ctx, req.(*DeleteTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TaskService_ServiceDesc is the grpc.ServiceDesc for TaskService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TaskService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "task.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTask",
			Handler:    _TaskService_GetTask_Handler,
		},
		{
			MethodName: "ListTasks",
			Handler:    _TaskService_ListTasks_Handler,
		},
		{
			MethodName: "CreateTask",
			Handler:    _TaskService_CreateTask_Handler,
		},
		{
			MethodName: "UpdateTask",
			Handler:    _TaskService_UpdateTask_Handler,
		},
		{
			MethodName: "DeleteTask",
			Handler:    _TaskService_DeleteTask_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "task.proto",
}
