// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import session "gophr.v2/session"

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *Repository) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: ctx, id
func (_m *Repository) Find(ctx context.Context, id string) (*session.Session, error) {
	ret := _m.Called(ctx, id)

	var r0 *session.Session
	if rf, ok := ret.Get(0).(func(context.Context, string) *session.Session); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*session.Session)
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

// Save provides a mock function with given fields: ctx, s
func (_m *Repository) Save(ctx context.Context, s *session.Session) error {
	ret := _m.Called(ctx, s)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *session.Session) error); ok {
		r0 = rf(ctx, s)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}