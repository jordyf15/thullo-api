// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "github.com/jordyf15/thullo-api/models"
	mock "github.com/stretchr/testify/mock"

	os "os"

	sync "sync"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// UploadFile provides a mock function with given fields: respond, wg, currentImage, file, metadata
func (_m *Storage) UploadFile(respond chan<- error, wg *sync.WaitGroup, currentImage *models.Image, file *os.File, metadata map[string]string) {
	_m.Called(respond, wg, currentImage, file, metadata)
}

type mockConstructorTestingTNewStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorage(t mockConstructorTestingTNewStorage) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}