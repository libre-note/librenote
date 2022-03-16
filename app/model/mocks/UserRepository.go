// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"
	model "librenote/app/model"

	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: tx, user
func (_m *UserRepository) CreateUser(tx context.Context, user *model.User) error {
	ret := _m.Called(tx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) error); ok {
		r0 = rf(tx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUser provides a mock function with given fields: tx, id
func (_m *UserRepository) GetUser(tx context.Context, id int32) (model.User, error) {
	ret := _m.Called(tx, id)

	var r0 model.User
	if rf, ok := ret.Get(0).(func(context.Context, int32) model.User); ok {
		r0 = rf(tx, id)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int32) error); ok {
		r1 = rf(tx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByEmail provides a mock function with given fields: tx, email
func (_m *UserRepository) GetUserByEmail(tx context.Context, email string) (model.User, error) {
	ret := _m.Called(tx, email)

	var r0 model.User
	if rf, ok := ret.Get(0).(func(context.Context, string) model.User); ok {
		r0 = rf(tx, email)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(tx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: tx, user
func (_m *UserRepository) UpdateUser(tx context.Context, user *model.User) error {
	ret := _m.Called(tx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) error); ok {
		r0 = rf(tx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
