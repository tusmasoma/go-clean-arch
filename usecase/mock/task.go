// Code generated by MockGen. DO NOT EDIT.
// Source: task.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	entity "github.com/tusmasoma/go-clean-arch/entity"
	usecase "github.com/tusmasoma/go-clean-arch/usecase"
)

// MockTaskUseCase is a mock of TaskUseCase interface.
type MockTaskUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockTaskUseCaseMockRecorder
}

// MockTaskUseCaseMockRecorder is the mock recorder for MockTaskUseCase.
type MockTaskUseCaseMockRecorder struct {
	mock *MockTaskUseCase
}

// NewMockTaskUseCase creates a new mock instance.
func NewMockTaskUseCase(ctrl *gomock.Controller) *MockTaskUseCase {
	mock := &MockTaskUseCase{ctrl: ctrl}
	mock.recorder = &MockTaskUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskUseCase) EXPECT() *MockTaskUseCaseMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockTaskUseCase) CreateTask(ctx context.Context, params *usecase.CreateTaskParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockTaskUseCaseMockRecorder) CreateTask(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockTaskUseCase)(nil).CreateTask), ctx, params)
}

// DeleteTask mocks base method.
func (m *MockTaskUseCase) DeleteTask(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockTaskUseCaseMockRecorder) DeleteTask(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockTaskUseCase)(nil).DeleteTask), ctx, id)
}

// GetTask mocks base method.
func (m *MockTaskUseCase) GetTask(ctx context.Context, id string) (*entity.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTask", ctx, id)
	ret0, _ := ret[0].(*entity.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask.
func (mr *MockTaskUseCaseMockRecorder) GetTask(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockTaskUseCase)(nil).GetTask), ctx, id)
}

// UpdateTask mocks base method.
func (m *MockTaskUseCase) UpdateTask(ctx context.Context, params *usecase.UpdateTaskParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockTaskUseCaseMockRecorder) UpdateTask(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockTaskUseCase)(nil).UpdateTask), ctx, params)
}
