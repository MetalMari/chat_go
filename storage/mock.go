package storage

import (
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Storage is an autogenerated mock type for the Storage type
type StorageMock struct {
	mock.Mock
}

// CreateMessage provides a mock function with given fields: m
func (_m *StorageMock) CreateMessage(m Message) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(Message) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateUser provides a mock function with given fields: u
func (_m *StorageMock) CreateUser(u User) error {
	ret := _m.Called(u)

	var r0 error
	if rf, ok := ret.Get(0).(func(User) error); ok {
		r0 = rf(u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteMessage provides a mock function with given fields: m
func (_m *StorageMock) DeleteMessage(m Message) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(Message) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetMessages provides a mock function with given fields: login
func (_m *StorageMock) GetMessages(login string) ([]Message, error) {
	ret := _m.Called(login)

	var r0 []Message
	if rf, ok := ret.Get(0).(func(string) []Message); ok {
		r0 = rf(login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUsers provides a mock function with given fields:
func (_m *StorageMock) GetUsers() ([]User, error) {
	ret := _m.Called()

	var r0 []User
	if rf, ok := ret.Get(0).(func() []User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStorage creates a new instance of Storage. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorage(t testing.TB) *StorageMock {
	mock := &StorageMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
