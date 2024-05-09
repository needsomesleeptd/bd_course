// Code generated by MockGen. DO NOT EDIT.
// Source: bl/anotattionTypeService/anottationTypeRepo/anotattionTypeRepo.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	models "annotater/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIAnotattionTypeRepository is a mock of IAnotattionTypeRepository interface.
type MockIAnotattionTypeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIAnotattionTypeRepositoryMockRecorder
}

// MockIAnotattionTypeRepositoryMockRecorder is the mock recorder for MockIAnotattionTypeRepository.
type MockIAnotattionTypeRepositoryMockRecorder struct {
	mock *MockIAnotattionTypeRepository
}

// NewMockIAnotattionTypeRepository creates a new mock instance.
func NewMockIAnotattionTypeRepository(ctrl *gomock.Controller) *MockIAnotattionTypeRepository {
	mock := &MockIAnotattionTypeRepository{ctrl: ctrl}
	mock.recorder = &MockIAnotattionTypeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAnotattionTypeRepository) EXPECT() *MockIAnotattionTypeRepositoryMockRecorder {
	return m.recorder
}

// AddAnottationType mocks base method.
func (m *MockIAnotattionTypeRepository) AddAnottationType(markUp *models.MarkupType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAnottationType", markUp)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAnottationType indicates an expected call of AddAnottationType.
func (mr *MockIAnotattionTypeRepositoryMockRecorder) AddAnottationType(markUp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAnottationType", reflect.TypeOf((*MockIAnotattionTypeRepository)(nil).AddAnottationType), markUp)
}

// DeleteAnotattionType mocks base method.
func (m *MockIAnotattionTypeRepository) DeleteAnotattionType(id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAnotattionType", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAnotattionType indicates an expected call of DeleteAnotattionType.
func (mr *MockIAnotattionTypeRepositoryMockRecorder) DeleteAnotattionType(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAnotattionType", reflect.TypeOf((*MockIAnotattionTypeRepository)(nil).DeleteAnotattionType), id)
}

// GetAnottationTypeByID mocks base method.
func (m *MockIAnotattionTypeRepository) GetAnottationTypeByID(id uint64) (*models.MarkupType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnottationTypeByID", id)
	ret0, _ := ret[0].(*models.MarkupType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnottationTypeByID indicates an expected call of GetAnottationTypeByID.
func (mr *MockIAnotattionTypeRepositoryMockRecorder) GetAnottationTypeByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnottationTypeByID", reflect.TypeOf((*MockIAnotattionTypeRepository)(nil).GetAnottationTypeByID), id)
}
