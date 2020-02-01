// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	user "gophr/v2/gophr.api/user"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// GetAndComparePassword provides a mock function with given fields: ctx, username, password
func (_m *Service) GetAndComparePassword(ctx context.Context, username string, password string) (*user.User, error) {
	ret := _m.Called(ctx, username, password)

	var r0 *user.User
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *user.User); ok {
		r0 = rf(ctx, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByEmail provides a mock function with given fields: ctx, id
func (_m *Service) GetByEmail(ctx context.Context, id string) (*user.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *user.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *user.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *Service) GetByID(ctx context.Context, id string) (*user.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *user.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *user.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUsername provides a mock function with given fields: ctx, uname
func (_m *Service) GetByUsername(ctx context.Context, uname string) (*user.User, error) {
	ret := _m.Called(ctx, uname)

	var r0 *user.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *user.User); ok {
		r0 = rf(ctx, uname)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, uname)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, _a1
func (_m *Service) Save(ctx context.Context, _a1 *user.User) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *user.User) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
