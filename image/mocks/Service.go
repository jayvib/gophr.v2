// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import image "gophr.v2/image"
import mock "github.com/stretchr/testify/mock"

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Find provides a mock function with given fields: ctx, id
func (_m *Service) Find(ctx context.Context, id string) (*image.Image, error) {
	ret := _m.Called(ctx, id)

	var r0 *image.Image
	if rf, ok := ret.Get(0).(func(context.Context, string) *image.Image); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*image.Image)
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

// FindAll provides a mock function with given fields: ctx, offset
func (_m *Service) FindAll(ctx context.Context, offset int) ([]*image.Image, error) {
	ret := _m.Called(ctx, offset)

	var r0 []*image.Image
	if rf, ok := ret.Get(0).(func(context.Context, int) []*image.Image); ok {
		r0 = rf(ctx, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*image.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllByUser provides a mock function with given fields: ctx, userId, offset
func (_m *Service) FindAllByUser(ctx context.Context, userId string, offset int) ([]*image.Image, error) {
	ret := _m.Called(ctx, userId, offset)

	var r0 []*image.Image
	if rf, ok := ret.Get(0).(func(context.Context, string, int) []*image.Image); ok {
		r0 = rf(ctx, userId, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*image.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int) error); ok {
		r1 = rf(ctx, userId, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, _a1
func (_m *Service) Save(ctx context.Context, _a1 *image.Image) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *image.Image) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
