// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "github.com/jordyf15/thullo-api/models"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// GetGoogleTokenInfo provides a mock function with given fields: token
func (_m *Repository) GetGoogleTokenInfo(token string) (*models.GoogleTokenInfo, error) {
	ret := _m.Called(token)

	var r0 *models.GoogleTokenInfo
	if rf, ok := ret.Get(0).(func(string) *models.GoogleTokenInfo); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GoogleTokenInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}