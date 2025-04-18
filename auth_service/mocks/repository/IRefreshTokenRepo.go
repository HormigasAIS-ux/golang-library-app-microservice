// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import (
	model "auth_service/domain/model"

	mock "github.com/stretchr/testify/mock"
)

// IRefreshTokenRepo is an autogenerated mock type for the IRefreshTokenRepo type
type IRefreshTokenRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: refresh_token
func (_m *IRefreshTokenRepo) Create(refresh_token *model.RefreshToken) error {
	ret := _m.Called(refresh_token)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.RefreshToken) error); ok {
		r0 = rf(refresh_token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: id
func (_m *IRefreshTokenRepo) Delete(id string) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByToken provides a mock function with given fields: token
func (_m *IRefreshTokenRepo) GetByToken(token string) (*model.RefreshToken, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for GetByToken")
	}

	var r0 *model.RefreshToken
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.RefreshToken, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(string) *model.RefreshToken); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RefreshToken)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvalidateManyByUserUUID provides a mock function with given fields: userUUID
func (_m *IRefreshTokenRepo) InvalidateManyByUserUUID(userUUID string) error {
	ret := _m.Called(userUUID)

	if len(ret) == 0 {
		panic("no return value specified for InvalidateManyByUserUUID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userUUID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: refresh_token
func (_m *IRefreshTokenRepo) Update(refresh_token *model.RefreshToken) error {
	ret := _m.Called(refresh_token)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.RefreshToken) error); ok {
		r0 = rf(refresh_token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIRefreshTokenRepo creates a new instance of IRefreshTokenRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIRefreshTokenRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *IRefreshTokenRepo {
	mock := &IRefreshTokenRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
