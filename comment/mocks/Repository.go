// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "github.com/jordyf15/thullo-api/models"
	mock "github.com/stretchr/testify/mock"

	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *Repository) Create(_a0 *models.Comment) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Comment) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteCommentByID provides a mock function with given fields: commentID
func (_m *Repository) DeleteCommentByID(commentID primitive.ObjectID) error {
	ret := _m.Called(commentID)

	var r0 error
	if rf, ok := ret.Get(0).(func(primitive.ObjectID) error); ok {
		r0 = rf(commentID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCommentByID provides a mock function with given fields: commentID
func (_m *Repository) GetCommentByID(commentID primitive.ObjectID) (*models.Comment, error) {
	ret := _m.Called(commentID)

	var r0 *models.Comment
	if rf, ok := ret.Get(0).(func(primitive.ObjectID) *models.Comment); ok {
		r0 = rf(commentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Comment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(primitive.ObjectID) error); ok {
		r1 = rf(commentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0
func (_m *Repository) Update(_a0 *models.Comment) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Comment) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
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
